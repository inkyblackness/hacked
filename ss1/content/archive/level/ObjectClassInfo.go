package level

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
)

// ObjectClassInfo contains meta information about a level object class in an archive.
type ObjectClassInfo struct {
	DataSize   int
	EntryCount int
}

var objectClassInfoList = []ObjectClassInfo{
	{DataSize: 2, EntryCount: 16},
	{DataSize: 0, EntryCount: 32},
	{DataSize: 34, EntryCount: 32},
	{DataSize: 6, EntryCount: 32},
	{DataSize: 0, EntryCount: 32},
	{DataSize: 1, EntryCount: 8},
	{DataSize: 3, EntryCount: 16},
	{DataSize: 10, EntryCount: 176},
	{DataSize: 10, EntryCount: 128},
	{DataSize: 24, EntryCount: 64},
	{DataSize: 8, EntryCount: 64},
	{DataSize: 4, EntryCount: 32},
	{DataSize: 22, EntryCount: 160},
	{DataSize: 15, EntryCount: 64},
	{DataSize: 40, EntryCount: 64},
}

// ObjectClassInfoFor returns the info entry for the corresponding object class.
func ObjectClassInfoFor(class object.Class) ObjectClassInfo {
	return objectClassInfoList[int(class)]
}
