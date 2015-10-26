package g53

import (
	"g53/util"
	"testing"
)

func matchRRsetRaw(t *testing.T, rawData string, rs *RRset) {
	wire, _ := util.HexStrToBytes(rawData)
	buffer := util.NewInputBuffer(wire)
	nrs, err := RRsetFromWire(buffer)
	Assert(t, err == nil, "err should be nil")
	matchRRset(t, nrs, rs)
	render := NewMsgRender()
	nrs.Rend(render)
	WireMatch(t, wire, render.Data())
}

func matchRRset(t *testing.T, nrs *RRset, rs *RRset) {
	Assert(t, nrs.Name.Equals(rs.Name), "name should equal")
	Equal(t, nrs.Type, rs.Type)
	Equal(t, nrs.Class, rs.Class)
	Equal(t, len(nrs.Rdatas), len(rs.Rdatas))
	for i := 0; i < len(rs.Rdatas); i++ {
		Equal(t, nrs.Rdatas[i].String(), rs.Rdatas[i].String())
	}
}

func TestRRsetFromToWire(t *testing.T) {
	n, _ := NameFromString("test.example.com.")
	ra, _ := AFromString("192.0.2.1")
	matchRRsetRaw(t, "0474657374076578616d706c6503636f6d000001000100000e100004c0000201", &RRset{
		Name:   n,
		Type:   RR_A,
		Class:  CLASS_IN,
		Ttl:    RRTTL(3600),
		Rdatas: []Rdata{ra},
	})
}
