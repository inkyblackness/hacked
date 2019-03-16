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
	err := suite.sequencer.Add(compression.TileColorOp{Offset: 0x100000})
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

func (suite *ControlWordSequencerSuite) TestControlWordsAreOrderedByOperationFrequency() {
	suite.givenRegisteredOperations(
		compression.TileColorOp{Offset: 1},
		compression.TileColorOp{}, compression.TileColorOp{})
	suite.whenSequenceIsCreated()
	suite.thenControlWordsShouldBe(compression.ControlWordOf(12, 0, 0), compression.ControlWordOf(12, 0, 1))
}

// TODO: Operation frequency might not be the only deciding factor.
// The further criteria is the arrangement of control word sequences that generate offsets beyond 17 bit.
// mixed into this then also the magic of reusing extension sequences.
// I believe the order of the entries is irrelevant for the words that have parameters within 17 bits.
// the order becomes relevant for those that need extension. Though, even then, what is the benefit
// of putting an entry with a high parameter value before those with a smaller value? It still needs
// the same amount of extensions to reach. Perhaps it is important whether the extensions are within
// the same bucket, or spanning several jumps of 15 indices...

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
