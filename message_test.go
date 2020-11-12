package g53

import (
	"testing"

	"github.com/ben-han-cn/g53/util"
)

func matchMessageRaw(t *testing.T, rawData string, m *Message) {
	wire, _ := util.HexStrToBytes(rawData)
	buf := util.NewInputBuffer(wire)
	nm, err := MessageFromWire(buf)
	Assert(t, err == nil, "err should be nil")

	Equal(t, nm.Header, m.Header)
	matchQuestion(t, nm.Question, m.Question)
	matchSection(t, nm.GetSection(AnswerSection), m.GetSection(AnswerSection))
	matchSection(t, nm.GetSection(AuthSection), m.GetSection(AuthSection))
	matchSection(t, nm.GetSection(AdditionalSection), m.GetSection(AdditionalSection))

	render := NewMsgRender()
	nm.Rend(render)
	WireMatch(t, wire, render.Data())
}

func matchSection(t *testing.T, ns Section, s Section) {
	Equal(t, len(ns), len(s))
	for i := 0; i < len(ns); i++ {
		matchRRset(t, ns[i], s[i])
	}
}

func TestSimpleMessageFromToWire(t *testing.T) {
	qn, _ := NameFromString("test.example.com.")
	ra1, _ := AFromString("192.0.2.2")
	ra2, _ := AFromString("192.0.2.1")
	ns, _ := NameFromString("example.com.")
	ra3, _ := NSFromString("ns1.example.com.")
	glue, _ := NameFromString("ns1.example.com.")
	ra4, _ := AFromString("2.2.2.2")
	question := &Question{
		Name:  *qn,
		Type:  RR_A,
		Class: CLASS_IN,
	}
	edns := &EDNS{
		UdpSize:     uint16(4096),
		DnssecAware: false,
	}

	msg := NewMsgBuilder(&Message{}).
		SetId(1200).
		SetHeaderFlag(FLAG_QR, true).
		SetHeaderFlag(FLAG_AA, true).
		SetHeaderFlag(FLAG_RD, true).
		SetOpcode(OP_QUERY).
		SetRcode(R_NOERROR).
		SetQuestion(question).
		AddRRset(AnswerSection, &RRset{
			Name:   *qn,
			Type:   RR_A,
			Class:  CLASS_IN,
			Ttl:    RRTTL(3600),
			Rdatas: []Rdata{ra1, ra2},
		}).
		AddRRset(AuthSection, &RRset{
			Name:   *ns,
			Type:   RR_NS,
			Class:  CLASS_IN,
			Ttl:    RRTTL(3600),
			Rdatas: []Rdata{ra3},
		}).
		AddRRset(AdditionalSection, &RRset{
			Name:   *glue,
			Type:   RR_A,
			Class:  CLASS_IN,
			Ttl:    RRTTL(3600),
			Rdatas: []Rdata{ra4},
		}).
		SetEdns(edns).
		Done()

	matchMessageRaw(t, "04b0850000010002000100020474657374076578616d706c6503636f6d0000010001c00c0001000100000e100004c0000202c00c0001000100000e100004c0000201c0110002000100000e100006036e7331c011c04e0001000100000e100004020202020000291000000000000000", msg)
}

