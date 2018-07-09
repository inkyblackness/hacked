package object

// StandardPropertiesTable returns a properties table based on the standard
// configuration of the existing objprop.dat file.
func StandardPropertiesTable() PropertiesTable {
	return NewPropertiesTable(StandardDescriptors())
}

// StandardProperties returns an array of class descriptors that represent the standard
// configuration of the existing objprop.dat file.
func StandardDescriptors() Descriptors {
	var result Descriptors

	{ // Guns
		subclasses := []SubclassDescriptor{
			{TypeCount: 5, SpecificDataSize: 1},
			{TypeCount: 2, SpecificDataSize: 1},
			{TypeCount: 2, SpecificDataSize: 16},
			{TypeCount: 2, SpecificDataSize: 13},
			{TypeCount: 3, SpecificDataSize: 13},
			{TypeCount: 2, SpecificDataSize: 18},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 2, Subclasses: subclasses})
	}
	{ // Ammo
		subclasses := []SubclassDescriptor{
			{TypeCount: 2, SpecificDataSize: 1},
			{TypeCount: 2, SpecificDataSize: 1},
			{TypeCount: 3, SpecificDataSize: 1},
			{TypeCount: 2, SpecificDataSize: 1},
			{TypeCount: 2, SpecificDataSize: 1},
			{TypeCount: 2, SpecificDataSize: 1},
			{TypeCount: 2, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 14, Subclasses: subclasses})
	}
	{ // Physics
		subclasses := []SubclassDescriptor{
			{TypeCount: 6, SpecificDataSize: 20},
			{TypeCount: 16, SpecificDataSize: 6},
			{TypeCount: 2, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 1, Subclasses: subclasses})
	}
	{ // Grenades
		subclasses := []SubclassDescriptor{
			{TypeCount: 5, SpecificDataSize: 1},
			{TypeCount: 3, SpecificDataSize: 3},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 15, Subclasses: subclasses})
	}
	{ // Drugs
		subclasses := []SubclassDescriptor{
			{TypeCount: 7, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 22, Subclasses: subclasses})
	}
	{ // Hardware
		subclasses := []SubclassDescriptor{
			{TypeCount: 5, SpecificDataSize: 1},
			{TypeCount: 10, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 9, Subclasses: subclasses})
	}
	{ // Software
		subclasses := []SubclassDescriptor{
			{TypeCount: 7, SpecificDataSize: 1},
			{TypeCount: 3, SpecificDataSize: 1},
			{TypeCount: 4, SpecificDataSize: 1},
			{TypeCount: 5, SpecificDataSize: 1},
			{TypeCount: 3, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 5, Subclasses: subclasses})
	}
	{ // BigStuff
		subclasses := []SubclassDescriptor{
			{TypeCount: 9, SpecificDataSize: 1},
			{TypeCount: 10, SpecificDataSize: 1},
			{TypeCount: 11, SpecificDataSize: 1},
			{TypeCount: 4, SpecificDataSize: 1},
			{TypeCount: 9, SpecificDataSize: 1},
			{TypeCount: 8, SpecificDataSize: 1},
			{TypeCount: 16, SpecificDataSize: 1},
			{TypeCount: 10, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 2, Subclasses: subclasses})
	}
	{ // SmallStuff
		subclasses := []SubclassDescriptor{
			{TypeCount: 8, SpecificDataSize: 1},
			{TypeCount: 10, SpecificDataSize: 1},
			{TypeCount: 15, SpecificDataSize: 1},
			{TypeCount: 6, SpecificDataSize: 1},
			{TypeCount: 12, SpecificDataSize: 1},
			{TypeCount: 12, SpecificDataSize: 6},
			{TypeCount: 9, SpecificDataSize: 1},
			{TypeCount: 8, SpecificDataSize: 2},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 2, Subclasses: subclasses})
	}
	{ // Fixtures
		subclasses := []SubclassDescriptor{
			{TypeCount: 9, SpecificDataSize: 1},
			{TypeCount: 7, SpecificDataSize: 1},
			{TypeCount: 3, SpecificDataSize: 1},
			{TypeCount: 11, SpecificDataSize: 1},
			{TypeCount: 2, SpecificDataSize: 0}, // skipped vending machines
			{TypeCount: 3, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 1, Subclasses: subclasses})
	}
	{ // Doors
		subclasses := []SubclassDescriptor{
			{TypeCount: 10, SpecificDataSize: 1},
			{TypeCount: 9, SpecificDataSize: 1},
			{TypeCount: 7, SpecificDataSize: 1},
			{TypeCount: 5, SpecificDataSize: 1},
			{TypeCount: 10, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 1, Subclasses: subclasses})
	}
	{ // Animating
		subclasses := []SubclassDescriptor{
			{TypeCount: 9, SpecificDataSize: 1},
			{TypeCount: 11, SpecificDataSize: 1},
			{TypeCount: 14, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 2, Subclasses: subclasses})
	}
	{ // Traps
		subclasses := []SubclassDescriptor{
			{TypeCount: 13, SpecificDataSize: 1},
			{TypeCount: 1, SpecificDataSize: 2},
			{TypeCount: 5, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 1, Subclasses: subclasses})
	}
	{ // Containers
		subclasses := []SubclassDescriptor{
			{TypeCount: 3, SpecificDataSize: 1},
			{TypeCount: 3, SpecificDataSize: 1},
			{TypeCount: 4, SpecificDataSize: 1},
			{TypeCount: 8, SpecificDataSize: 1},
			{TypeCount: 13, SpecificDataSize: 1},
			{TypeCount: 7, SpecificDataSize: 1},
			{TypeCount: 8, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 3, Subclasses: subclasses})
	}
	{ // Critters
		subclasses := []SubclassDescriptor{
			{TypeCount: 9, SpecificDataSize: 3},
			{TypeCount: 12, SpecificDataSize: 1},
			{TypeCount: 7, SpecificDataSize: 1},
			{TypeCount: 7, SpecificDataSize: 6},
			{TypeCount: 2, SpecificDataSize: 1},
		}

		result = append(result, ClassDescriptor{GenericDataSize: 75, Subclasses: subclasses})
	}

	return result
}
