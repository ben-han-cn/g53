package g53

import (
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"g53/util"
)

type RDFCodingType uint8
type RDFDisplayType uint8

const (
	RDF_C_NAME RDFCodingType = iota
	RDF_C_UINT8
	RDF_C_UINT16
	RDF_C_UINT32
	RDF_C_IPV4
	RDF_C_IPV6
	RDF_C_BINARY
	RDF_C_BYTE_BINARY
	RDF_C_TXT
)

const (
	RDF_D_NAME RDFDisplayType = iota
	RDF_D_INT
	RDF_D_IP
	RDF_D_TXT
	RDF_D_HEX
	RDF_D_B32
	RDF_D_B64
	RDF_D_STR
)

func fieldFromWire(ct RDFCodingType, buffer *util.InputBuffer, ll uint16) (interface{}, uint16, error) {
	switch ct {
	case RDF_C_NAME:
		pos := buffer.Position()
		n, err := NameFromWire(buffer, true)
		if err != nil {
			return nil, ll, err
		} else {
			return n, ll - uint16(buffer.Position()-pos), nil
		}

	case RDF_C_UINT8:
		d, err := buffer.ReadUint8()
		if err != nil {
			return nil, ll, err
		} else {
			return d, ll - 1, nil
		}

	case RDF_C_UINT16:
		d, err := buffer.ReadUint16()
		if err != nil {
			return nil, ll, err
		} else {
			return d, ll - 2, nil
		}

	case RDF_C_UINT32:
		d, err := buffer.ReadUint32()
		if err != nil {
			return nil, ll, err
		} else {
			return d, ll - 4, nil
		}

	case RDF_C_IPV4:
		d, err := buffer.ReadBytes(4)
		if err != nil {
			return nil, ll, err
		} else {
			return net.IP(d), ll - 4, nil
		}

	case RDF_C_IPV6:
		d, err := buffer.ReadBytes(16)
		if err != nil {
			return nil, ll, err
		} else {
			return net.IP(d), ll - 16, nil
		}

	case RDF_C_TXT:
		var ss []string
		var d interface{}
		var err error
		for ll > 0 {
			d, ll, err = fieldFromWire(RDF_C_BYTE_BINARY, buffer, ll)
			if err != nil {
				return nil, ll, err
			}
			bs, _ := d.([]uint8)
			ss = append(ss, string(bs))
		}
		return ss, 0, nil

	case RDF_C_BINARY:
		d, err := buffer.ReadBytes(uint(ll))
		if err != nil {
			return nil, ll, err
		}
		return d, 0, nil

	case RDF_C_BYTE_BINARY:
		l, err := buffer.ReadUint8()
		if err != nil {
			return nil, ll, err
		}
		ll -= 1
		if uint16(l) > ll {
			return nil, ll, errors.New("character string is too long")
		}
		d, err := buffer.ReadBytes(uint(l))
		if err != nil {
			return nil, ll, err
		}
		return d, ll - uint16(l), nil

	default:
		return nil, ll, errors.New("unknown rdata file type")
	}
}

func rendField(ct RDFCodingType, data interface{}, render *MsgRender) {
	switch ct {
	case RDF_C_NAME:
		n, _ := data.(*Name)
		n.Rend(render)
	case RDF_C_UINT8:
		d, _ := data.(uint8)
		render.WriteUint8(d)
	case RDF_C_UINT16:
		d, _ := data.(uint16)
		render.WriteUint16(d)
	case RDF_C_UINT32:
		d, _ := data.(uint32)
		render.WriteUint32(d)
	case RDF_C_IPV4, RDF_C_IPV6:
		d, _ := data.(net.IP)
		render.WriteData([]uint8(d))
	case RDF_C_BINARY:
		d, _ := data.([]uint8)
		render.WriteData(d)
	case RDF_C_TXT:
		ds, _ := data.([]string)
		for _, d := range ds {
			rendField(RDF_C_BYTE_BINARY, d, render)
		}
		d, _ := data.([]uint8)
		render.WriteData(d)
	case RDF_C_BYTE_BINARY:
		d, _ := data.([]uint8)
		render.WriteUint8(uint8(len(d)))
		render.WriteData(d)
	}
}

func fieldToWire(ct RDFCodingType, data interface{}, buffer *util.OutputBuffer) {
	switch ct {
	case RDF_C_NAME:
		n, _ := data.(*Name)
		n.ToWire(buffer)

	case RDF_C_UINT8:
		d, _ := data.(uint8)
		buffer.WriteUint8(d)

	case RDF_C_UINT16:
		d, _ := data.(uint16)
		buffer.WriteUint16(d)

	case RDF_C_UINT32:
		d, _ := data.(uint32)
		buffer.WriteUint32(d)

	case RDF_C_IPV4, RDF_C_IPV6:
		ip, _ := data.(net.IP)
		buffer.WriteData([]uint8(ip))

	case RDF_C_BINARY:
		d, _ := data.([]uint8)
		buffer.WriteData(d)

	case RDF_C_TXT:
		ds, _ := data.([]string)
		for _, d := range ds {
			fieldToWire(RDF_C_BYTE_BINARY, d, buffer)
		}
		d, _ := data.([]uint8)
		buffer.WriteData(d)

	case RDF_C_BYTE_BINARY:
		d, _ := data.([]uint8)
		buffer.WriteUint8(uint8(len(d)))
		buffer.WriteData(d)
	}
}

func fieldFromStr(dt RDFDisplayType, s string) (interface{}, error) {
	switch dt {
	case RDF_D_NAME:
		n, err := NameFromString(s)
		if err != nil {
			return nil, err
		} else {
			return n, nil
		}

	case RDF_D_INT:
		d, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		} else {
			return d, nil
		}

	case RDF_D_IP:
		ip := net.ParseIP(s)
		if ip == nil {
			return nil, errors.New("invalid ip address")
		} else {
			return ip, nil
		}

	case RDF_D_TXT:
		d := strings.Split(s, " ")
		return d, nil

	case RDF_D_HEX:
		d, err := util.HexStrToBytes(s)
		if err != nil {
			return nil, err
		} else {
			return d, nil
		}

	case RDF_D_B32:
		d, err := base32.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, err
		} else {
			return []uint8(d), nil
		}

	case RDF_D_B64:
		d, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, err
		} else {
			return []uint8(d), nil
		}

	case RDF_D_STR:
		return s, nil

	default:
		return nil, errors.New("unknown display type")
	}
}

func fieldToStr(dt RDFDisplayType, d interface{}) string {
	switch dt {
	case RDF_D_NAME:
		n, _ := d.(*Name)
		return n.String(false)
	case RDF_D_INT:
		return fmt.Sprintf("%v", d)
	case RDF_D_IP:
		ip, _ := d.(net.IP)
		return ip.String()
	case RDF_D_TXT:
		ss, _ := d.([]string)
		return strings.Join(ss, " ")
	case RDF_D_HEX:
		bs, _ := d.([]uint8)
		s := ""
		for _, b := range bs {
			s += fmt.Sprintf("%x", b)
		}
		return s
	case RDF_D_B32:
		bs, _ := d.([]uint8)
		return base32.StdEncoding.EncodeToString([]byte(bs))
	case RDF_D_B64:
		bs, _ := d.([]uint8)
		return base64.StdEncoding.EncodeToString([]byte(bs))
	case RDF_D_STR:
		s, _ := d.(string)
		return s
	default:
		return ""
	}
}
