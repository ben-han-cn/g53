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

	for _, s := range []string{"", "good", "ddddd, dddd"} {
		out := NewOutputBuffer(1024)
		out.WriteVariableLenBytes(StringToBytes(s))

		in := NewInputBuffer(out.Data())
		strRead, _ := in.ReadVariableLenBytes()
		if BytesToString(strRead) != s || string(strRead) != s {
			t.Errorf("for string %s failed\n", s)
			t.FailNow()
		}
	}
}
