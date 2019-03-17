package compression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
)

type ControlWordSequencerSuite struct {
	suite.Suite

	sequencer compression.ControlWordSequencer

	sequenceErr error
	sequence    compression.ControlWordSequence
}

func TestControlWordSequencer(t *testing.T) {
	suite.Run(t, new(ControlWordSequencerSuite))
}

func (suite *ControlWordSequencerSuite) SetupTest() {
	suite.sequencer = compression.ControlWordSequencer{}
	suite.sequenceErr = nil
	suite.sequence = compression.ControlWordSequence{}
}

func (suite *ControlWordSequencerSuite) TestEmptySequence() {
	suite.whenSequenceIsCreated()
	suite.thenControlWordsShouldHaveALenOf(0)
}

func (suite *ControlWordSequencerSuite) TestAddReturnsErrorOnWrongOffset() {
	err := suite.sequencer.Add(compression.TileColorOp{Offset: compression.ControlWordParamLimit + 1})
	assert.NotNil(suite.T(), err, "error expected adding an operation with too high offset")
}

func (suite *ControlWordSequencerSuite) TestSingleOperation() {
	suite.givenRegisteredOperations(compression.TileColorOp{})
	suite.whenSequenceIsCreated()
	suite.thenControlWordsShouldHaveALenOf(1)
}

func (suite *ControlWordSequencerSuite) TestDuplicatedOperationShouldResultInOneControlWord() {
	suite.givenRegisteredOperations(compression.TileColorOp{}, compression.TileColorOp{})
	suite.whenSequenceIsCreated()
	suite.thenControlWordsShouldHaveALenOf(1)
}

func (suite *ControlWordSequencerSuite) TestControlWordsAreOrderedByOperationFrequencyFirst() {
	suite.givenRegisteredOperations(
		compression.TileColorOp{Offset: 1},
		compression.TileColorOp{}, compression.TileColorOp{})
	suite.whenSequenceIsCreated()
	suite.thenControlWordsShouldBe(compression.ControlWordOf(12, 0, 0), compression.ControlWordOf(12, 0, 1))
}

func (suite *ControlWordSequencerSuite) TestControlWordsAreOrderedByOffsetSecond() {
	suite.givenRegisteredOperations(
		compression.TileColorOp{Offset: 1}, compression.TileColorOp{Offset: 1},
		compression.TileColorOp{}, compression.TileColorOp{})
	suite.whenSequenceIsCreated()
	suite.thenControlWordsShouldBe(compression.ControlWordOf(12, 0, 0), compression.ControlWordOf(12, 0, 1))
}

func (suite *ControlWordSequencerSuite) TestControlWordsAreOrderedByControlTypeThird() {
	suite.givenRegisteredOperations(
		compression.TileColorOp{Type: compression.CtrlColorTile16ColorsMasked},
		compression.TileColorOp{Type: compression.CtrlColorTile16ColorsMasked},
		compression.TileColorOp{}, compression.TileColorOp{})
	suite.whenSequenceIsCreated()
	suite.thenControlWordsShouldBe(
		compression.ControlWordOf(12, 0, 0),
		compression.ControlWordOf(12, compression.CtrlColorTile16ColorsMasked, 0))
}

func (suite *ControlWordSequencerSuite) TestControlWordLongOffsetsWhenBitstreamSpaceExhausted_FirstTime() {
	suite.givenSequencerBitstreamIndexLimitOf(1)
	suite.givenRegisteredOperations(
		compression.TileColorOp{Offset: 0},
		compression.TileColorOp{Offset: 1},
		compression.TileColorOp{Offset: 2},
		compression.TileColorOp{Offset: 3})
	suite.whenSequenceIsCreated()
	suite.thenControlWordsShouldBe(
		compression.ControlWordOf(12, 0, 0),
		compression.ControlWordOf(0, 0, 2), // first long offset jump to index 2
		compression.ControlWordOf(4, 0, 1),
		compression.ControlWordOf(4, 0, 2),
		compression.ControlWordOf(4, 0, 3))
}

