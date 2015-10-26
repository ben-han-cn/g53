package g53

import (
	"bytes"
	"g53/util"
)

type SectionType int

const (
	AnswerSection     SectionType = 0
	AuthSection                   = 1
	AdditionalSection             = 2
)

const SectionCount = 3

type Section []*RRset

type Message struct {
	Header   *Header
	Question *Question
	Sections [SectionCount]Section
}

func MessageFromWire(buffer *util.InputBuffer) (*Message, error) {
	h, err := HeaderFromWire(buffer)
	if err != nil {
		return nil, err
	}

	var q *Question
	if h.QDCount == 1 {
		q, err = QuestionFromWire(buffer)
		if err != nil {
			return nil, err
		}
	}

	answer, err := sectionFromWire(h.ANCount, buffer)
	if err != nil {
		return nil, err
	}

	authority, err := sectionFromWire(h.NSCount, buffer)
	if err != nil {
		return nil, err
	}

	additional, err := sectionFromWire(h.ARCount, buffer)
	if err != nil {
		return nil, err
	}

	sections := [SectionCount]Section{answer, authority, additional}
	return &Message{
		Header:   h,
		Question: q,
		Sections: sections,
	}, nil
}

func sectionFromWire(count uint16, buffer *util.InputBuffer) (Section, error) {
	var s Section
	for i := uint16(0); i < count; i++ {
		rrset, err := RRsetFromWire(buffer)
		if err != nil {
			return s, err
		}
		s = append(s, rrset)
	}
	return s, nil
}

func (m *Message) Rend(r *MsgRender) {
	m.Header.Rend(r)
	if m.Question != nil {
		m.Question.Rend(r)
	}

	for i := 0; i < SectionCount; i++ {
		m.Sections[i].Rend(r)
	}
}

func (s Section) Rend(r *MsgRender) {
	for _, rrset := range s {
		rrset.Rend(r)
	}
}

func (m *Message) ToWire(buffer *util.OutputBuffer) {
	m.Header.ToWire(buffer)
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
	var buf bytes.Buffer
	buf.WriteString(m.Header.String())
	buf.WriteString("\n")

	buf.WriteString(";; QUESTION SECTION:\n")
	if m.Question != nil {
		buf.WriteString(m.Question.String())
	}

	buf.WriteString("\n;; ANSWER SECTION:\n")
	buf.WriteString(m.Sections[AnswerSection].String())

	buf.WriteString("\n;; AUTHORITY SECTION:\n")
	buf.WriteString(m.Sections[AuthSection].String())

	buf.WriteString("\n;; ADDITIONAL SECTION:\n")
	buf.WriteString(m.Sections[AdditionalSection].String())
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
