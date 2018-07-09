package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardDescriptorsReturnsProperLength(t *testing.T) {
	descriptor := StandardDescriptors()
	totalLength := 4 // version prefix

	for _, classDesc := range descriptor {
		totalLength += classDesc.TotalDataSize()
	}
	assert.Equal(t, 17951, totalLength) // as taken from original CD
}

func TestStandardDescriptorsReturnsProperAmount(t *testing.T) {
	descriptor := StandardDescriptors()
	total := 0

	for _, classDesc := range descriptor {
		total += classDesc.TotalTypeCount()
	}
	assert.Equal(t, 476, total)
}
