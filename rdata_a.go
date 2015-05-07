package g53

import (
	"errors"
	"net"

	"g53/util"
)

type A struct {
	Host net.IP
}

func (a *A) Rend(r *MsgRender) {
	rendField(RDF_C_IPV4, a.Host, r)
}

func (a *A) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_IPV4, a.Host, buffer)
}

func (a *A) String() string {
	return fieldToStr(RDF_D_IP, a.Host)
}

func AFromWire(buffer *util.InputBuffer, ll uint16) (*A, error) {
	f, ll, err := fieldFromWire(RDF_C_IPV4, buffer, ll)
	if err != nil {
		return nil, err
	} else if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	} else {
		host, _ := f.(net.IP)
		return &A{host}, nil
	}
}

func AFromStr(s string) (*A, error) {
	f, err := fieldFromStr(RDF_D_IP, s)
	if err == nil {
		host, _ := f.(net.IP)
		return &A{host}, nil
	} else {
		return nil, err
	}
}
