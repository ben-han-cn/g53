package g53

import (
	"errors"
	"strings"

	"g53/util"
)

type MX struct {
	Preference uint16
	Exchange   *Name
}

func (mx *MX) Rend(r *MsgRender) {
	rendField(RDF_C_UINT16, mx.Preference, r)
	rendField(RDF_C_NAME, mx.Exchange, r)
}

func (mx *MX) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_UINT16, mx.Preference, buffer)
	fieldToWire(RDF_C_NAME, mx.Exchange, buffer)
}

func (mx *MX) String() string {
	return strings.Join([]string{
		fieldToStr(RDF_D_INT, mx.Preference),
		fieldToStr(RDF_D_NAME, mx.Exchange)}, " ")
}

func MXFromWire(buffer *util.InputBuffer, ll uint16) (*MX, error) {
	f, ll, err := fieldFromWire(RDF_C_UINT16, buffer, ll)
	if err != nil {
		return nil, err
	}
	preference, _ := f.(uint16)

	f, ll, err = fieldFromWire(RDF_C_NAME, buffer, ll)
	if err != nil {
		return nil, err
	}
	exchange, _ := f.(*Name)

	if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	}

	return &MX{preference, exchange}, nil
}

func MXFromStr(s string) (*MX, error) {
	fields := strings.Split(s, " ")
	if len(fields) != 2 {
		return nil, errors.New("fields count for mx isn't 2")
	}

	f, err := fieldFromStr(RDF_D_INT, fields[0])
	if err != nil {
		return nil, err
	}
	preference, _ := f.(int)

	f, err = fieldFromStr(RDF_D_NAME, fields[1])
	if err != nil {
		return nil, err
	}
	exchange, _ := f.(*Name)
	return &MX{uint16(preference), exchange}, nil
}
