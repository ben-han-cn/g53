package g53

import (
	"fmt"

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
		return nil, fmt.Errorf("extra data %d in opt rdata part", ll)
	} else {
		d, _ := f.([]uint8)
		return &OPT{d}, nil
	}
}

func OPTFromString(s string) (*OPT, error) {
	f, err := fieldFromStr(RDF_D_HEX, s)
	if err == nil {
		d, _ := f.([]uint8)
		return &OPT{d}, nil
	} else {
		return nil, err
	}
}
