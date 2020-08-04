package g53

import (
	"fmt"

	"github.com/ben-han-cn/g53/util"
)

type Rdata interface {
	Rend(r *MsgRender)
	FromWire(*util.InputBuffer, uint16) error
	ToWire(*util.OutputBuffer)
	Compare(Rdata) int
	String() string
}

func RdataFromWire(typ RRType, buf *util.InputBuffer) (Rdata, error) {
	rdlen, err := buf.ReadUint16()
	if err != nil {
		return nil, err
	}

	//RR_OPT or rr in UPDATE message may have empty rdlen
	if rdlen == 0 {
		return nil, nil
	}

	rdata := acquireRdata(typ)
	if rdata == nil {
		return nil, fmt.Errorf("unknown rr type %s", typ.String())
	} else if err := rdata.FromWire(buf, rdlen); err != nil {
		releaseRdata(typ, rdata)
		return nil, err
	} else {
		return rdata, nil
	}
}

func RdataFromString(t RRType, s string) (Rdata, error) {
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
	case RR_DNAME:
		return DNameFromString(s)
	case RR_RRSIG:
		return RRSigFromString(s)
	case RR_MX:
		return MXFromString(s)
	case RR_TXT:
		return TxtFromString(s)
	case RR_RP:
		return RPFromString(s)
	case RR_SPF:
		return SPFFromString(s)
	case RR_NSEC3:
		return NSEC3FromString(s)
	case RR_DS:
		return DSFromString(s)
	case RR_WA:
		return WAFromString(s)
	case RR_WAAAA:
		return WAAAAFromString(s)
	case RR_WCNAME:
		return WCNameFromString(s)
	default:
		return nil, fmt.Errorf("unimplement type: %v", t)
	}
}
