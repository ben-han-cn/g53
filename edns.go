package g53

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	"g53/util"
)

const (
	VERSION_SHIFT  = 16
	EXTRCODE_SHIFT = 24
	VERSION_MASK   = 0x00ff0000
	EXTFLAG_DO     = 0x00008000

	EDNS_SUBNET = 8
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

type SubnetOpt struct {
	family uint16
	mask   uint8
	scope  uint8
	ip     net.IP
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

	rdlen, _ := buffer.ReadUint16()
	options := []Option{}
	if rdlen != 0 {
		code, _ := buffer.ReadUint16()
		if code == EDNS_SUBNET {
			if opt, err := subnetOptFromWire(buffer); err == nil {
				options = append(options, opt)
			} else {
				return nil, err
			}
		}
	}

	return &EDNS{
		Version:       version,
		extendedRcode: extendedRcode,
		UdpSize:       uint16(udpSize),
		DnssecAware:   dnssecAware,
		Options:       options,
	}, nil
}

func EdnsFromRRset(rrset *RRset) *EDNS {
	util.Assert(rrset.Type == RR_OPT, "edns should generate from otp")
	udpSize := uint16(rrset.Class)
	flags := uint32(rrset.Ttl)
	dnssecAware := (flags & EXTFLAG_DO) != 0
	extendedRcode := uint8(flags >> EXTRCODE_SHIFT)
	version := uint8((flags & VERSION_MASK) >> VERSION_SHIFT)

	options := []Option{}
	if len(rrset.Rdatas) > 0 {
		for _, rdata := range rrset.Rdatas {
			opt := subnetOptFromRdata(rdata)
			if opt != nil {
				options = append(options, opt)
			}
		}
	}

	return &EDNS{
		Version:       version,
		extendedRcode: extendedRcode,
		UdpSize:       udpSize,
		DnssecAware:   dnssecAware,
		Options:       options,
	}
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

func (subnet *SubnetOpt) Rend(render *MsgRender) {
	render.WriteUint16(EDNS_SUBNET)
	ipLen := uint(subnet.mask / 8)
	if subnet.mask%8 != 0 {
		ipLen += 1
	}

	render.WriteUint16(uint16(2 + 2 + ipLen))
	render.WriteUint16(subnet.family)
	render.WriteUint8(subnet.mask)
	render.WriteUint8(subnet.scope)
	var ipToWrite net.IP
	if subnet.family == 1 {
		ipToWrite = subnet.ip.To4().Mask(net.CIDRMask(int(subnet.mask), net.IPv4len*8))
	} else {
		ipToWrite = subnet.ip.To16().Mask(net.CIDRMask(int(subnet.mask), net.IPv6len*8))
	}
	render.WriteData([]byte(ipToWrite)[0:ipLen])
}

func (subnet *SubnetOpt) String() string {
	return fmt.Sprintf("; CLIENT-SUBNET: %s/%d\n", subnet.ip.String(), subnet.mask)
}

//read from OPTION-LENGTH
func subnetOptFromWire(buffer *util.InputBuffer) (Option, error) {
	l, _ := buffer.ReadUint16()
	family, _ := buffer.ReadUint16()
	mask, _ := buffer.ReadUint8()
	scope, _ := buffer.ReadUint8()
	var ip net.IP
	switch family {
	case 1:
		addr := make([]byte, 4)
		addr_data, _ := buffer.ReadBytes(uint(l - 4))
		copy(addr, addr_data)
		ip = net.IPv4(addr[0], addr[1], addr[2], addr[3])
	case 2:
		addr := make([]byte, 16)
		addr_data, _ := buffer.ReadBytes(uint(l - 4))
		copy(addr, addr_data)
		ip = net.IP{addr[0], addr[1], addr[2], addr[3], addr[4],
			addr[5], addr[6], addr[7], addr[8], addr[9], addr[10],
			addr[11], addr[12], addr[13], addr[14], addr[15]}
	}

	if ip != nil {
		return &SubnetOpt{family: family,
			mask:  mask,
			scope: scope,
			ip:    ip}, nil
	} else {
		return nil, fmt.Errorf("unkown family")
	}
}

func subnetOptFromRdata(rdata Rdata) Option {
	opt := rdata.(*OPT)
	if len(opt.Data) == 0 {
		return nil
	}

	buffer := util.NewInputBuffer(opt.Data)
	code, _ := buffer.ReadUint16()
	if code != EDNS_SUBNET {
		return nil
	}

	option, err := subnetOptFromWire(buffer)
	if err == nil {
		return option
	} else {
		return nil
	}
}

func (e *EDNS) AddSubnetV4(ip_ string) error {
	if ip := net.ParseIP(ip_); ip != nil {
		e.Options = append(e.Options, &SubnetOpt{
			family: 1,
			mask:   32,
			scope:  0,
			ip:     ip,
		})
		return nil
	} else {
		return fmt.Errorf("invalid ip address:%s", ip_)
	}
}
