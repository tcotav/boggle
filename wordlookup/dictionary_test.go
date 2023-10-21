package wordlookup

import (
	"os"
	"path"
	"runtime"
	"testing"
)

// https://stackoverflow.com/a/70050794
// allow our tests to run from the root of the project
// and find the data files
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestNewDictionaryFromFile(t *testing.T) {
	lookup, err := NewDictionaryFromFile("data/words_alpha.txt")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if lookup == nil {
		t.Errorf("Expected a dictionary, got nil")
	}
}

// TODO add more tests
// check that we can do a lookup for various words -- both positive and negative
