//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
	"strings"
)

type MyReader struct {
	reader *bufio.Reader
}

type MyWriter struct {
	writer *bufio.Writer
}

func (r *MyReader) ReadLine() (string, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimRight(line, "\n"), err
}

func (w *MyWriter) Write(l string) error {
	_, err := w.writer.WriteString(l + "\n")
	if err != nil {
		return err
	}
	return w.writer.Flush()
}

type HeapElem struct {
	str string
	ind int
}

type MyHeap []HeapElem

func (h MyHeap) Len() int {
	return len(h)
}

func (h MyHeap) Less(i, j int) bool {
	return h[i].str < h[j].str
}

func (h MyHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MyHeap) Push(x interface{}) {
	*h = append(*h, x.(HeapElem))
}

func (h *MyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func NewReader(r io.Reader) LineReader {
	return &MyReader{reader: bufio.NewReader(r)}
}

func NewWriter(w io.Writer) LineWriter {
	return &MyWriter{writer: bufio.NewWriter(w)}
}

func Merge(w LineWriter, readers ...LineReader) error {
	h := &MyHeap{}
	heap.Init(h)
	for i, reader := range readers {
		line, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return err
		}
		if !(err == io.EOF && line == "") {
			heap.Push(h, HeapElem{str: line, ind: i})
		}
	}
	for h.Len() > 0 {
		elem := heap.Pop(h).(HeapElem)
		err := w.Write(elem.str)
		if err != nil {
			return err
		}
		line, err := readers[elem.ind].ReadLine()
		if err != nil && err != io.EOF {
			return err
		} else if !(err == io.EOF && line == "") {
			heap.Push(h, HeapElem{str: line, ind: elem.ind})
		}
	}
	return nil
}

func Sort(w io.Writer, in ...string) error {
	var readers []LineReader
	for _, file := range in {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		myReader := NewReader(f)
		var str []string
		for {
			line, errReader := myReader.ReadLine()

			if errReader != nil && errReader != io.EOF {
				return err
			}
			if !(errReader == io.EOF && line == "") {
				str = append(str, line)
			}
			if errReader == io.EOF {
				break
			}
		}
		err = f.Close()
		if err != nil {
			return err
		}
		sort.Strings(str)
		sortFile, err := os.Create(file)
		if err != nil {
			return err
		}
		writer := bufio.NewWriter(sortFile)
		for _, line := range str {
			_, err = writer.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
		if err = writer.Flush(); err != nil {
			return err
		}
		sortFile.Close()
	}
	for _, file := range in {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		myReader := NewReader(f)
		readers = append(readers, myReader)
	}
	return Merge(NewWriter(w), readers...)
}
