package main
//import "github.com/davecheney/profile"
import (
	"github.com/joernweissenborn/aursir4go"
	_ "log"
	"time"
)

func main(){
	//defer profile.Start(profile.CPUProfile).Stop()

	iface:=aursir4go.NewInterface("testex")


//	exp := iface.AddExport(aursir4go.HelloAurSirAppKey,[]string{})
//
//	for r := range exp.Request {
//		var sayhelloreq aursir4go.SayHelloReq
//		r.Decode(&sayhelloreq)
//		log.Println("Got",sayhelloreq.Greeting)
//		exp.Reply(&r,aursir4go.SayHelloRes{"MOINSEN, you said"+sayhelloreq.Greeting})
//	}
	iface.Close()

	time.Sleep(100*time.Second)
}
