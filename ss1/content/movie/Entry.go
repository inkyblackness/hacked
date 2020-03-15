package movie

import (
	"bytes"
	"io"

	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// EntryBucketPriority describes the general priority of the bucket.
// Lower numbers are put before higher numbers if the bucket has the same timestamp.
type EntryBucketPriority int

// List of entry bucket priorities.
const (
	EntryBucketPriorityAudio           = 0
	EntryBucketPriorityVideoControl    = 1
	EntryBucketPrioritySubtitleControl = 2
	EntryBucketPrioritySubtitle        = 3
	EntryBucketPriorityFrame           = 4
)

// EntryBucket is a set of entries that need to be together within the stream.
// This exists mainly to cover special cases of entry within a stream that have no valid timestamp.
type EntryBucket struct {
	// Priority helps in merging content from multiple buckets.
	Priority EntryBucketPriority
	// Timestamp specifies at which time the bucket should be inserted into the stream.
	// This is typically derived from one of the timestamps of the contained entries.
	Timestamp Timestamp
	// Entries is the list of entries within this bucket.
	Entries []Entry
}

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

// UnknownEntryData is an entry that is not know to this codebase.
type UnknownEntryData struct {
	DataType DataType
	Raw      []byte
}

// UnknownEntryFrom decodes an entry from given data.
func UnknownEntryFrom(dataType DataType, r io.Reader, dataSize int) (UnknownEntryData, error) {
	data := UnknownEntryData{
		DataType: dataType,
		Raw:      make([]byte, dataSize),
	}
	coder := serial.NewDecoder(r)
	coder.Code(data.Raw)
	return data, coder.FirstError()
}

// Type describes the entry.
func (data UnknownEntryData) Type() DataType {
	return data.DataType
}

// Bytes returns the raw bytes of the entry.
func (data UnknownEntryData) Bytes() []byte {
	return data.Raw
}

// LowResVideoEntryData is a compressed low-resolution video frame.
type LowResVideoEntryData struct {
	BoundingBox [4]uint16
	Packed      []byte
}

// LowResVideoEntryFrom decodes an entry from given data.
func LowResVideoEntryFrom(r io.Reader, dataSize int) (LowResVideoEntryData, error) {
	var data LowResVideoEntryData
	coder := serial.NewDecoder(r)
	coder.Code(&data.BoundingBox)
	data.Packed = make([]byte, dataSize-int(coder.CurPos()))
	coder.Code(data.Packed)
	return data, coder.FirstError()
}

// Type describes the entry.
func (data LowResVideoEntryData) Type() DataType {
	return DataTypeLowResVideo
}

// Bytes returns the raw bytes of the entry.
func (data LowResVideoEntryData) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	coder.Code(data.BoundingBox)
	coder.Code(data.Packed)
	return buf.Bytes()
}

// HighResVideoEntryData is a compressed high-resolution video frame.
type HighResVideoEntryData struct {
	Bitstream  []byte
	Maskstream []byte
}

// HighResVideoEntryFrom decodes an entry from given data.
func HighResVideoEntryFrom(r io.Reader, dataSize int) (HighResVideoEntryData, error) {
	var data HighResVideoEntryData
	coder := serial.NewDecoder(r)
	var maskstreamOffset uint16
	coder.Code(&maskstreamOffset)
	data.Bitstream = make([]byte, int(maskstreamOffset)-int(coder.CurPos()))
	coder.Code(data.Bitstream)
	data.Maskstream = make([]byte, dataSize-int(coder.CurPos()))
	coder.Code(data.Maskstream)
	return data, coder.FirstError()
}

// Type describes the entry.
func (data HighResVideoEntryData) Type() DataType {
	return DataTypeHighResVideo
}

// Bytes returns the raw bytes of the entry.
func (data HighResVideoEntryData) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	maskstreamOffset := uint16(len(data.Bitstream)) + 2
	coder.Code(&maskstreamOffset)
	coder.Code(data.Bitstream)
	coder.Code(data.Maskstream)
	return buf.Bytes()
}

// AudioEntryData is an entry with audio samples.
type AudioEntryData struct {
	Samples []byte
}

// AudioEntryFrom decodes an entry from given data.
func AudioEntryFrom(r io.Reader, dataSize int) (AudioEntryData, error) {
	var data AudioEntryData
	data.Samples = make([]byte, dataSize)
	coder := serial.NewDecoder(r)
	coder.Code(data.Samples)
	return data, coder.FirstError()
}

