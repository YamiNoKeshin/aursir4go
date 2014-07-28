package AurSir4Go

import (
	"testing"
	"time"
	"log"
)

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


	res,_ :=	imp.CallFunction(Testkey.Functions[0].Name,SayHelloReq{"AHOI"},ONE2ONE)
	req := <-exp.Request
    var SayHelloReq SayHelloReq
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI"{
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req,SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var result SayHelloRes
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


	imp2.CallFunction(Testkey.Functions[0].Name,SayHelloReq{"AHOI"},MANY2ONE)
	req := <-exp.Request
    var SayHelloReq SayHelloReq
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI"{
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req,SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var res1 SayHelloRes
	imp1.Listen().Decode(&res1)
	log.Println("res1",res1)
	if res1.Answer != "MOINSEN"{
		T.Error("got wrong result parameter")
	}
	var res2 SayHelloRes
	imp2.Listen().Decode(&res2)
	log.Println("res2",res2)
	if res2.Answer != "MOINSEN"{
		T.Error("got wrong result parameter")
	}
}


func TestDelayedExporter(T *testing.T){
	importer, imp := testimporter()
	defer importer.Close()



	res,_ :=	imp.CallFunction(Testkey.Functions[0].Name,SayHelloReq{"AHOI"},ONE2ONE)
	exporter, exp := testexporter()
	defer exporter.Close()
	req := <-exp.Request
	var SayHelloReq SayHelloReq
	req.Decode(&SayHelloReq)
	log.Println(SayHelloReq)

	if SayHelloReq.Greeting != "AHOI"{
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req,SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var result SayHelloRes
	(<-res).Decode(&result)
	log.Println(result)
	if result.Answer != "MOINSEN"{
		T.Error("got wrong result parameter")
	}
}

func TestTagging(T *testing.T){
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()
	time.Sleep(100*time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}

	imp.UpdateTags([]string{"testtag"})
	time.Sleep(300*time.Millisecond)
	if imp.Connected == true {
		T.Error("could not disconnect from appkey")
	}
	exp.UpdateTags([]string{"testtag"})

	time.Sleep(100*time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}
	exp.UpdateTags([]string{"testtag","anothertag"})

	time.Sleep(100*time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}

	exp.UpdateTags([]string{"anothertag"})
	time.Sleep(300*time.Millisecond)
	if imp.Connected == true {
		T.Error("could not disconnect from appkey")
	}
	imp.UpdateTags([]string{})
	time.Sleep(100*time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}
}
/*func replyOnce(exp ExportedAppKey) {

}*/




func testexporter() (*AurSirInterface, *ExportedAppKey){
	iface:=NewInterface("testex")


	exp:=iface.AddExport(Testkey,nil)
	return iface, exp

}
func testimporter() (*AurSirInterface,*ImportedAppKey){
	iface:=NewInterface("testimp")
	imp:=iface.AddImport(Testkey,nil)
	return iface, imp
}

