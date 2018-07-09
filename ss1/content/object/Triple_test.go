package object_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/object"

	"github.com/stretchr/testify/assert"
)

func TestClasses(t *testing.T) {
	classes := object.Classes()
	assert.Equal(t, 15, len(classes))
}

func TestClassString(t *testing.T) {
	assert.Equal(t, "Gun", object.Class(0).String())
	assert.Equal(t, "UnknownFF", object.Class(0xFF).String())
}

func TestTripleString(t *testing.T) {
	tt := []struct {
		triple   object.Triple
		expected string
	}{
		{object.TripleFrom(1, 2, 3), " 1/2/ 3"},
		{object.TripleFrom(20, 0, 10), "20/0/10"},
	}
	for _, tc := range tt {
		result := tc.triple.String()
		assert.Equal(t, tc.expected, result)
	}
}
