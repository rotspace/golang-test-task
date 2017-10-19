package tagc

import "io"

type sizeCounter struct {
	io.Reader
	readed int
}

func newSizeCounter(rdr io.Reader) *sizeCounter {
	return &sizeCounter{Reader: rdr}
}

func (tc *sizeCounter) Read(p []byte) (l int, err error) {
	l, err = tc.Reader.Read(p)
	tc.readed += l
	return
}
