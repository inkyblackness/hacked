package object

import (
	"errors"

	"github.com/inkyblackness/hacked/ss1/serial"
)

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

// Code serializes the table with given coder.
func (table PropertiesTable) Code(coder serial.Coder) {
	version := propertiesFileVersion
	coder.Code(&version)
	for _, class := range table {
		for _, subclass := range class {
			for _, objType := range subclass {
				coder.Code(objType.Generic)
			}
		}
		for _, subclass := range class {
			for _, objType := range subclass {
				coder.Code(objType.Specific)
			}
		}
	}
	for _, class := range table {
		for _, subclass := range class {
			for _, objType := range subclass {
				coder.Code(objType.Common)
			}
		}
	}
}

// TriplesInClass returns all triples that are available in the given class.
func (table PropertiesTable) TriplesInClass(class Class) []Triple {
	var triples []Triple
	if int(class) < len(table) {
		subclasses := table[class]
		for subclass, subclassEntry := range subclasses {
			for objType := range subclassEntry {
				triples = append(triples, TripleFrom(int(class), subclass, objType))
			}
		}
	}
	return triples
}

// TripleIndex returns the linear index of the given index.
func (table PropertiesTable) TripleIndex(triple Triple) int {
	if int(triple.Class) >= len(table) {
		return -1
	}
	counter := 0
	for class := Class(0); class < triple.Class; class++ {
		subclasses := table[class]
		for _, types := range subclasses {
			counter += len(types)
		}
	}
	subclasses := table[triple.Class]
	if int(triple.Subclass) >= len(subclasses) {
		return -1
	}
	for subclass := Subclass(0); subclass < triple.Subclass; subclass++ {
		counter += len(subclasses[subclass])
	}
	types := subclasses[triple.Subclass]
	if int(triple.Type) >= len(types) {
		return -1
	}
	return counter + int(triple.Type)
}
