package g53

import (
	"math"
	"path/filepath"
	"reflect"
	"runtime"
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
		t.Errorf("expect name %v but get %v\n", str, n.String(true))
	}
}

func Assert(t *testing.T, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d: "+msg+"\033[39m\n\n",
			append([]interface{}{filepath.Base(file), line}, v...)...)
		t.FailNow()
	}
}

func Equal(t *testing.T, act, exp interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n",
			filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

func Nequal(t *testing.T, act, exp interface{}) {
	if reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n",
			filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}
