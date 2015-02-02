package main
//import "github.com/davecheney/profile"
import (
	"github.com/joernweissenborn/aursir4go"
	_ "log"
	"github.com/joernweissenborn/aursir4go/Example/keys"
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
}
