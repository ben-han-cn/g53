package g53

import (
	"bytes"
	"errors"

	"g53/util"
)

var ErrQueryQuestionIsNotValid = errors.New("query should have exact one question")

type SectionType int

const (
	AnswerSection     SectionType = 0
	AuthSection                   = 1
	AdditionalSection             = 2
)

const SectionCount = 3

type Section []*RRset

func (section Section) rrCount() int {
	count := 0
	for _, rrset := range section {
		count += rrset.RrCount()
	}
	return count
}

type Message struct {
	Header   Header
	Question *Question
	Sections [SectionCount]Section
	Edns     *EDNS
}

func MakeQuery(name *Name, typ RRType, msgSize int, dnssec bool) *Message {
	h := Header{}
	h.SetFlag(FLAG_RD, true)
	h.Opcode = OP_QUERY

	q := &Question{
		Name:  name,
		Type:  typ,
		Class: CLASS_IN,
	}

	return &Message{
		Header:   h,
		Question: q,
		Edns: &EDNS{
			UdpSize:     uint16(msgSize),
			DnssecAware: dnssec,
		},
	}
}

func MessageFromWire(buffer *util.InputBuffer) (*Message, error) {
	m := Message{}
	h := &m.Header
	err := HeaderFromWire(h, buffer)
	if err != nil {
		return nil, err
	}

	if h.QDCount == 1 {
		q, err := QuestionFromWire(buffer)
		if err != nil {
			return nil, err
		}
		m.Question = q
	} else if h.Opcode == OP_QUERY {
		return nil, ErrQueryQuestionIsNotValid
	}

	for i := 0; i < SectionCount; i++ {
		if err := m.sectionFromWire(SectionType(i), buffer); err != nil {
			return nil, err
		}
	}

	return &m, nil
}

func (m *Message) sectionFromWire(st SectionType, buffer *util.InputBuffer) error {
	var s Section
	var count uint16

	switch st {
	case AnswerSection:
		count = m.Header.ANCount
	case AuthSection:
		count = m.Header.NSCount
	case AdditionalSection:
		count = m.Header.ARCount
	}

	var lastRrset *RRset
	for i := uint16(0); i < count; i++ {
		rrset, err := RRsetFromWire(buffer)
		if err != nil {
			return err
		}

		if lastRrset == nil {
			lastRrset = rrset
			continue
		}

		if lastRrset.IsSameRrset(rrset) {
			lastRrset.Rdatas = append(lastRrset.Rdatas, rrset.Rdatas[0])
		} else {
			if st == AdditionalSection && lastRrset.Type == RR_OPT {
				m.Edns = EdnsFromRRset(lastRrset)
			} else {
				s = append(s, lastRrset)
			}
			lastRrset = rrset
		}
	}

	if lastRrset != nil {
		if st == AdditionalSection && lastRrset.Type == RR_OPT {
			m.Edns = EdnsFromRRset(lastRrset)
		} else {
			s = append(s, lastRrset)
		}
	}

	m.Sections[st] = s
	return nil
}

func (m *Message) Rend(r *MsgRender) {
	m.recalculateSectionRrCount()
	(&m.Header).Rend(r)

	if m.Question != nil {
		m.Question.Rend(r)
	}

	for i := 0; i < SectionCount; i++ {
		m.Sections[i].Rend(r)
	}

	if m.Edns != nil {
		m.Edns.Rend(r)
	}
}

func (m *Message) recalculateSectionRrCount() {
	if m.Question == nil {
		m.Header.QDCount = 0
	} else {
		m.Header.QDCount = 1
	}

	m.Header.ANCount = uint16(m.Sections[AnswerSection].rrCount())
	m.Header.NSCount = uint16(m.Sections[AuthSection].rrCount())
	m.Header.ARCount = uint16(m.Sections[AdditionalSection].rrCount())

	if m.Edns != nil {
		m.Header.ARCount += 1
	}
}

func (s Section) Rend(r *MsgRender) {
	for _, rrset := range s {
		rrset.Rend(r)
	}
}

func (m *Message) ToWire(buffer *util.OutputBuffer) {
	(&m.Header).ToWire(buffer)
	if m.Question != nil {
		m.Question.ToWire(buffer)
	}

	for i := 0; i < SectionCount; i++ {
		m.Sections[i].ToWire(buffer)
	}
}