func TestCompliateMessageFromToWire(t *testing.T) {
	knet_cn := "04b08180000100010004000d03777777046b6e657402636e0000010001c00c00010001000002580004caad0b0ac01000020001000000c1001404676e7331097a646e73636c6f7564036e657400c01000020001000000c10014046c6e7332097a646e73636c6f75640362697a00c01000020001000000c1001504676e7332097a646e73636c6f7564036e6574c015c01000020001000000c10015046c6e7331097a646e73636c6f756404696e666f00c039000100010000262c000401089801c0790001000100000599000401089901c09a00010001000007c800046f012189c09a00010001000007c8000477a7e9e9c09a00010001000007c80004b683170bc09a00010001000007c80004010865fdc09a001c0001000007c8001024018d00000400000000000000000001c0590001000100002fea000477a7e9ebc0590001000100002fea0004b683170cc0590001000100002fea0004010865fcc0590001000100002fea00046f01218ac059001c00010000249f001024018d000006000000000000000000010000291000000000000000"
	/*
		;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 1200
		;; flags:  qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 4, ADDITIONAL: 13,

		;; OPT PSEUDOSECTION:
		; EDNS: version: 0, udp: 4096
		;; QUESTION SECTION:
		www.knet.cn. IN A

		;; ANSWER SECTION:
		www.knet.cn.	600	IN	A	202.173.11.10

		;; AUTHORITY SECTION:
		knet.cn.	193	IN	NS	gns1.zdnscloud.net.
		knet.cn.	193	IN	NS	lns2.zdnscloud.biz.
		knet.cn.	193	IN	NS	gns2.zdnscloud.net.cn.
		knet.cn.	193	IN	NS	lns1.zdnscloud.info.

		;; ADDITIONAL SECTION:
		gns1.zdnscloud.net.	9772	IN	A	1.8.152.1
		gns2.zdnscloud.net.cn.	1433	IN	A	1.8.153.1

		lns1.zdnscloud.info.	1992	IN	A	111.1.33.137
		lns1.zdnscloud.info.	1992	IN	A	119.167.233.233
		lns1.zdnscloud.info.	1992	IN	A	182.131.23.11
		lns1.zdnscloud.info.	1992	IN	A	1.8.101.253
		lns1.zdnscloud.info.	1992	IN	AAAA	2401:8d00:4::1

		lns2.zdnscloud.biz.	12266	IN	A	119.167.233.235
		lns2.zdnscloud.biz.	12266	IN	A	182.131.23.12
		lns2.zdnscloud.biz.	12266	IN	A	1.8.101.252
		lns2.zdnscloud.biz.	12266	IN	A	111.1.33.138
		lns2.zdnscloud.biz.	9375	IN	AAAA	2401:8d00:6::1
	*/

	qn, _ := NameFromString("www.knet.cn.")
	ra1, _ := AFromString("202.173.11.10")
	question := &Question{
		Name:  *qn,
		Type:  RR_A,
		Class: CLASS_IN,
	}
	nsName, _ := NameFromString("knet.cn.")
	ns1, _ := NSFromString("gns1.zdnscloud.net.")
	ns2, _ := NSFromString("lns2.zdnscloud.biz.")
	ns3, _ := NSFromString("gns2.zdnscloud.net.cn.")
	ns4, _ := NSFromString("lns1.zdnscloud.info.")
	glueName1, _ := NameFromString("gns1.zdnscloud.net.")
	glue1, _ := AFromString("1.8.152.1")
	glueName2, _ := NameFromString("gns2.zdnscloud.net.cn.")
	glue2, _ := AFromString("1.8.153.1")
	glueName3, _ := NameFromString("lns1.zdnscloud.info.")
	glue31, _ := AFromString("111.1.33.137")
	glue32, _ := AFromString("119.167.233.233")
	glue33, _ := AFromString("182.131.23.11")
	glue34, _ := AFromString("1.8.101.253")
	glue35, _ := AAAAFromString("2401:8d00:4::1")
	glueName4, _ := NameFromString("lns2.zdnscloud.biz.")
	glue41, _ := AFromString("119.167.233.235")
	glue42, _ := AFromString("182.131.23.12")
	glue43, _ := AFromString("1.8.101.252")
	glue44, _ := AFromString("111.1.33.138")
	glue45, _ := AAAAFromString("2401:8d00:6::1")
	edns := &EDNS{
		UdpSize:     uint16(4096),
		DnssecAware: false,
	}

	msg := NewMsgBuilder(&Message{}).
		SetId(1200).
		SetHeaderFlag(FLAG_QR, true).
		SetHeaderFlag(FLAG_RD, true).
		SetHeaderFlag(FLAG_RA, true).
		SetOpcode(OP_QUERY).
		SetRcode(R_NOERROR).
		SetQuestion(question).
		AddRRset(AnswerSection, &RRset{
			Name:   *qn,
			Type:   RR_A,
			Class:  CLASS_IN,
			Ttl:    RRTTL(600),
			Rdatas: []Rdata{ra1},
		}).
		AddRRset(AuthSection, &RRset{
			Name:   *nsName,
			Type:   RR_NS,
			Class:  CLASS_IN,
			Ttl:    RRTTL(193),
			Rdatas: []Rdata{ns1, ns2, ns3, ns4},
		}).
		ResizeSection(AdditionalSection, 6).
		AddRRset(AdditionalSection, &RRset{
			Name:   *glueName1,
			Type:   RR_A,
			Class:  CLASS_IN,
			Ttl:    RRTTL(9772),
			Rdatas: []Rdata{glue1},
		}).
		AddRRset(AdditionalSection, &RRset{
			Name:   *glueName2,
			Type:   RR_A,
			Class:  CLASS_IN,
			Ttl:    RRTTL(1433),
			Rdatas: []Rdata{glue2},
		}).
		AddRRset(AdditionalSection, &RRset{
			Name:   *glueName3,
			Type:   RR_A,
			Class:  CLASS_IN,
			Ttl:    RRTTL(1992),
			Rdatas: []Rdata{glue31, glue32, glue33, glue34},
		}).
		AddRRset(AdditionalSection, &RRset{
			Name:   *glueName3,
			Type:   RR_AAAA,
			Class:  CLASS_IN,
			Ttl:    RRTTL(1992),
			Rdatas: []Rdata{glue35},
		}).
		AddRRset(AdditionalSection, &RRset{
			Name:   *glueName4,
			Type:   RR_A,
			Class:  CLASS_IN,
			Ttl:    RRTTL(12266),
			Rdatas: []Rdata{glue41, glue42, glue43, glue44},
		}).
		AddRRset(AdditionalSection, &RRset{
			Name:   *glueName4,
			Type:   RR_AAAA,
			Class:  CLASS_IN,
			Ttl:    RRTTL(9375),
			Rdatas: []Rdata{glue45},
		}).
		SetEdns(edns).
		Done()

	matchMessageRaw(t, knet_cn, msg)
}

