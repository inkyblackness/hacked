package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClassDescriptorTotalDataLengthReturnsCompleteLength(t *testing.T) {
	var mainDesc ClassDescriptor

	mainDesc.GenericDataSize = 7
	mainDesc.Subclasses = append(mainDesc.Subclasses,
		SubclassDescriptor{TypeCount: 2, SpecificDataSize: 3},
		SubclassDescriptor{TypeCount: 1, SpecificDataSize: 20})

	assert.Equal(t, (7*(2+1))+(2*3)+(1*20)+CommonPropertiesSize*3, mainDesc.TotalDataSize())
}

func TestClassDescriptorTotalTypeCountReturnsTotalAmount(t *testing.T) {
	var mainDesc ClassDescriptor

	mainDesc.GenericDataSize = 7
	mainDesc.Subclasses = append(mainDesc.Subclasses,
		SubclassDescriptor{TypeCount: 2, SpecificDataSize: 3},
		SubclassDescriptor{TypeCount: 1, SpecificDataSize: 20})

	assert.Equal(t, 3, mainDesc.TotalTypeCount())
}
