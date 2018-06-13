package resource_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
)

func TestKeyOfReturnsKey(t *testing.T) {
	assert.Equal(t, resource.Key{ID: 0x1000, Lang: resource.LangGerman, Index: 123}, resource.KeyOf(0x1000, resource.LangGerman, 123))
	assert.Equal(t, resource.Key{ID: 0xA000, Lang: resource.LangAny, Index: 0}, resource.KeyOf(0xA000, resource.LangAny, 0))
}
