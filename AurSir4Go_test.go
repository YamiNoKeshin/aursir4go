package AurSir4Go

import (
	"testing"
	"time"
	"log"
)
var testkey = AppKey{
	"org.aursir.helloaursir",
	[]Function{
		Function{
			"SayHello",
			[]Data{
				Data{
					"Greeting",
					1}},
			[]Data{
				Data{
					"Answer",
					1}}}}}
func TestInitCloseIface(t *testing.T){
	iface:=NewInterface("test")
	defer iface.Close()
	for !iface.Connected(){
		time.Sleep(1* time.Millisecond)
	}

}
func TestExportKey(t *testing.T) {
	iface, exp := testexporter()
	defer iface.Close()
	log.Println("ExportId",exp)
	time.Sleep(1*time.Second)

}

func TestImportKey(t *testing.T) {
	iface, imp := testimporter()
	defer iface.Close()
	log.Println("map",*iface.imports)
	time.Sleep(1*time.Second)
	log.Println("Import",imp)


}

func TestKeyAvailable(T *testing.T){
	importer, imp := testimporter()
	defer importer.Close()
	exporter, _ := testexporter()
	time.Sleep(100*time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}
	exporter.Close()
	time.Sleep(100*time.Millisecond)
	if imp.Connected == true {
		T.Error("could not disconnect from appkey")
	}

}


func TestFunCall121(T *testing.T){
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()


	res,_ :=	imp.CallFunction(testkey.Functions[0].Name,sayHelloReq{"AHOI"},ONE2ONE)
	req := <-exp.Request
    var sayhelloreq sayHelloReq
	req.Decode(&sayhelloreq)
	if sayhelloreq.Greeting != "AHOI"{
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req,sayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var result sayHelloRes
	(<-res).Decode(&result)
	log.Println(result)
	if result.Answer != "MOINSEN"{
		T.Error("got wrong result parameter")
	}
}

func TestFunCallN21(T *testing.T){
	importer1, imp1 := testimporter()
	defer importer1.Close()
	imp1.ListenToFunction("SayHello")
	importer2, imp2 := testimporter()
	imp2.ListenToFunction("SayHello")
	defer importer2.Close()
	exporter, exp := testexporter()
	defer exporter.Close()


	imp2.CallFunction(testkey.Functions[0].Name,sayHelloReq{"AHOI"},MANY2ONE)
	req := <-exp.Request
    var sayhelloreq sayHelloReq
	req.Decode(&sayhelloreq)
	if sayhelloreq.Greeting != "AHOI"{
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req,sayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var res1 sayHelloRes
	imp1.Listen().Decode(&res1)
	log.Println("res1",res1)
	if res1.Answer != "MOINSEN"{
		T.Error("got wrong result parameter")
	}
	var res2 sayHelloRes
	imp2.Listen().Decode(&res2)
	log.Println("res2",res2)
	if res2.Answer != "MOINSEN"{
		T.Error("got wrong result parameter")
	}
}


func TestDelayedExporter(T *testing.T){
	importer, imp := testimporter()
	defer importer.Close()



	res,_ :=	imp.CallFunction(testkey.Functions[0].Name,sayHelloReq{"AHOI"},ONE2ONE)
	exporter, exp := testexporter()
	defer exporter.Close()
	req := <-exp.Request
	var sayhelloreq sayHelloReq
	req.Decode(&sayhelloreq)
	log.Println(sayhelloreq)

	if sayhelloreq.Greeting != "AHOI"{
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req,sayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var result sayHelloRes
	(<-res).Decode(&result)
	log.Println(result)
	if result.Answer != "MOINSEN"{
		T.Error("got wrong result parameter")
	}
}


/*func replyOnce(exp ExportedAppKey) {

}*/



type sayHelloReq struct {
	Greeting string
	}

type sayHelloRes struct {
	Answer string
}

func testexporter() (*AurSirInterface, *ExportedAppKey){
	iface:=NewInterface("testex")


	exp:=iface.AddExport(testkey,nil)
	return iface, exp

}
func testimporter() (*AurSirInterface,*ImportedAppKey){
	iface:=NewInterface("testimp")
	imp:=iface.AddImport(testkey,nil)
	return iface, imp
}

