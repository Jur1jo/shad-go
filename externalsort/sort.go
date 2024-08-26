//go:build !solution

package externalsort

import (
	"bytes"
	"io"
	"os"
	"sort"
)

func updateReader(i int, minElement []string, hasElement []bool, readers ...LineReader) {
	s, err := readers[i].ReadLine()
	minElement[i] = ""
	hasElement[i] = false
	if err == nil || len(s) != 0 {
		minElement[i] = s
		hasElement[i] = true
	}
}

func Merge(w LineWriter, readers ...LineReader) error {
	n := len(readers)
	minElement := make([]string, n)
	hasElement := make([]bool, n)
	for i := 0; i < n; i++ {
		updateReader(i, minElement, hasElement, readers...)
	}
	for true {
		minInd := -1
		for i := 0; i < n; i++ {
			if hasElement[i] {
				if minInd == -1 {
					minInd = i
				} else if minElement[minInd] > minElement[i] {
					minInd = i
				}
			}
		}
		if minInd == -1 {
			break
		}
		err := w.Write(minElement[minInd])
		if err != nil {
			return err
		}
		updateReader(minInd, minElement, hasElement, readers...)
	}
	return nil
}

func Sort(w io.Writer, in ...string) error {
	files := make([]*os.File, len(in))
	defer func() {
		for i := 0; i < len(in); i++ {
			files[i].Close()
		}
	}()
	for i, name := range in {
		file, err := os.Open(name)
		files[i] = file
		if err != nil {
			return err
		}
	}
	readers := make([]LineReader, len(in))
	for i := 0; i < len(in); i++ {
		arr := make([]string, 0)
		r := NewReader(files[i])
		for true {
			s, err := r.ReadLine()
			if s == "" && err != nil {
				break
			}
			arr = append(arr, s)
		}
		sort.Strings(arr)
		buf := bytes.NewBuffer(make([]byte, 0))
		for i := 0; i < len(arr); i++ {
			if i+1 != len(arr) || arr[i] == "" {
				buf.Write([]byte(arr[i] + "\n"))
			} else {
				buf.Write([]byte(arr[i]))
			}
		}
		readers[i] = NewReader(buf)
	}
	return Merge(NewWriter(w), readers...)
}
