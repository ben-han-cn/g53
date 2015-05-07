package g53

import (
	"errors"

	"g53/util"
)

type NS struct {
	Name *Name
}

func (ns *NS) Rend(r *MsgRender) {
	rendField(RDF_C_NAME, ns.Name, r)
}

func (ns *NS) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_NAME, ns.Name, buffer)
}

func (ns *NS) String() string {
	return fieldToStr(RDF_D_NAME, ns.Name)
}

func NSFromWire(buffer *util.InputBuffer, ll uint16) (*NS, error) {
	n, ll, err := fieldFromWire(RDF_C_NAME, buffer, ll)
	if err != nil {
		return nil, err
	} else if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	} else {
		name, _ := n.(*Name)
		return &NS{name}, nil
	}
}

func NSFromStr(s string) (*NS, error) {
	n, err := fieldFromStr(RDF_D_NAME, s)
	if err == nil {
		name, _ := n.(*Name)
		return &NS{name}, nil
	} else {
		return nil, err
	}
}
