package main
//import "github.com/davecheney/profile"
import (
	"github.com/joernweissenborn/aursir4go"
	"log"
	"github.com/joernweissenborn/aursir4go/Example/keys"
<<<<<<< HEAD
	"log"
)

func main(){
	iface, err := aursir4go.NewInterface("ExampleExporter")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Docking to runtime...")

	iface.WaitUntilDocked()

	log.Println("...Done")

	exp := iface.AddExport(keys.HelloAurSirAppKey, []string{"Tag1", "Tag2"})

	for request := range exp.Request {
		var sayhelloreq keys.SayHelloReq
		request.Decode(&sayhelloreq)
		log.Println("Got greeting:",sayhelloreq.Greeting)
		exp.Reply(&request,keys.SayHelloRes{"Greetings back from AurSir4Go!"})

	}
=======
	"fmt"
	"time"
)

func main(){
	//defer profile.Start(profile.CPUProfile).Stop()
	buf := []byte{123, 34, 83, 97, 121, 72, 101, 108, 108, 111, 34, 58, 34, 70, 114, 111, 109, 34, 125}
	log.Println(string(buf))
	iface, _:=aursir4go.NewInterface("testex")
//	iface.AddExport(keys.HelloAurSirAppKey, nil)


		exp := iface.AddExport(keys.HelloAurSirAppKey,[]string{})

	for r := range exp.Request {
		var sayhelloreq keys.SayHelloReq
		r.Decode(&sayhelloreq)
		log.Println(r)
		log.Println(sayhelloreq.Greeting)
		time.Sleep(1*time.Second)
		exp.Reply(r,keys.SayHelloRes{"MOINSEN, you said"+sayhelloreq.Greeting})
		
	}
	iface.Close()
	fmt.Println("done")
>>>>>>> v0.2devel
}
