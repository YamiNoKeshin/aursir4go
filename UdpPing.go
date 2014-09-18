package aursir4go

import "net"
import (
	"log"
	"time"
)

func pingUdp(UUID string, killFlag *bool){

	var pingtime = 8*time.Second

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5556")
	if err != nil {
		log.Fatal("DOCKERZMQ",err)
	}
	con, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatal("DOCKERZMQ",err)
	}
	t := time.NewTimer(pingtime)

	for _ = range t.C{
		if (*killFlag){
			break
		}
		con.Write([]byte(UUID))
		t.Reset(pingtime)
	}
	log.Println("Stopping UDP")
}
