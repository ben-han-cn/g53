package g53

import (
	"github.com/ben-han-cn/g53/util"
)

func MakeAXFR(zone *Name, tsig *Tsig) *Message {
	q := &Question{
		Name:  *zone,
		Type:  RR_AXFR,
		Class: CLASS_IN,
	}
	return NewMsgBuilder(&Message{}).
		SetId(util.GenMessageId()).
		SetOpcode(OP_QUERY).
		SetQuestion(q).
		SetTsig(tsig).
		Done()
}

func MakeIXFR(zone *Name, currentSOA *RRset, tsig *Tsig) *Message {
	q := &Question{
		Name:  *zone,
		Type:  RR_IXFR,
		Class: CLASS_IN,
	}

	return NewMsgBuilder(&Message{}).
		SetId(util.GenMessageId()).
		SetOpcode(OP_QUERY).
		SetQuestion(q).
		SetTsig(tsig).
		AddRRset(AuthSection, currentSOA).
		Done()
}
