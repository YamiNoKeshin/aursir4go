package main
//import "github.com/davecheney/profile"
import (
	"github.com/joernweissenborn/aursir4go"
	_ "log"
	"github.com/joernweissenborn/aursir4go/Example/keys"
	"log"
	"github.com/joernweissenborn/aursir4go/calltypes"
)

func main(){

	iface, err := aursir4go.NewInterface("ExampleImporter")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Docking to runtime...")

	iface.WaitUntilDocked()

	log.Println("...Done")


	imp := iface.AddImport(keys.HelloAurSirAppKey, []string{})

	sayhelloreq := keys.SayHelloReq{"Hello AurSir"}
	log.Println("Sending 'Hello AurSir'")

	request, err := imp.CallFunction("SayHello",sayhelloreq,calltypes.ONE2ONE)
	if err != nil {
		log.Fatal(err)
	}

	result := <-request

	var answer keys.SayHelloRes
	err = result.Decode(&answer)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Got answer:",answer.Answer)

	iface.Close()
}
