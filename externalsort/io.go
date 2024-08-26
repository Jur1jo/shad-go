//go:build !change

package externalsort

import (
	"io"
)

const bufSize = 10

type LineReader interface {
	ReadLine() (string, error)
}

type LineWriter interface {
	Write(l string) error
}

type Reader struct {
	r        io.Reader
	cntEnter int
	buf      []byte
	indBuf   int
}

type Writer struct {
	w io.Writer
}

func (r *Reader) ReadLine() (string, error) {
	res := make([]byte, 0)
	for r.cntEnter == 0 {
		for _, b := range r.buf {
			res = append(res, b)
		}
		r.indBuf = 0
		r.buf = make([]byte, bufSize)
		n, err := r.r.Read(r.buf)
		r.buf = r.buf[:r.indBuf+n]
		if n == 0 && err != nil {
			return string(res), err
		}
		for i := 0; i < n; i++ {
			if rune(r.buf[i]) == '\n' {
				r.cntEnter++
			}
		}
	}

	for ; rune(r.buf[r.indBuf]) != '\n' && r.indBuf < len(r.buf); r.indBuf++ {
		res = append(res, r.buf[r.indBuf])
	}
	r.buf = r.buf[r.indBuf+1:]
	r.indBuf = 0
	r.cntEnter -= 1
	return string(res), nil
}

func (w *Writer) Write(l string) error {
	buf := make([]byte, len(l))
	for i, b := range l {
		buf[i] = byte(b)
	}
	for len(buf) > 0 {
		n, err := w.w.Write(buf)
		buf = buf[n:]
		if err != nil {
			return err
		}
	}
	_, err := w.w.Write([]byte("\n"))
	return err
}

func NewReader(r io.Reader) LineReader {
	return &Reader{r: r}
}

func NewWriter(w io.Writer) LineWriter {
	return &Writer{w: w}
}
