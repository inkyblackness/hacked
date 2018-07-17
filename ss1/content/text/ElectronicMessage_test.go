package text_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ElectronicMessageSuite struct {
	suite.Suite
	cp text.Codepage
}

func TestElectronicMessageSuite(t *testing.T) {
	suite.Run(t, new(ElectronicMessageSuite))
}

func (suite *ElectronicMessageSuite) SetupTest() {
	suite.cp = text.DefaultCodepage()
}

func (suite *ElectronicMessageSuite) TestEncodeBasicMessage() {
	message := text.EmptyElectronicMessage()

	message.Title = "1"
	message.Sender = "2"
	message.Subject = "3"
	message.VerboseText = "4"
	message.TerseText = "5"

	encoded := message.Encode(suite.cp)

	suite.verifyBlock(0, encoded, []byte{0x00})
	suite.verifyBlock(1, encoded, []byte{0x31, 0x00})
	suite.verifyBlock(2, encoded, []byte{0x32, 0x00})
	suite.verifyBlock(3, encoded, []byte{0x33, 0x00})
	suite.verifyBlock(4, encoded, []byte{0x34, 0x00})
	suite.verifyBlock(5, encoded, []byte{0x00})
	suite.verifyBlock(6, encoded, []byte{0x35, 0x00})
	suite.verifyBlock(7, encoded, []byte{0x00})
}

func (suite *ElectronicMessageSuite) TestEncodeMeta_A() {
	message := text.EmptyElectronicMessage()

	message.NextMessage = 0x20
	message.ColorIndex = 0x13
	message.LeftDisplay = 30
	message.RightDisplay = 40

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, len(encoded) > 0)
	suite.verifyBlock(0, encoded, suite.cp.Encode("i20 c13 30,40"))
}

func (suite *ElectronicMessageSuite) TestEncodeMeta_B() {
	message := text.EmptyElectronicMessage()

	message.IsInterrupt = true
	message.LeftDisplay = 31

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, len(encoded) > 0)
	suite.verifyBlock(0, encoded, suite.cp.Encode("t 31"))
}

func (suite *ElectronicMessageSuite) TestEncodeMeta_C() {
	message := text.EmptyElectronicMessage()

	message.IsInterrupt = true

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, len(encoded) > 0)
	suite.verifyBlock(0, encoded, suite.cp.Encode("t"))
}

func (suite *ElectronicMessageSuite) TestEncodeCreatesNewBlocksPerNewLine() {
	message := text.EmptyElectronicMessage()

	message.VerboseText = "line1\n\n\nline2"
	message.TerseText = "terse1\n\n\nterse2"

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, len(encoded) > 0)
	suite.verifyBlock(4, encoded, suite.cp.Encode("line1\n\n\nline2"))
	suite.verifyBlock(6, encoded, suite.cp.Encode("terse1\n\n\nterse2"))
}

func (suite *ElectronicMessageSuite) TestEncodeBreaksUpLinesAfterLimitCharacters() {
	message := text.EmptyElectronicMessage()

	message.VerboseText = "aaaaaaaaa bbbbbbbbb ccccccccc ddddddddd eeeeeeeee fffffffff ggggggggg hhhhhhhhh iiiiiiiii jjjjjjjjj kkkkk"

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, len(encoded) > 0)
	suite.verifyBlock(4, encoded,
		suite.cp.Encode("aaaaaaaaa bbbbbbbbb ccccccccc ddddddddd eeeeeeeee fffffffff ggggggggg hhhhhhhhh "))
	suite.verifyBlock(5, encoded,
		suite.cp.Encode("iiiiiiiii jjjjjjjjj kkkkk"))
}

func (suite *ElectronicMessageSuite) TestDecodeMeta() {
	message, err := text.DecodeElectronicMessage(suite.cp, suite.holderWithMeta("i20 c13 30,40"))

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), 0x20, message.NextMessage)
	assert.Equal(suite.T(), 0x13, message.ColorIndex)
	assert.Equal(suite.T(), 30, message.LeftDisplay)
	assert.Equal(suite.T(), 40, message.RightDisplay)
}

func (suite *ElectronicMessageSuite) TestDecodeMeta_Failure() {
	_, err := text.DecodeElectronicMessage(suite.cp, suite.holderWithMeta("i20 c 13 30,40"))

	assert.NotNil(suite.T(), err)
}

func (suite *ElectronicMessageSuite) TestDecodeMetaColorIs8BitUnsigned() {
	message, err := text.DecodeElectronicMessage(suite.cp, suite.holderWithMeta("cD1"))

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), 0xD1, message.ColorIndex)
}

