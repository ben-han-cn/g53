package g53

import (
	"errors"

	"github.com/ben-han-cn/g53/util"
)

type NS struct {
	Name *Name
}

func (ns *NS) Rend(r *MsgRender) {
	rendField(RDF_C_NAME, ns.Name, r)
}

func (ns *NS) ToWire(buf *util.OutputBuffer) {
	fieldToWire(RDF_C_NAME, ns.Name, buf)
}

func (ns *NS) Compare(other Rdata) int {
	return fieldCompare(RDF_C_NAME, ns.Name, other.(*NS).Name)
}

func (ns *NS) String() string {
	return fieldToString(RDF_D_NAME, ns.Name)
}

func (ns *NS) FromWire(buf *util.InputBuffer, ll uint16) error {
	n, ll, err := nameFieldFromWire(ns.Name, buf, ll)
	if err != nil {
		return err
	} else if ll != 0 {
		return errors.New("extra data in cname rdata part")
	} else {
		ns.Name = n
		return nil
	}
}

func NSFromString(s string) (*NS, error) {
	n, err := fieldFromString(RDF_D_NAME, s)
	if err == nil {
		name, _ := n.(*Name)
		return &NS{name}, nil
	} else {
		return nil, err
	}
}
