.PHONY: all test clean glide fast release build install script/build script/install install script/test-all man-pages fmt deps-ensure build version

################################################################################################
## local - runtime

ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  RUNTIME_OS_SLUG			:= osx
else
  RUNTIME_OS_SLUG 			:= nix
endif
RUNTIME_OS_VERSION 			?= $(shell uname -r)
RUNTIME_OS_ARCH 			?= $(shell uname -m)
RUNTIME_OS_INFO 			?= $(shell uname -a)
RUNTIME_OS_NAME 			?= $(shell uname -s)

################################################################################################
## local - build

## program
PROG_NAME 					:= hub
PROG_NAME_SUFFIX 			:= 
PROG_SRCS 					:= $(shell git ls-files '*.go' | grep -v '^vendor/')
PROG_BINS 					:= $(shell ls -1 $(CURDIR)/cmd)

## local build
BIN_PREFIX_DIR 				:= ./bin
BIN_BASE_NAME 				:= $(PROG_NAME_SUFFIX)$(PROG_NAME)
BIN_FILE_PATH 				:= $(BIN_PREFIX_DIR)/$(BIN_BASE_NAME)

## local dist
DIST_PREFIX_DIR 			:= ./dist
DIST_BASE_NAME 				:= $(PROG_NAME_SUFFIX)$(PROG_NAME)
DIST_FILE_PATH 				:= $(DIST_PREFIX_DIR)/$(DIST_BASE_NAME)
DIST_ARCHS 					:= "linux darwin"
DIST_OSS 					:= "amd64"

################################################################################################
## docker

#### build
DOCKER_PREFIX_DIR 			:= ./docker
DOCKER_BIN_FILE_PATH 		:= $(DOCKER_PREFIX_DIR)/$(BIN_BASE_NAME)

#### image
DOCKER_IMAGE_OWNER 			:= sniperkit
DOCKER_IMAGE_BASENAME 		:= hub
DOCKER_IMAGE_TAG 			:= 3.7-alpine
DOCKER_IMAGE 				:= $(DOCKER_IMAGE_OWNER)/$(DOCKER_IMAGE_BASENAME):$(DOCKER_IMAGE_TAG)
DOCKER_MULTI_STAGE_IMAGE 	:= $(DOCKER_IMAGE_OWNER)/$(DOCKER_IMAGE_BASENAME)-multi:$(DOCKER_IMAGE_TAG)

################################################################################################
## version

# vcs
REPO_VCS 					:= github.com
REPO_OWNER 					:= sniperkit
REPO_NAME 					:= hub
REPO_BRANCH_DEV 			:= sniperkit
REPO_IS_FORK 				:= true

REPO_URI 					:= $(REPO_VCS)/$(REPO_OWNER)/$(REPO_NAME)
REPO_REMOTE_ORIGIN_URL 		:= $(shell git config --get remote.origin.url)
REPO_BRANCH_CURRENT 		:= $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))
REPO_BRANCH_EXPECTED 		?= $(REPO_BRANCH_DEV)

# vcs - orign
REPO_ORIG_VCS 				:= github.com
REPO_ORIG_OWNER 			:= github
REPO_ORIG_NAME 				:= hub
REPO_ORIG_BRANCH 			:= master

REPO_ORIG_URI 				:= $(REPO_ORIG_VCS)/$(REPO_ORIG_OWNER)/$(REPO_ORIG_NAME)
REPO_ORIG_REMOTE_URL 		:= git://$(REPO_ORIG_URI).git

ifeq ($(REPO_BRANCH_CURRENT), $(REPO_BRANCH_EXPECTED))
  REPO_BRANCH_MISMATCH 		:= false
else
  REPO_BRANCH_MISMATCH 		:= true
endif

VERSION ?= $(shell git describe --tags)
VERSION_INCODE = $(shell perl -ne '/^var version.*"([^"]+)".*$$/ && print "v$$1\n"' main.go)
VERSION_INCHANGELOG = $(shell perl -ne '/^\# Release (\d+(\.\d+)+) / && print "$$1\n"' CHANGELOG.md | head -n1)

#### vcs - commit 
COMMIT_ID   				?= $(shell git describe --tags --always --dirty=-dev)
COMMIT_UNIX 				?= $(shell git show -s --format=%ct HEAD)
COMMIT_HASH 				?= $(shell git rev-parse HEAD)