// Type describes the entry.
func (data AudioEntryData) Type() DataType {
	return DataTypeAudio
}

// Bytes returns the raw bytes of the entry.
func (data AudioEntryData) Bytes() []byte {
	return data.Samples
}

// SubtitleEntryData is an entry with Subtitle samples.
type SubtitleEntryData struct {
	Control SubtitleControl
	Text    []byte
}

// SubtitleEntryFrom decodes an entry from given data.
func SubtitleEntryFrom(r io.Reader, dataSize int) (SubtitleEntryData, error) {
	var data SubtitleEntryData
	coder := serial.NewDecoder(r)
	var header SubtitleHeader
	coder.Code(&header)
	data.Control = header.Control
	if header.TextOffset > coder.CurPos() {
		coder.Code(make([]byte, header.TextOffset-coder.CurPos()))
	}
	data.Text = make([]byte, dataSize-int(coder.CurPos()))
	coder.Code(data.Text)
	return data, coder.FirstError()
}

// Type describes the entry.
func (data SubtitleEntryData) Type() DataType {
	return DataTypeSubtitle
}

// Bytes returns the raw bytes of the entry.
func (data SubtitleEntryData) Bytes() []byte {
	header := SubtitleHeader{
		Control:    data.Control,
		TextOffset: SubtitleDefaultTextOffset,
	}
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	coder.Code(&header)
	coder.Code(make([]byte, int(header.TextOffset-coder.CurPos())))
	coder.Code(data.Text)
	return buf.Bytes()
}

// PaletteEntryData is an entry with a palette.
type PaletteEntryData struct {
	Colors bitmap.Palette
}

// PaletteEntryFrom decodes an entry from given data.
func PaletteEntryFrom(r io.Reader) (PaletteEntryData, error) {
	var data PaletteEntryData
	coder := serial.NewDecoder(r)
	coder.Code(&data.Colors)
	return data, coder.FirstError()
}

// Type describes the entry.
func (data PaletteEntryData) Type() DataType {
	return DataTypePalette
}

// Bytes returns the raw bytes of the entry.
func (data PaletteEntryData) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	coder := serial.NewEncoder(buf)
	coder.Code(&data.Colors)
	return buf.Bytes()
}

// PaletteResetEntryData marks a reset of the palette.
type PaletteResetEntryData struct {
}

// PaletteResetEntryFrom decodes an entry from given data.
func PaletteResetEntryFrom() (PaletteResetEntryData, error) {
	return PaletteResetEntryData{}, nil
}

// Type describes the entry.
func (data PaletteResetEntryData) Type() DataType {
	return DataTypePaletteReset
}

// Code serializes the entry.
func (data PaletteResetEntryData) Bytes() []byte {
	return nil
}

// PaletteLookupEntryData contains a list of palette indices.
type PaletteLookupEntryData struct {
	List []byte
}

// PaletteLookupEntryFrom decodes an entry from given data.
func PaletteLookupEntryFrom(r io.Reader, dataSize int) (PaletteLookupEntryData, error) {
	var data PaletteLookupEntryData
	data.List = make([]byte, dataSize)
	coder := serial.NewDecoder(r)
	coder.Code(data.List)
	return data, coder.FirstError()
}

// Type describes the entry.
func (data PaletteLookupEntryData) Type() DataType {
	return DataTypePaletteLookupList
}

// Code serializes the entry.
func (data PaletteLookupEntryData) Bytes() []byte {
	return data.List
}

// ControlDictionaryEntryData contains control words for compression.
type ControlDictionaryEntryData struct {
	Words []compression.ControlWord
}

// ControlDictionaryEntryFrom decodes an entry from given data.
func ControlDictionaryEntryFrom(r io.Reader, dataSize int) (ControlDictionaryEntryData, error) {
	var data ControlDictionaryEntryData
	coder := serial.NewDecoder(r)
	b := make([]byte, dataSize)
	coder.Code(b)
	if coder.FirstError() != nil {
		return data, coder.FirstError()
	}
	var err error
	data.Words, err = compression.UnpackControlWords(b)
	return data, err
}

// Type describes the entry.
func (data ControlDictionaryEntryData) Type() DataType {
	return DataTypeControlDictionary
}

// Code serializes the entry.
func (data ControlDictionaryEntryData) Bytes() []byte {
	return compression.PackControlWords(data.Words)
}
