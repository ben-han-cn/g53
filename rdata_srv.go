package g53

import (
	"errors"
	"strings"

	"g53/util"
)

type SRV struct {
	Priority uint16
	Weight   uint16
	Port     uint16
	Target   *Name
}

func (srv *SRV) Rend(r *MsgRender) {
	rendField(RDF_C_UINT16, srv.Priority, r)
	rendField(RDF_C_UINT16, srv.Weight, r)
	rendField(RDF_C_UINT16, srv.Port, r)
	rendField(RDF_C_NAME_UNCOMPRESS, srv.Target, r)
}

func (srv *SRV) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_UINT16, srv.Priority, buffer)
	fieldToWire(RDF_C_UINT16, srv.Weight, buffer)
	fieldToWire(RDF_C_UINT16, srv.Port, buffer)
	fieldToWire(RDF_C_NAME, srv.Target, buffer)
}

func (srv *SRV) Compare(other Rdata) int {
	otherSRV := other.(*SRV)
	order := fieldCompare(RDF_C_UINT16, srv.Priority, otherSRV.Priority)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT16, srv.Weight, otherSRV.Weight)
	if order != 0 {
		return order
	}

	order = fieldCompare(RDF_C_UINT16, srv.Port, otherSRV.Port)
	if order != 0 {
		return order
	}

	return fieldCompare(RDF_C_NAME, srv.Target, otherSRV.Target)
}

func (srv *SRV) String() string {
	var ss []string
	ss = append(ss, fieldToStr(RDF_D_INT, srv.Priority))
	ss = append(ss, fieldToStr(RDF_D_INT, srv.Weight))
	ss = append(ss, fieldToStr(RDF_D_INT, srv.Port))
	ss = append(ss, fieldToStr(RDF_D_NAME, srv.Target))
	return strings.Join(ss, " ")
}

func SRVFromWire(buffer *util.InputBuffer, ll uint16) (*SRV, error) {
	p, ll, err := fieldFromWire(RDF_C_UINT16, buffer, ll)
	if err != nil {
		return nil, err
	}

	w, ll, err := fieldFromWire(RDF_C_UINT16, buffer, ll)
	if err != nil {
		return nil, err
	}

	port, ll, err := fieldFromWire(RDF_C_UINT16, buffer, ll)
	if err != nil {
		return nil, err
	}

	t, ll, err := fieldFromWire(RDF_C_NAME, buffer, ll)
	if err != nil {
		return nil, err
	}

	if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	}

	return &SRV{p.(uint16), w.(uint16), port.(uint16), t.(*Name)}, nil
}

func SRVFromString(s string) (*SRV, error) {
	fields := strings.Split(s, " ")
	if len(fields) != 4 {
		return nil, errors.New("short of fields for srv")
	}

	p, err := fieldFromStr(RDF_D_INT, fields[0])
	if err != nil {
		return nil, err
	}
	priority, _ := p.(uint16)

	w, err := fieldFromStr(RDF_D_INT, fields[1])
	if err != nil {
		return nil, err
	}
	weight, _ := w.(uint16)

	p, err = fieldFromStr(RDF_D_INT, fields[2])
	if err != nil {
		return nil, err
	}
	port, _ := p.(uint16)

	t, err := fieldFromStr(RDF_D_NAME, fields[3])
	if err != nil {
		return nil, err
	}
	target, _ := t.(*Name)

	return &SRV{priority, weight, port, target}, nil
}