#### semantic version 
BUILD_COUNT 				?= $(shell git rev-list --count HEAD)
BUILD_UNIX  				?= $(shell date +%s)
BUILD_VERSION 				:= $(shell cat $(CURDIR)/VERSION)
BUILD_TIME 					:= $(shell date)

################################################################################################
## golang

GO15VENDOREXPERIMENT	= 1
BUILD_LDFLAGS 			= 	\
							-X '$(REPO_URI)/pkg/version.Version=$(VERSION)' \
							-X '$(REPO_URI)/pkg/version.CranchName=$(REPO_BRANCH)' \
							-X '$(REPO_URI)/pkg/version.CommitHash=$(COMMIT_HASH)' \
							-X '$(REPO_URI)/pkg/version.CommitID=$(COMMIT_ID)' \
							-X '$(REPO_URI)/pkg/version.CommitUnix=$(COMMIT_UNIX)' \
							-X '$(REPO_URI)/pkg/version.BuildVersion=$(BUILD_VERSION)' \
							-X '$(REPO_URI)/pkg/version.BuildCount=$(BUILD_COUNT)' \
							-X '$(REPO_URI)/pkg/version.BuildUnix=$(BUILD_UNIX)'

SOURCES 				= $(shell shared/scripts/build files)
SOURCES_FMT 			= $(shell shared/scripts/build files | cut -d/ -f1-2 | sort -u)

## gox - cross-build
GOX_OSARCH_LIST 		:= darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm freebsd/386 freebsd/amd64 freebsd/arm netbsd/386 netbsd/amd64 netbsd/arm openbsd/386 openbsd/amd64 openbsd/arm windows/386 windows/amd64
GOX_CMD_LIST 			:= $(shell ls -1 $(CURDIR)/cmd)

################################################################################################
## makefile
INFO_BREAKLINE := "\n"
INFO_HEADER := "$(INFO_BREAKLINE)------------------------------------------------------------------------------------------"

INFO_FOOTER := "$(INFO_BREAKLINE)------------------------------$(INFO_BREAKLINE)"

default: help

# all: deps-ensure test build install version

all: deps test build install version dist ## Trigger targets for generating a new release: deps, test, build, install, version and dist targets

info: clear info-runtime info-vcs info-docker info-footer ## Print all Makefile related variables

.PHONY: commit
commit: ensure-branch-dev
	@git add .
	@git commit -am "commit changes for..." 2>/dev/null; true

.PHONY: ensure-branch-dev
ensure-branch-dev: info-vcs
	@if [ $(REPO_BRANCH_CURRENT) == "master" ];then \
		git branch $(REPO_BRANCH_EXPECTED) 2>/dev/null; true; \
		git checkout $(REPO_BRANCH_EXPECTED) 2>/dev/null; true; \
	fi

.PHONY: update-fork-master
update-fork-master: commit
	@echo "current branch: $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))"
	@if [ $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)) != $(REPO_BRANCH_EXPECTED) ]; then \
		git checkout $(REPO_BRANCH_EXPECTED); \
		echo "chekout branch: $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))" ; \
	fi
	git remote add upstream $(REPO_ORIG_REMOTE_URL) 2>/dev/null; true
	git fetch upstream 2>/dev/null; true
	git pull upstream $(REPO_ORIG_BRANCH)
	git checkout $(REPO_BRANCH_CURRENT)
# git merge -Xours origin/$(REPO_ORIG_BRANCH)

.PHONY: update-local-master
update-local-master: commit
	git remote add upstream $(REPO_ORIG_REMOTE_URL) 2>/dev/null; true
	git fetch upstream 2>/dev/null; true
	git checkout $(REPO_ORIG_BRANCH)
	git pull upstream $(REPO_ORIG_BRANCH)
	git merge -Xours origin/$(REPO_ORIG_BRANCH)
	git checkout $(REPO_BRANCH_EXPECTED)

# git reset --hard HEAD

clear: ## Clear terminal screen 
	@clear

info-header:
	@echo ""
	@echo "------------------------------"

info-footer:
	@echo "$(INFO_FOOTER)"

.PHONY: info-runtime
info-runtime:  ## Print local runtime env variables
	@echo "$(INFO_HEADER)"
	@echo "Runtime:"
	@echo " - RUNTIME_OS_NAME: $(RUNTIME_OS_NAME)"
	@echo " - RUNTIME_ARCH: $(RUNTIME_OS_ARCH)"
	@echo " - RUNTIME_OS_VERSION: $(RUNTIME_OS_VERSION)"
	@echo " - RUNTIME_OS_SLUG: $(RUNTIME_OS_SLUG)"
	@echo " - RUNTIME_OS_INFO: $(RUNTIME_OS_INFO)"

