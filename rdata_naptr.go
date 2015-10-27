package g53

import (
	"bytes"
	"errors"
	"strings"

	"g53/util"
)

type NAPTR struct {
	Order       uint16
	Preference  uint16
	Flags       string
	Services    string
	Regexp      string
	Replacement *Name
}

func (naptr *NAPTR) Rend(r *MsgRender) {
	rendField(RDF_C_UINT16, naptr.Order, r)
	rendField(RDF_C_UINT16, naptr.Preference, r)
	rendField(RDF_C_BYTE_BINARY, []byte(naptr.Flags), r)
	rendField(RDF_C_BYTE_BINARY, []byte(naptr.Services), r)
	rendField(RDF_C_BYTE_BINARY, []byte(naptr.Regexp), r)
	rendField(RDF_C_NAME, naptr.Replacement, r)
}

func (naptr *NAPTR) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_UINT16, naptr.Order, buffer)
	fieldToWire(RDF_C_UINT16, naptr.Preference, buffer)
	fieldToWire(RDF_C_BYTE_BINARY, []byte(naptr.Flags), buffer)
	fieldToWire(RDF_C_BYTE_BINARY, []byte(naptr.Services), buffer)
	fieldToWire(RDF_C_BYTE_BINARY, []byte(naptr.Regexp), buffer)
	fieldToWire(RDF_C_NAME, naptr.Replacement, buffer)
}

func (naptr *NAPTR) String() string {
	var buf bytes.Buffer
	buf.WriteString(fieldToStr(RDF_D_INT, naptr.Order))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_INT, naptr.Preference))
	buf.WriteString(" ")
	buf.WriteString(strings.Join([]string{"\"", fieldToStr(RDF_D_STR, naptr.Flags), "\""}, ""))
	buf.WriteString(" ")
	buf.WriteString(strings.Join([]string{"\"", fieldToStr(RDF_D_STR, naptr.Services), "\""}, ""))
	buf.WriteString(" ")
	buf.WriteString(strings.Join([]string{"\"", fieldToStr(RDF_D_STR, naptr.Regexp), "\""}, ""))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_NAME, naptr.Replacement))
	return buf.String()
}

func NAPTRFromWire(buffer *util.InputBuffer, ll uint16) (*NAPTR, error) {
	o, ll, err := fieldFromWire(RDF_C_UINT16, buffer, ll)
	if err != nil {
		return nil, err
	}
	order, _ := o.(uint16)

	p, ll, err := fieldFromWire(RDF_C_UINT16, buffer, ll)
	if err != nil {
		return nil, err
	}
	preference, _ := p.(uint16)

	f, ll, err := fieldFromWire(RDF_C_BYTE_BINARY, buffer, ll)
	if err != nil {
		return nil, err
	}
	f_, _ := f.([]uint8)
	flags := string(f_)

	s, ll, err := fieldFromWire(RDF_C_BYTE_BINARY, buffer, ll)
	if err != nil {
		return nil, err
	}
	s_, _ := s.([]uint8)
	service := string(s_)

	r, ll, err := fieldFromWire(RDF_C_BYTE_BINARY, buffer, ll)
	if err != nil {
		return nil, err
	}
	r_, _ := r.([]uint8)
	regex := string(r_)

	n, ll, err := fieldFromWire(RDF_C_NAME, buffer, ll)
	if err != nil {
		return nil, err
	}
	replacement, _ := n.(*Name)

	if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	}

	return &NAPTR{order, preference, flags, service, regex, replacement}, nil
}

func NAPTRFromString(s string) (*NAPTR, error) {
	fields := strings.Split(s, " ")
	if len(fields) != 6 {
		return nil, errors.New("short of fields for naptr")
	}

	o, err := fieldFromStr(RDF_D_INT, fields[0])
	if err != nil {
		return nil, err
	}
	order, _ := o.(uint16)

	p, err := fieldFromStr(RDF_D_INT, fields[1])
	if err != nil {
		return nil, err
	}
	preference, _ := p.(uint16)

	f, err := fieldFromStr(RDF_D_STR, fields[2])
	if err != nil {
		return nil, err
	}
	flags, _ := f.(string)

	se, err := fieldFromStr(RDF_D_STR, fields[3])
	if err != nil {
		return nil, err
	}
	service, _ := se.(string)

	r, err := fieldFromStr(RDF_D_STR, fields[4])
	if err != nil {
		return nil, err
	}
	regex, _ := r.(string)

	n, err := fieldFromStr(RDF_D_NAME, fields[5])
	if err != nil {
		return nil, err
	}
	replacement, _ := n.(*Name)

	return &NAPTR{order, preference, flags, service, regex, replacement}, nil
}
