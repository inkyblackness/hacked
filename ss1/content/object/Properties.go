package object

import "errors"

// PropertiesTable is a collection of class-specific properties.
type PropertiesTable []ClassProperties

// ClassProperties is a collection of subclass-specific properties.
type ClassProperties []SubclassProperties

// SubclassProperties is a collection of type-specific properties.
type SubclassProperties []Properties

// Properties contains the object-specific properties.
type Properties struct {
	Common   []byte
	Generic  []byte
	Specific []byte
}

// NewPropertiesTable returns a new instance based on given descriptors.
func NewPropertiesTable(desc Descriptors) PropertiesTable {
	classCount := len(desc)
	table := make([]ClassProperties, classCount)
	for class := 0; class < classCount; class++ {
		classDesc := desc[class]
		subclassCount := len(classDesc.Subclasses)
		subclassProperties := make([]SubclassProperties, subclassCount)
		table[class] = subclassProperties
		for subclass := 0; subclass < subclassCount; subclass++ {
			subclassDesc := classDesc.Subclasses[subclass]
			typeProperties := make([]Properties, subclassDesc.TypeCount)
			subclassProperties[subclass] = typeProperties
			for objType := 0; objType < subclassDesc.TypeCount; objType++ {
				prop := &typeProperties[objType]
				prop.Common = make([]byte, CommonPropertiesSize)
				prop.Generic = make([]byte, classDesc.GenericDataSize)
				prop.Specific = make([]byte, subclassDesc.SpecificDataSize)
			}
		}
	}
	return table
}

// ForObject returns the object-specific properties by given triple.
func (table PropertiesTable) ForObject(triple Triple) (Properties, error) {
	if int(triple.Class) >= len(table) {
		return Properties{}, errors.New("invalid class")
	}
	classEntry := table[triple.Class]
	if int(triple.Subclass) >= len(classEntry) {
		return Properties{}, errors.New("invalid subclass")
	}
	subclassEntry := classEntry[triple.Subclass]
	if int(triple.Type) >= len(subclassEntry) {
		return Properties{}, errors.New("invalid type")
	}
	return subclassEntry[triple.Type], nil
}
