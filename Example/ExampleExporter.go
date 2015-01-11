package main
//import "github.com/davecheney/profile"
import (
	"github.com/joernweissenborn/aursir4go"
	_ "log"
	"github.com/joernweissenborn/aursir4go/Example/keys"
)

func main(){
	//defer profile.Start(profile.CPUProfile).Stop()

	iface, _:=aursir4go.NewInterface("testex")
	iface.AddExport(keys.HelloAurSirAppKey, nil)


	//	exp := iface.AddExport(aursir4go.HelloAurSirAppKey,[]string{})
//
//	for r := range exp.Request {
//		var sayhelloreq aursir4go.SayHelloReq
//		r.Decode(&sayhelloreq)
//		log.Println("Gox
// t",sayhelloreq.Greeting)
//		exp.Reply(&r,aursir4go.SayHelloRes{"MOINSEN, you said"+sayhelloreq.Greeting})
//	}
	iface.WaitUntilDocked()
	iface.Close()

}
