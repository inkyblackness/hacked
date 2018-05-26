package lgres

const (
	headerString                   = "LG Res File v2\r\n"
	commentTerminator              = byte(0x1A)
	resourceDirectoryFileOffsetPos = 0x7C
	boundarySize                   = 4

	resourceTypeFlagCompressed = byte(0x01)
	resourceTypeFlagCompound   = byte(0x02)
)
