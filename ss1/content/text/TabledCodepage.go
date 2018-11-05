package text

type tabledCodepage struct {
	tableToRune []rune
	tableToByte map[rune]byte
}

func (cp *tabledCodepage) Encode(value string) []byte {
	result := make([]byte, 0, len(value)+1)

	for _, c := range value {
		mapped, known := cp.tableToByte[c]
		if !known {
			mapped = cp.tableToByte['?']
		}
		result = append(result, mapped)
	}
	result = append(result, 0x00)

	return result
}

func (cp *tabledCodepage) Decode(data []byte) string {
	runes := make([]rune, 0, len(data))

	for _, value := range data {
		if value != 0x00 {
			runes = append(runes, cp.tableToRune[value])
		}
	}

	return string(runes)
}
