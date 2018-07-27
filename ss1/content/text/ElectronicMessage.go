package text

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/inkyblackness/hacked/ss1/resource"
)

const (
	metaExpressionInterrupt    = "(t)"
	metaExpressionNextMessage  = "(?:i([0-9a-fA-F]{2}))"
	metaExpressionColorIndex   = "(?:c([0-9a-fA-F]{2}))"
	metaExpressionLeftDisplay  = "([0-9]+)"
	metaExpressionRightDisplay = "(?:,[ ]*([0-9]+))"
)

var metaExpression = regexp.MustCompile("^[ ]*(?:(?:" +
	metaExpressionInterrupt + "|" +
	metaExpressionNextMessage + "|" +
	metaExpressionColorIndex + "|" +
	metaExpressionLeftDisplay + "|" +
	metaExpressionRightDisplay + ")[ ]*)*$")

// ElectronicMessage describes one message.
type ElectronicMessage struct {
	// NextMessage identifies the message that should follow this one. -1 for none.
	NextMessage int
	// IsInterrupt marks messages that are the ones following another.
	IsInterrupt bool
	// ColorIndex is into the palette for the sender text. -1 for default.
	ColorIndex int
	// LeftDisplay indicates the image to show in the left MFD. -1 for none.
	LeftDisplay int
	// RightDisplay indicates the image to show in the right MFD. -1 for none.
	RightDisplay int

	// Title of the message.
	Title string
	// Sender of the message.
	Sender string
	// Subject of the message.
	Subject string
	// VerboseText the long form.
	VerboseText string
	// TerseText the short form.
	TerseText string
}

// EmptyElectronicMessage returns an instance of an empty electronic message.
func EmptyElectronicMessage() ElectronicMessage {
	message := ElectronicMessage{
		NextMessage:  -1,
		ColorIndex:   -1,
		LeftDisplay:  -1,
		RightDisplay: -1,
	}

	return message
}

// DecodeElectronicMessage tries to decode a message from given block holder.
func DecodeElectronicMessage(cp Codepage, provider resource.BlockProvider) (message ElectronicMessage, err error) {
	blockIndex := 0
	nextBlockString := func() (line string) {
		if err != nil {
			return
		}
		if blockIndex < provider.BlockCount() {
			blockReader, blockErr := provider.Block(blockIndex)
			if blockErr != nil {
				err = fmt.Errorf("failed to access block %v: %v", blockIndex, blockErr)
				return
			}
			data, dataErr := ioutil.ReadAll(blockReader)
			if dataErr != nil {
				err = fmt.Errorf("failed to read data from block %v: %v", blockIndex, dataErr)
			}
			line = cp.Decode(data)
			blockIndex++
		}
		return
	}
	nextText := func() (text string) {
		for line := nextBlockString(); len(line) > 0; line = nextBlockString() {
			text += line
		}
		return
	}

	metaString := nextBlockString()
	metaData := metaExpression.FindStringSubmatch(metaString)
	parseInt := func(metaIndex int, base, bits int) (result int) {
		var value uint64
		result = -1
		if (err == nil) && (len(metaData[metaIndex]) > 0) {
			value, err = strconv.ParseUint(metaData[metaIndex], base, bits)
			if err == nil {
				result = int(value)
			}
		}
		return
	}

	message = EmptyElectronicMessage()
	if (metaData == nil) || (len(metaData[0]) != len(metaString)) {
		err = fmt.Errorf("failed to parse meta string: <%v>", metaString)
	}
	if (err == nil) && (len(metaData[1]) > 0) {
		message.IsInterrupt = true
	}
	message.NextMessage = parseInt(2, 16, 16)
	message.ColorIndex = parseInt(3, 16, 8)
	message.LeftDisplay = parseInt(4, 10, 16)
	message.RightDisplay = parseInt(5, 10, 16)
	message.Title = nextBlockString()
	message.Sender = nextBlockString()
	message.Subject = nextBlockString()
	message.VerboseText = nextText()
	message.TerseText = nextText()

	return
}

// Encode serializes the message into a block holder.
func (message ElectronicMessage) Encode(cp Codepage) [][]byte {
	var blocks [][]byte

	blocks = append(blocks, cp.Encode(message.metaString()))
	blocks = append(blocks, cp.Encode(message.Title))
	blocks = append(blocks, cp.Encode(message.Sender))
	blocks = append(blocks, cp.Encode(message.Subject))
	for _, line := range Blocked(message.VerboseText) {
		blocks = append(blocks, cp.Encode(line))
	}
	for _, line := range Blocked(message.TerseText) {
		blocks = append(blocks, cp.Encode(line))
	}

	return blocks
}

func (message ElectronicMessage) metaString() string {
	result := ""
	add := func(sep, part string) {
		if len(result) > 0 {
			result += sep
		}
		result += part
	}

	if message.IsInterrupt {
		add("", "t")
	}
	if message.NextMessage >= 0 {
		add(" ", fmt.Sprintf("i%02X", message.NextMessage))
	}
	if message.ColorIndex >= 0 {
		add(" ", fmt.Sprintf("c%02X", message.ColorIndex))
	}
	if message.LeftDisplay >= 0 {
		add(" ", fmt.Sprintf("%d", message.LeftDisplay))
	}
	if message.RightDisplay >= 0 {
		add("", fmt.Sprintf(",%d", message.RightDisplay))
	}

	return result
}
