package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/zdnscloud/g53"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "s", "10.43.26.57:5553", "dns server port default to 53")
}

func main() {
	flag.Parse()

	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic("resolver failed:" + addr)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		panic("connect to server failed")
	}
	defer conn.Close()

	qn, err := g53.NewName("www.zdns.cn", false)
	if err != nil {
		panic("invalid name to query:" + err.Error())
	}
	qtype, err := g53.TypeFromString("a")
	if err != nil {
		panic("invalid type to query:" + err.Error())
	}
	msg := g53.MakeQuery(qn, qtype, 4096, false)
	render := g53.NewMsgRender()

	for {
		msg.Header.Id = uint16(rand.Intn(4096))
		msg.Rend(render)
		n, err := conn.Write(render.Data())
		if err != nil {
			fmt.Printf("send get err %v\n", err)
		} else {
			fmt.Printf("send %d query to %s\n", n, serverAddr)
		}
		<-time.After(time.Second)
		render.Clear()
	}
}
