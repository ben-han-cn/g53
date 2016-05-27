package g53

import (
	"fmt"
	"g53/util"
)

const (
	EDNS_VIEW = 53
)

type ViewOpt struct {
	view string
}

func (vo *ViewOpt) Rend(render *MsgRender) {
	render.WriteUint16(EDNS_VIEW)
	viewLen := len(vo.view)
	render.WriteUint16(uint16(1 + viewLen))
	render.WriteUint8(uint8(len(vo.view)))
	render.WriteData([]byte(vo.view))
}

func (vo *ViewOpt) String() string {
	return fmt.Sprintf("; CLIENT-VIEW: %s\n", vo.view)
}

//read from OPTION-LENGTH
func viewOptFromWire(buffer *util.InputBuffer) (Option, error) {
	l, err := buffer.ReadUint8()
	if err != nil {
		return nil, err
	}

	view, err := buffer.ReadBytes(uint(l))
	if err != nil {
		return nil, err
	}

	return &ViewOpt{
		view: string(view),
	}, nil
}

func viewOptFromRdata(rdata Rdata) Option {
	opt := rdata.(*OPT)
	if len(opt.Data) == 0 {
		return nil
	}

	buffer := util.NewInputBuffer(opt.Data)
	code, _ := buffer.ReadUint16()
	if code != EDNS_VIEW {
		return nil
	}

	option, err := viewOptFromWire(buffer)
	if err == nil {
		return option
	} else {
		return nil
	}
}

func (e *EDNS) AddSubnetView(view string) error {
	e.Options = append(e.Options, &ViewOpt{
		view: view,
	})
	return nil
}
