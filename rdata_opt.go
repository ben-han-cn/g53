package g53

import (
	"errors"

	"g53/util"
)

type OPT struct {
	Data []uint8
}

func (opt *OPT) Rend(r *MsgRender) {
	rendField(RDF_C_BINARY, opt.Data, r)
}

func (opt *OPT) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_BINARY, opt.Data, buffer)
}

func (opt *OPT) String() string {
	return fieldToStr(RDF_D_HEX, opt.Data)
}

func OPTFromWire(buffer *util.InputBuffer, ll uint16) (*OPT, error) {
	f, ll, err := fieldFromWire(RDF_C_BINARY, buffer, ll)

	if err != nil {
		return nil, err
	} else if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	} else {
		d, _ := f.([]uint8)
		return &OPT{d}, nil
	}
}

func OPTFromStr(s string) (*OPT, error) {
	f, err := fieldFromStr(RDF_D_HEX, s)
	if err == nil {
		d, _ := f.([]uint8)
		return &OPT{d}, nil
	} else {
		return nil, err
	}
}
