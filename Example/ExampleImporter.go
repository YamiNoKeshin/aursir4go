package main
import (
	"github.com/joernweissenborn/aursir4go"
	_ "log"
	"github.com/joernweissenborn/aursir4go/Example/keys"
<<<<<<< HEAD
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

=======
	"fmt"
	"log"
)

func main(){
	//defer profile.Start(profile.CPUProfile).Stop()

	iface, _:=aursir4go.NewInterface("testex")
	imp := iface.AddImport(keys.HelloAurSirAppKey, nil)
	req,_ := imp.Call("SayHello",keys.SayHelloReq{"Hello from go"})

	var res keys.SayHelloRes
	(<-req).Decode(&res)
	log.Println("Gor result", res)
>>>>>>> v0.2devel
	iface.Close()
}
