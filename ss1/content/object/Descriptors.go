package object

// Descriptors describes a set of object classes.
type Descriptors []ClassDescriptor

// ClassDescriptor describes a single object class.
type ClassDescriptor struct {
	// GenericDataSize specifies the length, in bytes, of one generic type entry.
	GenericDataSize int
	// Subclasses contains descriptions of the subclasses of this class.
	// The index into the array is the subclass ID.
	Subclasses []SubclassDescriptor
}

// TotalDataSize returns the total length the class requires.
func (desc ClassDescriptor) TotalDataSize() int {
	total := 0
	for _, subclass := range desc.Subclasses {
		total += (desc.GenericDataSize * subclass.TypeCount) + subclass.TotalDataSize()
	}
	return total
}

// TotalTypeCount returns the total number of types in this class.
func (desc ClassDescriptor) TotalTypeCount() int {
	total := 0
	for _, subclass := range desc.Subclasses {
		total += subclass.TypeCount
	}
	return total
}

// SubclassDescriptor describes one subclass.
type SubclassDescriptor struct {
	// TypeCount specifies how many types exist in this subclass.
	TypeCount int
	// SpecificDataSize specifies the length, in bytes, of one specific type entry.
	SpecificDataSize int
}

// TotalDataSize returns the total length, in bytes, the subclass requires in the properties file.
func (desc SubclassDescriptor) TotalDataSize() int {
	return (desc.SpecificDataSize + CommonPropertiesSize) * desc.TypeCount
}
