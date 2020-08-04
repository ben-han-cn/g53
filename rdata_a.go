package g53

import (
	"errors"
	"net"

	"github.com/ben-han-cn/g53/util"
)

type A struct {
	Host net.IP
}

func (a *A) Rend(r *MsgRender) {
	rendField(RDF_C_IPV4, a.Host, r)
}

func (a *A) ToWire(buf *util.OutputBuffer) {
	fieldToWire(RDF_C_IPV4, a.Host, buf)
}

func (a *A) Compare(other Rdata) int {
	return fieldCompare(RDF_C_IPV4, a.Host, other.(*A).Host)
}

func (a *A) String() string {
	return fieldToString(RDF_D_IP, a.Host)
}

func (a *A) FromWire(buf *util.InputBuffer, ll uint16) error {
	host, ll, err := ipv4FieldFromWire(a.Host, buf, ll)
	if err != nil {
		return err
	} else if ll != 0 {
		return errors.New("extra data in a rdata part")
	} else {
		a.Host = host
		return nil
	}
}

func AFromString(s string) (*A, error) {
	f, err := fieldFromString(RDF_D_IP, s)
	if err == nil {
		host, _ := f.(net.IP)
		return &A{host.To4()}, nil
	} else {
		return nil, err
	}
}
