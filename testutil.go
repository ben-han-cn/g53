package g53

import (
	"math"
	"testing"
)

func WireMatch(t *testing.T, expectData []uint8, actualData []uint8) {
	if len(expectData) != len(actualData) {
		t.Errorf("want len [%v] but get [%v]", len(expectData), len(actualData))
	}

	minLen := int(math.Min(float64(len(expectData)), float64(len(actualData))))
	for i := 0; i < minLen; i++ {
		if expectData[i] != actualData[i] {
			t.Errorf("match part %v", expectData[:i])
			t.Errorf("right:%v", expectData[i:])
			t.Errorf("wrong:%v", actualData[i:])
			break
		}
	}
}

func NameEqToStr(t *testing.T, n *Name, str string) {
	s, _ := NewName(str, true)
	if n.Equals(s) == false {
		actualStr, _ := n.String(true)
		t.Errorf("expect name %v but get %v\n", str, actualStr)
	}
}
