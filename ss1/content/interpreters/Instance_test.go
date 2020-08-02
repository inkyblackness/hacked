package interpreters_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/interpreters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type InstanceSuite struct {
	suite.Suite

	data []byte
	inst *interpreters.Instance
}

func TestInstanceSuite(t *testing.T) {
	suite.Run(t, new(InstanceSuite))
}

func (suite *InstanceSuite) SetupTest() {
	sub1 := interpreters.New().
		With("subField0", 0, 1).
		With("subField1", 1, 2)

	sub2 := interpreters.New().
		With("subFieldA", 0, 2).
		With("subFieldB", 2, 1)

	desc := interpreters.New().
		With("field0", 0, 1).
		With("field1", 1, 1).
		With("field2", 2, 2).
		With("field3", 4, 4).
		With("misaligned", 9, 2).
		With("beyond", 256, 16).
		Refining("sub1", 3, 3, sub1, interpreters.Always).
		Refining("sub2", 6, 3, sub2, func(inst *interpreters.Instance) bool { return inst.Get("field0") == 0 }).
		Refining("sub3", 7, 1, interpreters.New(), interpreters.Always)

	suite.data = []byte{0x01, 0x5A, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A}
	suite.inst = desc.For(suite.data)
}

func (suite *InstanceSuite) TestGetReturnsZeroForUnknownKey() {
	result := suite.inst.Get("unknown")

	assert.Equal(suite.T(), uint32(0), result)
}

func (suite *InstanceSuite) TestGetReturnsValueLittleEndian() {
	result := suite.inst.Get("field2")

	assert.Equal(suite.T(), uint32(0x0403), result)
}

func (suite *InstanceSuite) TestGetReturnsZeroForKeyBeyondSize() {
	result := suite.inst.Get("beyond")

	assert.Equal(suite.T(), uint32(0), result)
}

func (suite *InstanceSuite) TestSetIgnoresMisalignedFields() {
	suite.inst.Set("misaligned", 0xEEFF)

	assert.Equal(suite.T(), []byte{0x09, 0x0A}, suite.data[8:10])
}

func (suite *InstanceSuite) TestSetStoresValue() {
	suite.inst.Set("field3", 0xAABBCCDD)

	assert.Equal(suite.T(), []byte{0xDD, 0xCC, 0xBB, 0xAA}, suite.data[4:8])
}

func (suite *InstanceSuite) TestRefinedForUnknownKeyReturnsNullObject() {
	refined := suite.inst.Refined("unknown")

	require.NotNil(suite.T(), refined)
	assert.Equal(suite.T(), uint32(0), refined.Get("something"))
}

func (suite *InstanceSuite) TestRefinedReturnsInstanceForSubsection() {
	refined := suite.inst.Refined("sub1")

	assert.Equal(suite.T(), uint32(0x0605), refined.Get("subField1"))
}

func (suite *InstanceSuite) TestRefinedAllowsModificationOfOriginalData() {
	refined := suite.inst.Refined("sub1")
	refined.Set("subField0", 0xAB)

	assert.Equal(suite.T(), byte(0xAB), suite.data[3])
}

func (suite *InstanceSuite) TestRefinedReturnsInstanceEvenIfNotActive() {
	refined := suite.inst.Refined("sub2")

	assert.Equal(suite.T(), uint32(0x0807), refined.Get("subFieldA"))
}

func (suite *InstanceSuite) TestKeysReturnsListOfKeysSortedByStartIndex() {
	keys := suite.inst.Keys()

	assert.Equal(suite.T(), []string{"field0", "field1", "field2", "field3", "misaligned", "beyond"}, keys)
}

func (suite *InstanceSuite) TestActiveRefinementsReturnsListOfActiveKeysSortedByStartIndex() {
	keys := suite.inst.ActiveRefinements()

	assert.Equal(suite.T(), []string{"sub1", "sub3"}, keys)
}

func (suite *InstanceSuite) TestActiveRefinementsCanChange() {
	suite.inst.Set("field0", 0)
	keys := suite.inst.ActiveRefinements()

	assert.Equal(suite.T(), []string{"sub1", "sub2", "sub3"}, keys)
}

func (suite *InstanceSuite) TestRawReturnsData() {
	data := suite.inst.Refined("sub2").Raw()

	assert.Equal(suite.T(), []byte{0x07, 0x08, 0x09}, data)
}

func (suite *InstanceSuite) TestRawAllowsModificationOfOriginalData() {
	suite.inst.Refined("sub2").Raw()[1] = 0xEF

	assert.Equal(suite.T(), byte(0xEF), suite.data[7])
}

func (suite *InstanceSuite) TestUndefinedReturnsUndefinedBits() {
	data := suite.inst.Undefined()

	assert.Equal(suite.T(), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF}, data)
}

func (suite *InstanceSuite) TestUndefinedConsidersActiveRefinements() {
	suite.inst.Set("field0", 0)
	data := suite.inst.Undefined()

	assert.Equal(suite.T(), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF}, data)
}
