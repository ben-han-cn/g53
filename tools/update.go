package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net"

	"github.com/ben-han-cn/g53"
	"github.com/ben-han-cn/g53/util"
)

var (
	addr   string
	key    string
	secret string
	rr     string
	zone   string
)

func init() {
	flag.StringVar(&addr, "addr", "", "dns server port")
	flag.StringVar(&key, "key", "", "key name")
	flag.StringVar(&secret, "secret", "", "key secret")
	flag.StringVar(&zone, "zone", "", "zone to add")
	flag.StringVar(&rr, "rr", "", "rr to add")
}

func main() {
	flag.Parse()

	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("addr %s isn't valid:%s", addr, err.Error())
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatalf("connect to server %s failed:%s", addr, err.Error())
	}
	defer conn.Close()

	zone_, err := g53.NewName(zone, false)
	if err != nil {
		log.Fatalf("zone isn't valid:%s", err.Error())
	}

	rrset, err := g53.RRsetFromString(rr)
	if err != nil {
		log.Fatalf("rr isn't valid:%s", err.Error())
	}

	msg := g53.MakeUpdate(zone_)
	msg.UpdateAddRRset(rrset)
	msg.Header.Id = 1200

	log.Printf("secret is %s\n", secret)
	secret := base64.StdEncoding.EncodeToString([]byte(secret))
	log.Printf("after encode, secret is %s\n", secret)
	tsig, err := g53.NewTSIG(key, secret, "hmac-md5")
	if err != nil {
		log.Fatalf("create tsig failed:%s", err.Error())
	}
	msg.SetTSIG(tsig)
	msg.RecalculateSectionRRCount()

	render := g53.NewMsgRender()
	msg.Rend(render)
	conn.Write(render.Data())

	answerBuffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(answerBuffer)
	if err != nil {
		log.Fatalf("read from server failed:%s", err.Error())
	}

	if n > 0 {
		answer, err := g53.MessageFromWire(util.NewInputBuffer(answerBuffer))
		if err == nil {
			log.Printf(answer.String())
		} else {
			log.Fatalf("get err %s\n", err.Error())
		}
	}
}
