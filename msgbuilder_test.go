package g53

import (
	"testing"

	"github.com/ben-han-cn/g53/util"
)

func TestMessageBuilderFilter(t *testing.T) {
	knet_cn := "04b08180000100010004000d03777777046b6e657402636e0000010001c00c00010001000002580004caad0b0ac01000020001000000c1001404676e7331097a646e73636c6f7564036e657400c01000020001000000c10014046c6e7332097a646e73636c6f75640362697a00c01000020001000000c1001504676e7332097a646e73636c6f7564036e6574c015c01000020001000000c10015046c6e7331097a646e73636c6f756404696e666f00c039000100010000262c000401089801c0790001000100000599000401089901c09a00010001000007c800046f012189c09a00010001000007c8000477a7e9e9c09a00010001000007c80004b683170bc09a00010001000007c80004010865fdc09a001c0001000007c8001024018d00000400000000000000000001c0590001000100002fea000477a7e9ebc0590001000100002fea0004b683170cc0590001000100002fea0004010865fcc0590001000100002fea00046f01218ac059001c00010000249f001024018d000006000000000000000000010000291000000000000000"
	wire, _ := util.HexStrToBytes(knet_cn)
	msg, _ := MessageFromWire(util.NewInputBuffer(wire))
	zone := NameFromStringUnsafe("zdnscloud.biz")
	msg = NewMsgBuilder(msg).FilterRRset(AdditionalSection, func(rrset *RRset) bool {
		return rrset.Name.IsSubDomain(zone)
	}).Done()
	Equal(t, msg.SectionRRsetCount(AdditionalSection), 3)

	msg = NewMsgBuilder(msg).FilterRRset(AnswerSection, func(rrset *RRset) bool {
		return rrset.Type == RR_AAAA
	}).Done()
	Equal(t, msg.SectionRRsetCount(AnswerSection), 0)

	msg = NewMsgBuilder(msg).FilterRRset(AnswerSection, func(rrset *RRset) bool {
		return rrset.Type == RR_AAAA
	}).Done()
	Equal(t, msg.SectionRRsetCount(AnswerSection), 0)

	Equal(t, msg.SectionRRsetCount(AuthSection), 1)
	msg = NewMsgBuilder(msg).FilterRRset(AuthSection, func(rrset *RRset) bool {
		return rrset.Type == RR_NS
	}).Done()
	Equal(t, msg.SectionRRsetCount(AuthSection), 1)
	ns := msg.GetSection(AuthSection)
	Equal(t, ns[0].RRCount(), 4)
}
