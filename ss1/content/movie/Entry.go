package movie

import (
	"bytes"
	"io"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// Entry describes a timestamped block from a MOVI container.
type Entry struct {
	// Timestamp marks the beginning time of the entry.
	Timestamp Timestamp
	// Data is a type specific information of a movie.
	Data EntryData
}

// EntryData describes one entry in a MOVI container.
type EntryData interface {
	// Type describes the content type of the data.
	Type() DataType
	// Bytes returns the raw bytes of the entry.
	Bytes() []byte
}

// UnknownEntry is an entry that is not know to this codebase.
type UnknownEntry struct {
	DataType DataType
	Raw      []byte
}

// UnknownEntryFrom decodes an entry from given data.
func UnknownEntryFrom(dataType DataType, r io.Reader, dataSize int) (UnknownEntry, error) {
	entry := UnknownEntry{
		DataType: dataType,
		Raw:      make([]byte, dataSize),
	}
	coder := serial.NewDecoder(r)
	coder.Code(entry.Raw)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry UnknownEntry) Type() DataType {
	return entry.DataType
}

// Bytes returns the raw bytes of the entry.
func (entry UnknownEntry) Bytes() []byte {
	return entry.Raw
}

// LowResVideoEntry is a compressed low-resolution video frame.
type LowResVideoEntry struct {
	BoundingBox [4]uint16
	Packed      []byte
}

// LowResVideoEntryFrom decodes an entry from given data.
func LowResVideoEntryFrom(r io.Reader, dataSize int) (LowResVideoEntry, error) {
	var entry LowResVideoEntry
	coder := serial.NewDecoder(r)
	coder.Code(&entry.BoundingBox)
	entry.Packed = make([]byte, dataSize-int(coder.CurPos()))
	coder.Code(entry.Packed)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry LowResVideoEntry) Type() DataType {
	return DataTypeLowResVideo
}

// Bytes returns the raw bytes of the entry.
func (entry LowResVideoEntry) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	coder.Code(entry.BoundingBox)
	coder.Code(entry.Packed)
	return buf.Bytes()
}

// HighResVideoEntry is a compressed high-resolution video frame.
type HighResVideoEntry struct {
	Bitstream  []byte
	Maskstream []byte
}

// HighResVideoEntryFrom decodes an entry from given data.
func HighResVideoEntryFrom(r io.Reader, dataSize int) (HighResVideoEntry, error) {
	var entry HighResVideoEntry
	coder := serial.NewDecoder(r)
	var maskstreamOffset uint16
	coder.Code(&maskstreamOffset)
	entry.Bitstream = make([]byte, int(maskstreamOffset)-int(coder.CurPos()))
	coder.Code(entry.Bitstream)
	entry.Maskstream = make([]byte, dataSize-int(coder.CurPos()))
	coder.Code(entry.Maskstream)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry HighResVideoEntry) Type() DataType {
	return DataTypeHighResVideo
}

// Bytes returns the raw bytes of the entry.
func (entry HighResVideoEntry) Bytes() []byte {
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
	Samples []byte
}

// AudioEntryFrom decodes an entry from given data.
func AudioEntryFrom(r io.Reader, dataSize int) (AudioEntry, error) {
	var entry AudioEntry
	entry.Samples = make([]byte, dataSize)
	coder := serial.NewDecoder(r)
	coder.Code(entry.Samples)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry AudioEntry) Type() DataType {
	return DataTypeAudio
}

// Bytes returns the raw bytes of the entry.
func (entry AudioEntry) Bytes() []byte {
	return entry.Samples
}

// SubtitleEntry is an entry with Subtitle samples.
type SubtitleEntry struct {
	Control SubtitleControl
	Text    []byte
}

// SubtitleEntryFrom decodes an entry from given data.
func SubtitleEntryFrom(r io.Reader, dataSize int) (SubtitleEntry, error) {
	var entry SubtitleEntry
	coder := serial.NewDecoder(r)
	var header SubtitleHeader
	coder.Code(&header)
	entry.Control = header.Control
	if header.TextOffset > coder.CurPos() {
		coder.Code(make([]byte, header.TextOffset-coder.CurPos()))
	}
	entry.Text = make([]byte, dataSize-int(coder.CurPos()))
	coder.Code(entry.Text)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry SubtitleEntry) Type() DataType {
	return DataTypeSubtitle
}

// Bytes returns the raw bytes of the entry.
func (entry SubtitleEntry) Bytes() []byte {
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
	Colors bitmap.Palette
}

// PaletteEntryFrom decodes an entry from given data.
func PaletteEntryFrom(r io.Reader) (PaletteEntry, error) {
	var entry PaletteEntry
	coder := serial.NewDecoder(r)
	coder.Code(&entry.Colors)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry PaletteEntry) Type() DataType {
	return DataTypePalette
}

// Bytes returns the raw bytes of the entry.
func (entry PaletteEntry) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	coder.Code(&entry.Colors)
	return buf.Bytes()
}

// PaletteResetEntry marks a reset of the palette.
type PaletteResetEntry struct {
}

// PaletteResetEntryFrom decodes an entry from given data.
func PaletteResetEntryFrom() (PaletteResetEntry, error) {
	var entry PaletteResetEntry
	return entry, nil
}

// Type describes the entry.
func (entry PaletteResetEntry) Type() DataType {
	return DataTypePaletteReset
}

// Code serializes the entry.
func (entry PaletteResetEntry) Bytes() []byte {
	return nil
}

// PaletteLookupEntry contains a list of palette indices.
type PaletteLookupEntry struct {
	List []byte
}

// PaletteLookupEntryFrom decodes an entry from given data.
func PaletteLookupEntryFrom(r io.Reader, dataSize int) (PaletteLookupEntry, error) {
	var entry PaletteLookupEntry
	entry.List = make([]byte, dataSize)
	coder := serial.NewDecoder(r)
	coder.Code(entry.List)
	return entry, coder.FirstError()
}

// Type describes the entry.
func (entry PaletteLookupEntry) Type() DataType {
	return DataTypePaletteLookupList
}

// Code serializes the entry.
func (entry PaletteLookupEntry) Bytes() []byte {
	return entry.List
}

// ControlDictionaryEntry contains control words for compression.
type ControlDictionaryEntry struct {
	Words []compression.ControlWord
}

// ControlDictionaryEntryFrom decodes an entry from given data.
func ControlDictionaryEntryFrom(r io.Reader, dataSize int) (ControlDictionaryEntry, error) {
	var entry ControlDictionaryEntry
	coder := serial.NewDecoder(r)
	data := make([]byte, dataSize)
	coder.Code(data)
	if coder.FirstError() != nil {
		return entry, coder.FirstError()
	}
	var err error
	entry.Words, err = compression.UnpackControlWords(data)
	return entry, err
}

// Type describes the entry.
func (entry ControlDictionaryEntry) Type() DataType {
	return DataTypeControlDictionary
}

// Code serializes the entry.
func (entry ControlDictionaryEntry) Bytes() []byte {
	return compression.PackControlWords(entry.Words)
}
