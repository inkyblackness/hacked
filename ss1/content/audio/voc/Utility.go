package voc

const (
	fileHeader         string  = "Creative Voice File\u001A"
	standardHeaderSize uint16  = 0x1A
	baseVersion        uint16  = 0x010A
	versionCheckValue  uint16  = 0x1234
	rateBase           float32 = 1000000.0
)

func lengthFromBlockStart(blockStart []byte) int {
	return int(blockStart[3])<<16 + int(blockStart[2])<<8 + int(blockStart[1])
}

func divisorToSampleRate(divisor byte) float32 {
	return rateBase / float32(256-int(divisor))
}

func sampleRateToDivisor(sampleRate float32) byte {
	return byte(256 - int(rateBase/sampleRate))
}
