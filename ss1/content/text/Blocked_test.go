package text_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/text"

	"github.com/stretchr/testify/assert"
)

func TestSplitText(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{input: "", expected: []string{""}},
		{input: "a", expected: []string{"a", ""}},
		{input: "a", expected: []string{"a", ""}},
		{input: "b\n", expected: []string{"b\n", ""}},
		{input: "   spacing!   ", expected: []string{"   spacing!   ", ""}},
		{input: "terse1\n\n\nterse2", expected: []string{"terse1\n\n\nterse2", ""}},
		{
			input: "aaaaaaaaa bbbbbbbbb ccccccccc ddddddddd eeeeeeeee fffffffff ggggggggg hhhhhhhhh iiiiiiiii jjjjjjjjj kkkkk",
			expected: []string{
				"aaaaaaaaa bbbbbbbbb ccccccccc ddddddddd eeeeeeeee fffffffff ggggggggg hhhhhhhhh ",
				"iiiiiiiii jjjjjjjjj kkkkk",
				"",
			},
		},
	}

	for _, tc := range tt {
		result := text.Blocked(tc.input)
		assert.Equal(t, tc.expected, result, "wrong result for <"+tc.input+">")
	}
}