.PHONY: info-vcs
info-vcs:  ## Print source-control related variables
	@echo "$(INFO_HEADER)"
	@echo "Source-Control:"
	@echo " - REPO_URI: $(REPO_URI)"
	@echo " - REPO_REMOTE_ORIGIN_URL: $(REPO_REMOTE_ORIGIN_URL)"
	@echo " - REPO_BRANCH_CURRENT: $(REPO_BRANCH_CURRENT)"
	@echo " - REPO_BRANCH_EXPECTED: $(REPO_BRANCH_EXPECTED)"
	@if [ $(REPO_BRANCH_MISMATCH) == "true" ]; then \
		echo " - !!! Warning !!! REPO_BRANCH_MISMATCH $(REPO_BRANCH_MISMATCH)"; \
	 	echo " - REPO_BRANCH_CURRENT=$(REPO_BRANCH_CURRENT) not equal to REPO_BRANCH_EXPECTED=$(REPO_BRANCH_EXPECTED)"; \
		echo ""; \
	fi
	@if [ $(REPO_IS_FORK) == "true" ]; then \
		echo " - REPO_ORIG_VCS: $(REPO_ORIG_VCS)" ; \
		echo " - REPO_ORIG_OWNER: $(REPO_ORIG_OWNER)" ; \
		echo " - REPO_ORIG_NAME: $(REPO_ORIG_NAME)" ; \
		echo " - REPO_ORIG_BRANCH: $(REPO_ORIG_BRANCH)" ; \
		echo " - REPO_ORIG_URI: $(REPO_ORIG_URI)" ; \
		echo " - REPO_ORIG_REMOTE_URL: $(REPO_ORIG_REMOTE_URL)" ; \
	fi
	@echo " - COMMIT_ID: $(COMMIT_ID)"
	@echo " - COMMIT_UNIX: $(COMMIT_UNIX)"
	@echo " - COMMIT_HASH: $(COMMIT_HASH)"
	@echo " - BUILD_COUNT: $(BUILD_COUNT)"
	@echo " - BUILD_UNIX: $(BUILD_UNIX)"
	@echo " - BUILD_VERSION: $(BUILD_VERSION)"
	@echo " - BUILD_TIME: $(BUILD_TIME)"

.PHONY: ls-cmd
ls-cmd:
	@echo "list of binaries available: \n"
	@echo "$(PROG_BINS)"

.PHONY: build
build: ## Build binary for local operating system 
	@$(foreach cmd, $(shell ls -1c ./cmd), $(call go_build_cmd,$(cmd)))

