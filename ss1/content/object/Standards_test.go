package object_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/object"
)

func TestStandardDescriptorsReturnsProperLength(t *testing.T) {
	descriptor := object.StandardDescriptors()
	totalLength := 4 // version prefix

	for _, classDesc := range descriptor {
		totalLength += classDesc.TotalDataSize()
	}
	assert.Equal(t, 17951, totalLength) // as taken from original CD
}

func TestStandardDescriptorsReturnsProperAmount(t *testing.T) {
	descriptor := object.StandardDescriptors()
	total := 0

	for _, classDesc := range descriptor {
		total += classDesc.TotalTypeCount()
	}
	assert.Equal(t, 476, total)
}
