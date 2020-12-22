package g53

import (
	"bytes"
	"testing"

	"github.com/ben-han-cn/g53/util"
)

func TestVerify(t *testing.T) {
	reqRaw := []byte{0x74, 0xdc, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x3, 0x63, 0x6f, 0x6d, 0x0, 0x0, 0x6, 0x0, 0x1, 0x0, 0x0, 0x29, 0x4, 0xd0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x9, 0x0, 0x0, 0x7, 0x61, 0x6c, 0x69, 0x62, 0x61, 0x62, 0x61, 0x0, 0x0, 0xfa, 0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3a, 0x8, 0x68, 0x6d, 0x61, 0x63, 0x2d, 0x6d, 0x64, 0x35, 0x7, 0x73, 0x69, 0x67, 0x2d, 0x61, 0x6c, 0x67, 0x3, 0x72, 0x65, 0x67, 0x3, 0x69, 0x6e, 0x74, 0x0, 0x0, 0x0, 0x5f, 0xd7, 0x70, 0x4c, 0x1, 0x2c, 0x0, 0x10, 0x24, 0x5c, 0xd2, 0x97, 0x1b, 0xc, 0xb9, 0xfe, 0x39, 0x64, 0x85, 0x9a, 0x53, 0x5, 0x9a, 0xb7, 0x74, 0xdc, 0x0, 0x0, 0x0, 0x0}

	req, _ := MessageFromWire(util.NewInputBuffer(reqRaw))
	key, _ := NewTsigKey("alibaba.",
		"z08GzEnlCDGy/W3Zw/2NHg==",
		"hmac-md5")
	Assert(t, key.VerifyMAC(req, nil) == nil, "")

	buf := util.NewOutputBuffer(1024)
	key.ToWire(buf)

	data := buf.Data()
	key, _ = TsigKeyFromWire(util.NewInputBuffer(data))
	//key should hold its own memory, so clean the underlaying data won't change its content
	for i := range data {
		data[i] = 0
	}
	Equal(t, key.Name, "alibaba.")
	rawSecret, _ := fromBase64([]byte("z08GzEnlCDGy/W3Zw/2NHg=="))
	Assert(t, bytes.Equal(key.rawSecret, rawSecret), "")
}