func (suite *ControlWordSequencerSuite) TestControlWordLongOffsetsWhenBitstreamSpaceExhausted_SecondTime() {
	suite.givenSequencerBitstreamIndexLimitOf(1)
	suite.givenRegisteredOperations(
		compression.TileColorOp{Offset: 0}, compression.TileColorOp{Offset: 1},
		compression.TileColorOp{Offset: 2}, compression.TileColorOp{Offset: 3},
		compression.TileColorOp{Offset: 4}, compression.TileColorOp{Offset: 5},
		compression.TileColorOp{Offset: 6}, compression.TileColorOp{Offset: 7},
		compression.TileColorOp{Offset: 8}, compression.TileColorOp{Offset: 9},
		compression.TileColorOp{Offset: 10}, compression.TileColorOp{Offset: 11},
		compression.TileColorOp{Offset: 12}, compression.TileColorOp{Offset: 13},
		compression.TileColorOp{Offset: 14}, compression.TileColorOp{Offset: 15},
		compression.TileColorOp{Offset: 16}, compression.TileColorOp{Offset: 17})
	suite.whenSequenceIsCreated()
	suite.thenControlWordsShouldBe(
		compression.ControlWordOf(12, 0, 0),
		compression.ControlWordOf(0, 0, 2), // first long offset jump to index 2
		compression.ControlWordOf(4, 0, 1), compression.ControlWordOf(4, 0, 2),
		compression.ControlWordOf(4, 0, 3), compression.ControlWordOf(4, 0, 4),
		compression.ControlWordOf(4, 0, 5), compression.ControlWordOf(4, 0, 6),
		compression.ControlWordOf(4, 0, 7), compression.ControlWordOf(4, 0, 8),
		compression.ControlWordOf(4, 0, 9), compression.ControlWordOf(4, 0, 10),
		compression.ControlWordOf(4, 0, 11), compression.ControlWordOf(4, 0, 12),
		compression.ControlWordOf(4, 0, 13), compression.ControlWordOf(4, 0, 14),
		compression.ControlWordOf(4, 0, 15),
		compression.ControlWordOf(0, 0, 18), // second long offset jump to index 18
		compression.ControlWordOf(4, 0, 16), compression.ControlWordOf(4, 0, 17))
}

type bitstreamExpectation struct {
	bits  int
	value uint32
}

func bits12(value uint32) bitstreamExpectation {
	return bitstreamExpectation{bits: 12, value: value}
}

func bits4(value uint32) bitstreamExpectation {
	return bitstreamExpectation{bits: 4, value: value}
}

func (suite *ControlWordSequencerSuite) TestBitstreamForSimpleSequence() {
	ops := []compression.TileColorOp{{Offset: 1}, {}, {}}
	suite.givenRegisteredOperations(ops...)
	suite.whenSequenceIsCreated()
	suite.thenBitstreamShouldBeFor(ops,
		[]bitstreamExpectation{bits12(1), bits12(0), bits12(0)})
}

func (suite *ControlWordSequencerSuite) givenSequencerBitstreamIndexLimitOf(value uint32) {
	suite.sequencer.BitstreamIndexLimit = value
}

func (suite *ControlWordSequencerSuite) givenRegisteredOperations(ops ...compression.TileColorOp) {
	suite.T().Helper()
	for index, op := range ops {
		err := suite.sequencer.Add(op)
		require.Nil(suite.T(), err, "no error expected adding operation ", index)
	}
}

func (suite *ControlWordSequencerSuite) whenSequenceIsCreated() {
	suite.sequence, suite.sequenceErr = suite.sequencer.Sequence()
}

func (suite *ControlWordSequencerSuite) thenControlWordsShouldHaveALenOf(expected int) {
	suite.T().Helper()
	require.Nil(suite.T(), suite.sequenceErr, "no error expected in this verification")
	words := suite.sequence.ControlWords()
	assert.Equal(suite.T(), expected, len(words))
}

func (suite *ControlWordSequencerSuite) thenControlWordsShouldBe(expected ...compression.ControlWord) {
	suite.T().Helper()
	words := suite.sequence.ControlWords()
	assert.Equal(suite.T(), len(expected), len(words), "length mismatch")
	assert.Equal(suite.T(), expected, words, "words mismatch")
}

func (suite *ControlWordSequencerSuite) thenBitstreamShouldBeFor(
	ops []compression.TileColorOp,
	expected []bitstreamExpectation) {
	data, err := suite.sequence.BitstreamFor(ops)
	suite.T().Helper()
	require.Nil(suite.T(), err, "no error expected extracting bitstream")
	bitstream := compression.NewBitstreamReader(data)
	for index, exp := range expected {
		value := bitstream.Read(exp.bits)
		assert.Equal(suite.T(), exp.value, value, "Value mismatch at index ", index)
		bitstream.Advance(exp.bits)
	}
}
