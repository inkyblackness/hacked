package numbers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/numbers"
)

func TestBinaryCodedDecimal(t *testing.T) {
	tt := []struct {
		int uint16
		bcd uint16
	}{
		{int: 0, bcd: 0x0000},
		{int: 1, bcd: 0x0001},
		{int: 20, bcd: 0x0020},
		{int: 300, bcd: 0x0300},
		{int: 4000, bcd: 0x4000},
		{int: 50000, bcd: 0x0000},
		{int: 9999, bcd: 0x9999},
		{int: 451, bcd: 0x0451},
	}

	for _, tc := range tt {
		toBCD := numbers.ToBinaryCodedDecimal(tc.int)
		assert.Equal(t, tc.bcd, toBCD, "Wrong result converting to BCD for %d", tc.int)

		fromBCD := numbers.FromBinaryCodedDecimal(tc.bcd)
		if tc.int <= 9999 {
			assert.Equal(t, tc.int, fromBCD, "Wrong result converting from BCD for %d", tc.int)
		} else {
			assert.Equal(t, uint16(0), fromBCD, "Wrong result converting from BCD for %d", tc.int)
		}
	}
}
