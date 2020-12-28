package resource_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ResourceSelectorSuite struct {
	suite.Suite

	resources []resource.View
	view      resource.View

	isCompoundList bool
	viewStrategy   resource.ViewStrategy
}

func TestResourceSelectorSuite(t *testing.T) {
	suite.Run(t, new(ResourceSelectorSuite))
}

func (suite *ResourceSelectorSuite) SetupTest() {
	suite.resources = nil
	suite.view = nil
	suite.isCompoundList = false
	suite.viewStrategy = nil
}

func (suite *ResourceSelectorSuite) TestBlockReturnsData() {
	suite.givenResource(resource.BlocksFrom([][]byte{{0x01}}))
	suite.whenInstanceIsCreated()
	suite.thenResourceBlockShouldBe(0, []byte{0x01})
}

func (suite *ResourceSelectorSuite) TestBlockReturnsDataFromLastEntryByDefault() {
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Blocks = resource.BlocksFrom([][]byte{{0xAA}})
	})
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Blocks = resource.BlocksFrom([][]byte{{0xBB}})
	})
	suite.whenInstanceIsCreated()
	suite.thenResourceBlockShouldBe(0, []byte{0xBB})
}

func (suite *ResourceSelectorSuite) TestBlockReturnsDataFromLastEntryForCompoundNonListResources() {
	suite.givenViewStrategyIsSet()
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = true
		res.Blocks = resource.BlocksFrom([][]byte{{0xAA}})
	})
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = true
		res.Blocks = resource.BlocksFrom([][]byte{{0xBB}, {0xCC}})
	})
	suite.whenInstanceIsCreated()
	suite.thenResourceBlockCountShouldBe(2)
	suite.thenResourceBlockShouldBe(0, []byte{0xBB})
	suite.thenResourceBlockShouldBe(1, []byte{0xCC})
}

func (suite *ResourceSelectorSuite) TestBlockReturnsDataFromLastNonEmptyEntryIfACompoundList() {
	suite.givenResourceIsACompoundList()
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = true
		res.Blocks = resource.BlocksFrom([][]byte{{0xAA}})
	})
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = true
		res.Blocks = resource.BlocksFrom([][]byte{{}})
	})
	suite.whenInstanceIsCreated()
	suite.thenResourceBlockShouldBe(0, []byte{0xAA})
}

func (suite *ResourceSelectorSuite) TestBlockCountOfACompoundListReturnsHighestCount() {
	suite.givenResourceIsACompoundList()
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = true
		res.Blocks = resource.BlocksFrom([][]byte{{}, {}, {0x11}})
	})
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = true
		res.Blocks = resource.BlocksFrom([][]byte{{}, {0x22}})
	})
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = true
		res.Blocks = resource.BlocksFrom([][]byte{{}, {}, {}, {0x33}})
	})
	suite.whenInstanceIsCreated()
	suite.thenResourceBlockCountShouldBe(4)
}

func (suite *ResourceSelectorSuite) TestMetaOfCompoundListIsThatOfFirst() {
	suite.givenResourceIsACompoundList()
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = true
		res.Properties.Compressed = true
		res.Properties.ContentType = resource.Movie
		res.Blocks = resource.BlocksFrom([][]byte{{}, {}, {0x11}})
	})
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = false
		res.Properties.Compressed = false
		res.Properties.ContentType = resource.Text
		res.Blocks = resource.BlocksFrom([][]byte{{}, {0x22}})
	})
	suite.whenInstanceIsCreated()
	suite.thenResourceShouldHaveMeta(true, resource.Movie, true)
}

func (suite *ResourceSelectorSuite) TestMetaOfNonCompoundListIsThatOfLast() {
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = true
		res.Properties.Compressed = true
		res.Properties.ContentType = resource.Movie
		res.Blocks = resource.BlocksFrom([][]byte{{}, {}, {0x11}})
	})
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Properties.Compound = false
		res.Properties.Compressed = false
		res.Properties.ContentType = resource.Text
		res.Blocks = resource.BlocksFrom([][]byte{{}, {0x22}})
	})
	suite.whenInstanceIsCreated()
	suite.thenResourceShouldHaveMeta(false, resource.Text, false)
}

