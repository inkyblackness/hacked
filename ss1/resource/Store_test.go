package resource_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite

	store           resource.Store
	resourceCounter int
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}

func (suite *StoreSuite) SetupTest() {
	suite.store = resource.Store{}
}

func (suite *StoreSuite) TestNewInstanceIsEmpty() {
	suite.whenInstanceIsCreated()
	suite.thenIDsShouldBeEmpty()
}

func (suite *StoreSuite) TestResourceReturnsErrorForUnknownResource() {
	suite.whenInstanceIsCreated()
	suite.thenViewShouldReturnErrorFor(resource.ID(10))
}

func (suite *StoreSuite) TestDelWillHaveStoreIgnorePreviousEntry() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.whenResourceIsDeleted(resource.ID(2))
	suite.thenIDsShouldBe([]resource.ID{resource.ID(1)})
	suite.thenViewShouldReturnErrorFor(resource.ID(2))
}

func (suite *StoreSuite) TestDelWillHaveStoreReportFewerIDs() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.givenStoredResource(resource.ID(3), suite.aResource())
	suite.givenStoredResource(resource.ID(4), suite.aResource())
	suite.whenResourceIsDeleted(resource.ID(2))
	suite.thenIDsShouldBe([]resource.ID{resource.ID(1), resource.ID(3), resource.ID(4)})
}

func (suite *StoreSuite) TestPutOverridesPreviousResources() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	newRes := suite.aResource()
	suite.whenResourceIsPut(resource.ID(2), newRes)
	suite.thenIDsShouldBe([]resource.ID{resource.ID(2), resource.ID(1)})
	suite.thenReturnedViewShouldBe(resource.ID(2), newRes)
}

func (suite *StoreSuite) TestDelWillHaveStoreIgnoreResourceEvenIfPutMultipleTimes() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.whenResourceIsDeleted(resource.ID(2))
	suite.thenIDsShouldBe([]resource.ID{resource.ID(1)})
	suite.thenViewShouldReturnErrorFor(resource.ID(2))
}

func (suite *StoreSuite) TestPutAddsNewResourcesAtEnd() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	newRes := suite.aResource()
	suite.whenResourceIsPut(resource.ID(3), newRes)
	suite.thenIDsShouldBe([]resource.ID{resource.ID(2), resource.ID(1), resource.ID(3)})
	suite.thenReturnedViewShouldBe(resource.ID(3), newRes)
}

func (suite *StoreSuite) TestPutRestoresIDAtOldPosition() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	suite.givenResourceWasDeleted(resource.ID(2))
	newRes := suite.aResource()
	suite.whenResourceIsPut(resource.ID(2), newRes)
	suite.thenIDsShouldBe([]resource.ID{resource.ID(2), resource.ID(1)})
}

func (suite *StoreSuite) givenAnInstance() {
	suite.whenInstanceIsCreated()
}

func (suite *StoreSuite) givenStoredResource(id resource.ID, res *resource.Resource) {
	err := suite.store.Put(id, res)
	require.Nil(suite.T(), err, "No error expected storing resource")
}

func (suite *StoreSuite) givenResourceWasDeleted(id resource.ID) {
	suite.whenResourceIsDeleted(id)
}

func (suite *StoreSuite) whenInstanceIsCreated() {
	suite.store = resource.Store{}
}

func (suite *StoreSuite) whenResourceIsDeleted(id resource.ID) {
	suite.store.Del(id)
}

func (suite *StoreSuite) whenResourceIsPut(id resource.ID, res *resource.Resource) {
	err := suite.store.Put(id, res)
	require.Nil(suite.T(), err, "No error expected putting resource")
}

func (suite *StoreSuite) thenIDsShouldBeEmpty() {
	assert.Equal(suite.T(), 0, len(suite.store.IDs()))
}

func (suite *StoreSuite) thenIDsShouldBe(expected []resource.ID) {
	assert.Equal(suite.T(), expected, suite.store.IDs())
}

func (suite *StoreSuite) thenReturnedViewShouldBe(id resource.ID, expected *resource.Resource) {
	res, err := suite.store.View(id)
	assert.Nil(suite.T(), err, "No error expected for ID %v", id)
	assert.Equal(suite.T(), expected, res, "Different res returned for ID %v", id)
}

func (suite *StoreSuite) thenViewShouldReturnErrorFor(id resource.ID) {
	_, err := suite.store.View(id)
	assert.Error(suite.T(), err, "Error expected for ID %v ", id) // nolint: vet
}

func (suite *StoreSuite) aResource() *resource.Resource {
	suite.resourceCounter++
	return &resource.Resource{Blocks: resource.BlocksFrom([][]byte{{byte(suite.resourceCounter)}})}
}
