package g53

import (
	"bytes"
	"errors"
	"fmt"

	"g53/util"
)

type HeaderFlag uint16
type FlagField uint16

const (
	FLAG_QR FlagField = 0x8000
	FLAG_AA           = 0x0400
	FLAG_TC           = 0x0200
	FLAG_RD           = 0x0100
	FLAG_RA           = 0x0080
	FLAG_AD           = 0x0020
	FLAG_CD           = 0x0010
)

const (
	HEADERFLAG_MASK uint16 = 0x87b0
	OPCODE_MASK            = 0x7800
	OPCODE_SHIFT           = 11
	RCODE_MASK             = 0x000f
)

type Header struct {
	Id      uint16
	Flag    HeaderFlag
	Opcode  Opcode
	Rcode   Rcode
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func (h *Header) Clear() {
	h.Flag = 0
	h.QDCount = 0
	h.ANCount = 0
	h.NSCount = 0
	h.ARCount = 0
}

func (h *Header) GetFlag(ff FlagField) bool {
	return (uint16(h.Flag) & uint16(ff)) != 0
}

func (h *Header) SetFlag(ff FlagField, set bool) {
	if set {
		h.Flag = HeaderFlag(uint16(h.Flag) | uint16(ff))
	} else {
		h.Flag = HeaderFlag(uint16(h.Flag) & uint16(^ff))
	}
}

func HeaderFromWire(buffer *util.InputBuffer) (*Header, error) {
	if buffer.Len() < 12 {
		return nil, errors.New("too short wire data for message header")
	}
	id, _ := buffer.ReadUint16()
	flag, _ := buffer.ReadUint16()
	qdcount, _ := buffer.ReadUint16()
	ancount, _ := buffer.ReadUint16()
	nscount, _ := buffer.ReadUint16()
	arcount, _ := buffer.ReadUint16()
	return &Header{
		Id:      id,
		Flag:    HeaderFlag(flag & HEADERFLAG_MASK),
		Opcode:  Opcode((flag & OPCODE_MASK) >> OPCODE_SHIFT),
		Rcode:   Rcode(flag & RCODE_MASK),
		QDCount: qdcount,
		ANCount: ancount,
		NSCount: nscount,
		ARCount: arcount,
	}, nil
}

func (h *Header) Rend(r *MsgRender) {
	r.WriteUint16(h.Id)
	r.WriteUint16(h.flag())
	r.WriteUint16(h.QDCount)
	r.WriteUint16(h.ANCount)
	r.WriteUint16(h.NSCount)
	r.WriteUint16(h.ARCount)
}

func (h *Header) flag() uint16 {
	flag := (uint16(h.Opcode) << OPCODE_SHIFT) & OPCODE_MASK
	flag |= uint16(h.Rcode) & RCODE_MASK
	flag |= uint16(h.Flag) & HEADERFLAG_MASK
	return flag
}

func (h *Header) ToWire(buffer *util.OutputBuffer) {
	buffer.WriteUint16(h.Id)
	buffer.WriteUint16(h.flag())
	buffer.WriteUint16(h.QDCount)
	buffer.WriteUint16(h.ANCount)
	buffer.WriteUint16(h.NSCount)
	buffer.WriteUint16(h.ARCount)
}

func (h *Header) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(";; ->>HEADER<<- opcode: %s, status: %s, id: %d\n", h.Opcode.String(), h.Rcode.String(), h.Id))
	buf.WriteString(";; flags: ")
	if h.GetFlag(FLAG_QR) {
		buf.WriteString(" qr")
	}

	if h.GetFlag(FLAG_AA) {
		buf.WriteString(" aa")
	}

	if h.GetFlag(FLAG_TC) {
		buf.WriteString(" tc")
	}

	if h.GetFlag(FLAG_RD) {
		buf.WriteString(" rd")
	}

	if h.GetFlag(FLAG_RA) {
		buf.WriteString(" ra")
	}

	if h.GetFlag(FLAG_AD) {
		buf.WriteString(" ad")
	}

	if h.GetFlag(FLAG_CD) {
		buf.WriteString(" cd")
	}
	buf.WriteString("; ")

	buf.WriteString(fmt.Sprintf("QUERY: %d, ", h.QDCount))
	buf.WriteString(fmt.Sprintf("ANSWER: %d, ", h.ANCount))
	buf.WriteString(fmt.Sprintf("AUTHORITY: %d, ", h.NSCount))
	buf.WriteString(fmt.Sprintf("ADDITIONAL: %d, ", h.ARCount))
	buf.WriteString("\n")
	return buf.String()
}