func (suite *ResourceSelectorSuite) TestBlockReturnsErrorForIndexOutOfRange() {
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Blocks = resource.BlocksFrom(nil)
	})
	suite.whenInstanceIsCreated()
	suite.thenBlockShouldReturnErrorFor(0)
	suite.thenBlockShouldReturnErrorFor(-1)
}

func (suite *ResourceSelectorSuite) TestBlockReturnsEmptyReaderIfAllLayersAreEmpty() {
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Blocks = resource.BlocksFrom([][]byte{{}})
	})
	suite.whenInstanceIsCreated()
	suite.thenBlockShouldReturnFor(0, []byte{})
}

func (suite *ResourceSelectorSuite) givenResourceIsACompoundList() {
	suite.isCompoundList = true
	suite.viewStrategy = suite
}

func (suite *ResourceSelectorSuite) givenViewStrategyIsSet() {
	suite.viewStrategy = suite
}

func (suite *ResourceSelectorSuite) givenResource(blocks resource.Blocks) {
	suite.givenSpecificResource(func(res *resource.Resource) {
		res.Blocks = blocks
	})
}

func (suite *ResourceSelectorSuite) givenSpecificResource(modifier func(*resource.Resource)) {
	res := resource.Resource{}
	modifier(&res)
	suite.resources = append(suite.resources, res)
}

func (suite *ResourceSelectorSuite) whenInstanceIsCreated() {
	selector := resource.Selector{
		Lang: resource.LangAny,
		From: suite,
		As:   suite.viewStrategy,
	}

	var err error
	suite.view, err = selector.Select(resource.ID(1000))
	require.Nil(suite.T(), err, "No error expected creating view for %v", 1000)
}

func (suite *ResourceSelectorSuite) thenResourceBlockCountShouldBe(expected int) {
	assert.Equal(suite.T(), expected, suite.view.BlockCount(), "Proper block count expected")
}

func (suite *ResourceSelectorSuite) thenResourceBlockShouldBe(block int, expected []byte) {
	reader, err := suite.view.Block(block)
	require.Nil(suite.T(), err, fmt.Sprintf("No error expected retrieving block %v", block))
	data, readErr := ioutil.ReadAll(reader)
	require.Nil(suite.T(), readErr, fmt.Sprintf("No error expected reading block %v", block))
	assert.Equal(suite.T(), expected, data, "Block data expected")
}

func (suite *ResourceSelectorSuite) thenResourceShouldHaveMeta(compound bool, contentType resource.ContentType, compressed bool) {
	assert.Equal(suite.T(), compound, suite.view.Compound(), "Compound state is wrong")
	assert.Equal(suite.T(), contentType, suite.view.ContentType(), "Content type is wrong")
	assert.Equal(suite.T(), compressed, suite.view.Compressed(), "Compression is wrong")
}

func (suite *ResourceSelectorSuite) thenBlockShouldReturnErrorFor(index int) {
	_, err := suite.view.Block(index)
	assert.NotNil(suite.T(), err, "Error expected")
}

func (suite *ResourceSelectorSuite) thenBlockShouldReturnFor(index int, expected []byte) {
	reader, readerErr := suite.view.Block(index)
	require.Nil(suite.T(), readerErr, "No error expected retrieving block")
	require.NotNil(suite.T(), reader, "Reader expected")
	data, dataErr := ioutil.ReadAll(reader)
	require.Nil(suite.T(), dataErr, "No error expected reading block")
	assert.Equal(suite.T(), expected, data, "Data mismatch")
}

func (suite *ResourceSelectorSuite) IsCompoundList(id resource.ID) bool {
	assert.Equal(suite.T(), resource.ID(1000), id, "Unknown resource queried")
	return suite.isCompoundList
}

func (suite *ResourceSelectorSuite) Filter(lang resource.Language, id resource.ID) resource.List {
	return suite.resources
}
