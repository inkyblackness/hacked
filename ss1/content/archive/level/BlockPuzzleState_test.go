package level_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"

	"github.com/stretchr/testify/assert"
)

func TestBlockPuzzleCellValueReturnsTheCellValue_A(t *testing.T) {
	var data [16]byte
	state := level.NewBlockPuzzleState(data[:], 1, 1)

	assert.Equal(t, 0, state.CellValue(0, 0))
}

func TestBlockPuzzleCellValueReturnsTheCellValue_B(t *testing.T) {
	var data [16]byte
	data[16-4] = 2
	state := level.NewBlockPuzzleState(data[:], 1, 1)

	assert.Equal(t, 2, state.CellValue(0, 0))
}

func TestBlockPuzzleCellValueReturnsTheCellValue_C(t *testing.T) {
	var data [16]byte
	data[16-4] = 0x18
	state := level.NewBlockPuzzleState(data[:], 2, 1)

	assert.Equal(t, 3, state.CellValue(1, 0))
}

func TestBlockPuzzleCellValueReturnsTheCellValue_D(t *testing.T) {
	var data [16]byte
	data[15] = 0x40
	data[8] = 0x01
	state := level.NewBlockPuzzleState(data[:], 4, 3)

	assert.Equal(t, 5, state.CellValue(3, 1))
}

func TestBlockPuzzleCellValueReturnsTheCellValue_E(t *testing.T) {
	var data [16]byte
	data[3] = 0x38
	state := level.NewBlockPuzzleState(data[:], 6, 7)

	assert.Equal(t, 7, state.CellValue(5, 6))
}

func TestBlockPuzzleCellValueReturnsZeroOutOfBounds_A(t *testing.T) {
	var data [16]byte
	for index := 0; index < len(data); index++ {
		data[index] = 0xFF
	}
	state := level.NewBlockPuzzleState(data[:], 1, 1)

	assert.Equal(t, 0, state.CellValue(0, 1))
	assert.Equal(t, 0, state.CellValue(1, 0))
	assert.Equal(t, 0, state.CellValue(5, 5))
	assert.Equal(t, 0, state.CellValue(8, 8))
}

func TestBlockPuzzleCellValueReturnsZeroOutOfBounds_B(t *testing.T) {
	var data [16]byte
	for index := 0; index < len(data); index++ {
		data[index] = 0xFF
	}
	state := level.NewBlockPuzzleState(data[:], 7, 6)

	assert.Equal(t, 0, state.CellValue(6, 7))
	assert.Equal(t, 0, state.CellValue(10, 10))
}

func TestBlockPuzzleCellValueReturnsZeroOutOfBounds_C(t *testing.T) {
	var data [16]byte
	for index := 0; index < len(data); index++ {
		data[index] = 0xFF
	}
	state := level.NewBlockPuzzleState(data[:], 7, 7)

	assert.Equal(t, 0, state.CellValue(6, 0))
}

func TestBlockPuzzleSetCellValue_A(t *testing.T) {
	var data [16]byte
	state := level.NewBlockPuzzleState(data[:], 1, 1)

	state.SetCellValue(0, 0, 1)
	assert.Equal(t, [16]byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00}, data)
}

func TestBlockPuzzleSetCellValue_B(t *testing.T) {
	var data [16]byte
	state := level.NewBlockPuzzleState(data[:], 2, 2)

	state.SetCellValue(1, 0, 7)
	assert.Equal(t, [16]byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xC0, 0x01, 0x00, 0x00}, data)
}

func TestBlockPuzzleSetCellValue_C(t *testing.T) {
	var data [16]byte
	state := level.NewBlockPuzzleState(data[:], 6, 7)

	state.SetCellValue(5, 6, 5)
	assert.Equal(t, [16]byte{
		0x00, 0x00, 0x00, 0x28,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00}, data)
}

func TestBlockPuzzleSetCellValueIgnoredOutOfBounds(t *testing.T) {
	var data [16]byte
	state := level.NewBlockPuzzleState(data[:], 6, 7)

	state.SetCellValue(8, 8, 7)
	assert.Equal(t, [16]byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00}, data)
}
