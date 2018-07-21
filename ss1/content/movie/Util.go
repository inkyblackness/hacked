package movie

const fractionDivisor = float32(0x10000)

func timeFromRaw(second byte, fraction uint16) float32 {
	return float32(second) + float32(fraction)/fractionDivisor
}

func timeToRaw(time float32) (second byte, fraction uint16) {
	second = byte(time)
	fraction = uint16((time - float32(second)) * fractionDivisor)
	return
}
