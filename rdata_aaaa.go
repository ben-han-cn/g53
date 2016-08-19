package g53

import (
	"errors"
	"net"

	"g53/util"
)

type AAAA struct {
	Host net.IP
}

func (aaaa *AAAA) Rend(r *MsgRender) {
	rendField(RDF_C_IPV6, aaaa.Host, r)
}

func (aaaa *AAAA) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_IPV6, aaaa.Host, buffer)
}

func (aaaa *AAAA) String() string {
	return fieldToStr(RDF_D_IP, aaaa.Host)
}

func AAAAFromWire(buffer *util.InputBuffer, ll uint16) (*AAAA, error) {
	f, ll, err := fieldFromWire(RDF_C_IPV6, buffer, ll)
	if err != nil {
		return nil, err
	} else if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	} else {
		host, _ := f.(net.IP)
		return &AAAA{host.To16()}, nil
	}
}

func AAAAFromString(s string) (*AAAA, error) {
	f, err := fieldFromStr(RDF_D_IP, s)
	if err == nil {
		host, _ := f.(net.IP)
		return &AAAA{host}, nil
	} else {
		return nil, err
	}
}
