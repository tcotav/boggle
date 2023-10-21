package solver

import (
	"github.com/tcotav/boggle/wordlookup"
)

// BoggleSolver represents the Boggle game solver.
type BoggleSolver struct {
	board      [][]rune
	dictionary *wordlookup.Dictionary
	visited    [][]bool
}

// NewBoggleSolver creates a new BoggleSolver instance.
func NewBoggleSolver(board [][]rune, dictionary *wordlookup.Dictionary) *BoggleSolver {
	return &BoggleSolver{
		board:      board,
		dictionary: dictionary,
		visited:    make([][]bool, len(board)),
	}
}

// FindWords finds all valid words on the Boggle board.
func (bs *BoggleSolver) FindWords() []string {
	var result []string
	for i := range bs.board {
		bs.visited[i] = make([]bool, len(bs.board[i]))
	}

	for i := range bs.board {
		for j := range bs.board[i] {
			// initialize currentword to empty string as we start with no letters
			bs.dfs(i, j, "", &result)
		}
	}
	return result
}

func (bs *BoggleSolver) isWord(word string) bool {
	return bs.dictionary.Contains(word)
}

func (bs *BoggleSolver) dfs(x, y int, currentWord string, result *[]string) {
	bs.visited[x][y] = true

	// add current letter to currentWord
	currentWord += string(bs.board[x][y])

	// check if currentWord is a valid word
	if bs.isWord(currentWord) {
		// TODO: we should check for dupes or change from slice to map
		// equiv of a python set()
		// map[string]bool or map[string]struct{} bc empty struct is 0 bytes
		*result = append(*result, currentWord)
	}

	// dx and dy are the 8 directions we can move in the matrix
	// corresponding to the coordinates around our current position
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	for k := 0; k < 8; k++ {
		newX, newY := x+dx[k], y+dy[k]
		if bs.isValid(newX, newY) && !bs.visited[newX][newY] {
			bs.dfs(newX, newY, currentWord, result)
		}
	}

	bs.visited[x][y] = false
}

func (bs *BoggleSolver) isValid(x, y int) bool {
	return x >= 0 && x < len(bs.board) && y >= 0 && y < len(bs.board[x])
}
