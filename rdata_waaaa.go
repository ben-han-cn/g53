package g53

import (
	"errors"
	"math"
	"net"
	"regexp"
	"strings"

	"github.com/ben-han-cn/g53/util"
)

type WAAAA struct {
	Priority uint16
	Weight   uint16
	Host     net.IP
}

func (aaaa *WAAAA) Rend(r *MsgRender) {
	rendField(RDF_C_UINT16, aaaa.Priority, r)
	rendField(RDF_C_UINT16, aaaa.Weight, r)
	rendField(RDF_C_IPV6, aaaa.Host, r)
}

func (aaaa *WAAAA) ToWire(buf *util.OutputBuffer) {
	fieldToWire(RDF_C_UINT16, aaaa.Priority, buf)
	fieldToWire(RDF_C_UINT16, aaaa.Weight, buf)
	fieldToWire(RDF_C_IPV6, aaaa.Host, buf)
}

func (aaaa *WAAAA) Compare(other_ Rdata) int {
	other := other_.(*WAAAA)
	order := fieldCompare(RDF_C_UINT16, aaaa.Priority, other.Priority)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT16, aaaa.Weight, other.Weight)
	if order != 0 {
		return order
	}

	return fieldCompare(RDF_C_IPV6, aaaa.Host, other.Host)
}

func (aaaa *WAAAA) String() string {
	return strings.Join([]string{
		fieldToString(RDF_D_INT, aaaa.Priority),
		fieldToString(RDF_D_INT, aaaa.Weight),
		fieldToString(RDF_D_IP, aaaa.Host)}, " ")
}

func WAAAAFromWire(buf *util.InputBuffer, ll uint16) (*WAAAA, error) {
	f, ll, err := fieldFromWire(RDF_C_UINT16, buf, ll)
	if err != nil {
		return nil, err
	}
	priority, _ := f.(uint16)

	f, ll, err = fieldFromWire(RDF_C_UINT16, buf, ll)
	if err != nil {
		return nil, err
	}
	weight, _ := f.(uint16)

	f, ll, err = fieldFromWire(RDF_C_IPV6, buf, ll)
	if err != nil {
		return nil, err
	} else if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	} else {
		host, _ := f.(net.IP)
		return &WAAAA{priority, weight, host.To16()}, nil
	}
}

var waaaaRdataTemplate = regexp.MustCompile(`^\s*(\S+)\s+(\S+)\s+(\S+)\s*$`)

func WAAAAFromString(s string) (*WAAAA, error) {
	fields := waaaaRdataTemplate.FindStringSubmatch(s)
	if len(fields) != 4 {
		return nil, errors.New("fields count for wa isn't 3")
	}

	fields = fields[1:]
	f, err := fieldFromString(RDF_D_INT, fields[0])
	if err != nil {
		return nil, err
	}
	priority, _ := f.(int)
	if priority > math.MaxUint16 {
		return nil, ErrOutOfRange
	}

	fields = fields[2:]
	f, err = fieldFromString(RDF_D_INT, fields[1])
	if err != nil {
		return nil, err
	}
	weight, _ := f.(int)
	if weight > math.MaxUint16 {
		return nil, ErrOutOfRange
	}

	f, err = fieldFromString(RDF_D_IP, fields[2])
	if err == nil {
		host, _ := f.(net.IP)
		return &WAAAA{uint16(priority), uint16(weight), host.To4()}, nil
	} else {
		return nil, err
	}

}
