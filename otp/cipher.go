//go:build !solution

package otp

import (
	"io"
)

type otpReader struct {
	r    io.Reader
	prng io.Reader
}

type otpWriter struct {
	w    io.Writer
	prng io.Reader
}

func (w *otpReader) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	encRes := make([]byte, 1, 1)
	rndRes := make([]byte, 1, 1)
	encN, err := w.r.Read(encRes)
	if encN == 0 {
		return 0, err
	}
	w.prng.Read(rndRes)
	b[0] = encRes[0] ^ rndRes[0]
	return 1, err
}

func (w *otpWriter) Write(b []byte) (int, error) {
	resBytes := make([]byte, len(b), len(b))
	for i := 0; ; {
		n, err := w.prng.Read(resBytes[i:])
		i += n
		if err != nil || i == len(b) {
			break
		}
	}
	for i := range b {
		resBytes[i] = resBytes[i] ^ b[i]
	}
	for i := 0; ; {
		n, err := w.w.Write(resBytes)
		i += n
		if err != nil || i == len(b) {
			return i, err
		}
	}
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &otpReader{r: r, prng: prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &otpWriter{w: w, prng: prng}
}
