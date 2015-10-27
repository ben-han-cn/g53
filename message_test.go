package g53

import (
	"g53/util"
	"testing"
)

func buildHeader(id uint16, setFlag []FlagField, counts []uint16, opcode Opcode, rcode Rcode) *Header {
	h := &Header{
		Id:      id,
		Opcode:  opcode,
		Rcode:   rcode,
		QDCount: counts[0],
		ANCount: counts[1],
		NSCount: counts[2],
		ARCount: counts[3],
	}

	for _, f := range setFlag {
		h.SetFlag(f, true)
	}

	return h
}

func matchMessageRaw(t *testing.T, rawData string, m *Message) {
	wire, _ := util.HexStrToBytes(rawData)
	buffer := util.NewInputBuffer(wire)
	nm, err := MessageFromWire(buffer)
	Assert(t, err == nil, "err should be nil")

	Equal(t, *(nm.Header), *(m.Header))
	matchQuestion(t, nm.Question, m.Question)
	matchSection(t, nm.GetSection(AnswerSection), m.GetSection(AnswerSection))
	matchSection(t, nm.GetSection(AuthSection), m.GetSection(AuthSection))
	matchSection(t, nm.GetSection(AdditionalSection), m.GetSection(AdditionalSection))

	render := NewMsgRender()
	nm.Mode = RENDER
	nm.Rend(render)
	WireMatch(t, wire, render.Data())
}

func matchSection(t *testing.T, ns Section, s Section) {
	Equal(t, len(ns), len(s))
	for i := 0; i < len(ns); i++ {
		matchRRset(t, ns[i], s[i])
	}
}

func TestMessageFromToWire(t *testing.T) {
	qn, _ := NameFromString("test.example.com.")
	ra1, _ := AFromString("192.0.2.1")
	ra2, _ := AFromString("192.0.2.2")

	var answer Section
	answer = append(answer, &RRset{
		Name:   qn,
		Type:   RR_A,
		Class:  CLASS_IN,
		Ttl:    RRTTL(3600),
		Rdatas: []Rdata{ra1},
	})

	answer = append(answer, &RRset{
		Name:   qn,
		Type:   RR_A,
		Class:  CLASS_IN,
		Ttl:    RRTTL(7200),
		Rdatas: []Rdata{ra2},
	})

	matchMessageRaw(t, "1035850000010002000000000474657374076578616d706c6503636f6d0000010001c00c0001000100000e100004c0000201c00c0001000100001c200004c0000202", &Message{
		Header: buildHeader(uint16(0x1035), []FlagField{FLAG_QR, FLAG_AA, FLAG_RD}, []uint16{1, 2, 0, 0}, OP_QUERY, R_NOERROR),
		Question: &Question{
			Name:  qn,
			Type:  RR_A,
			Class: CLASS_IN,
		},
		Sections: [...]Section{answer, []*RRset{}, []*RRset{}},
	})
}
