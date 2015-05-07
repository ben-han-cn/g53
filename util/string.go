package util

import (
	"strconv"
)

func HexStrToBytes(s string) (result []uint8, err error) {
	for i := 0; i < len(s); i += 2 {
		d, err := strconv.ParseUint(s[i:i+2], 16, 8)
		if err != nil {
			break
		}
		result = append(result, uint8(d))
	}
	return
}
