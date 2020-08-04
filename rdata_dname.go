package g53

import (
	"errors"

	"github.com/ben-han-cn/g53/util"
)

type DName struct {
	Target *Name
}

func (c *DName) Rend(r *MsgRender) {
	rendField(RDF_C_NAME, c.Target, r)
}

func (c *DName) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_NAME, c.Target, buffer)
}

func (c *DName) Compare(other Rdata) int {
	return fieldCompare(RDF_C_NAME, c.Target, other.(*DName).Target)
}

func (c *DName) String() string {
	return fieldToString(RDF_D_NAME, c.Target)
}

func DNameFromWire(buf *util.InputBuffer, ll uint16) (*DName, error) {
	var d DName
	if err := d.FromWire(buf, ll); err != nil {
		return nil, err
	} else {
		return &d, nil
	}
}

func (d *DName) FromWire(buf *util.InputBuffer, ll uint16) error {
	n, ll, err := nameFieldFromWire(d.Target, buf, ll)
	if err != nil {
		return err
	} else if ll != 0 {
		return errors.New("extra data in cname rdata part")
	} else {
		d.Target = n
		return nil
	}
}

func DNameFromString(s string) (*DName, error) {
	n, err := fieldFromString(RDF_D_NAME, s)
	if err == nil {
		name, _ := n.(*Name)
		return &DName{name}, nil
	} else {
		return nil, err
	}
}
