package g53

import (
	"errors"
	"net"

	"github.com/ben-han-cn/g53/util"
)

type AAAA struct {
	Host net.IP
}

func (aaaa *AAAA) Rend(r *MsgRender) {
	rendField(RDF_C_IPV6, aaaa.Host, r)
}

func (aaaa *AAAA) ToWire(buf *util.OutputBuffer) {
	fieldToWire(RDF_C_IPV6, aaaa.Host, buf)
}

func (aaaa *AAAA) Compare(other Rdata) int {
	return fieldCompare(RDF_C_IPV6, aaaa.Host, other.(*AAAA).Host)
}

func (aaaa *AAAA) String() string {
	return fieldToString(RDF_D_IP, aaaa.Host)
}

func (aaaa *AAAA) FromWire(buf *util.InputBuffer, ll uint16) error {
	host, ll, err := ipv6FieldFromWire(aaaa.Host, buf, ll)
	if err != nil {
		return err
	} else if ll != 0 {
		return errors.New("extra data in rdata part")
	} else {
		aaaa.Host = host
		return nil
	}
}

func AAAAFromString(s string) (*AAAA, error) {
	f, err := fieldFromString(RDF_D_IP, s)
	if err == nil {
		host, _ := f.(net.IP)
		return &AAAA{host}, nil
	} else {
		return nil, err
	}
}
