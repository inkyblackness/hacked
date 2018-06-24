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
	{2, 16},
	{0, 32},
	{34, 32},
	{6, 32},
	{0, 32},
	{1, 8},
	{3, 16},
	{10, 176},
	{10, 128},
	{24, 64},
	{8, 64},
	{4, 32},
	{22, 160},
	{15, 64},
	{40, 64},
}

// ObjectClassInfoFor returns the info entry for the corresponding object class.
func ObjectClassInfoFor(class object.Class) ObjectClassInfo {
	return objectClassInfoList[int(class)]
}
