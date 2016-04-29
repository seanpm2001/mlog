package mlog

import (
	"io"
	"sync"
)

type sliceWriter interface {
	Write([]byte) (int, error)
	WriteByte(byte) error
	WriteString(string) (int, error)
	AppendIntWidth(int, int)
	Truncate(int)
}

type sliceBuffer struct {
	data []byte
}

func (sb *sliceBuffer) AppendIntWidth(i int, wid int) {
	digits := 0
	// write digits backwards (easier/faster)
	for i >= 10 {
		q := i / 10
		sb.data = append(sb.data, byte('0'+i-q*10))
		i = q
		digits++
	}
	sb.data = append(sb.data, byte('0'+i))
	digits++

	for j := wid - digits; j > 0; j-- {
		sb.data = append(sb.data, '0')
		digits++
	}

	// reverse to proper order
	sblen := len(sb.data)
	for i, j := sblen-digits, sblen-1; i < j; i, j = i+1, j-1 {
		sb.data[i], sb.data[j] = sb.data[j], sb.data[i]
	}
}

func (sb *sliceBuffer) Write(b []byte) (int, error) {
	sb.data = append(sb.data, b...)
	return len(b), nil
}

func (sb *sliceBuffer) WriteByte(c byte) error {
	sb.data = append(sb.data, c)
	return nil
}

func (sb *sliceBuffer) WriteString(s string) (int, error) {
	sb.data = append(sb.data, s...)
	return len(s), nil
}

func (sb *sliceBuffer) WriteTo(w io.Writer) (int, error) {
	return w.Write(sb.data)
}

func (sb *sliceBuffer) Bytes() []byte {
	return sb.data
}

func (sb *sliceBuffer) String() string {
	return string(sb.data)
}

func (sb *sliceBuffer) Len() int {
	return len(sb.data)
}

func (sb *sliceBuffer) Reset() {
	sb.Truncate(0)
}

func (sb *sliceBuffer) Truncate(i int) {
	sb.data = sb.data[:i]
}

type sliceBufferPool struct {
	*sync.Pool
}

func newSliceBufferPool() *sliceBufferPool {
	return &sliceBufferPool{
		&sync.Pool{New: func() interface{} {
			return &sliceBuffer{make([]byte, 0, 1024)}
		}},
	}
}

func (sp *sliceBufferPool) Get() *sliceBuffer {
	return (sp.Pool.Get()).(*sliceBuffer)
}

func (sp *sliceBufferPool) Put(c *sliceBuffer) {
	c.Truncate(0)
	sp.Pool.Put(c)
}
