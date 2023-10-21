package main

import (
	"encoding/json"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tcotav/boggle/middleware"
	"github.com/tcotav/boggle/solver"
	"github.com/tcotav/boggle/wordlookup"
)

type Server struct {
	Lookup *wordlookup.Dictionary
}

func main() {
	// TODO - set up config, currently everything is hard coded
	lookup, err := wordlookup.NewDictionaryFromFile("data/words_alpha.txt")
	if err != nil {
		// no dictionary means we want to stop the service
		// log.Panic vs. log.Fatal - panic will allow any defered to finish, fatal will not. it will os.Exit(1)
		log.Panic().Str("App", "bogglesvc").Str("Call", "main").Err(err).Msg("")
	}
	// set up our server
	s := Server{Lookup: lookup}

	ginMode := "release"
	gin.SetMode(ginMode) // need to set this to turn off default DEBUG noise from gin logging
	router := gin.New()
	// set our custom request logger here w/gin
	router.Use(middleware.RequestLogger())
	// lose this when we switched from gin.Default() to gin.New()
	// so add back, recovers from panics and 500s if there is one
	router.Use(gin.Recovery())

	// hack in some dummy handler to test basics
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// here's our boggle
	router.POST("/boggle", s.PostBoggle)
	log.Info().Str("App", "bogglesvc").Str("Call", "main").Msg("Starting server on port 8080")
	router.Run(":8080")
}

// BoggleRequest is the request body for the POST /boggle endpoint
type BoggleRequest struct {
	Matrix [][]rune `json:"matrix"`
}

// utility function to make sure our 4x4 matrix has valid characters
func validateInput(matrix [][]rune) bool {
	// make sure we have a 4x4 matrix
	if len(matrix) != 4 {
		return false
	}

	// all the rows are length 4
	for i := 0; i < len(matrix); i++ {
		if len(matrix[i]) != 4 {
			return false
		}
	}

	// then test the runes to ensure they're valid
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			// check for valid characters - matches upper and lower A-Z
			if !unicode.IsLetter(matrix[i][j]) {
				return false
			}
		}
	}
	return true
}

func (s *Server) PostBoggle(c *gin.Context) {
	// get the request body
	var boggleRequest BoggleRequest
	// Decode the JSON request
	// TODO -- this unmarshal fails w/ "cannot unmarshal string into Go struct field BoggleRequest.matrix of type int32"
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&boggleRequest); err != nil {
		log.Error().Str("App", "bogglesvc").Str("Call", "postBoggle").Err(err).Msg("Error parsing request body")
		c.JSON(400, gin.H{"message": "Error parsing request body"})
		return
	}
	/*  variation tried -- same error as above but uses a gin helper method

	//err := c.BindJSON(&boggleRequest)
	if err != nil {
		log.Error().Str("App", "bogglesvc").Str("Call", "postBoggle").Err(err).Msg("Error parsing request body")
		c.JSON(400, gin.H{"message": "Error parsing request body"})
		return
	}*/
	if !validateInput(boggleRequest.Matrix) {
		log.Error().Str("App", "bogglesvc").Str("Call", "postBoggle").Msg("Invalid input")
		c.JSON(400, gin.H{"message": "Invalid input"})
		return
	}

	// get the matrix
	matrix := boggleRequest.Matrix

	// process the matrix
	b := solver.NewBoggleSolver(matrix, s.Lookup)
	words := b.FindWords()

	// return the words
	c.JSON(200, gin.H{"words": words})
}
