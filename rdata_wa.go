package g53

import (
	"errors"
	"math"
	"net"
	"regexp"
	"strings"

	"github.com/ben-han-cn/g53/util"
)

type WA struct {
	Weight uint16
	Host   net.IP
}

func (a *WA) Rend(r *MsgRender) {
	rendField(RDF_C_UINT16, a.Weight, r)
	rendField(RDF_C_IPV4, a.Host, r)
}

func (a *WA) ToWire(buf *util.OutputBuffer) {
	fieldToWire(RDF_C_UINT16, a.Weight, buf)
	fieldToWire(RDF_C_IPV4, a.Host, buf)
}

func (a *WA) Compare(other Rdata) int {
	return fieldCompare(RDF_C_IPV4, a.Host, other.(*WA).Host)
}

func (a *WA) String() string {
	return strings.Join([]string{
		fieldToString(RDF_D_INT, a.Weight),
		fieldToString(RDF_D_IP, a.Host)}, " ")
}

func WAFromWire(buf *util.InputBuffer, ll uint16) (*WA, error) {
	f, ll, err := fieldFromWire(RDF_C_UINT16, buf, ll)
	if err != nil {
		return nil, err
	}
	weight, _ := f.(uint16)

	f, ll, err = fieldFromWire(RDF_C_IPV4, buf, ll)
	if err != nil {
		return nil, err
	} else if ll != 0 {
		return nil, errors.New("extra data in a rdata part")
	} else {
		host, _ := f.(net.IP)
		return &WA{weight, host}, nil
	}
}

var waRdataTemplate = regexp.MustCompile(`^\s*(\S+)\s+(\S+)\s*$`)

func WAFromString(s string) (*WA, error) {
	fields := waRdataTemplate.FindStringSubmatch(s)
	if len(fields) != 3 {
		return nil, errors.New("fields count for wa isn't 3")
	}

	fields = fields[1:]
	f, err := fieldFromString(RDF_D_INT, fields[0])
	if err != nil {
		return nil, err
	}
	weight, _ := f.(int)
	if weight > math.MaxUint16 {
		return nil, ErrOutOfRange
	}

	f, err = fieldFromString(RDF_D_IP, fields[1])
	if err == nil {
		host, _ := f.(net.IP)
		return &WA{uint16(weight), host.To4()}, nil
	} else {
		return nil, err
	}
}