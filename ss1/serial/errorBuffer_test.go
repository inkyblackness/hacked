package serial_test

import (
	"fmt"
	"io"

	"github.com/inkyblackness/hacked/ss1"
)

type errorBuffer struct {
	position        int64
	callCounter     int
	errorOnNextCall bool
	errorByteCount  int
}

func (buf *errorBuffer) check() (err error) {
	buf.callCounter++
	if buf.errorOnNextCall {
		buf.errorOnNextCall = false
		err = ss1.StringError(fmt.Sprintf("errorBuffer on call number %v", buf.callCounter))
	}
	return
}

func (buf *errorBuffer) Seek(offset int64, whence int) (int64, error) {
	err := buf.check()
	if err != nil {
		return 0, err
	}
	switch whence {
	case io.SeekStart:
		buf.position = offset
	case io.SeekCurrent:
		buf.position += offset
	case io.SeekEnd:
		panic("SeekEnd not supported")
	}
	return buf.position, nil
}

func (buf *errorBuffer) Read(p []byte) (n int, err error) {
	err = buf.check()
	if err != nil {
		return buf.errorByteCount, err
	}
	return len(p), nil
}

func (buf *errorBuffer) Write(p []byte) (n int, err error) {
	err = buf.check()
	if err != nil {
		return buf.errorByteCount, err
	}
	return len(p), nil
}
