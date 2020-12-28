package compression_test

import (
	"bytes"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource/lgres/internal/compression"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WordReaderSuite struct {
	suite.Suite
	reader *compression.WordReader
}

func TestWordReaderSuite(t *testing.T) {
	suite.Run(t, new(WordReaderSuite))
}

func (suite *WordReaderSuite) TestReadCanReadFirstWord() {
	suite.givenReaderFrom()

	assert.Equal(suite.T(), compression.EndOfStream, suite.reader.Read())
}

func (suite *WordReaderSuite) TestReadCanReadSeveralWords() {
	input := []compression.Word{
		compression.Word(0x3FFF), compression.Word(0x0000), compression.Word(0x3FFF),
		compression.Word(0x0000), compression.Word(0x2001), compression.Word(0x1234),
	}
	suite.givenReaderFrom(input...)

	output := make([]compression.Word, len(input)+1)
	for i := 0; i < len(output); i++ {
		output[i] = suite.reader.Read()
	}

	assert.Equal(suite.T(), append(input, compression.EndOfStream), output)
}

func (suite *WordReaderSuite) givenReaderFrom(words ...compression.Word) {
	store := serial.NewByteStore()
	encoder := serial.NewEncoder(store)
	writer := compression.NewWordWriter(encoder)

	for _, value := range words {
		writer.Write(value)
	}
	writer.Close()

	source := bytes.NewReader(store.Data())
	suite.reader = compression.NewWordReader(serial.NewDecoder(source))
}
