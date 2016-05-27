package g53

import (
	"fmt"
	"net"

	"g53/util"
)

const (
	EDNS_SUBNET = 8
)

type SubnetOpt struct {
	family uint16
	mask   uint8
	scope  uint8
	ip     net.IP
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
