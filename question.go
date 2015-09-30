package g53

import (
	"strings"
)

type Question struct {
	name Name
	typ  RRType
	cls  RRClass
}

func QuestionFromWire(buffer *util.InputBuffer) (*Question, error) {
	n, err := NameFromWire(buffer, true)
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
		name: n,
		typ:  t,
		cls:  cls,
	}, nil
}

func (q *Question) Rend(render *MsgRender) {
	render.WriteName(q.name, false)
	render.WriteUint16(q.typ)
	render.WriteUint16(q.cls)
}

func (q *Question) ToWire(buffer *util.OutputBuffer) {
	q.name.ToWire(buffer)
	buffer.WriteUint16(q.typ)
	buffer.WriteUint16(q.cls)
}

func (q *Question) String() string {
	return strings.Join([]string{q.name.String(), q.typ.String(), q.cls.String()}, " ")
}
