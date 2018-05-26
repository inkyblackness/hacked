package serial

import "io"

// SeekingReadCloser combines the interfaces of a Reader, a Seeker and a Closer.
type SeekingReadCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
