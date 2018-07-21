package wav

type riffChunkType uint32
type riffContentType uint32

const (
	riffChunkTypeRiff = riffChunkType(0x46464952)
	riffChunkTypeFmt  = riffChunkType(0x20746d66)
	riffChunkTypeData = riffChunkType(0x61746164)
)

const (
	riffContentTypeWave = riffContentType(0x45564157)
)

type riffChunkTag struct {
	ChunkType riffChunkType
	Size      uint32
}
