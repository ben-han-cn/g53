package g53

import (
	"errors"

	"g53/util"
)

type SPF struct {
	Data []string
}

func (spf *SPF) Rend(r *MsgRender) {
	rendField(RDF_C_TXT, spf.Data, r)
}

func (spf *SPF) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_TXT, spf.Data, buffer)
}

func (spf *SPF) String() string {
	return fieldToStr(RDF_D_TXT, spf.Data)
}

func SPFFromWire(buffer *util.InputBuffer, ll uint16) (*SPF, error) {
	f, ll, err := fieldFromWire(RDF_C_TXT, buffer, ll)
	if err != nil {
		return nil, err
	} else if ll != 0 {
		return nil, errors.New("extra data in rdata part when parse spf")
	} else {
		data, _ := f.([]string)
		return &SPF{data}, nil
	}
}

func SPFFromString(s string) (*SPF, error) {
	f, err := fieldFromStr(RDF_D_TXT, s)
	if err != nil {
		return nil, err
	} else {
		return &SPF{f.([]string)}, nil
	}
}
