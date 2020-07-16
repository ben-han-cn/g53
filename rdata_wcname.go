package g53

import (
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/ben-han-cn/g53/util"
)

type WCName struct {
	Priority uint16
	Weight   uint16
	Name     *Name
}

func (c *WCName) Rend(r *MsgRender) {
	rendField(RDF_C_UINT16, c.Priority, r)
	rendField(RDF_C_UINT16, c.Weight, r)
	rendField(RDF_C_NAME, c.Name, r)
}

func (c *WCName) ToWire(buf *util.OutputBuffer) {
	fieldToWire(RDF_C_UINT16, c.Priority, buf)
	fieldToWire(RDF_C_UINT16, c.Weight, buf)
	fieldToWire(RDF_C_NAME, c.Name, buf)
}

func (c *WCName) String() string {
	return strings.Join([]string{
		fieldToString(RDF_D_INT, c.Priority),
		fieldToString(RDF_D_INT, c.Weight),
		fieldToString(RDF_D_NAME, c.Name)}, " ")
}

func (c *WCName) Compare(other_ Rdata) int {
	other := other_.(*WCName)
	order := fieldCompare(RDF_C_UINT16, c.Priority, other.Priority)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT16, c.Weight, other.Weight)
	if order != 0 {
		return order
	}

	return fieldCompare(RDF_C_NAME, c.Name, other.Name)
}

func WCNameFromWire(buf *util.InputBuffer, ll uint16) (*WCName, error) {
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

	n, ll, err := fieldFromWire(RDF_C_NAME, buf, ll)
	if err != nil {
		return nil, err
	} else if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	} else {
		name, _ := n.(*Name)
		return &WCName{priority, weight, name}, nil
	}
}

var wcnameRdataTemplate = regexp.MustCompile(`^\s*(\S+)\s+(\S+)\s+(\S+)\s*$`)

func WCNameFromString(s string) (*WCName, error) {
	fields := wcnameRdataTemplate.FindStringSubmatch(s)
	if len(fields) != 4 {
		return nil, errors.New("fields count for wcname isn't 3")
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

	n, err := fieldFromString(RDF_D_NAME, fields[2])
	if err == nil {
		name, _ := n.(*Name)
		return &WCName{uint16(priority), uint16(weight), name}, nil
	} else {
		return nil, err
	}
}
