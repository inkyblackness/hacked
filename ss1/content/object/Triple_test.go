package object_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/object"

	"github.com/stretchr/testify/assert"
)

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
