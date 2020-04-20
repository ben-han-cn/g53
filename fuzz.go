package g53

import (
	"github.com/ben-han-cn/g53/util"
)

func Fuzz(data []byte) int {
	buf := util.NewInputBuffer(data)
	_, err := MessageFromWire(buf)
	if err != nil {
		return 0
	} else {
		return 1
	}
}
