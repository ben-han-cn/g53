package g53

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/ben-han-cn/g53/util"
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
	Options       []Option
}

type Option interface {
	Rend(*MsgRender)
	String() string
}

func EdnsFromWire(buf *util.InputBuffer) (*EDNS, error) {
	e := EDNS{}
	if err := e.FromWire(buf); err != nil {
		return nil, err
	} else {
		return &e, nil
	}
}

func (e *EDNS) FromWire(buf *util.InputBuffer) error {
	if _, err := buf.ReadUint8(); err != nil {
		return err
	}

	if t, err := TypeFromWire(buf); err != nil {
		return err
	} else if t != RR_OPT {
		return errors.New("edns rr type isn't opt")
	}

	udpSize, err := ClassFromWire(buf)
	if err != nil {
		return err
	}

	flags_, err := TTLFromWire(buf)
	dnssecAware := (uint32(flags_) & EXTFLAG_DO) != 0
	extendedRcode := uint8(uint32(flags_) >> EXTRCODE_SHIFT)
	version := uint8((uint32(flags_) & VERSION_MASK) >> VERSION_SHIFT)

	rdlen, _ := buf.ReadUint16()
	opts := e.Options[:0]
	if rdlen != 0 {
		code, _ := buf.ReadUint16()
		switch code {
		case EDNS_SUBNET:
			if opt, err := subnetOptFromWire(buf); err == nil {
				opts = append(opts, opt)
			} else {
				return err
			}
		case EDNS_VIEW:
			if opt, err := viewOptFromWire(buf); err == nil {
				opts = append(opts, opt)
			} else {
				return err
			}
		case EDNS_EXPIRE:
			if opt, err := expireOptFromWire(buf); err == nil {
				opts = append(opts, opt)
			} else {
				return err
			}
		}
	}

	e.Version = version
	e.extendedRcode = extendedRcode
	e.UdpSize = uint16(udpSize)
	e.DnssecAware = dnssecAware
	e.Options = opts
	return nil
}

func EdnsFromRRset(rrset *RRset) *EDNS {
	var e EDNS
	e.FromRRset(rrset)
	return &e
}

func (e *EDNS) FromRRset(rrset *RRset) error {
	util.Assert(rrset.Type == RR_OPT, "edns should generate from otp")

	udpSize := uint16(rrset.Class)
	flags := uint32(rrset.Ttl)
	dnssecAware := (flags & EXTFLAG_DO) != 0
	extendedRcode := uint8(flags >> EXTRCODE_SHIFT)
	version := uint8((flags & VERSION_MASK) >> VERSION_SHIFT)

	opts := e.Options[:0]
	if len(rrset.Rdatas) > 0 {
		for _, rdata := range rrset.Rdatas {
			opt := rdata.(*OPT)
			if len(opt.Data) == 0 {
				continue
			}

			buf := util.NewInputBuffer(opt.Data)
			code, err := buf.ReadUint16()
			if err != nil {
				return err
			}

			if code == EDNS_SUBNET {
				if option, err := subnetOptFromWire(buf); err == nil {
					opts = append(opts, option)
				} else {
					return err
				}
			} else if code == EDNS_VIEW {
				if option, err := viewOptFromWire(buf); err == nil {
					opts = append(opts, option)
				} else {
					return err
				}
			} else if code == EDNS_EXPIRE {
				if option, err := expireOptFromWire(buf); err == nil {
					opts = append(opts, option)
				} else {
					return err
				}
			}
		}
	}

	e.Version = version
	e.extendedRcode = extendedRcode
	e.UdpSize = udpSize
	e.DnssecAware = dnssecAware
	e.Options = opts
	return nil
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
	if len(e.Options) == 0 {
		r.WriteUint16(0)
	} else {
		pos := r.Len()
		r.Skip(2)
		for _, opt := range e.Options {
			opt.Rend(r)
		}
		r.WriteUint16At(uint16(r.Len()-pos-2), pos)
	}
}

func (e *EDNS) ToWire(buf *util.OutputBuffer) {
	flags := uint32(e.extendedRcode) << EXTRCODE_SHIFT
	flags |= (uint32(e.Version) << VERSION_SHIFT) & VERSION_MASK
	if e.DnssecAware {
		flags |= EXTFLAG_DO
	}

	Root.ToWire(buf)
	RRType(RR_OPT).ToWire(buf)
	RRClass(e.UdpSize).ToWire(buf)
	RRTTL(flags).ToWire(buf)
	buf.WriteUint16(0)
}

func (e *EDNS) String() string {
	var header bytes.Buffer
	header.WriteString(fmt.Sprintf("; EDNS: version: %d, ", e.Version))
	if e.DnssecAware {
		header.WriteString("flags: do; ")
	}
	header.WriteString(fmt.Sprintf("udp: %d", e.UdpSize))
	desc := []string{header.String()}
	for _, opt := range e.Options {
		desc = append(desc, opt.String())
	}
	return strings.Join(desc, "\n") + "\n"
}

func (e *EDNS) CleanOption() {
	e.Options = []Option{}
}
