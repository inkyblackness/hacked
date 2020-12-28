package object_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/object"
)

func TestClassDescriptorTotalDataLengthReturnsCompleteLength(t *testing.T) {
	var mainDesc object.ClassDescriptor

	mainDesc.GenericDataSize = 7
	mainDesc.Subclasses = append(mainDesc.Subclasses,
		object.SubclassDescriptor{TypeCount: 2, SpecificDataSize: 3},
		object.SubclassDescriptor{TypeCount: 1, SpecificDataSize: 20})

	assert.Equal(t, (7*(2+1))+(2*3)+(1*20)+object.CommonPropertiesSize*3, mainDesc.TotalDataSize())
}

func TestClassDescriptorTotalTypeCountReturnsTotalAmount(t *testing.T) {
	var mainDesc object.ClassDescriptor

	mainDesc.GenericDataSize = 7
	mainDesc.Subclasses = append(mainDesc.Subclasses,
		object.SubclassDescriptor{TypeCount: 2, SpecificDataSize: 3},
		object.SubclassDescriptor{TypeCount: 1, SpecificDataSize: 20})

	assert.Equal(t, 3, mainDesc.TotalTypeCount())
}