func (s Section) ToWire(buffer *util.OutputBuffer) {
	for _, rrset := range s {
		rrset.ToWire(buffer)
	}
}

func (m *Message) String() string {
	m.recalculateSectionRrCount()

	var buf bytes.Buffer
	buf.WriteString(m.Header.String())
	buf.WriteString("\n")

	if m.Edns != nil {
		buf.WriteString(";; OPT PSEUDOSECTION:\n")
		buf.WriteString(m.Edns.String())
	}

	buf.WriteString(";; QUESTION SECTION:\n")
	if m.Question != nil {
		buf.WriteString(m.Question.String())
		buf.WriteString("\n")
	}

	if len(m.Sections[AnswerSection]) > 0 {
		buf.WriteString("\n;; ANSWER SECTION:\n")
		buf.WriteString(m.Sections[AnswerSection].String())
	}

	if len(m.Sections[AuthSection]) > 0 {
		buf.WriteString("\n;; AUTHORITY SECTION:\n")
		buf.WriteString(m.Sections[AuthSection].String())
	}

	if len(m.Sections[AdditionalSection]) > 0 {
		buf.WriteString("\n;; ADDITIONAL SECTION:\n")
		buf.WriteString(m.Sections[AdditionalSection].String())
	}
	return buf.String()
}

func (s Section) String() string {
	var buf bytes.Buffer
	for _, rrset := range s {
		buf.WriteString(rrset.String())
	}
	return buf.String()
}

func (m *Message) GetSection(st SectionType) Section {
	return m.Sections[st]
}

func (m *Message) Clear() {
	m.Header.Clear()
	m.Question = nil
	for i := 0; i < SectionCount; i++ {
		m.Sections[i] = nil
	}
}

func (m *Message) AddRRset(st SectionType, rrset *RRset) {
	m.Sections[st] = append(m.Sections[st], rrset)
}

func (m *Message) AddRr(st SectionType, name *Name, typ RRType, class RRClass, ttl RRTTL, rdata Rdata, merge bool) {
	if merge {
		if i := m.rrsetIndex(st, name, typ, class); i != -1 {
			m.Sections[st][i].AddRdata(rdata)
			m.Sections[st][i].Ttl = ttl
			return
		}
	}
	newRRset := &RRset{
		Name:   name,
		Type:   typ,
		Class:  class,
		Ttl:    ttl,
		Rdatas: []Rdata{rdata},
	}
	m.AddRRset(st, newRRset)
}

func (m *Message) HasRRset(st SectionType, rrset *RRset) bool {
	return m.rrsetIndex(st, rrset.Name, rrset.Type, rrset.Class) != -1
}

func (m *Message) rrsetIndex(st SectionType, name *Name, typ RRType, class RRClass) int {
	for i, rrset := range m.Sections[st] {
		if rrset.Class == class &&
			rrset.Type == typ &&
			rrset.Name.Equals(name) {
			return i
		}
	}
	return -1
}

func (m *Message) HasRRsetWithNameType(st SectionType, n *Name, t RRType) bool {
	return false
}

func (m *Message) MakeResponse() *Message {
	h := Header{
		Id:      m.Header.Id,
		Opcode:  OP_QUERY,
		QDCount: m.Header.QDCount,
	}

	h.SetFlag(FLAG_QR, true)
	h.SetFlag(FLAG_RD, m.Header.GetFlag(FLAG_RD))

	return &Message{
		Header:   h,
		Question: m.Question,
	}
}

func (m *Message) ClearSection(s SectionType) {
	m.Sections[s] = nil
	switch s {
	case AnswerSection:
		m.Header.ANCount = 0
	case AuthSection:
		m.Header.NSCount = 0
	case AdditionalSection:
		m.Edns = nil
		m.Header.ARCount = 0
	default:
		panic("question section couldn't be cleared")
	}
}

func (m *Message) SectionRrCount(s SectionType) int {
	return m.Sections[s].rrCount()
}

func (m *Message) SectionRRsetCount(s SectionType) int {
	return len(m.Sections[s])
}
