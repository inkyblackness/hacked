package serial

import "io"

// SeekingWriteCloser combines the interfaces of a Writer, a Seeker and a Closer.
type SeekingWriteCloser interface {
	io.Writer
	io.Seeker
	io.Closer
}
