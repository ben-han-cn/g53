package g53

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"g53/util"
)

type RRTTL uint32
type RRClass uint16
type RRType uint16

const (
	CLASS_IN   RRClass = 1
	CLASS_CH           = 3
	CLASS_HS           = 4
	CLASS_NONE         = 254
	CLASS_ANY          = 255
)

const (
	/** a host address */
	RR_A RRType = 1
	/** an authoritative name server */
	RR_NS = 2
	/** the canonical name for an alias */
	RR_CNAME = 5
	/**  marks the start of a zone of authority */
	RR_SOA = 6
	/**  a mailbox domain name (EXPERIMENTAL) */
	RR_MB = 7
	/**  a mail group member (EXPERIMENTAL) */
	RR_MG = 8
	/**  a mail rename domain name (EXPERIMENTAL) */
	RR_MR = 9
	/**  a null RR (EXPERIMENTAL) */
	RR_NULL = 10
	/**  a well known service description */
	RR_WKS = 11
	/**  a domain name pointer */
	RR_PTR = 12
	/**  host information */
	RR_HINFO = 13
	/**  mailbox or mail list information */
	RR_MINFO = 14
	/**  mail exchange */
	RR_MX = 15
	/**  text strings */
	RR_TXT = 16
	/**  RFC1183 */
	RR_RP = 17
	/**  RFC1183 */
	RR_AFSDB = 18
	/**  RFC1183 */
	RR_X25 = 19
	/**  RFC1183 */
	RR_ISDN = 20
	/**  RFC1183 */
	RR_RT = 21
	/**  RFC1706 */
	RR_NSAP = 22
	/**  RFC1348 */
	RR_NSAP_PTR = 23
	/**  2535typecode */
	RR_SIG = 24
	/**  2535typecode */
	RR_KEY = 25
	/**  RFC2163 */
	RR_PX = 26
	/**  RFC1712 */
	RR_GPOS = 27
	/**  ipv6 address */
	RR_AAAA = 28
	/**  LOC record  RFC1876 */
	RR_LOC = 29
	/**  2535typecode */
	RR_NXT = 30
	/**  draft-ietf-nimrod-dns-01.txt */
	RR_EID = 31
	/**  draft-ietf-nimrod-dns-01.txt */
	RR_NIMLOC = 32
	/**  SRV record RFC2782 */
	RR_SRV = 33
	/**  http://www.jhsoft.com/rfc/af-saa-0069.000.rtf */
	RR_ATMA = 34
	/**  RFC2915 */
	RR_NAPTR = 35
	/**  RFC2230 */
	RR_KX = 36
	/**  RFC2538 */
	RR_CERT = 37
	/**  RFC2874 */
	RR_A6 = 38
	/**  RFC2672 */
	RR_DNAME = 39
	/**  dnsind-kitchen-sink-02.txt */
	RR_SINK = 40
	/**  Pseudo OPT record... */
	RR_OPT = 41
	/**  RFC3123 */
	RR_APL = 42
	/**  RFC4034 RFC3658 */
	RR_DS = 43
	/**  SSH Key Fingerprint */
	RR_SSHFP = 44 /* RFC 4255 */
	/**  IPsec Key */
	RR_IPSECKEY = 45 /* RFC 4025 */
	/**  DNSSEC */
	RR_RRSIG  = 46 /* RFC 4034 */
	RR_NSEC   = 47 /* RFC 4034 */
	RR_DNSKEY = 48 /* RFC 4034 */

	RR_DHCID = 49 /* RFC 4701 */
	/* NSEC3 */
	RR_NSEC3      = 50 /* RFC 5155 */
	RR_NSEC3PARAM = 51 /* RFC 5155 */
	RR_TLSA       = 52 /* RFC 6698 */

	RR_HIP = 55 /* RFC 5205 */

	/** draft-reid-dnsext-zs */
	RR_NINFO = 56
	/** draft-reid-dnsext-rkey */
	RR_RKEY = 57
	/** draft-ietf-dnsop-trust-history */
	RR_TALINK = 58
	/** draft-barwood-dnsop-ds-publis */
	RR_CDS = 59

	RR_SPF = 99 /* RFC 4408 */

	RR_UINFO  = 100
	RR_UID    = 101
	RR_GID    = 102
	RR_UNSPEC = 103

	RR_NID = 104 /* RFC 6742 */
	RR_L32 = 105 /* RFC 6742 */
	RR_L64 = 106 /* RFC 6742 */
	RR_LP  = 107 /* RFC 6742 */

	RR_EUI48 = 108 /* RFC 7043 */
	RR_EUI64 = 109 /* RFC 7043 */

	RR_TKEY = 249 /* RFC 2930 */
	RR_TSIG = 250
	RR_IXFR = 251
	RR_AXFR = 252
	/**  A request for mailbox-related records (MB MG or MR) */
	RR_MAILB = 253
	/**  A request for mail agent RRs (Obsolete - see MX) */
	RR_MAILA = 254
	/**  any type (wildcard) */
	RR_ANY = 255
	/** draft-faltstrom-uri-06 */
	RR_URI = 256
	RR_CAA = 257 /* RFC 6844 */

	/** DNSSEC Trust Authorities */
	RR_TA = 32768
	/* RFC 4431 5074 DNSSEC Lookaside Validation */
	RR_DLV = 32769
)

