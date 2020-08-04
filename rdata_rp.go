package g53

import (
	"bytes"
	"errors"
	"regexp"

	"github.com/ben-han-cn/g53/util"
)

type RP struct {
	Mbox *Name
	Txt  *Name
}

func (rp *RP) Rend(r *MsgRender) {
	rendField(RDF_C_NAME, rp.Mbox, r)
	rendField(RDF_C_NAME, rp.Txt, r)
}

func (rp *RP) ToWire(buf *util.OutputBuffer) {
	fieldToWire(RDF_C_NAME, rp.Mbox, buf)
	fieldToWire(RDF_C_NAME, rp.Txt, buf)
}

func (rp *RP) Compare(other Rdata) int {
	ord := fieldCompare(RDF_C_NAME, rp.Mbox, other.(*RP).Mbox)
	if ord == 0 {
		return fieldCompare(RDF_C_NAME, rp.Txt, other.(*RP).Txt)
	} else {
		return ord
	}
}

func (rp *RP) String() string {
	var buf bytes.Buffer
	buf.WriteString(fieldToString(RDF_D_NAME, rp.Mbox))
	buf.WriteByte(' ')
	buf.WriteString(fieldToString(RDF_D_NAME, rp.Txt))
	return buf.String()
}

func (rp *RP) FromWire(buf *util.InputBuffer, ll uint16) error {
	n, ll, err := nameFieldFromWire(rp.Mbox, buf, ll)
	if err != nil {
		return err
	} else {
		rp.Mbox = n
	}

	n, ll, err = nameFieldFromWire(rp.Txt, buf, ll)
	if err != nil {
		return err
	} else if ll != 0 {
		return errors.New("extra data in cname rdata part")
	} else {
		rp.Txt = n
		return nil
	}
}

var rpRdataTemplate = regexp.MustCompile(`^\s*(\S+)\s+(\S+)\s*$`)

func RPFromString(s string) (*RP, error) {
	fields := rpRdataTemplate.FindStringSubmatch(s)
	if len(fields) != 3 {
		return nil, errors.New("short of fields for rp")
	}

	fields = fields[1:]
	mbox, err := fieldFromString(RDF_D_NAME, fields[0])
	if err != nil {
		return nil, err
	}

	txt, err := fieldFromString(RDF_D_NAME, fields[1])
	if err != nil {
		return nil, err
	}

	return &RP{mbox.(*Name), txt.(*Name)}, nil
}
