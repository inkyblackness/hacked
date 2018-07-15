package level

// BlockPuzzleState describes the state of a block puzzle.
type BlockPuzzleState struct {
	width  int
	height int
	data   []byte
}

// NewBlockPuzzleState returns a new instance of a block puzzle state modifier.
// The returned instance works with the passed data slice directly.
func NewBlockPuzzleState(data []byte, height, width int) *BlockPuzzleState {
	return &BlockPuzzleState{
		data:   data,
		width:  width,
		height: height}
}

// CellValue returns the value of the identified cell.
func (state *BlockPuzzleState) CellValue(row, col int) int {
	scratch := uint16(0)

	if state.positionOk(row, col) {
		byteOffset, byteShift := state.byteOffset(row, col)
		bitsAvailable := 8 - byteShift

		scratch = uint16(state.data[state.mappedByteIndex(byteOffset)]) >> byteShift
		if bitsAvailable < 3 {
			scratch |= uint16(state.data[state.mappedByteIndex(byteOffset+1)]) << bitsAvailable
		}
	}

	return int(scratch & 7)
}

// SetCellValue sets the value of the identified cell.
func (state *BlockPuzzleState) SetCellValue(row, col int, value int) {
	if state.positionOk(row, col) {
		byteOffset, byteShift := state.byteOffset(row, col)
		needsSecondByte := (8 - byteShift) < 3
		scratch := uint16(state.data[state.mappedByteIndex(byteOffset)])
		mask := uint16(0x0007) << byteShift

		if needsSecondByte {
			scratch |= uint16(state.data[state.mappedByteIndex(byteOffset+1)]) << 8
		}
		scratch = (scratch & ^mask) | (uint16(value) << byteShift)
		state.data[state.mappedByteIndex(byteOffset)] = byte((scratch >> 0) & 0xFF)
		if needsSecondByte {
			state.data[state.mappedByteIndex(byteOffset+1)] = byte((scratch >> 8) & 0xFF)
		}
	}
}

func (state *BlockPuzzleState) bitOffset(row, col int) int {
	return (row*state.width + col) * 3
}

func (state *BlockPuzzleState) positionOk(row, col int) bool {
	bitOffset := state.bitOffset(row, col)

	return (row < state.height) && (col < state.width) && (bitOffset <= (128 - 3))
}

func (state *BlockPuzzleState) byteOffset(row, col int) (byteOffset int, byteShift uint) {
	bitOffset := state.bitOffset(row, col)
	byteOffset = bitOffset / 8
	byteShift = uint(bitOffset % 8)
	return
}

func (state *BlockPuzzleState) mappedByteIndex(index int) int {
	remainder := index % 4
	return 15 - (index - remainder + (3 - remainder))
}
