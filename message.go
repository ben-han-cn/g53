package g53

import (
	"bytes"
	"fmt"

	"github.com/ben-han-cn/g53/util"
)

type SectionType int

const (
	AnswerSection     SectionType = 0
	AuthSection       SectionType = 1
	AdditionalSection SectionType = 2
)

const SectionCount = 3

type Section []*RRset

func (s Section) rrCount() int {
	count := 0
	for i := 0; i < len(s); i++ {
		rrCount := s[i].RRCount()
		//for empty rdata, just count as 1
		if rrCount == 0 {
			rrCount = 1
		}
		count += rrCount
	}
	return count
}

func (s Section) Rend(r *MsgRender) {
	for i := 0; i < len(s); i++ {
		s[i].Rend(r)
	}
}

func (s Section) ToWire(buf *util.OutputBuffer) {
	for i := 0; i < len(s); i++ {
		s[i].ToWire(buf)
	}
}

func (s Section) String() string {
	var buf bytes.Buffer
	for i := 0; i < len(s); i++ {
		buf.WriteString(s[i].String())
	}
	return buf.String()
}

func (s Section) Clear() Section {
	for i := range s {
		s[i] = nil
	}
	return s[:0]
}

type Message struct {
	noCopy

	Header   Header
	Question *Question
	question Question
	sections [SectionCount]Section
	edns     EDNS
	Edns     *EDNS
	Tsig     *Tsig
}

