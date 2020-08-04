package g53

import (
	"fmt"
	"sync"
)

var (
	rrsetPool      sync.Pool
	aRdataPool     sync.Pool
	aaaaRdataPool  sync.Pool
	cnameRdataPool sync.Pool
	soaRdataPool   sync.Pool
	nsRdataPool    sync.Pool
)

func acquireRRset() *RRset {
	rrset := rrsetPool.Get()
	if rrset != nil {
		target := rrset.(*RRset)
		if len(target.Rdatas) != 0 {
			panic(fmt.Sprintf("acquire return non-empty rdatas %s", target.String()))
		}
		return rrset.(*RRset)
	} else {
		rrset := &RRset{}
		if len(rrset.Rdatas) != 0 {
			panic(fmt.Sprintf("empty rrset with non-empty rdatas %s", rrset.String()))
		}
		return rrset
	}
}

func releaseRRset(rrset *RRset) {
	for i, rdata := range rrset.Rdatas {
		releaseRdata(rrset.Type, rdata)
		rrset.Rdatas[i] = nil
	}
	rrset.Rdatas = rrset.Rdatas[:0]

	if len(rrset.Rdatas) != 0 {
		panic(fmt.Sprintf("release rrset with non-empty rdatas %s", rrset.String()))
	}

	fmt.Printf("----> release rrset:%s\n", rrset.String())
	rrsetPool.Put(rrset)
}

func acquireRdata(typ RRType) Rdata {
	switch typ {
	case RR_A:
		a := aRdataPool.Get()
		if a != nil {
			return a.(*A)
		} else {
			return &A{}
		}
	case RR_AAAA:
		aaaa := aaaaRdataPool.Get()
		if aaaa != nil {
			return aaaa.(*AAAA)
		} else {
			return &AAAA{}
		}
	case RR_CNAME:
		cname := cnameRdataPool.Get()
		if cname != nil {
			return cname.(*CName)
		} else {
			return &CName{}
		}
	case RR_SOA:
		soa := soaRdataPool.Get()
		if soa != nil {
			return soa.(*SOA)
		} else {
			return &SOA{}
		}
	case RR_NS:
		ns := nsRdataPool.Get()
		if ns != nil {
			return ns.(*NS)
		} else {
			return &NS{}
		}
	case RR_OPT:
		return &OPT{}
	case RR_PTR:
		return &PTR{}
	case RR_SRV:
		return &SRV{}
	case RR_NAPTR:
		return &NAPTR{}
	case RR_DNAME:
		return &DName{}
	case RR_RRSIG:
		return &RRSig{}
	case RR_MX:
		return &MX{}
	case RR_TXT:
		return &Txt{}
	case RR_RP:
		return &RP{}
	case RR_SPF:
		return &SPF{}
	case RR_TSIG:
		return &TSIG{}
	case RR_NSEC3:
		return &NSEC3{}
	case RR_DS:
		return &DS{}
	case RR_WA:
		return &WA{}
	case RR_WAAAA:
		return &WAAAA{}
	case RR_WCNAME:
		return &WCName{}
	default:
		return nil
	}
}

func releaseRdata(typ RRType, rdata Rdata) {
	switch typ {
	case RR_A:
		if _, ok := rdata.(*A); !ok {
			panic(fmt.Sprintf("put not a %#v", rdata))
		}
		aRdataPool.Put(rdata)
	case RR_AAAA:
		if _, ok := rdata.(*AAAA); !ok {
			panic(fmt.Sprintf("put not aaaa %#v", rdata))
		}
		aaaaRdataPool.Put(rdata)
	case RR_CNAME:
		if _, ok := rdata.(*CName); !ok {
			panic(fmt.Sprintf("put not cname %#v", rdata))
		}
		cnameRdataPool.Put(rdata)
	case RR_SOA:
		if _, ok := rdata.(*SOA); !ok {
			panic(fmt.Sprintf("put not soa %#v", rdata))
		}
		soaRdataPool.Put(rdata)
	case RR_NS:
		if _, ok := rdata.(*NS); !ok {
			panic(fmt.Sprintf("put not ns %#v", rdata))
		}
		nsRdataPool.Put(rdata)
	default:
	}
}
