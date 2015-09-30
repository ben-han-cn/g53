package g53

import (
	"errors"
)

const (
	DefaultMaxUdpSize uint16 = 512
	VersionShift      uint32 = 16
	VersionMask       uint32 = 0x00ff0000
	SupportedVersion  uint8  = 0
	ExtFlagDo         uint32 = 0x00008000
	ExtRcodeShift     uint32 = 24
)

type Edns struct {
	version       uint8
	udpSize       uint16
	dnssecAware   bool
	extendedRcode uint8
}

func FromRRset(rrset *RRset) (*Edns, error) {
	if rrset.rrtype != OPT {
		return nil, errors.New("edns type isn't opt")
	}

	version := uint8((uint32(rrset.rrttl) & VersionMask) >> VersionShift)
	if version > SupportedVersion {
		return nil, errors.New("edns version isn't supported")
	}

	if rrset.name.Equals(Root) == false {
		return nil, errors.New("edns name should be root")
	}

	dnssecAware := (uint32(rrset.rrttl) & ExtFlagDo) != 0
	udpSize := uint16(rrset.rrclass)
	extendedRcode := uint8(uint32(rrset.rrttl) >> ExtRcodeShift)
	return &Edns{
		version:       version,
		udpSize:       udpSize,
		dnssecAware:   dnssecAware,
		extendedRcode: extendedRcode,
	}, nil
}

func (edns *Edns) Rend(render *MsgRender) {
	extRcodeFlags = ends.extendedRcode << ExtRcodeShift
	extRcodeFlags |= (edns.version << VersionShift) & VersionMask
	if edns.dnssecAware {
		extRcodeFlags |= ExtFlagDo
	}

	Root.Rend(render)
	OPT.Rend(render)
	RRClass(edns.udpSize).Render(render)
	RRTTL(extRcodeFlags).Render(render)
	render.WriteUint16(0)
}

func (edns *Edns) ToWire(buffer *util.OutputBuffer) {
	extRcodeFlags = ends.extendedRcode << ExtRcodeShift
	extRcodeFlags |= (edns.version << VersionShift) & VersionMask
	if edns.dnssecAware {
		extRcodeFlags |= ExtFlagDo
	}

	Root.ToWire(buffer)
	OPT.ToWire(buffer)
	RRClass(edns.udpSize).ToWire(buffer)
	RRTTL(extRcodeFlags).ToWire(buffer)
	buffer.WriteUint16(0)
}

func (edns *Edns) String() string {
	ret := "; EDNS: version: "
	ret += strconv.Itoa(int(edns.version))
	ret += ", flags:"
	if edns.dnssecAware {
		ret += " do"
	}
	ret += "; udp: " + strconv.Itoa(int(edns.udpSize)) + "\n"
	return ret
}
