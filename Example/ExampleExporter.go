package main

import (
	"github.com/joernweissenborn/aursir4go"
	"log"
)

func main(){
	iface:=aursir4go.NewInterface("testex")


	exp:=iface.AddExport(aursir4go.Testkey,nil)

	for r := range exp.Request {
		var sayhelloreq aursir4go.SayHelloReq
		r.Decode(&sayhelloreq)
		log.Println("Got",sayhelloreq.Greeting)
		exp.Reply(&r,aursir4go.SayHelloRes{"MOINSEN, you said"+sayhelloreq.Greeting})
	}
}
