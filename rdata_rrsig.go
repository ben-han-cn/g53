package g53

import (
	"bytes"
	"errors"
	"strings"

	"g53/util"
)

type RRSig struct {
	Covered     RRType
	Algorithm   uint8
	Labels      uint8
	OriginalTtl uint32
	SigExpire   uint32
	Inception   uint32
	Tag         uint16
	Signer      *Name
	Signature   []uint8
}

func (rrsig *RRSig) Rend(r *MsgRender) {
	rendField(RDF_C_UINT16, uint16(rrsig.Covered), r)
	rendField(RDF_C_UINT8, rrsig.Algorithm, r)
	rendField(RDF_C_UINT8, rrsig.Labels, r)
	rendField(RDF_C_UINT32, rrsig.OriginalTtl, r)
	rendField(RDF_C_UINT32, rrsig.SigExpire, r)
	rendField(RDF_C_UINT32, rrsig.Inception, r)
	rendField(RDF_C_UINT16, rrsig.Tag, r)
	rendField(RDF_C_NAME, rrsig.Signer, r)
	rendField(RDF_C_BINARY, rrsig.Signature, r)
}

func (rrsig *RRSig) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_UINT16, uint16(rrsig.Covered), buffer)
	fieldToWire(RDF_C_UINT8, rrsig.Algorithm, buffer)
	fieldToWire(RDF_C_UINT8, rrsig.Labels, buffer)
	fieldToWire(RDF_C_UINT32, rrsig.OriginalTtl, buffer)
	fieldToWire(RDF_C_UINT32, rrsig.SigExpire, buffer)
	fieldToWire(RDF_C_UINT32, rrsig.Inception, buffer)
	fieldToWire(RDF_C_UINT16, rrsig.Tag, buffer)
	fieldToWire(RDF_C_NAME, rrsig.Signer, buffer)
	fieldToWire(RDF_C_BINARY, rrsig.Signature, buffer)
}

func (rrsig *RRSig) Compare(other Rdata) int {
	otherRRSig := other.(*RRSig)
	order := fieldCompare(RDF_C_UINT16, uint16(rrsig.Covered), uint16(otherRRSig.Covered))
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT8, rrsig.Algorithm, otherRRSig.Algorithm)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT8, rrsig.Labels, otherRRSig.Labels)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT32, rrsig.OriginalTtl, otherRRSig.OriginalTtl)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT32, rrsig.SigExpire, otherRRSig.SigExpire)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT32, rrsig.Inception, otherRRSig.Inception)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT16, rrsig.Tag, otherRRSig.Tag)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_NAME, rrsig.Signer, otherRRSig.Signer)
	if order != 0 {
		return order
	}

	return fieldCompare(RDF_C_BINARY, rrsig.Signature, otherRRSig.Signature)
}

func (rrsig *RRSig) String() string {
	var buf bytes.Buffer
	buf.WriteString(fieldToStr(RDF_D_STR, rrsig.Covered.String()))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_INT, rrsig.Algorithm))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_INT, rrsig.Labels))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_INT, rrsig.OriginalTtl))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_INT, rrsig.SigExpire))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_INT, rrsig.Inception))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_INT, rrsig.Tag))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_NAME, rrsig.Signer))
	buf.WriteString(" ")
	buf.WriteString(fieldToStr(RDF_D_B64, rrsig.Signature))
	return buf.String()
}

func RRSigFromWire(buffer *util.InputBuffer, ll uint16) (*RRSig, error) {
	covered, ll, err := fieldFromWire(RDF_C_UINT16, buffer, ll)
	if err != nil {
		return nil, err
	}

	algorithm, ll, err := fieldFromWire(RDF_C_UINT8, buffer, ll)
	if err != nil {
		return nil, err
	}

	labels, ll, err := fieldFromWire(RDF_C_UINT8, buffer, ll)
	if err != nil {
		return nil, err
	}

	originalTtl, ll, err := fieldFromWire(RDF_C_UINT32, buffer, ll)
	if err != nil {
		return nil, err
	}

	sigExpire, ll, err := fieldFromWire(RDF_C_UINT32, buffer, ll)
	if err != nil {
		return nil, err
	}

	inception, ll, err := fieldFromWire(RDF_C_UINT32, buffer, ll)
	if err != nil {
		return nil, err
	}

	tag, ll, err := fieldFromWire(RDF_C_UINT16, buffer, ll)
	if err != nil {
		return nil, err
	}

	signer, ll, err := fieldFromWire(RDF_C_NAME, buffer, ll)
	if err != nil {
		return nil, err
	}

	signature, ll, err := fieldFromWire(RDF_C_BINARY, buffer, ll)
	if err != nil {
		return nil, err
	}

	if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	}

	return &RRSig{RRType(covered.(uint16)), algorithm.(uint8), labels.(uint8), originalTtl.(uint32), sigExpire.(uint32), inception.(uint32), tag.(uint16), signer.(*Name), signature.([]uint8)}, nil
}

func RRSigFromString(s string) (*RRSig, error) {
	fields := strings.Split(s, " ")
	if len(fields) != 9 {
		return nil, errors.New("short of fields for rrsig")
	}

	covered, err := TypeFromString(fields[0])
	if err != nil {
		return nil, err
	}

	algorithm, err := fieldFromStr(RDF_D_INT, fields[1])
	if err != nil {
		return nil, err
	}

	labels, err := fieldFromStr(RDF_D_INT, fields[2])
	if err != nil {
		return nil, err
	}

	originalTtl, err := fieldFromStr(RDF_D_INT, fields[3])
	if err != nil {
		return nil, err
	}

	sigExpire, err := fieldFromStr(RDF_D_INT, fields[4])
	if err != nil {
		return nil, err
	}

	inception, err := fieldFromStr(RDF_D_INT, fields[5])
	if err != nil {
		return nil, err
	}

	tag, err := fieldFromStr(RDF_D_INT, fields[6])
	if err != nil {
		return nil, err
	}

	signer, err := fieldFromStr(RDF_D_NAME, fields[7])
	if err != nil {
		return nil, err
	}

	signature, err := fieldFromStr(RDF_D_B64, fields[8])
	if err != nil {
		return nil, err
	}

	return &RRSig{covered, uint8(algorithm.(int)), uint8(labels.(int)), uint32(originalTtl.(int)), uint32(sigExpire.(int)), uint32(inception.(int)), uint16(tag.(int)), signer.(*Name), signature.([]uint8)}, nil
}
