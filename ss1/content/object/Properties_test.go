package object_test

import (
	"bytes"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
)

func TestPropertiesForObjectReturnErrorForWrongIDs(t *testing.T) {
	tt := []struct {
		name   string
		triple object.Triple
	}{
		{"wrong class", object.TripleFrom(20, 0, 2)},
		{"wrong subclass", object.TripleFrom(1, 40, 0)},
		{"wrong type", object.TripleFrom(1, 0, 20)},
	}
	table := object.StandardPropertiesTable()

	for _, tc := range tt {
		_, err := table.ForObject(tc.triple)
		assert.Error(t, err, "error expected for "+tc.name)
	}
}

func TestPropertiesForObjectData(t *testing.T) {
	tt := []struct {
		triple   object.Triple
		genSize  int
		specSize int
	}{
		{object.TripleFrom(0, 0, 0), 2, 1},
		{object.TripleFrom(14, 4, 1), 75, 1},
		{object.TripleFrom(10, 3, 4), 1, 1},
	}
	table := object.StandardPropertiesTable()

	for _, tc := range tt {
		prop, err := table.ForObject(tc.triple)
		assert.Nil(t, err, "no error expected for "+tc.triple.String())
		assert.Equal(t, tc.genSize, len(prop.Generic), "wrong generic size for "+tc.triple.String())
		assert.Equal(t, tc.specSize, len(prop.Specific), "wrong specific size for "+tc.triple.String())
	}
}

func TestPropertiesEncoding(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	table := object.StandardPropertiesTable()

	table.Code(encoder)
	result := buf.Bytes()
	assert.Equal(t, 17951, len(result)) // as taken from original CD
}