func MessageFromWire(buf *util.InputBuffer) (*Message, error) {
	m := Message{}
	if err := m.FromWire(buf); err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

func (m *Message) FromWire(buf *util.InputBuffer) error {
	if err := m.Header.FromWire(buf); err != nil {
		return err
	}

	if m.Header.QDCount == 1 {
		if err := m.question.FromWire(buf); err != nil {
			return err
		}
		m.Question = &m.question
	} else {
		m.Question = nil //in axfr message, question could be nil
	}

	m.Edns = nil
	m.Tsig = nil
	for i := 0; i < SectionCount; i++ {
		if err := m.sectionFromWire(SectionType(i), buf); err != nil {
			return err
		}
	}

	return nil
}

func (m *Message) sectionFromWire(st SectionType, buf *util.InputBuffer) error {
	var s Section
	var count uint16
	switch st {
	case AnswerSection:
		count = m.Header.ANCount
		s = m.sections[0]
	case AuthSection:
		count = m.Header.NSCount
		s = m.sections[1]
	case AdditionalSection:
		count = m.Header.ARCount
		s = m.sections[2]
	}

	if count == 0 {
		return nil
	}

	var (
		lastRRset *RRset
		opt       *RRset
		tsig      *RRset
	)
	for i := uint16(0); i < count; i++ {
		var rrset RRset
		if err := rrset.FromWire(buf); err != nil {
			return err
		}

		if rrset.Type == RR_OPT {
			if st != AdditionalSection {
				return fmt.Errorf("opt rr exist in non-addtional section")
			}
			if opt != nil {
				return fmt.Errorf("opt can only has one rr")
			}
			opt = &rrset
			continue
		}

		if rrset.Type == RR_TSIG {
			if st != AdditionalSection {
				return fmt.Errorf("tsig rr exist in non-addtional section")
			} else if i != count-1 {
				return fmt.Errorf("tsig rr isn't the last rr")
			} else if tsig != nil {
				return fmt.Errorf("tsig should has only one rr")
			}
			tsig = &rrset
			continue
		}

		if lastRRset == nil {
			lastRRset = &rrset
			continue
		}

		if lastRRset.IsSameRRset(&rrset) {
			if rrset.Type == RR_OPT || rrset.Type == RR_TSIG {
				return fmt.Errorf("opt or tsig should has only one rdata")
			}
			if len(rrset.Rdatas) == 0 {
				return fmt.Errorf("duplicate rrset with empty rdata")
			}
			lastRRset.Rdatas = append(lastRRset.Rdatas, rrset.Rdatas[0])
		} else {
			s = append(s, lastRRset)
			lastRRset = &rrset
		}
	}

	if lastRRset != nil {
		s = append(s, lastRRset)
	}

	if opt != nil {
		if err := m.edns.FromRRset(opt); err != nil {
			return err
		}
		m.Edns = &m.edns
	} else {
		m.Edns = nil
	}

	if tsig != nil {
		if tsig, err := TsigFromRRset(tsig); err != nil {
			return err
		} else {
			m.Tsig = tsig
		}
	} else {
		m.Tsig = nil
	}

	m.sections[st] = s
	return nil
}

func (m *Message) Rend(r *MsgRender) {
	if m.Tsig != nil {
		m.Header.ARCount -= 1
	}

	(&m.Header).Rend(r)

	if m.Question != nil {
		m.Question.Rend(r)
	}

	for i := 0; i < SectionCount; i++ {
		m.sections[i].Rend(r)
	}

	if m.Edns != nil {
		m.Edns.Rend(r)
	}

	if m.Tsig != nil {
		m.Tsig.Rend(r)
		m.Header.ARCount += 1
		r.WriteUint16At(uint16(m.Header.ARCount), 10)
	}
}

func (m *Message) ToWire(buf *util.OutputBuffer) {
	(&m.Header).ToWire(buf)
	if m.Question != nil {
		m.Question.ToWire(buf)
	}

	for i := 0; i < SectionCount; i++ {
		m.sections[i].ToWire(buf)
	}
}

func (m *Message) String() string {
	var buf bytes.Buffer
	buf.WriteString(m.Header.String())
	buf.WriteByte('\n')

	if m.Edns != nil {
		buf.WriteString(";; OPT PSEUDOSECTION:\n")
		buf.WriteString(m.Edns.String())
	}

	buf.WriteString(";; QUESTION SECTION:\n")
	if m.Question != nil {
		buf.WriteString(m.Question.String())
		buf.WriteByte('\n')
	}

	if len(m.sections[AnswerSection]) > 0 {
		buf.WriteString("\n;; ANSWER SECTION:\n")
		buf.WriteString(m.sections[AnswerSection].String())
	}

	if len(m.sections[AuthSection]) > 0 {
		buf.WriteString("\n;; AUTHORITY SECTION:\n")
		buf.WriteString(m.sections[AuthSection].String())
	}

	if len(m.sections[AdditionalSection]) > 0 {
		buf.WriteString("\n;; ADDITIONAL SECTION:\n")
		buf.WriteString(m.sections[AdditionalSection].String())
	}

	if m.Tsig != nil {
		buf.WriteString("\n;; Tsig PSEUDOSECTION:\n")
		buf.WriteString(m.Tsig.String())
	}

	return buf.String()
}

func (m *Message) GetSection(st SectionType) Section {
	return m.sections[st]
}

func (m *Message) Clear() {
	m.Header.Clear()
	m.Question = nil
	//this will reuse the backend array, this may cause
	//memory leak if there is a big section but after that
	//the section has very few rrset
	for i := 0; i < SectionCount; i++ {
		m.sections[i] = m.sections[i].Clear()
	}
	m.Edns = nil
	m.Tsig = nil
}

func (m *Message) HasRRset(st SectionType, rrset *RRset) bool {
	return m.rrsetIndex(st, &rrset.Name, rrset.Type, rrset.Class) != -1
}

func (m *Message) SectionRRCount(st SectionType) int {
	c := m.sections[st].rrCount()
	if st == AdditionalSection {
		if m.Edns != nil {
			c += 1
		}
		if m.Tsig != nil {
			c += 1
		}
	}
	return c
}

func (m *Message) SectionRRsetCount(st SectionType) int {
	c := len(m.sections[st])
	if st == AdditionalSection {
		if m.Edns != nil {
			c += 1
		}
		if m.Tsig != nil {
			c += 1
		}
	}
	return c
}
