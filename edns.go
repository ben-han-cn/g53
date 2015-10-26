package g53

import (
	"bytes"
	"errors"
	"fmt"

	"g53/util"
)

const (
	VERSION_SHIFT  = 16
	EXTRCODE_SHIFT = 24
	VERSION_MASK   = 0x00ff0000
	EXTFLAG_DO     = 0x00008000
)

type EDNS struct {
	Version       uint8
	extendedRcode uint8
	UdpSize       uint16
	DnssecAware   bool
}

func EdnsFromWire(buffer *util.InputBuffer) (*EDNS, error) {
	buffer.ReadUint8()

	t, err := TypeFromWire(buffer)
	if err != nil {
		return nil, err
	} else if t != RR_OPT {
		return nil, errors.New("edns rr type isn't opt")
	}

	udpSize, err := ClassFromWire(buffer)
	if err != nil {
		return nil, err
	}

	flags_, err := TTLFromWire(buffer)
	dnssecAware := (uint32(flags_) & EXTFLAG_DO) != 0
	extendedRcode := uint8(uint32(flags_) >> EXTRCODE_SHIFT)
	version := uint8((uint32(flags_) & VERSION_MASK) >> VERSION_SHIFT)
	return &EDNS{
		Version:       version,
		extendedRcode: extendedRcode,
		UdpSize:       uint16(udpSize),
		DnssecAware:   dnssecAware,
	}, nil
}

func (e *EDNS) Rend(r *MsgRender) {
	flags := uint32(e.extendedRcode) << EXTRCODE_SHIFT
	flags |= (uint32(e.Version) << VERSION_SHIFT) & VERSION_MASK
	if e.DnssecAware {
		flags |= EXTFLAG_DO
	}

	Root.Rend(r)
	RRType(RR_OPT).Rend(r)
	RRClass(e.UdpSize).Rend(r)
	RRTTL(flags).Rend(r)
	r.WriteUint16(0)
}

func (e *EDNS) ToWire(buffer *util.OutputBuffer) {
	flags := uint32(e.extendedRcode) << EXTRCODE_SHIFT
	flags |= (uint32(e.Version) << VERSION_SHIFT) & VERSION_MASK
	if e.DnssecAware {
		flags |= EXTFLAG_DO
	}

	Root.ToWire(buffer)
	RRType(RR_OPT).ToWire(buffer)
	RRClass(e.UdpSize).ToWire(buffer)
	RRTTL(flags).ToWire(buffer)
	buffer.WriteUint16(0)
}

func (e *EDNS) String() string {
	var desc bytes.Buffer
	desc.WriteString(fmt.Sprintf("; EDNS: version: %d, ", e.Version))
	if e.DnssecAware {
		desc.WriteString("flags: do; ")
	}
	desc.WriteString(fmt.Sprintf("udp: %d\n", e.UdpSize))
	return desc.String()
}