var typeNameMap = map[RRType]string{
	RR_A:          "a",
	RR_NS:         "ns",
	RR_CNAME:      "cname",
	RR_SOA:        "soa",
	RR_MB:         "mb",
	RR_MG:         "mg",
	RR_MR:         "mr",
	RR_NULL:       "null",
	RR_WKS:        "wks",
	RR_PTR:        "ptr",
	RR_HINFO:      "hinfo",
	RR_MINFO:      "minfo",
	RR_MX:         "mx",
	RR_TXT:        "txt",
	RR_RP:         "rp",
	RR_AFSDB:      "afsdb",
	RR_X25:        "x25",
	RR_ISDN:       "isdn",
	RR_RT:         "rt",
	RR_NSAP:       "nsap",
	RR_NSAP_PTR:   "ptr",
	RR_SIG:        "sig",
	RR_KEY:        "key",
	RR_PX:         "px",
	RR_GPOS:       "gpos",
	RR_AAAA:       "aaaa",
	RR_LOC:        "loc",
	RR_NXT:        "nxt",
	RR_EID:        "eid",
	RR_NIMLOC:     "nimloc",
	RR_SRV:        "srv",
	RR_ATMA:       "atma",
	RR_NAPTR:      "naptr",
	RR_KX:         "kx",
	RR_CERT:       "cert",
	RR_A6:         "a6",
	RR_DNAME:      "dname",
	RR_SINK:       "sink",
	RR_OPT:        "opt",
	RR_APL:        "apl",
	RR_DS:         "ds",
	RR_SSHFP:      "sshfp",
	RR_IPSECKEY:   "ipseckey",
	RR_RRSIG:      "rrsig",
	RR_NSEC:       "nsec",
	RR_DNSKEY:     "dnskey",
	RR_DHCID:      "dhcid",
	RR_NSEC3:      "nsec3",
	RR_NSEC3PARAM: "nsec3param",
	RR_TLSA:       "tlsa",
	RR_HIP:        "hip",
	RR_NINFO:      "ninfo",
	RR_RKEY:       "pkey",
	RR_TALINK:     "talink",
	RR_CDS:        "cds",
	RR_SPF:        "spf",
	RR_UINFO:      "uinfo",
	RR_UID:        "uid",
	RR_GID:        "gid",
	RR_UNSPEC:     "unspec",
	RR_NID:        "nid",
	RR_L32:        "l32",
	RR_L64:        "l64",
	RR_LP:         "lp",
	RR_EUI48:      "eui48",
	RR_EUI64:      "eui64",

	RR_TKEY:  "tkey",
	RR_TSIG:  "tsig",
	RR_IXFR:  "ixfr",
	RR_AXFR:  "axfr",
	RR_MAILB: "mailb",
	RR_MAILA: "maila",
	RR_ANY:   "any",
	RR_URI:   "uri",
	RR_CAA:   "caa",
	RR_TA:    "ta",
	RR_DLV:   "dlv",
}

func TTLFromWire(buffer *util.InputBuffer) (RRTTL, error) {
	ttl, err := buffer.ReadUint32()
	if err != nil {
		return RRTTL(0), err
	}

	return RRTTL(ttl), nil
}

func TTLFromStr(s string) (RRTTL, error) {
	ttl, err := strconv.Atoi(s)
	if err != nil {
		return RRTTL(0), err
	}

	return RRTTL(ttl), nil
}

func (ttl RRTTL) Rend(render *MsgRender) {
	render.WriteUint32(uint32(ttl))
}

func (ttl RRTTL) ToWire(buffer *util.OutputBuffer) {
	buffer.WriteUint32(uint32(ttl))
}

func (ttl RRTTL) String() string {
	return strconv.Itoa(int(ttl))
}

func ClassFromWire(buffer *util.InputBuffer) (RRClass, error) {
	cls, err := buffer.ReadUint16()
	if err != nil {
		return RRClass(0), err
	}

	return RRClass(cls), nil
}

func ClassFromStr(s string) (RRClass, error) {
	s = strings.ToUpper(s)
	switch s {
	case "IN":
		return CLASS_IN, nil
	case "CH":
		return CLASS_CH, nil
	case "HS":
		return CLASS_HS, nil
	case "NONE":
		return CLASS_NONE, nil
	case "ANY":
		return CLASS_ANY, nil
	default:
		return RRClass(0), errors.New("unknownclass")
	}
}

