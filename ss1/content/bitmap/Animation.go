package bitmap

import (
	"fmt"
	"io"

	"github.com/inkyblackness/hacked/ss1"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial"
)

const (
	animationEndTag     = 0x0C
	animationEndTagData = 0x01

	animationEntryTag = 0x04

	errEndTagDataWrong ss1.StringError = "end tag data wrong"
)

// Animation describes a sequence of bitmaps to form a small movie.
type Animation struct {
	Width      int16
	Height     int16
	ResourceID resource.ID
	IntroFlag  uint16
	Entries    []AnimationEntry
}

// AnimationEntry describes one part of an animation with common frame time.
type AnimationEntry struct {
	FirstFrame byte
	LastFrame  byte
	FrameTime  int16
}

// ReadAnimation extracts an animation from a serialized stream.
func ReadAnimation(reader io.Reader) (anim Animation, err error) {
	decoder := serial.NewDecoder(reader)
	codeAnimationBase(decoder, &anim)
	done := false
	err = decoder.FirstError()
	for !done && (err == nil) {
		var tag byte
		decoder.Code(&tag)
		err = decoder.FirstError()
		switch tag {
		case animationEntryTag:
			var entry AnimationEntry
			decoder.Code(&entry)
			anim.Entries = append(anim.Entries, entry)
			err = decoder.FirstError()
		case animationEndTag:
			var endTagData byte
			decoder.Code(&endTagData)
			if endTagData != animationEndTagData {
				err = errEndTagDataWrong
			}
			done = true
		default:
			if err == nil {
				err = fmt.Errorf("unknown entry tag 0x%02X", int(tag))
			}
		}
	}

	return
}

// WriteAnimation serializes an animation to a writer.
func WriteAnimation(writer io.Writer, anim Animation) error {
	encoder := serial.NewEncoder(writer)
	codeAnimationBase(encoder, &anim)
	for _, entry := range anim.Entries {
		encoder.Code([]byte{animationEntryTag})
		encoder.Code(&entry) // nolint: scopelint
	}
	encoder.Code([]byte{animationEndTag, animationEndTagData})

	return encoder.FirstError()
}

func codeAnimationBase(coder serial.Coder, anim *Animation) {
	coder.Code(&anim.Width)
	coder.Code(&anim.Height)
	coder.Code(&anim.ResourceID)
	var unknown [6]byte
	coder.Code(&unknown)
	coder.Code(&anim.IntroFlag)
}
