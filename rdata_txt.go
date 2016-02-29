package g53

import (
	"errors"

	"g53/util"
)

type Txt struct {
	Data []string
}

func (txt *Txt) Rend(r *MsgRender) {
	rendField(RDF_C_TXT, txt.Data, r)
}

func (txt *Txt) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_TXT, txt.Data, buffer)
}

func (txt *Txt) String() string {
	return fieldToStr(RDF_D_TXT, txt.Data)
}

func TxtFromWire(buffer *util.InputBuffer, ll uint16) (*Txt, error) {
	f, ll, err := fieldFromWire(RDF_C_TXT, buffer, ll)
	if err != nil {
		return nil, err
	} else if ll != 0 {
		return nil, errors.New("extra data in rdata part when parse txt")
	} else {
		data, _ := f.([]string)
		return &Txt{data}, nil
	}
}

func TxtFromString(s string) (*Txt, error) {
	f, err := fieldFromStr(RDF_D_TXT, s)
	if err != nil {
		return nil, err
	} else {
		return &Txt{f.([]string)}, nil
	}
}
