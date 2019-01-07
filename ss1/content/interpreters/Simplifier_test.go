package interpreters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplifierRawValueCallsHandler(t *testing.T) {
	var calledMinValue int64
	var calledMaxValue int64
	rawHandler := func(minValue, maxValue int64, formatter RawValueFormatter) {
		calledMinValue, calledMaxValue = minValue, maxValue
	}
	simpl := NewSimplifier(rawHandler)

	simpl.rawValue(&entry{count: 2})

	assert.Equal(t, int64(-1), calledMinValue)
	assert.Equal(t, int64(32767), calledMaxValue)
}

func TestSimplifierEnumValueReturnsFalseIfNoHandlerRegistered(t *testing.T) {
	simpl := NewSimplifier(func(minValue, maxValue int64, formatter RawValueFormatter) {})

	result := simpl.enumValue(map[uint32]string{})

	assert.False(t, result)
}

func TestSimplifierEnumValueCallsRegisteredHandler(t *testing.T) {
	result := map[uint32]string{}
	simpl := NewSimplifier(func(minValue, maxValue int64, formatter RawValueFormatter) {})
	simpl.SetEnumValueHandler(func(values map[uint32]string) {
		result = values
	})

	expected := map[uint32]string{1: "value"}
	simpl.enumValue(expected)

	assert.Equal(t, expected, result)
}