func TestMessageAllocate(t *testing.T) {
	resp := []byte{
		4, 176, 133, 0, 0, 1, 0, 4, 0, 0, 0, 7, 3, 105, 115, 99, 3, 111, 114, 103, 0, 0, 2, 0, 1, 192, 12, 0, 2, 0, 1, 0, 0, 28, 32, 0, 25, 2, 110, 115, 3, 105, 115, 99, 11, 97, 102, 105, 108, 105, 97, 115, 45, 110, 115, 116, 4, 105, 110, 102, 111, 0, 192, 12, 0, 2, 0, 1, 0, 0, 28, 32, 0, 13, 3, 111, 114, 100, 6, 115, 110, 115, 45, 112, 98, 192, 12, 192, 12, 0, 2, 0, 1, 0, 0, 28, 32, 0, 6, 3, 97, 109, 115, 192, 78, 192, 12, 0, 2, 0, 1, 0, 0, 28, 32, 0, 7, 4, 115, 102, 98, 97, 192, 78, 192, 99, 0, 1, 0, 1, 0, 0, 28, 32, 0, 4, 199, 6, 1, 30, 192, 99, 0, 28, 0, 1, 0, 0, 28, 32, 0, 16, 32, 1, 5, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 192, 74, 0, 1, 0, 1, 0, 0, 28, 32, 0, 4, 199, 6, 0, 30, 192, 74, 0, 28, 0, 1, 0, 0, 28, 32, 0, 16, 32, 1, 5, 0, 0, 113, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 192, 117, 0, 1, 0, 1, 0, 0, 28, 32, 0, 4, 149, 20, 64, 3, 192, 117, 0, 28, 0, 1, 0, 0, 28, 32, 0, 16, 32, 1, 4, 248, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 25, 0, 0, 41, 16, 0, 0, 0, 0, 0, 0, 0,
	}
	msg_, _ := MessageFromWire(util.NewInputBuffer(resp))
	str := msg_.String()

	var msg Message
	for i := 0; i < 20; i++ {
		msg.Clear()
		err := msg.FromWire(util.NewInputBuffer(resp))
		Assert(t, err == nil, "")
		Assert(t, msg.String() == str, "")
	}

	req := []byte{
		4, 176, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 3, 105, 115, 99, 3, 111, 114, 103, 0, 0, 2, 0, 1, 0, 0, 41, 16, 0, 0, 0, 0, 0, 0, 0}
	msg.FromWire(util.NewInputBuffer(req))
	allocs := testing.AllocsPerRun(10, func() {
		msg.Clear()
		msg.FromWire(util.NewInputBuffer(req))
	})
	Assert(t, allocs == 3, "allocate %v", allocs)
}

func benchmarkParseMessage(b *testing.B, raw string) {
	wire, _ := util.HexStrToBytes(raw)
	buf := util.NewInputBuffer(wire)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MessageFromWire(buf)
		buf.SetPosition(0)
	}
}

func BenchmarkParseKnetMessage(b *testing.B) {
	benchmarkParseMessage(b, "04b08180000100010004000d03777777046b6e657402636e0000010001c00c00010001000002580004caad0b0ac01000020001000000c1001404676e7331097a646e73636c6f7564036e657400c01000020001000000c10014046c6e7332097a646e73636c6f75640362697a00c01000020001000000c1001504676e7332097a646e73636c6f7564036e6574c015c01000020001000000c10015046c6e7331097a646e73636c6f756404696e666f00c039000100010000262c000401089801c0790001000100000599000401089901c09a00010001000007c800046f012189c09a00010001000007c8000477a7e9e9c09a00010001000007c80004b683170bc09a00010001000007c80004010865fdc09a001c0001000007c8001024018d00000400000000000000000001c0590001000100002fea000477a7e9ebc0590001000100002fea0004b683170cc0590001000100002fea0004010865fcc0590001000100002fea00046f01218ac059001c00010000249f001024018d000006000000000000000000010000291000000000000000")
}

func BenchmarkParseTestExample(b *testing.B) {
	benchmarkParseMessage(b, "04b0850000010002000100020474657374076578616d706c6503636f6d0000010001c00c0001000100000e100004c0000202c00c0001000100000e100004c0000201c0110002000100000e100006036e7331c011c04e0001000100000e100004020202020000291000000000000000")
}
