package wordlookup

import (
	"bufio"
	"os"
	"strings"
	"unicode"

	"github.com/dghubble/trie"
	"github.com/rs/zerolog/log"
)

// Dictionary is a trie of words
// chose trie because prefix matching is a common use case
// searches are O(k) where k is the length of the word
type Dictionary struct {
	trie *trie.RuneTrie
}

func NewDictionary() *Dictionary {
	return &Dictionary{trie.NewRuneTrie()}
}

func NewDictionaryFromFile(path string) (*Dictionary, error) {
	d := Dictionary{trie.NewRuneTrie()}
	err := d.LoadFromFile(path)
	return &d, err
}

// Add adds a word to the dictionary
func (d *Dictionary) Add(word string) {
	d.trie.Put(word, true)
}

func (d Dictionary) Contains(word string) bool {
	val := d.trie.Get(word)
	return val != nil
}

// isAlphaWord returns true if the word is all letters
// utility function
func isAlphaWord(word string) bool {
	for _, r := range word {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// LoadFromFile loads a dictionary from a file
// also does some data cleaning that arguably should be done elsewhere
// as a preprocessing step prior to using it in the service as this delays startup
func (d *Dictionary) LoadFromFile(path string) error {
	// open file for reading
	readfile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer readfile.Close()

	// read line by line
	fileScanner := bufio.NewScanner(readfile)

	// set the split function for fileScanner, specifically here to split on new lines
	fileScanner.Split(bufio.ScanLines)
	count := 0
	for fileScanner.Scan() {
		// add each line to the dictionary
		word := fileScanner.Text()
		// apply the boggle rules
		// all words must be at least 3 characters long
		if len(word) < 3 {
			continue
		}
		// word list we used has some words with apostrophes, dashes, and numbers
		// we only want to process words that are all letters
		if isAlphaWord(word) {
			// then lowercase the word as our word list has a mix
			word = strings.ToLower(word)
			count += 1
			// and add it to the dictionary
			d.Add(word)
		}
	}
	log.Info().Int("count", count).Msg("Loaded dictionary")
	return nil
}
