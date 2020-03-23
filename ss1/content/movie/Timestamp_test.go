package movie_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/movie"
)

func TestTimestampFromSecond(t *testing.T) {
	ts := movie.TimestampFromSeconds(7.25)

	assert.Equal(t, byte(7), ts.Second)
	assert.Equal(t, uint16(0x4000), ts.Fraction)
}

func TestTimeToRawB(t *testing.T) {
	ts := movie.TimestampFromSeconds(255.75)

	assert.Equal(t, byte(255), ts.Second)
	assert.Equal(t, uint16(0xC000), ts.Fraction)
}

func TestDurationConversion(t *testing.T) {
	for fraction := 0; fraction <= math.MaxUint16; fraction++ {
		ts := movie.Timestamp{Second: 0, Fraction: uint16(fraction)}
		dur := ts.ToDuration()
		ts2 := movie.TimestampFromDuration(dur)

		assert.Equal(t, ts, ts2, fmt.Sprintf("Mismatch for fraction %v", fraction))
	}
}
