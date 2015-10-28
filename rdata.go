package g53

import (
	"errors"
	"fmt"

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
	case RR_PTR:
		return PTRFromWire(buffer, rdlen)
	case RR_SRV:
		return SRVFromWire(buffer, rdlen)
	case RR_NAPTR:
		return NAPTRFromWire(buffer, rdlen)
	case RR_RRSIG:
		return RRSigFromWire(buffer, rdlen)
	default:
		return nil, fmt.Errorf("unimplement type: %v", t)
	}
}

func RdataFromStr(t RRType, s string) (Rdata, error) {
	switch t {
	case RR_A:
		return AFromString(s)
	case RR_AAAA:
		return AAAAFromString(s)
	case RR_CNAME:
		return CNameFromString(s)
	case RR_SOA:
		return SOAFromString(s)
	case RR_NS:
		return NSFromString(s)
	case RR_OPT:
		return OPTFromString(s)
	case RR_PTR:
		return PTRFromString(s)
	case RR_SRV:
		return SRVFromString(s)
	case RR_NAPTR:
		return NAPTRFromString(s)
	case RR_RRSIG:
		return RRSigFromString(s)
	default:
		return nil, errors.New("unimplement type")
	}
}
