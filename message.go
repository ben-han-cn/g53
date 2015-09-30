package g53

import (
	"log"
	"os"

	"g53/util"
)

type QID uint16
type Section []RRset

const (
	QDSection  = "QUESTION"
	ANSection  = "ANSWER"
	NSSection  = "AUTHORITY"
	ARSection  = "ADDITIONAL"
	NumSection = 4
)

type Message struct {
	qid      QID
	rcode    Rcode
	opcode   Opcode
	header   MsgHeader
	question Question
	sections [NumSection]Section
}
