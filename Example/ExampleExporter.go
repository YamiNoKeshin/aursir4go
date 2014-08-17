package main

import (
	"github.com/joernweissenborn/aursir4go"
	_ "log"
	"time"
)

func main(){
	iface:=aursir4go.NewInterface("testex")


	_=iface.AddExport(aursir4go.Testkey,nil)

	/*for r := range exp.Request {
		var sayhelloreq aursir4go.SayHelloReq
		r.Decode(&sayhelloreq)
		log.Println("Got",sayhelloreq.Greeting)
		exp.Reply(&r,aursir4go.SayHelloRes{"MOINSEN, you said"+sayhelloreq.Greeting})
	}*/
	time.Sleep(10*time.Second)
	iface.Close()
}
