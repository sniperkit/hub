package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

func TestDirIsNotEmpty(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)
	ioutil.TempFile(dir, "gh-utils-test-")

	assert.T(t, !isEmptyDir(dir))
}

func TestDirIsEmpty(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	assert.T(t, isEmptyDir(dir))
}

func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "gh-utils-test-")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}