package movie

import (
	"bytes"
	"io"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// Entry describes a block from a MOVI container.
type Entry interface {
	// Timestamp marks the beginning time of the entry.
	Timestamp() Timestamp
	// SetTimestamp changes the point in time of the entry.
	SetTimestamp(value Timestamp)
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

// SetTimestamp changes the timestamp of the entry.
func (entry *EntryBase) SetTimestamp(value Timestamp) {
	entry.Time = value
}

// UnknownEntry is an entry that is not know to this codebase.
type UnknownEntry struct {
	EntryBase
	DataType DataType
	Bytes    []byte
}

// UnknownEntryFrom decodes an entry from given data.
func UnknownEntryFrom(timestamp Timestamp, dataType DataType, r io.Reader, dataSize int) (*UnknownEntry, error) {
	entry := UnknownEntry{
		EntryBase: EntryBase{Time: timestamp},
		DataType:  dataType,
		Bytes:     make([]byte, dataSize),
	}
	coder := serial.NewDecoder(r)
	coder.Code(entry.Bytes)
	return &entry, coder.FirstError()
}

// Type describes the entry.
func (entry UnknownEntry) Type() DataType {
	return entry.DataType
}

// Data returns the raw bytes of the entry.
func (entry UnknownEntry) Data() []byte {
	return entry.Bytes
}

// LowResVideoEntry is a compressed low-resolution video frame.
type LowResVideoEntry struct {
	EntryBase

	BoundingBox [4]uint16
	Packed      []byte
}

// LowResVideoEntryFrom decodes an entry from given data.
func LowResVideoEntryFrom(timestamp Timestamp, r io.Reader, dataSize int) (*LowResVideoEntry, error) {
	var entry LowResVideoEntry
	entry.Time = timestamp
	coder := serial.NewDecoder(r)
	coder.Code(&entry.BoundingBox)
	entry.Packed = make([]byte, dataSize-int(coder.CurPos()))
	coder.Code(entry.Packed)
	return &entry, coder.FirstError()
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
func HighResVideoEntryFrom(timestamp Timestamp, r io.Reader, dataSize int) (*HighResVideoEntry, error) {
	var entry HighResVideoEntry
	entry.Time = timestamp
	coder := serial.NewDecoder(r)
	var maskstreamOffset uint16
	coder.Code(&maskstreamOffset)
	entry.Bitstream = make([]byte, int(maskstreamOffset)-int(coder.CurPos()))
	coder.Code(entry.Bitstream)
	entry.Maskstream = make([]byte, dataSize-int(coder.CurPos()))
	coder.Code(entry.Maskstream)
	return &entry, coder.FirstError()
}

// Type describes the entry.
func (entry HighResVideoEntry) Type() DataType {
	return DataTypeHighResVideo
}

// Data returns the raw bytes of the entry.
func (entry HighResVideoEntry) Data() []byte {
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	maskstreamOffset := uint16(len(entry.Bitstream)) + 2
	coder.Code(&maskstreamOffset)
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
func AudioEntryFrom(timestamp Timestamp, r io.Reader, dataSize int) (*AudioEntry, error) {
	var entry AudioEntry
	entry.Time = timestamp
	entry.Samples = make([]byte, dataSize)
	coder := serial.NewDecoder(r)
	coder.Code(entry.Samples)
	return &entry, coder.FirstError()
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
func SubtitleEntryFrom(timestamp Timestamp, r io.Reader, dataSize int) (*SubtitleEntry, error) {
	var entry SubtitleEntry
	entry.Time = timestamp
	coder := serial.NewDecoder(r)
	var header SubtitleHeader
	coder.Code(&header)
	entry.Control = header.Control
	if header.TextOffset > coder.CurPos() {
		coder.Code(make([]byte, header.TextOffset-coder.CurPos()))
	}
	entry.Text = make([]byte, dataSize-int(coder.CurPos()))
	coder.Code(entry.Text)
	return &entry, coder.FirstError()
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
func PaletteEntryFrom(timestamp Timestamp, r io.Reader) (*PaletteEntry, error) {
	var entry PaletteEntry
	entry.Time = timestamp
	coder := serial.NewDecoder(r)
	coder.Code(&entry.Colors)
	return &entry, coder.FirstError()
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
func PaletteResetEntryFrom(timestamp Timestamp) (*PaletteResetEntry, error) {
	var entry PaletteResetEntry
	entry.Time = timestamp
	return &entry, nil
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
func PaletteLookupEntryFrom(timestamp Timestamp, r io.Reader, dataSize int) (*PaletteLookupEntry, error) {
	var entry PaletteLookupEntry
	entry.Time = timestamp
	entry.List = make([]byte, dataSize)
	coder := serial.NewDecoder(r)
	coder.Code(entry.List)
	return &entry, coder.FirstError()
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

// ControlDictionaryEntryFrom decodes an entry from given data.
func ControlDictionaryEntryFrom(timestamp Timestamp, r io.Reader, dataSize int) (*ControlDictionaryEntry, error) {
	var entry ControlDictionaryEntry
	entry.Time = timestamp
	coder := serial.NewDecoder(r)
	data := make([]byte, dataSize)
	coder.Code(data)
	if coder.FirstError() != nil {
		return nil, coder.FirstError()
	}
	var err error
	entry.Words, err = compression.UnpackControlWords(data)
	return &entry, err
}

// Type describes the entry.
func (entry ControlDictionaryEntry) Type() DataType {
	return DataTypeControlDictionary
}

// Code serializes the entry.
func (entry ControlDictionaryEntry) Data() []byte {
	return compression.PackControlWords(entry.Words)
}
