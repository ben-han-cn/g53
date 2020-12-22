package util

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestBuffer(t *testing.T) {
	for _, l := range []int{0, 32, 129, 20003, 60006, 300009, 3100000} {
		bs := make([]byte, l)
		rand.Read(bs)
		out := NewOutputBuffer(1024)
		out.WriteVariableLenBytes(bs)

		in := NewInputBuffer(out.Data())
		bsRead, _ := in.ReadVariableLenBytes()
		if bytes.Compare(bs, bsRead) != 0 {
			t.Errorf("for len %d failed\n", l)
			t.FailNow()
		}
	}
}
