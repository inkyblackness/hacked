package movie

import (
	"bytes"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// Entry describes a block from a MOVI container.
type Entry interface {
	// Timestamp marks the beginning time of the entry, in seconds.
	Timestamp() Timestamp
	// Type describes the content type of the data.
	Type() DataType
	// Data returns the raw bytes of the entry.
	Data() []byte
}

// EntryBase is the core structure for all entries.
type EntryBase struct {
	Time Timestamp
}

// Timestamp helps implement the Entry interface.
func (entry EntryBase) Timestamp() Timestamp {
	return entry.Time
}

// UnknownEntry is an entry that is not know to this codebase.
type UnknownEntry struct {
	EntryBase
	DataType DataType
	Bytes    []byte
}

// UnknownEntryFrom decodes an entry from given data.
func UnknownEntryFrom(timestamp Timestamp, dataType DataType, data []byte) (UnknownEntry, error) {
	entry := UnknownEntry{
		EntryBase: EntryBase{Time: timestamp},
		DataType:  dataType,
		Bytes:     data,
	}
	return entry, nil
}

// Type describes the entry.
func (entry UnknownEntry) Type() DataType {
	return entry.DataType
}

// Data returns the raw bytes of the entry.
func (entry UnknownEntry) Data() []byte {
	return entry.Bytes
}

// EndOfMediaEntry marks the end of a list of entries.
type EndOfMediaEntry struct {
	EntryBase
}

// Type describes the entry.
func (entry EndOfMediaEntry) Type() DataType {
	return dataTypeEndOfMedia
}

// Code serializes the entry.
func (entry EndOfMediaEntry) Data() []byte {
	return nil
}

// LowResVideoEntry is a compressed low-resolution video frame.
type LowResVideoEntry struct {
	EntryBase

	BoundingBox [4]uint16
	Packed      []byte
}

// LowResVideoEntryFrom decodes an entry from given data.
func LowResVideoEntryFrom(timestamp Timestamp, data []byte) (LowResVideoEntry, error) {
	var entry LowResVideoEntry
	entry.Time = timestamp
	coder := serial.NewPositioningDecoder(bytes.NewReader(data))
	coder.Code(&entry.BoundingBox)
	entry.Packed = make([]byte, len(data)-int(coder.CurPos()))
	coder.Code(entry.Packed)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry LowResVideoEntry) Type() DataType {
	return DataTypeLowResVideo
}

// Data returns the raw bytes of the entry.
func (entry LowResVideoEntry) Data() []byte {
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	coder.Code(entry.BoundingBox)
	coder.Code(entry.Packed)
	return buf.Bytes()
}

// HighResVideoEntry is a compressed high-resolution video frame.
type HighResVideoEntry struct {
	EntryBase

	Bitstream  []byte
	Maskstream []byte
}

// HighResVideoEntryFrom decodes an entry from given data.
func HighResVideoEntryFrom(timestamp Timestamp, data []byte) (HighResVideoEntry, error) {
	var entry HighResVideoEntry
	entry.Time = timestamp
	coder := serial.NewDecoder(bytes.NewReader(data))
	var bitstreamLen uint16
	coder.Code(&bitstreamLen)
	entry.Bitstream = make([]byte, bitstreamLen)
	entry.Maskstream = make([]byte, len(data)-int(bitstreamLen))
	coder.Code(entry.Bitstream)
	coder.Code(entry.Maskstream)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry HighResVideoEntry) Type() DataType {
	return DataTypeHighResVideo
}

// Data returns the raw bytes of the entry.
func (entry HighResVideoEntry) Data() []byte {
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	bitstreamLen := uint16(len(entry.Bitstream))
	coder.Code(&bitstreamLen)
	coder.Code(entry.Bitstream)
	coder.Code(entry.Maskstream)
	return buf.Bytes()
}

// AudioEntry is an entry with audio samples.
type AudioEntry struct {
	EntryBase

	Samples []byte
}

// AudioEntryFrom decodes an entry from given data.
func AudioEntryFrom(timestamp Timestamp, data []byte) (AudioEntry, error) {
	var entry AudioEntry
	entry.Time = timestamp
	entry.Samples = data
	return entry, nil
}

// Type describes the entry.
func (entry AudioEntry) Type() DataType {
	return DataTypeAudio
}

// Data returns the raw bytes of the entry.
func (entry AudioEntry) Data() []byte {
	return entry.Samples
}

// SubtitleEntry is an entry with Subtitle samples.
type SubtitleEntry struct {
	EntryBase

	Control SubtitleControl
	Text    []byte
}

// SubtitleEntryFrom decodes an entry from given data.
func SubtitleEntryFrom(timestamp Timestamp, data []byte) (SubtitleEntry, error) {
	var entry SubtitleEntry
	entry.Time = timestamp
	coder := serial.NewDecoder(bytes.NewReader(data))
	var header SubtitleHeader
	coder.Code(&header)
	if header.TextOffset > coder.CurPos() {
		coder.Code(make([]byte, header.TextOffset-coder.CurPos()))
	}
	entry.Text = make([]byte, len(data)-int(coder.CurPos()))
	coder.Code(entry.Text)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry SubtitleEntry) Type() DataType {
	return DataTypeSubtitle
}

// Data returns the raw bytes of the entry.
func (entry SubtitleEntry) Data() []byte {
	header := SubtitleHeader{
		Control:    entry.Control,
		TextOffset: SubtitleDefaultTextOffset,
	}
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	coder.Code(&header)
	coder.Code(make([]byte, int(header.TextOffset-coder.CurPos())))
	coder.Code(entry.Text)
	return buf.Bytes()
}

// PaletteEntry is an entry with a palette.
type PaletteEntry struct {
	EntryBase

	Colors bitmap.Palette
}

// PaletteEntryFrom decodes an entry from given data.
func PaletteEntryFrom(timestamp Timestamp, data []byte) (PaletteEntry, error) {
	var entry PaletteEntry
	entry.Time = timestamp
	coder := serial.NewDecoder(bytes.NewReader(data))
	coder.Code(&entry.Colors)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry PaletteEntry) Type() DataType {
	return DataTypePalette
}

// Data returns the raw bytes of the entry.
func (entry PaletteEntry) Data() []byte {
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	coder.Code(&entry.Colors)
	return buf.Bytes()
}

// PaletteResetEntry marks a reset of the palette.
type PaletteResetEntry struct {
	EntryBase
}

// PaletteResetEntryFrom decodes an entry from given data.
func PaletteResetEntryFrom(timestamp Timestamp, data []byte) (PaletteResetEntry, error) {
	var entry PaletteResetEntry
	entry.Time = timestamp
	return entry, nil
}

// Type describes the entry.
func (entry PaletteResetEntry) Type() DataType {
	return DataTypePaletteReset
}

// Code serializes the entry.
func (entry PaletteResetEntry) Data() []byte {
	return nil
}

// PaletteLookupEntry contains a list of palette indices.
type PaletteLookupEntry struct {
	EntryBase

	List []byte
}

// PaletteLookupEntryFrom decodes an entry from given data.
func PaletteLookupEntryFrom(timestamp Timestamp, data []byte) (PaletteLookupEntry, error) {
	var entry PaletteLookupEntry
	entry.Time = timestamp
	entry.List = data
	return entry, nil
}

// Type describes the entry.
func (entry PaletteLookupEntry) Type() DataType {
	return DataTypePaletteLookupList
}

// Code serializes the entry.
func (entry PaletteLookupEntry) Data() []byte {
	return entry.List
}

// ControlDictionaryEntry contains control words for compression.
type ControlDictionaryEntry struct {
	EntryBase

	Words []compression.ControlWord
}

// PaletteDictionaryEntryFrom decodes an entry from given data.
func PaletteDictionaryEntryFrom(timestamp Timestamp, data []byte) (ControlDictionaryEntry, error) {
	var entry ControlDictionaryEntry
	entry.Time = timestamp
	var err error
	entry.Words, err = compression.UnpackControlWords(data)
	return entry, err
}

// Type describes the entry.
func (entry ControlDictionaryEntry) Type() DataType {
	return DataTypeControlDictionary
}

// Code serializes the entry.
func (entry ControlDictionaryEntry) Data() []byte {
	return compression.PackControlWords(entry.Words)
}
