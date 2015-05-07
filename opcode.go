package g53

import ()

type Opcode uint8

const (
	OP_QUERY      Opcode = 0  ///< 0: Standard query (RFC1035)
	OP_IQUERY            = 1  ///< 1: Inverse query (RFC1035)
	OP_STATUS            = 2  ///< 2: Server status request (RFC1035)
	OP_RESERVED3         = 3  ///< 3: Reserved for future use (RFC1035)
	OP_NOTIFY            = 4  ///< 4: Notify (RFC1996)
	OP_UPDATE            = 5  ///< 5: Dynamic update (RFC2136)
	OP_RESERVED6         = 6  ///< 6: Reserved for future use (RFC1035)
	OP_RESERVED7         = 7  ///< 7: Reserved for future use (RFC1035)
	OP_RESERVED8         = 8  ///< 8: Reserved for future use (RFC1035)
	OP_RESERVED9         = 9  ///< 9: Reserved for future use (RFC1035)
	OP_RESERVED10        = 10 ///< 10: Reserved for future use (RFC1035)
	OP_RESERVED11        = 11 ///< 11: Reserved for future use (RFC1035)
	OP_RESERVED12        = 12 ///< 12: Reserved for future use (RFC1035)
	OP_RESERVED13        = 13 ///< 13: Reserved for future use (RFC1035)
	OP_RESERVED14        = 14 ///< 14: Reserved for future use (RFC1035)
	OP_RESERVED15        = 15 ///< 15: Reserved for future use (RFC1035)
)

var OpcodeStr = map[Opcode]string{
	OP_QUERY:      "QUERY",
	OP_IQUERY:     "IQUERY",
	OP_STATUS:     "STATUS",
	OP_RESERVED3:  "RESERVED3",
	OP_NOTIFY:     "NOTIFY",
	OP_UPDATE:     "UPDATE",
	OP_RESERVED6:  "RESERVED6",
	OP_RESERVED7:  "RESERVED7",
	OP_RESERVED8:  "RESERVED8",
	OP_RESERVED9:  "RESERVED9",
	OP_RESERVED10: "RESERVED10",
	OP_RESERVED11: "RESERVED11",
	OP_RESERVED12: "RESERVED12",
	OP_RESERVED13: "RESERVED13",
	OP_RESERVED14: "RESERVED14",
	OP_RESERVED15: "RESERVED15",
}

func (c Opcode) String() string {
	return OpcodeStr[c]
}
