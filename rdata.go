package g53

import (
	"errors"

	"g53/util"
)

type Rdata interface {
	Rend(r *MsgRender)
	ToWire(buffer *util.OutputBuffer)
	String() string
}

/*
var registor = map[RRType][]RDFieldType{
	RR_SRV:        []RDFieldType{RDF_UINT16, RDF_UINT16, RDF_UINT16, RDF_NAME},
	RR_NAPTR:      []RDFieldType{RDF_UINT16, RDF_UINT16, RDF_STR, RDF_STR, RDF_STR, RDF_NAME},
	RR_DS:         []RDFieldType{RDF_UINT16, RDF_UINT8, RDF_UINT8, RDF_B64},
	RR_RRSIG:      []RDFieldType{RDF_UINT16, RDF_UINT8, RDF_UINT8, RDF_UINT32, RDF_UINT32, RDF_UINT32, RDF_UINT16, RDF_NAME, RDF_B64},
	RR_DNSKEY:     []RDFieldType{RDF_UINT16, RDF_UINT8, RDF_UINT8, RDF_B64},
	RR_NSEC3:      []RDFieldType{RDF_UINT8, RDF_UINT8, RDF_UINT16, RDF_MID_BINARY, RDF_MID_BINARY, RDF_B64},
	RR_NSEC3PARAM: []RDFieldType{RDF_UINT8, RDF_UINT8, RDF_UINT16, RDF_MID_BINARY},
}
*/

func RdataFromWire(t RRType, buffer *util.InputBuffer) (Rdata, error) {
	rdlen, err := buffer.ReadUint16()
	if err != nil {
		return nil, err
	}

	switch t {
	case RR_A:
		return AFromWire(buffer, rdlen)
	case RR_AAAA:
		return AAAAFromWire(buffer, rdlen)
	case RR_CNAME:
		return CNameFromWire(buffer, rdlen)
	case RR_SOA:
		return SOAFromWire(buffer, rdlen)
	case RR_NS:
		return NSFromWire(buffer, rdlen)
	case RR_OPT:
		return OPTFromWire(buffer, rdlen)
	default:
		return nil, errors.New("unimplement type")
	}
}

func RdataFromStr(t RRType, s string) (Rdata, error) {
	switch t {
	case RR_A:
		return AFromStr(s)
	case RR_AAAA:
		return AAAAFromStr(s)
	case RR_CNAME:
		return CNameFromStr(s)
	case RR_SOA:
		return SOAFromStr(s)
	case RR_NS:
		return NSFromStr(s)
	case RR_OPT:
		return OPTFromStr(s)
	default:
		return nil, errors.New("unimplement type")
	}
}
