package sheet

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSearchNoMatch ensures that the expected output is returned when no
// matches are found
func TestSearchNoMatch(t *testing.T) {

	// mock a cheatsheet
	sheet := Sheet{
		Text: "The quick brown fox\njumped over\nthe lazy dog.",
	}

	// compile the search regex
	reg, err := regexp.Compile("(?i)foo")
	assert.NoError(t, err, "failed to compile regex")

	// search the sheet
	matches := sheet.Search(reg)
	assert.Equal(t, "", matches, "failure: expected no matches")
}

// TestSearchSingleMatch asserts that the expected output is returned
// when a single match is returned
func TestSearchSingleMatch(t *testing.T) {

	// mock a cheatsheet
	sheet := Sheet{
		Text: "The quick brown fox\njumped over\n\nthe lazy dog.",
	}

	// compile the search regex
	reg, err := regexp.Compile("(?i)fox")
	assert.NoError(t, err, "failed to compile regex")

	// search the sheet
	matches := sheet.Search(reg)
	want := "The quick brown fox\njumped over"
	assert.Equal(t, want, matches, "failed to return expected matches")
}

// TestSearchMultiMatch asserts that the expected output is returned
// when a multiple matches are returned
func TestSearchMultiMatch(t *testing.T) {

	// mock a cheatsheet
	sheet := Sheet{
		Text: "The quick brown fox\n\njumped over\n\nthe lazy dog.",
	}

	// compile the search regex
	reg, err := regexp.Compile("(?i)the")
	assert.NoError(t, err, "failed to compile regex")

	// search the sheet
	matches := sheet.Search(reg)

	// specify the expected results
	want := "The quick brown fox\n\nthe lazy dog."
	assert.Equal(t, want, matches, "failed to return expected matches")
}
