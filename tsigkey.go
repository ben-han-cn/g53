package g53

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"hash"
	"strings"
	"time"

	"github.com/ben-han-cn/g53/util"
)

type TsigAlgorithm string

func AlgorithmFromString(name string) (TsigAlgorithm, error) {
	switch strings.ToLower(name) {
	case "hmac-md5", "hmac-md5.sig-alg.reg.int.":
		return HmacMD5, nil
	case "hmac-sha1", "hmac-sha1.":
		return HmacSHA1, nil
	case "hmac-sha256", "hmac-sha256.":
		return HmacSHA256, nil
	case "hmac-sha512", "hmac-sha512.":
		return HmacSHA512, nil
	default:
		return "", errors.New("No such algorothm")
	}
}

const (
	HmacMD5    TsigAlgorithm = "hmac-md5.sig-alg.reg.int."
	HmacSHA1   TsigAlgorithm = "hmac-sha1."
	HmacSHA256 TsigAlgorithm = "hmac-sha256."
	HmacSHA512 TsigAlgorithm = "hmac-sha512."
)

var ErrSig = errors.New("signature error")
var ErrTime = errors.New("tsig time expired")

type TsigKey struct {
	Name      string
	algo      TsigAlgorithm
	rawSecret []byte
}

func NewTsigKey(name, secret, alg string) (TsigKey, error) {
	algo, err := AlgorithmFromString(alg)
	if err != nil {
		return TsigKey{}, err
	}

	rawSecret, err := fromBase64([]byte(secret))
	if err != nil {
		return TsigKey{}, err
	}

	return TsigKey{
		Name:      name,
		algo:      algo,
		rawSecret: rawSecret,
	}, nil
}

func (k TsigKey) VerifyMAC(msg *Message, requestMac []byte) error {
	tsig := msg.Tsig
	msg.Tsig = nil
	render := NewMsgRender()
	msg.Header.ARCount -= 1
	msg.Rend(render)
	msg.Tsig = tsig
	msg.Header.ARCount += 1

	mac := k.genMessageHash(tsig, render, requestMac, false)
	if !hmac.Equal(mac, tsig.MAC) {
		return ErrSig
	} else {
		return nil
	}
}

func fromBase64(s []byte) ([]byte, error) {
	buflen := base64.StdEncoding.DecodedLen(len(s))
	buf := make([]byte, buflen)
	if n, err := base64.StdEncoding.Decode(buf, s); err != nil {
		return nil, err
	} else {
		return buf[:n], nil
	}
}

func (k TsigKey) hashSelect() hash.Hash {
	switch k.algo {
	case HmacMD5:
		return hmac.New(md5.New, k.rawSecret)
	case HmacSHA1:
		return hmac.New(sha1.New, k.rawSecret)
	case HmacSHA256:
		return hmac.New(sha256.New, k.rawSecret)
	case HmacSHA512:
		return hmac.New(sha512.New, k.rawSecret)
	default:
		panic("unreachable branch")
	}
}

func (k TsigKey) GenerateTsig(msgId uint16, render *MsgRender, requestMac []byte, timerOnly bool) *Tsig {
	tsig := &Tsig{
		Header: TsigHeader{
			Name:     *NameFromStringUnsafe(k.Name),
			Rrtype:   RR_TSIG,
			Class:    CLASS_ANY,
			Ttl:      0,
			Rdlength: 0,
		},
		Algorithm:  k.algo,
		TimeSigned: uint64(time.Now().Unix()),
		Fudge:      300,
		Error:      0,
		OtherLen:   0,
	}
	tsig.OrigId = msgId
	mac := k.genMessageHash(tsig, render, requestMac, timerOnly)
	tsig.MAC = mac
	tsig.MACSize = uint16(len(mac))
	return tsig
}

//render has rend the message for this tsig to generate hash
func (key TsigKey) genMessageHash(tsig *Tsig, render *MsgRender, requestMac []byte, timerOnly bool) []byte {
	h := key.hashSelect()
	if requestMac != nil {
		l := len(requestMac)
		lenBuf := []byte{uint8((l & 0xff00) >> 8), uint8(l & 0x00ff)}
		h.Write(lenBuf)
		h.Write(requestMac)
	}

	h.Write(render.Data())

	if timerOnly {
		buf := util.NewOutputBuffer(8)
		ts1 := uint16((tsig.TimeSigned & 0x0000ffff00000000) >> 32)
		ts2 := uint32(tsig.TimeSigned & 0x00000000ffffffff)
		buf.WriteUint16(ts1)
		buf.WriteUint32(ts2)
		buf.WriteUint16(tsig.Fudge)
		h.Write(buf.Data())
	} else {
		buf := util.NewOutputBuffer(64)
		tsig.Header.Name.ToWire(buf)
		CLASS_ANY.ToWire(buf)
		tsig.Header.Ttl.ToWire(buf)
		NameFromStringUnsafe(string(tsig.Algorithm)).ToWire(buf)
		ts1 := uint16((tsig.TimeSigned & 0x0000ffff00000000) >> 32)
		ts2 := uint32(tsig.TimeSigned & 0x00000000ffffffff)
		buf.WriteUint16(ts1)
		buf.WriteUint32(ts2)
		buf.WriteUint16(tsig.Fudge)
		buf.WriteUint16(tsig.Error)
		buf.WriteUint16(tsig.OtherLen)
		buf.WriteData(tsig.OtherData)
		h.Write(buf.Data())
	}

	return h.Sum(nil)
}
