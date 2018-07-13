package interpreters_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/interpreters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescriptionWithReturnsANewDescription(t *testing.T) {
	second := interpreters.New().With("test", 0, 32)

	require.NotNil(t, second)
	assert.NotEqual(t, interpreters.New(), second)
}

func TestDescriptionWithCopiesPreviousFields(t *testing.T) {
	first := interpreters.New().With("field1", 0, 32)
	second := first.With("field2", 32, 32)
	onlyField2 := interpreters.New().With("field2", 32, 32)

	assert.NotEqual(t, onlyField2, second)
}

func TestDescriptionWithLeavesOriginalAlone(t *testing.T) {
	first := interpreters.New().With("field1", 0, 32)
	onlyField1 := interpreters.New().With("field1", 0, 32)

	first.With("field2", 32, 32)

	assert.Equal(t, onlyField1, first)
}

func TestDescriptionForCreatesNewInstance(t *testing.T) {
	data := make([]byte, 10)
	inst := interpreters.New().For(data)

	assert.NotNil(t, inst)
}

func TestDescriptionRefiningReturnsANewDescription(t *testing.T) {
	refined := interpreters.New().With("field1", 0, 16)
	second := interpreters.New().Refining("test", 0, 4, refined, interpreters.Always)

	require.NotNil(t, second)
	assert.NotEqual(t, interpreters.New(), second)
}

func TestDescriptionRefiningCopiesPreviousFields(t *testing.T) {
	first := interpreters.New().With("fieldA", 0, 8)
	refined := interpreters.New().With("field1", 0, 16)
	second := first.Refining("test", 0, 4, refined, interpreters.Always)
	secondMissing := interpreters.New().Refining("test", 0, 4, refined, interpreters.Always)

	assert.NotEqual(t, secondMissing, second)
}

func TestDescriptionRefiningCopiesPreviousRefinements(t *testing.T) {
	first := interpreters.New().Refining("sub1", 0, 1, interpreters.New(), interpreters.Always)
	second := first.Refining("sub2", 0, 1, interpreters.New(), interpreters.Always)

	assert.NotContains(t, first.For(nil).ActiveRefinements(), "sub2")
	assert.Contains(t, second.For(nil).ActiveRefinements(), "sub1")
}

func TestDescriptionAsPanicsIfNoFieldIsActive(t *testing.T) {
	assert.PanicsWithValue(t, "No field active", func() {
		interpreters.New().As(func(simpl *interpreters.Simplifier) bool { return false })
	})
}

func TestDescriptionAsRegistersRangeFunctionForLastField(t *testing.T) {
	rangeFuncCalled := false
	rangeFunc := func(simpl *interpreters.Simplifier) bool {
		rangeFuncCalled = true
		return true
	}
	desc := interpreters.New().With("fieldA", 0, 2).As(rangeFunc)

	inst := desc.For([]byte{0x00, 0x01})
	simplifier := interpreters.NewSimplifier(func(minValue, maxValue int64, formatter interpreters.RawValueFormatter) {})
	inst.Describe("fieldA", simplifier)

	assert.True(t, rangeFuncCalled)
}
