package main

import (
	"github.com/joernweissenborn/AurSir4Go"
	"log"
)

func main(){
	iface:=AurSir4Go.NewInterface("testex")


	exp:=iface.AddExport(AurSir4Go.Testkey,nil)

	for r := range exp.Request {
		var sayhelloreq AurSir4Go.SayHelloReq
		r.Decode(&sayhelloreq)
		log.Println("Got",sayhelloreq.Greeting)
		exp.Reply(&r,AurSir4Go.SayHelloRes{"MOINSEN, you said"+sayhelloreq.Greeting})
	}
}
