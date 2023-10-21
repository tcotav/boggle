package solver

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/tcotav/boggle/wordlookup"
)

// https://stackoverflow.com/a/70050794
// allow our tests to run from the root of the project
// get the current executing and add ..
// and find the data files
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestBoggleSolver(t *testing.T) {
	// todo add more board variations
	boards := [][][]rune{
		{
			{'a', 'n', 't', 'z'},
			{'z', 'z', 'z', 'z'},
			{'z', 'z', 'z', 'z'},
			{'z', 'z', 'z', 'z'},
		},
		{
			{'l', 'z', 'z', 'z'},
			{'z', 'o', 's', 'z'},
			{'z', 't', 'e', 'z'},
			{'z', 'u', 's', 'z'},
		},
	}
	dict := wordlookup.NewDictionary()
	dict.Add("ant")
	dict.Add("lotus")
	dict.Add("lotuses")
	dict.Add("lots")
	dict.Add("lost")

	solution_count := []int{1, 4}
	for i, board := range boards {
		bs := NewBoggleSolver(board, dict)
		// this returns dupes because we're not checking for them
		// there's a todo
		words := bs.FindWords()
		// TODO remove logging from tests
		log.Info().Msgf("%v", words)

		// we fail because dupes -- need to remove dupes
		if len(words) != solution_count[i] {
			t.Errorf("Expected %d word, got %d", solution_count[i], len(words))
		}
	}
}