func (suite *ElectronicMessageSuite) TestDecodeMessage() {
	message, err := text.DecodeElectronicMessage(suite.cp, suite.holderWithMeta("10"))

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), "title", message.Title)
	assert.Equal(suite.T(), "sender", message.Sender)
	assert.Equal(suite.T(), "subject", message.Subject)
	assert.Equal(suite.T(), "verbose", message.VerboseText)
	assert.Equal(suite.T(), "terse", message.TerseText)
}

func (suite *ElectronicMessageSuite) TestDecodeMessageIsPossibleForVanillaDummyMails() {
	message, err := text.DecodeElectronicMessage(suite.cp, suite.vanillaStubMail())

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), "", message.Title)
	assert.Equal(suite.T(), "", message.Sender)
	assert.Equal(suite.T(), "", message.Subject)
	assert.Equal(suite.T(), "stub emailstub email", message.VerboseText)
	assert.Equal(suite.T(), "", message.TerseText)
}

func (suite *ElectronicMessageSuite) TestDecodeMessageIsPossibleForMissingTerminatingLine() {
	message, err := text.DecodeElectronicMessage(suite.cp, suite.holderWithMissingTerminatingLine())

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), "title", message.Title)
	assert.Equal(suite.T(), "sender", message.Sender)
	assert.Equal(suite.T(), "subject", message.Subject)
	assert.Equal(suite.T(), "verbose text", message.VerboseText)
	assert.Equal(suite.T(), "terse text", message.TerseText)
}

func (suite *ElectronicMessageSuite) TestRecodeMessage() {
	inMessage := text.EmptyElectronicMessage()
	inMessage.IsInterrupt = true
	inMessage.NextMessage = 0x10
	inMessage.ColorIndex = 0x20
	inMessage.LeftDisplay = 40
	inMessage.RightDisplay = 50
	inMessage.VerboseText = "abcd\nefgh\nsome"
	inMessage.TerseText = "\n"

	blocks := inMessage.Encode(suite.cp)
	outMessage, err := text.DecodeElectronicMessage(suite.cp, resource.MemoryBlockProvider(blocks))

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), outMessage)
	assert.Equal(suite.T(), inMessage, outMessage)
}

func (suite *ElectronicMessageSuite) TestRecodeMessageWithMultipleNewLines() {
	inMessage := text.EmptyElectronicMessage()
	inMessage.IsInterrupt = true
	inMessage.NextMessage = 0x10
	inMessage.ColorIndex = 0x20
	inMessage.LeftDisplay = 40
	inMessage.RightDisplay = 50
	inMessage.VerboseText = "first\n\n\nsecond"
	inMessage.TerseText = "terse\n"

	blocks := inMessage.Encode(suite.cp)
	outMessage, err := text.DecodeElectronicMessage(suite.cp, resource.MemoryBlockProvider(blocks))

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), outMessage)
	assert.Equal(suite.T(), inMessage, outMessage)
}

func (suite *ElectronicMessageSuite) verifyBlock(index int, blocks [][]byte, expected []byte) {
	require.True(suite.T(), len(blocks) > index, "too few blocks")
	assert.Equal(suite.T(), expected, blocks[index])
}

func (suite *ElectronicMessageSuite) holderWithMeta(meta string) resource.BlockProvider {
	blocks := [][]byte{
		suite.cp.Encode(meta),
		suite.cp.Encode("title"),
		suite.cp.Encode("sender"),
		suite.cp.Encode("subject"),
		suite.cp.Encode("verbose"),
		suite.cp.Encode(""),
		suite.cp.Encode("terse"),
		suite.cp.Encode("")}

	return resource.MemoryBlockProvider(blocks)
}

func (suite *ElectronicMessageSuite) vanillaStubMail() resource.BlockProvider {
	// The string resources contain a few mails which aren't used.
	// They are missing the terminating line for the verbose text.
	blocks := [][]byte{
		suite.cp.Encode(""),
		suite.cp.Encode(""),
		suite.cp.Encode(""),
		suite.cp.Encode(""),
		suite.cp.Encode("stub email"),
		suite.cp.Encode("stub email"),
		suite.cp.Encode("")}

	return resource.MemoryBlockProvider(blocks)
}

func (suite *ElectronicMessageSuite) holderWithMissingTerminatingLine() resource.BlockProvider {
	// This case is encountered once in gerstrng.res
	blocks := [][]byte{
		suite.cp.Encode(""),
		suite.cp.Encode("title"),
		suite.cp.Encode("sender"),
		suite.cp.Encode("subject"),
		suite.cp.Encode("verbose "),
		suite.cp.Encode("text"),
		suite.cp.Encode(""),
		suite.cp.Encode("terse "),
		suite.cp.Encode("text")}

	return resource.MemoryBlockProvider(blocks)
}