func (cls RRClass) Rend(render *MsgRender) {
	render.WriteUint16(uint16(cls))
}

func (cls RRClass) ToWire(buffer *util.OutputBuffer) {
	buffer.WriteUint16(uint16(cls))
}

func (cls RRClass) String() string {
	switch cls {
	case CLASS_IN:
		return "IN"
	case CLASS_CH:
		return "CH"
	case CLASS_HS:
		return "HS"
	case CLASS_NONE:
		return "NONE"
	case CLASS_ANY:
		return "ANY"
	default:
		return "unknownclass"
	}
}

func TypeFromWire(buffer *util.InputBuffer) (RRType, error) {
	t, err := buffer.ReadUint16()
	if err != nil {
		return RRType(0), err
	}

	return RRType(t), nil
}

func TypeFromString(s string) (RRType, error) {
	s = strings.ToLower(s)
	for t, ts := range typeNameMap {
		if ts == s {
			return t, nil
		}
	}
	return RRType(0), errors.New("unknown rr type")
}

func (t RRType) Rend(render *MsgRender) {
	render.WriteUint16(uint16(t))
}

func (t RRType) ToWire(buffer *util.OutputBuffer) {
	buffer.WriteUint16(uint16(t))
}

func (t RRType) String() string {
	s := typeNameMap[t]
	if s == "" {
		return fmt.Sprintf("unknowntype:%d", t)
	} else {
		return strings.ToUpper(s)
	}
}

type RRset struct {
	Name   Name
	Type   RRType
	Class  RRClass
	Ttl    RRTTL
	Rdatas []Rdata
}

func RRsetFromWire(buffer *util.InputBuffer) (*RRset, error) {
	rrset := &RRset{}

	err := rrset.Name.FromWire(buffer, true)
	if err != nil {
		return nil, err
	}

	rrset.Type, err = TypeFromWire(buffer)
	if err != nil {
		return nil, err
	}

	rrset.Class, err = ClassFromWire(buffer)
	if err != nil {
		return nil, err
	}

	rrset.Ttl, err = TTLFromWire(buffer)
	if err != nil {
		return nil, err
	}

	rdata, err := RdataFromWire(rrset.Type, buffer)
	if err != nil {
		return nil, err
	}

	rrset.Rdatas = []Rdata{rdata}
	return rrset, nil
}

func (rrset *RRset) Rend(r *MsgRender) {
	for _, rdata := range rrset.Rdatas {
		rrset.Name.Rend(r)
		rrset.Type.Rend(r)
		rrset.Class.Rend(r)
		rrset.Ttl.Rend(r)
		pos := r.Len()
		r.Skip(2)
		rdata.Rend(r)
		r.WriteUint16At(uint16(r.Len()-pos-2), pos)
	}
}

func (rrset *RRset) ToWire(buffer *util.OutputBuffer) {
	for _, rdata := range rrset.Rdatas {
		rrset.Name.ToWire(buffer)
		rrset.Type.ToWire(buffer)
		rrset.Class.ToWire(buffer)
		rrset.Ttl.ToWire(buffer)

		pos := buffer.Len()
		buffer.Skip(2)
		rdata.ToWire(buffer)
		buffer.WriteUint16At(uint16(buffer.Len()-pos-2), pos)
	}
}

func (rrset *RRset) String() string {
	header := strings.Join([]string{rrset.Name.String(false), rrset.Ttl.String(), rrset.Class.String(), rrset.Type.String()}, "\t")
	var buf bytes.Buffer
	for _, rdata := range rrset.Rdatas {
		buf.WriteString(header)
		buf.WriteString("\t")
		buf.WriteString(rdata.String())
		buf.WriteString("\n")
	}
	return buf.String()
}

func (rrset *RRset) RrCount() int {
	return len(rrset.Rdatas)
}

func (rrset *RRset) IsSameRrset(other *RRset) bool {
	return (rrset.Type == other.Type) && rrset.Name.Equals(&other.Name)
}

func (rrset *RRset) AddRdata(rdata Rdata) {
	rrset.Rdatas = append(rrset.Rdatas, rdata)
}

func (rrset *RRset) RotateRdata() {
	rrCount := rrset.RrCount()
	if rrCount < 2 {
		return
	}

	rrset.Rdatas = append([]Rdata{rrset.Rdatas[rrCount-1]}, rrset.Rdatas[0:rrCount-1]...)
}

func (rrset *RRset) MarshalJSON() ([]byte, error) {
	rrs := []map[string]interface{}{}
	for _, rdata := range rrset.Rdatas {
		rrs = append(rrs, map[string]interface{}{
			"name":  rrset.Name.String(true),
			"type":  rrset.Type.String(),
			"class": rrset.Class.String(),
			"ttl":   rrset.Ttl,
			"rdata": rdata.String(),
		})
	}
	return json.Marshal(rrs)
}