define go_build_cmd
    echo "## Building binary: $$(basename $1)" ;
    @go build -v -ldflags "$(BUILD_LDFLAGS)" -o ./bin/$$(basename $1) ./cmd/$(1)/*.go ;
    @./bin/$$(basename $1) version ;
    echo "" ;
endef

.PHONY: install
install: ## Install binary in your GOBIN path
	@go install -ldflags "$(BUILD_LDFLAGS)" $(REPO_URI)/cmd/...

.PHONY: dist
dist: ## Build all dist binaries for linux, darwin in amd64 arch.
	@gox -ldflags="$(BUILD_LDFLAGS)" -osarch="$(GOX_OSARCH_LIST)" -output="$(DIST_PREFIX_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}" $(REPO_URI)/cmd/...

.PHONY: version
version-current: ## Check current version of command build
	@which $(BIN_BASE_NAME)
	@$(BIN_BASE_NAME) --version

.PHONY: clean
clean: ## Clean previous build outputs 
	@go clean
	@git clean -fdx bin shared/man
	@rm -fr ./bin/$(BIN_FILE_PATH)
	@rm -fr ./dist/$(BIN_FILE_PATH)*

release: $(PROG_NAME) ## Push a new release version to the remote repository
	@git tag -a `cat VERSION`
	@git push origin `cat VERSION`

cover: ## Execute coverage tests
	@rm -rf *.coverprofile
	@go test -coverprofile=$(PROG_NAME).coverprofile ./pkg/...
	@gover
	@go tool cover -html=$(PROG_NAME).coverprofile ./pkg/...

deps: deps-ensure deps-dev deps-test ## Ensure all required dependencies and helpers

deps-all: deps-create deps-ensure deps-dev deps-test ## Re-create all dependencies list and ensure all locally

deps-create: ## Create program's dependencies list
	@rm -f glide.*
	@rm -f *Gopkg*
	@yes no | glide create

deps-ensure: ## Ensure locally all external dependencies required (package manager: glide)
	@glide install --strip-vendor

DEPS_DEV := github.com/sniperkit/crane/cmd/crane \
			github.com/sniperkit/gox/cmd/gox

.PHONY: deps-dev
deps-dev: ## Install required build helpers in GOBIN 
	@$(foreach cmd, $(DEPS_DEV), $(call go_install_cmd,$(cmd)))

DEPS_TEST := 	github.com/go-playground/overalls \
			 	github.com/mattn/goveralls \
			 	golang.org/x/tools/cmd/cover \
			 	github.com/alexkohler/prealloc \
			 	github.com/FiloSottile/vendorcheck \
			 	github.com/golang/dep/cmd/dep \
			 	github.com/golang/lint/golint \
			 	github.com/kisielk/errcheck \
			 	github.com/mdempsky/unconvert \
			 	github.com/opennota/check/... \
			 	honnef.co/go/tools/cmd/... \
			 	mvdan.cc/interfacer \
				github.com/dominikh/go-tools/...

.PHONY: deps-test
deps-test:  ## Install required program testing an ci helpers in GOBIN
	@$(foreach cmd, $(DEPS_TEST), $(call go_install_cmd,$(cmd)))

define go_install_cmd
    echo " [INSTALL] Package: $(1)";
	@errors=$$(go get -u $(1)); if [ "$${errors}" != "" ]; then echo \"$${errors}\"; exit 1; fi;
endef

.PHONY: lint
lint: ## Lint program's source code
	@errors=$$(gofmt -l .); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi
	@errors=$$(glide novendor | xargs -n 1 golint -min_confidence=0.3); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

.PHONY: vet
vet: ## Vet program's source code
	@go vet $$(glide novendor)

.PHONY: errcheck
errcheck: ## Check for errors
	@errcheck $(PACKAGES)

.PHONY: interfacer
interfacer: ## Suggest interface types
	@interfacer $(PACKAGES)

.PHONY: aligncheck
aligncheck: ## Find inefficiently packed structs
	@aligncheck $(PACKAGES)

.PHONY: structcheck
structcheck: ## Find unused struct fields
	@structcheck $(PACKAGES)

.PHONY: varcheck
varcheck: ## Find unused global variables and constants
	@varcheck $(PACKAGES)

.PHONY: unconvert
unconvert: ## Remove unnecessary type conversions from Go source
	@unconvert -v $(PACKAGES)

.PHONY: gosimple
gosimple: ## Suggest code simplifications
	@gosimple $(PACKAGES)

.PHONY: staticcheck
staticcheck: ## Execute a ton of static analysis checks
	@staticcheck $(PACKAGES)

.PHONY: unused
unused: ## Find for unused constants, variables, functions and types. 
	@unused $(PACKAGES)

.PHONY: vendorcheck
vendorcheck: ## Check that all Go dependencies are properly vendored
	@vendorcheck $(PACKAGES)
	@vendorcheck -u $(PACKAGES)

.PHONY: prealloc
prealloc: ## Find slice declarations that could potentially be preallocated.
	@prealloc $(PACKAGES)

.PHONY: test
test: ## Execute cover tests on program's sources
	@go test -cover $(PACKAGES)

coverage: ## Execute all coverage tests
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	@go tool cover -html=coverage-all.out

.PHONY: info-docker
info-docker:
	@echo "$(INFO_HEADER)"
	@echo "Docker:"
	@echo " - DOCKER_PREFIX_DIR: $(DOCKER_PREFIX_DIR)"
	@echo " - DOCKER_BIN_FILE_PATH: $(DOCKER_BIN_FILE_PATH)"
	@echo " - DOCKER_IMAGE_OWNER: $(DOCKER_IMAGE_OWNER)"
	@echo " - DOCKER_IMAGE_TAG: $(DOCKER_IMAGE_TAG)"
	@echo " - DOCKER_IMAGE: $(DOCKER_IMAGE)"
	@echo " - DOCKER_MULTI_STAGE_IMAGE: $(DOCKER_MULTI_STAGE_IMAGE)"

docker: docker-build  # docker-tag docker-commit docker-push ## Generate, tag and push a new docker image for this program.

docker-quick: docker-build docker-run ## Build and run quickly a docker container for this program

#docker-multistage: ## Build docker multi-stage container
#	@cd $(DOCKER_PREFIX_DIR) 
#	@docker build --force-rm -t $(DOCKER_MULTI_STAGE_IMAGE) --no-cache -f $(CURDIR)/docker/dockerfile-multi-stage-alpine3.7 .

.PHONY: docker-build
docker-build: ## Build docker container
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(BUILD_LDFLAGS)" -o $(DOCKER_BIN_FILE_PATH)-linux -v *.go
	@cd $(DOCKER_PREFIX_DIR) && docker build --force-rm -t $(DOCKER_IMAGE) --no-cache -f dockerfile-alpine3.7 .

.PHONY: docker-run
docker-run: ## Run docker container locally
	@docker run -ti --rm $(DOCKER_IMAGE)

.PHONY: docker-info
docker-info: ## Get docker client info and env variables
	@echo "'docker-info' is not implemented yet..."

docker-summary: ## Get docker image(s)/container(s) summary 
	@echo "'docker-summary' is not implemented yet..."

docker-commit: ## Commit latest docker image for this program
	@echo "'docker-commit' is not implemented yet..."

docker-tag: ## Tag latest docker image for this program
	@echo "'docker-push' is not implemented yet..."

docker-push: ## Push docker image to image registry
	@echo "'docker-push' is not implemented yet..."


help: ## Display the list of available targets.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

#generate-webui: build-webui ## Build the web front-end
#	if [ ! -d "static" ]; then \
#		mkdir -p static; \
#		docker run --rm -v "$$PWD/static":'/src/static' dtk-webui npm run build; \
#		echo 'For more informations show `webui/readme.md`' > $$PWD/static/DONT-EDIT-FILES-IN-THIS-DIRECTORY.md; \
#	fi

################################################################################################
## hub - help
MIN_COVERAGE = 89.4

HELP_CMD = \
	shared/man/man1/hub-alias.1 \
	shared/man/man1/hub-browse.1 \
	shared/man/man1/hub-ci-status.1 \
	shared/man/man1/hub-compare.1 \
	shared/man/man1/hub-create.1 \
	shared/man/man1/hub-delete.1 \
	shared/man/man1/hub-fork.1 \
	shared/man/man1/hub-pr.1 \
	shared/man/man1/hub-pull-request.1 \
	shared/man/man1/hub-release.1 \
	shared/man/man1/hub-issue.1 \
	shared/man/man1/hub-sync.1 \

HELP_EXT = \
	shared/man/man1/hub-am.1 \
	shared/man/man1/hub-apply.1 \
	shared/man/man1/hub-checkout.1 \
	shared/man/man1/hub-cherry-pick.1 \
	shared/man/man1/hub-clone.1 \
	shared/man/man1/hub-fetch.1 \
	shared/man/man1/hub-help.1 \
	shared/man/man1/hub-init.1 \
	shared/man/man1/hub-merge.1 \
	shared/man/man1/hub-push.1 \
	shared/man/man1/hub-remote.1 \
	shared/man/man1/hub-submodule.1 \

HELP_ALL = shared/man/man1/hub.1 $(HELP_CMD) $(HELP_EXT)

TEXT_WIDTH = 87

########################################################################################################
### Legacy

script/build/hub: $(SOURCES)
	@shared/scripts/build -o $@

script/test:
	@shared/scripts/build test

script/test-all: bin/cucumber
ifdef CI
	@shared/scripts/test --coverage $(MIN_COVERAGE)
else
	@shared/scripts/test
endif

bin/ronn bin/cucumber:
	@shared/scripts/bootstrap

fmt:
	go fmt ./...

man-pages: $(HELP_ALL:=.ronn) $(HELP_ALL) $(HELP_ALL:=.txt)

%.txt: %.ronn
	@groff -Wall -mtty-char -mandoc -Tutf8 -rLL=$(TEXT_WIDTH)n $< | col -b >$@

%.1: %.1.ronn bin/ronn
	@bin/ronn --organization=GITHUB --manual="Hub Manual" shared/man/man1/*.ronn

%.1.ronn: bin/hub
	@bin/hub help $(*F) --plain-text | shared/scripts/format-ronn $(*F) $@

shared/man/man1/hub.1.ronn:
	true

script/install: bin/hub man-pages
	@bash < shared/scripts/install.sh
