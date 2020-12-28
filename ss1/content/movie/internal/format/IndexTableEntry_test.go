package format_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
)

func TestIndexTableSizeFor_ExistingSizes(t *testing.T) {
	// These sample values are always the minimum and maximum amount of index entries
	// found for a given index size.

	assert.Equal(t, 0x0400, format.IndexTableSizeFor(3))
	assert.Equal(t, 0x0400, format.IndexTableSizeFor(127))

	assert.Equal(t, 0x0C00, format.IndexTableSizeFor(130))
	assert.Equal(t, 0x0C00, format.IndexTableSizeFor(218))

	assert.Equal(t, 0x1C00, format.IndexTableSizeFor(738))
	assert.Equal(t, 0x1C00, format.IndexTableSizeFor(755))

	assert.Equal(t, 0x3400, format.IndexTableSizeFor(1475))
	assert.Equal(t, 0x3400, format.IndexTableSizeFor(1523))
}
