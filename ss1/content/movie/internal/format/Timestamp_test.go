package format_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
)

func TestTimestampFromSecond(t *testing.T) {
	ts := format.TimestampFromSeconds(7.25)

	assert.Equal(t, byte(7), ts.Second)
	assert.Equal(t, uint16(0x4000), ts.Fraction)
}

func TestTimeToRawB(t *testing.T) {
	ts := format.TimestampFromSeconds(255.75)

	assert.Equal(t, byte(255), ts.Second)
	assert.Equal(t, uint16(0xC000), ts.Fraction)
}

func TestDurationConversion(t *testing.T) {
	for fraction := 0; fraction <= math.MaxUint16; fraction++ {
		ts := format.Timestamp{Second: 0, Fraction: uint16(fraction)}
		dur := ts.ToDuration()
		ts2 := format.TimestampFromDuration(dur)

		assert.Equal(t, ts, ts2, fmt.Sprintf("Mismatch for fraction %v", fraction))
	}
}

func TestTimestampPlus(t *testing.T) {
	tt := []struct {
		a        format.Timestamp
		b        format.Timestamp
		expected format.Timestamp
	}{
		{
			a:        format.Timestamp{Second: 0, Fraction: 0},
			b:        format.Timestamp{Second: 1, Fraction: 0x8000},
			expected: format.Timestamp{Second: 1, Fraction: 0x8000},
		},
		{
			a:        format.Timestamp{Second: 2, Fraction: 0x8000},
			b:        format.Timestamp{Second: 1, Fraction: 0x8000},
			expected: format.Timestamp{Second: 4, Fraction: 0x0000},
		},
		{
			a:        format.Timestamp{Second: 0xFF, Fraction: 0x0000},
			b:        format.Timestamp{Second: 2, Fraction: 0x8000},
			expected: format.Timestamp{Second: 0xFF, Fraction: 0xFFFF},
		},
	}
	for index, tc := range tt {
		tc := tc
		t.Run(fmt.Sprintf("case %d", index), func(t *testing.T) {
			result := tc.a.Plus(tc.b)
			assert.Equal(t, tc.expected, result)
		})
	}
}
