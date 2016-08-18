package g53

import (
	"fmt"
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
	Assert(t, nrs.Name.Equals(rs.Name), fmt.Sprintf("%s != %s", nrs.Name.String(false), rs.Name.String(false)))
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

func TestRRsetRoateRdata(t *testing.T) {
	ra1, _ := AFromString("1.1.1.1")
	ra2, _ := AFromString("2.2.2.2")
	ra3, _ := AFromString("3.3.3.3")
	n, _ := NameFromString("test.example.com.")
	rrset := &RRset{
		Name:   n,
		Type:   RR_A,
		Class:  CLASS_IN,
		Ttl:    RRTTL(3600),
		Rdatas: []Rdata{ra1},
	}
	rrset.RotateRdata()
	Equal(t, rrset.Rdatas[0].String(), ra1.String())

	rrset.AddRdata(ra2)
	rrset.RotateRdata()
	Equal(t, rrset.Rdatas[0].String(), ra2.String())
	Equal(t, rrset.Rdatas[1].String(), ra1.String())
	rrset.RotateRdata()
	Equal(t, rrset.Rdatas[0].String(), ra1.String())
	Equal(t, rrset.Rdatas[1].String(), ra2.String())

	rrset.AddRdata(ra3)
	rrset.RotateRdata()
	Equal(t, rrset.Rdatas[0].String(), ra3.String())
	Equal(t, rrset.Rdatas[1].String(), ra1.String())
	Equal(t, rrset.Rdatas[2].String(), ra2.String())
	rrset.RotateRdata()
	Equal(t, rrset.Rdatas[0].String(), ra2.String())
	Equal(t, rrset.Rdatas[1].String(), ra3.String())
	Equal(t, rrset.Rdatas[2].String(), ra1.String())
	rrset.RotateRdata()
	Equal(t, rrset.Rdatas[0].String(), ra1.String())
	Equal(t, rrset.Rdatas[1].String(), ra2.String())
	Equal(t, rrset.Rdatas[2].String(), ra3.String())
}
