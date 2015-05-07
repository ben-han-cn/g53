package g53

import (
	"errors"

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

func (hf HeaderFlag) GetFlag(ff FlagField) bool {
	return (uint16(hf) & uint16(ff)) != 0
}

func (hf HeaderFlag) SetFlag(ff FlagField, set bool) {
	if set {
		hf = HeaderFlag(uint16(hf) | uint16(ff))
	} else {
		hf = HeaderFlag(uint16(hf) & uint16(^ff))
	}
}

type MsgHeader struct {
	Id      uint16
	Flag    HeaderFlag
	Op      Opcode
	Rc      Rcode
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func MsgHeaderFromWire(buffer *util.InputBuffer) (*MsgHeader, error) {
	if buffer.Len() < 12 {
		return nil, errors.New("too short wire data for message header")
	}
	id, _ := buffer.ReadUint16()
	flag, _ := buffer.ReadUint16()
	qdcount, _ := buffer.ReadUint16()
	ancount, _ := buffer.ReadUint16()
	nscount, _ := buffer.ReadUint16()
	arcount, _ := buffer.ReadUint16()
	return &MsgHeader{
		Id:      id,
		Flag:    HeaderFlag(flag & HEADERFLAG_MASK),
		Op:      Opcode((flag & OPCODE_MASK) >> OPCODE_SHIFT),
		Rc:      Rcode(flag & RCODE_MASK),
		QDCount: qdcount,
		ANCount: ancount,
		NSCount: nscount,
		ARCount: arcount,
	}, nil
}

func (h *MsgHeader) Rend(r *MsgRender) {
	r.WriteUint16(h.Id)
	flag := (uint16(h.Op) << OPCODE_SHIFT) & OPCODE_MASK
	flag |= uint16(h.Rc) & RCODE_MASK
	flag |= uint16(h.Flag) | HEADERFLAG_MASK
	r.WriteUint16(flag)
	r.WriteUint16(h.QDCount)
	r.WriteUint16(h.ANCount)
	r.WriteUint16(h.NSCount)
	r.WriteUint16(h.ARCount)
}

func (h *MsgHeader) String() string {
	return ""
}
