package g53

import (
	"errors"

	"github.com/ben-han-cn/g53/util"
)

type PTR struct {
	Name *Name
}

func (p *PTR) Rend(r *MsgRender) {
	rendField(RDF_C_NAME, p.Name, r)
}

func (p *PTR) ToWire(buf *util.OutputBuffer) {
	fieldToWire(RDF_C_NAME, p.Name, buf)
}

func (p *PTR) Compare(other Rdata) int {
	return fieldCompare(RDF_C_NAME, p.Name, other.(*PTR).Name)
}

func (p *PTR) String() string {
	return fieldToString(RDF_D_NAME, p.Name)
}

func (p *PTR) FromWire(buf *util.InputBuffer, ll uint16) error {
	n, ll, err := nameFieldFromWire(p.Name, buf, ll)
	if err != nil {
		return err
	} else if ll != 0 {
		return errors.New("extra data in cname rdata part")
	} else {
		p.Name = n
		return nil
	}
}

func PTRFromString(s string) (*PTR, error) {
	n, err := fieldFromString(RDF_D_NAME, s)
	if err == nil {
		name, _ := n.(*Name)
		return &PTR{name}, nil
	} else {
		return nil, err
	}
}
