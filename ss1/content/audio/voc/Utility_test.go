package voc // nolint: testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrequencyDivisorFor11111(t *testing.T) {
	divisor := byte(0xA6)
	sampleRate := divisorToSampleRate(divisor)
	result := sampleRateToDivisor(sampleRate)

	assert.Equal(t, divisor, result)
}

func TestFrequencyDivisorFor222222(t *testing.T) {
	divisor := byte(0x2D)
	sampleRate := divisorToSampleRate(divisor)
	result := sampleRateToDivisor(sampleRate)

	assert.Equal(t, divisor, result)
}
