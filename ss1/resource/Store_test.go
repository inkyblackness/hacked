package resource_test

import (
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IdentifiedResource struct {
	id  resource.ID
	res *resource.Resource
}
type ResourceList []IdentifiedResource

func (list ResourceList) IDs() []resource.ID {
	ids := make([]resource.ID, 0, len(list))
	for _, entry := range list {
		ids = append(ids, entry.id)
	}
	return ids
}

func (list ResourceList) Resource(id resource.ID) (res resource.View, err error) {
	for _, entry := range list {
		if entry.id.Value() == id.Value() {
			res = entry.res
		}
	}
	if res == nil {
		err = fmt.Errorf("unknown id %v", id)
	}
	return res, err
}

type StoreSuite struct {
	suite.Suite

	store           *resource.Store
	resourceCounter int
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}

func (suite *StoreSuite) SetupTest() {
	suite.store = nil
}

func (suite *StoreSuite) TestNewInstanceIsEmpty() {
	suite.whenInstanceIsCreated()
	suite.thenIDsShouldBeEmpty()
}

func (suite *StoreSuite) TestResourceReturnsErrorForUnknownResource() {
	suite.whenInstanceIsCreated()
	suite.thenResourceShouldReturnErrorFor(resource.ID(10))
}

func (suite *StoreSuite) TestDelWillHaveStoreIgnorePreviousEntry() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.whenResourceIsDeleted(resource.ID(2))
	suite.thenIDsShouldBe([]resource.ID{resource.ID(1)})
	suite.thenResourceShouldReturnErrorFor(resource.ID(2))
}

func (suite *StoreSuite) TestPutOverridesPreviousResources() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	newRes := suite.aResource()
	suite.whenResourceIsPut(resource.ID(2), newRes)
	suite.thenIDsShouldBe([]resource.ID{resource.ID(2), resource.ID(1)})
	suite.thenReturnedResourceShouldBe(resource.ID(2), newRes)
}

func (suite *StoreSuite) TestDelWillHaveStoreIgnoreResourceEvenIfPutMultipleTimes() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.whenResourceIsDeleted(resource.ID(2))
	suite.thenIDsShouldBe([]resource.ID{resource.ID(1)})
	suite.thenResourceShouldReturnErrorFor(resource.ID(2))
}

func (suite *StoreSuite) TestPutAddsNewResourcesAtEnd() {
	suite.givenAnInstance()
	suite.givenStoredResource(resource.ID(2), suite.aResource())
	suite.givenStoredResource(resource.ID(1), suite.aResource())
	newRes := suite.aResource()
	suite.whenResourceIsPut(resource.ID(3), newRes)
	suite.thenIDsShouldBe([]resource.ID{resource.ID(2), resource.ID(1), resource.ID(3)})
	suite.thenReturnedResourceShouldBe(resource.ID(3), newRes)
}

func (suite *StoreSuite) givenAnInstance() {
	suite.whenInstanceIsCreated()
}

func (suite *StoreSuite) givenStoredResource(id resource.ID, res *resource.Resource) {
	suite.store.Put(id, res)
}

func (suite *StoreSuite) whenInstanceIsCreated() {
	suite.store = resource.NewStore()
}

func (suite *StoreSuite) whenResourceIsDeleted(id resource.ID) {
	suite.store.Del(id)
}

func (suite *StoreSuite) whenResourceIsPut(id resource.ID, res *resource.Resource) {
	suite.store.Put(id, res)
}

func (suite *StoreSuite) thenIDsShouldBeEmpty() {
	assert.Equal(suite.T(), 0, len(suite.store.IDs()))
}

func (suite *StoreSuite) thenIDsShouldBe(expected []resource.ID) {
	assert.Equal(suite.T(), expected, suite.store.IDs())
}

func (suite *StoreSuite) thenReturnedResourceShouldBe(id resource.ID, expected *resource.Resource) {
	res, err := suite.store.Resource(id)
	assert.Nil(suite.T(), err, "No error expected for ID %v", id)
	assert.Equal(suite.T(), expected, res, "Different res returned for ID %v", id)
}

func (suite *StoreSuite) thenResourceShouldReturnErrorFor(id resource.ID) {
	_, err := suite.store.Resource(id)
	assert.Error(suite.T(), err, "Error expected for ID %v ", id) // nolint: vet
}

func (suite *StoreSuite) aResource() *resource.Resource {
	suite.resourceCounter++
	return &resource.Resource{Blocks: resource.BlocksFrom([][]byte{{byte(suite.resourceCounter)}})}
}
