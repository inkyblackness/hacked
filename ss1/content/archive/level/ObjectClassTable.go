package level

// ObjectClassEntry describes an entry in a object-class specific list.
type ObjectClassEntry struct {
	ObjectID ObjectID
	Next     int16
	Prev     int16
	Data     []byte
}

// ObjectClassTable is a list of entries.
type ObjectClassTable []ObjectClassEntry
