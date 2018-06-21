// !build lib
package ghs

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type parseTestReulst struct {
	version  bool
	exitCode int
	sOpt     *SearchOpt
}

func TestOption_Parse(t *testing.T) {
	assert := func(result interface{}, want interface{}) {
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Returned %+v, want %+v", result, want)
		}
	}

	// want :exit, exitCode , sOpt, url, token
	assert(testParse("ghs -v"), &parseTestReulst{true, ExitCodeOK, nil})
	assert(testParse("ghs -h"), &parseTestReulst{false, ExitCodeError, nil})

	defaultOpt := SearchOpt{
		sort:    "best match",
		order:   "desc",
		query:   "SEARCH_WORD",
		max:     100,
		perPage: 100,
		baseURL: nil,
		token:   "",
	}

	// normal query test
	wantOpt := defaultOpt
	assert(testParse("ghs SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.sort = "stars"
	assert(testParse("ghs -s stars SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.order = "asc"
	assert(testParse("ghs -o asc SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.max = 1000
	assert(testParse("ghs -m 1000 SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.baseURL, _ = url.Parse("http://test.exmaple/")
	assert(testParse("ghs -e http://test.exmaple/ SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.token = "abcdefg"
	assert(testParse("ghs -t abcdefg SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "in:name SEARCH_WORD"
	assert(testParse("ghs -f name SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "user:sona-tar"
	assert(testParse("ghs -u sona-tar"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "repo:sona-tar/ghs"
	assert(testParse("ghs -r sona-tar/ghs"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "language:golang"
	assert(testParse("ghs -l golang"), &parseTestReulst{false, ExitCodeOK, &wantOpt})

	// no args test
	assert(testParse("ghs"), &parseTestReulst{false, ExitCodeError, nil})
	assert(testParse("ghs -o asc"), &parseTestReulst{false, ExitCodeError, nil})
	wantOpt = defaultOpt
	wantOpt.query = "user:sona-tar"
	assert(testParse("ghs -u sona-tar"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "repo:sona-tar/ghs"
	assert(testParse("ghs -r sona-tar/ghs"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "language:golang"
	assert(testParse("ghs -l golang"), &parseTestReulst{false, ExitCodeOK, &wantOpt})

	// invalid option value test
	assert(testParse("ghs -m 1001 SEARCH_WORD"), &parseTestReulst{false, ExitCodeError, nil})
	assert(testParse("ghs -m 0 SEARCH_WORD"), &parseTestReulst{false, ExitCodeError, nil})
	assert(testParse("ghs -e : SEARCH_WORD"), &parseTestReulst{false, ExitCodeError, nil})
}

func testParse(args_string string) *parseTestReulst {
	args := strings.Split(args_string, " ")[1:]
	flags, _ := NewFlags(args)
	version, exitCode, sOpt := flags.ParseOption()
	return &parseTestReulst{version, exitCode, sOpt}
}
