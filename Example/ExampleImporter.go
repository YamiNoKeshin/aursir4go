package main
import (
	"github.com/joernweissenborn/aursir4go"
	_ "log"
	"github.com/joernweissenborn/aursir4go/Example/keys"
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
	iface.Close()
	fmt.Println("done")
}
