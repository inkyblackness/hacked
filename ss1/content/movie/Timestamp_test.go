package movie_test

import (
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
