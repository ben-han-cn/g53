package util

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestBuffer(t *testing.T) {
	out := NewOutputBuffer(1024)
	sls := []int{0, 32, 220, 20003, 60006, 300009, 3100000}
	bss := make([][]byte, 0, len(sls))
	for _, l := range sls {
		bs := make([]byte, l)
		rand.Read(bs)
		out.WriteVariableLenBytes(bs)
		bss = append(bss, bs)
	}

	in := NewInputBuffer(out.Data())
	for i, bs := range bss {
		bsRead, _ := in.ReadVariableLenBytes()
		if bytes.Compare(bs, bsRead) != 0 {
			t.Errorf("for len %d failed\n", sls[i])
			t.FailNow()
		}
	}

	out = NewOutputBuffer(1024)
	ss := []string{"", "good", "ddddd, dddd"}
	for _, s := range ss {
		out.WriteVariableLenBytes(StringToBytes(s))
	}

	in = NewInputBuffer(out.Data())
	for _, s := range ss {
		strRead, _ := in.ReadVariableLenBytes()
		if BytesToString(strRead) != s || string(strRead) != s {
			t.Errorf("for string %s failed\n", s)
			t.FailNow()
		}
	}
}
