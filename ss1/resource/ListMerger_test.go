package resource

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ListMergerSuite struct {
	suite.Suite

	merger listMerger
}

func TestListMergerSuite(t *testing.T) {
	suite.Run(t, new(ListMergerSuite))
}

func (suite *ListMergerSuite) SetupTest() {
	suite.merger = listMerger{}
}

func (suite *ListMergerSuite) TestBlockReturnsErrorForIndexOutOfRange() {
	suite.thenBlockShouldReturnErrorFor(0)
	suite.thenBlockShouldReturnErrorFor(-1)
}

func (suite *ListMergerSuite) TestBlockReturnsEmptyReaderIfAllLayersAreEmpty() {
	suite.whenEntryStores([]byte{})
	suite.thenBlockShouldReturnFor(0, []byte{})
}

func (suite *ListMergerSuite) thenBlockShouldReturnErrorFor(index int) {
	_, err := suite.merger.Block(index)
	assert.NotNil(suite.T(), err, "Error expected")
}

func (suite *ListMergerSuite) whenEntryStores(data ...[]byte) {
	resource := &Resource{
		Properties: Properties{Compound: len(data) != 1},
		Blocks:     BlocksFrom(data),
	}
	suite.merger.list = append(suite.merger.list, resource.ToView())
}

func (suite *ListMergerSuite) thenBlockShouldReturnFor(index int, expected []byte) {
	reader, readerErr := suite.merger.Block(index)
	require.Nil(suite.T(), readerErr, "No error expected retrieving block")
	require.NotNil(suite.T(), reader, "Reader expected")
	data, dataErr := ioutil.ReadAll(reader)
	require.Nil(suite.T(), dataErr, "No error expected reading block")
	assert.Equal(suite.T(), expected, data, "Data mismatch")
}
