package resource_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
)

func TestNullProvider(t *testing.T) {
	provider := resource.NullProvider()
	verifyError := func(id uint16) {
		_, err := provider.Resource(resource.ID(id))
		assert.Error(t, err)
	}

	assert.Equal(t, 0, len(provider.IDs()))
	verifyError(0)
	verifyError(1)
}
