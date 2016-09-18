package g53

import (
	"strings"

	"g53/util"
)

type Question struct {
	Name  Name
	Type  RRType
	Class RRClass
}

func QuestionFromWire(buffer *util.InputBuffer) (*Question, error) {
	q := &Question{}
	if err := q.FromWire(buffer); err == nil {
		return q, nil
	} else {
		return nil, err
	}
}

func (q *Question) FromWire(buffer *util.InputBuffer) error {
	err := q.Name.FromWire(buffer, false)
	if err != nil {
		return err
	}

	q.Type, err = TypeFromWire(buffer)
	if err != nil {
		return err
	}

	q.Class, err = ClassFromWire(buffer)
	if err != nil {
		return err
	}

	return nil
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
	return strings.Join([]string{q.Name.String(false), q.Class.String(), q.Type.String()}, " ")
}
