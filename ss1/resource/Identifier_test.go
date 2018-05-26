package resource_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
)

func TestIDValueReturnsOwnValue(t *testing.T) {
	assert.Equal(t, uint16(0), resource.ID(0).Value())
	assert.Equal(t, uint16(123), resource.ID(123).Value())
	assert.Equal(t, uint16(math.MaxUint16), resource.ID(math.MaxUint16).Value(), "maximum of uint16 should be supported")
}

func TestIDImplementsStringer(t *testing.T) {
	assert.Equal(t, "0FA0", fmt.Sprintf("%v", resource.ID(0x0FA0)))
}
