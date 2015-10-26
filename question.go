package g53

import (
	"strings"

	"g53/util"
)

type Question struct {
	Name  *Name
	Type  RRType
	Class RRClass
}

func QuestionFromWire(buffer *util.InputBuffer) (*Question, error) {
	n, err := NameFromWire(buffer, false)
	if err != nil {
		return nil, err
	}

	t, err := TypeFromWire(buffer)
	if err != nil {
		return nil, err
	}

	cls, err := ClassFromWire(buffer)
	if err != nil {
		return nil, err
	}

	return &Question{
		Name:  n,
		Type:  t,
		Class: cls,
	}, nil
}

func (q *Question) Rend(r *MsgRender) {
	q.Name.Rend(r)
	q.Type.Rend(r)
	q.Class.Rend(r)
}

func (q *Question) ToWire(buffer *util.OutputBuffer) {
	q.Name.ToWire(buffer)
	q.Type.ToWire(buffer)
	q.Class.ToWire(buffer)
}

func (q *Question) String() string {
	return strings.Join([]string{q.Name.String(false), q.Type.String(), q.Class.String()}, " ")
}
