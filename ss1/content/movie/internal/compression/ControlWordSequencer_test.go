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
