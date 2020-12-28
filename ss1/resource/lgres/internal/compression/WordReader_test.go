package compression

import (
	"bytes"
	"testing"

	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WordReaderSuite struct {
	suite.Suite
	reader *wordReader
}

func TestWordReaderSuite(t *testing.T) {
	suite.Run(t, new(WordReaderSuite))
}

func (suite *WordReaderSuite) TestReadCanReadFirstWord() {
	suite.givenReaderFrom()

	assert.Equal(suite.T(), endOfStream, suite.reader.read())
}

func (suite *WordReaderSuite) TestReadCanReadSeveralWords() {
	input := []word{word(0x3FFF), word(0x0000), word(0x3FFF), word(0x0000), word(0x2001), word(0x1234)}
	suite.givenReaderFrom(input...)

	output := make([]word, len(input)+1)
	for i := 0; i < len(output); i++ {
		output[i] = suite.reader.read()
	}

	assert.Equal(suite.T(), append(input, endOfStream), output)
}

func (suite *WordReaderSuite) givenReaderFrom(words ...word) {
	store := serial.NewByteStore()
	encoder := serial.NewEncoder(store)
	writer := newWordWriter(encoder)

	for _, value := range words {
		writer.write(value)
	}
	writer.close()

	source := bytes.NewReader(store.Data())
	suite.reader = newWordReader(serial.NewDecoder(source))
}
