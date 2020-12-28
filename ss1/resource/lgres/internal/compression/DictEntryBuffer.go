package compression

// dictEntryBucketSizeBits governs how big one bucket in the buffer is.
const dictEntryBucketSizeBits = 8

type dictEntryBucket [1 << dictEntryBucketSizeBits]dictEntry

type dictEntryBuffer struct {
	buckets []*dictEntryBucket
}

func (buffer *dictEntryBuffer) entry(key word) *dictEntry {
	bucketIndex := int(key) >> dictEntryBucketSizeBits
	index := int(key & ((1 << dictEntryBucketSizeBits) - 1))

	for bucketIndex >= len(buffer.buckets) {
		buffer.buckets = append(buffer.buckets, new(dictEntryBucket))
	}

	return &buffer.buckets[bucketIndex][index]
}
