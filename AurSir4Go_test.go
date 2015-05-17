package aursir4go

import (
	"log"
	"testing"
	"time"
)

func TestInitCloseIface(t *testing.T) {
	iface := NewInterface("test")
	defer iface.Close()
	for !iface.Connected() {
		time.Sleep(1 * time.Millisecond)
	}

}
func TestExportKey(t *testing.T) {
	iface, exp := testexporter()
	defer iface.Close()
	log.Println("ExportId", exp)
	time.Sleep(1 * time.Second)

}

func TestImportKey(t *testing.T) {
	iface, imp := testimporter()
	defer iface.Close()
	log.Println("map", *iface.imports)
	time.Sleep(1 * time.Second)
	log.Println("Import", imp)

}

func TestKeyAvailable(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, _ := testexporter()
	time.Sleep(100 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}
	exporter.Close()
	time.Sleep(100 * time.Millisecond)
	if imp.Connected == true {
		T.Error("could not disconnect from appkey")
	}

}

func TestFunCall121(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()

	res, _ := imp.CallFunction(Testkey.Functions[0].Name, SayHelloReq{"AHOI"}, ONE2ONE)
	req := <-exp.Request
	var SayHelloReq SayHelloReq
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req, SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var result SayHelloRes
	(<-res).Decode(&result)
	log.Println(result)
	if result.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
}

func TestFunCallN21(T *testing.T) {
	importer1, imp1 := testimporter()
	defer importer1.Close()
	imp1.ListenToFunction("SayHello")
	importer2, imp2 := testimporter()
	imp2.ListenToFunction("SayHello")
	defer importer2.Close()
	exporter, exp := testexporter()
	defer exporter.Close()

	imp2.CallFunction(Testkey.Functions[0].Name, SayHelloReq{"AHOI"}, MANY2ONE)
	req := <-exp.Request
	var SayHelloReq SayHelloReq
	req.Decode(&SayHelloReq)
	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req, SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var res1 SayHelloRes
	imp1.Listen().Decode(&res1)
	log.Println("res1", res1)
	if res1.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
	var res2 SayHelloRes
	imp2.Listen().Decode(&res2)
	log.Println("res2", res2)
	if res2.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
}

func TestDelayedExporter(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()

	res, _ := imp.CallFunction(Testkey.Functions[0].Name, SayHelloReq{"AHOI"}, ONE2ONE)
	exporter, exp := testexporter()
	defer exporter.Close()
	req := <-exp.Request
	var SayHelloReq SayHelloReq
	req.Decode(&SayHelloReq)
	log.Println(SayHelloReq)

	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}
	err := exp.Reply(&req, SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var result SayHelloRes
	(<-res).Decode(&result)
	log.Println(result)
	if result.Answer != "MOINSEN" {
		T.Error("got wrong result parameter")
	}
}

func TestTagging(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()
	time.Sleep(100 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}

	imp.UpdateTags([]string{"testtag"})
	time.Sleep(300 * time.Millisecond)
	if imp.Connected == true {
		T.Error("could not disconnect from appkey")
	}
	exp.UpdateTags([]string{"testtag"})

	time.Sleep(300 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}
	exp.UpdateTags([]string{"testtag", "anothertag"})

	time.Sleep(300 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}

	exp.UpdateTags([]string{"anothertag"})
	time.Sleep(300 * time.Millisecond)
	if imp.Connected == true {
		T.Error("could not disconnect from appkey")
	}
	imp.UpdateTags([]string{})
	time.Sleep(300 * time.Millisecond)
	if imp.Connected == false {
		T.Error("could not connect to appkey")
	}
}

func TestCallChain(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()

	cc, _ := imp.NewCallChain(Testkey.Functions[0].Name, SayHelloReq{"AHOI"}, ONE2ONE)
	paramap := map[string]string{}
	paramap["String"] = "Answer"

	cc.AddCall("org.aursir.countstring", "CountString", paramap, ONE2ONE, []string{})
	err := cc.Finalize()
	log.Println(err)
	if err == nil {
		T.Error("Finalize should have thrown err now")
	}
	exporter1, exp1 := testexporterctrstr()
	defer exporter1.Close()

	err = cc.Finalize()
	if err != nil {
		T.Error("Finalize should not have thrown err now")
	}
	exporter2, exp2 := testexporter()
	defer exporter2.Close()

	req := <-exp2.Request
	var SayHelloReq SayHelloReq
	req.Decode(&SayHelloReq)
	log.Println(SayHelloReq)

	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}

	err = exp2.Reply(&req, SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var csr CountStringReq
	req = <-exp1.Request

	req.Decode(&csr)
	log.Println(csr)

	if csr.String != "MOINSEN" {
		T.Error("got wrong request parameter")
	}
	err = exp2.Reply(&req, CountStringRes{int64(len([]byte(csr.String)))})
	if err != nil {
		T.Error(err)
	}
}

func TestCallChainFinalize(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	imp2 := importer.AddImport(CountStringKey, []string{})
	cc, _ := imp.NewCallChain(Testkey.Functions[0].Name, SayHelloReq{"AHOI"}, ONE2ONE)
	paramap := map[string]string{}
	paramap["String"] = "Answer"

	exporter1, exp1 := testexporterctrstr()
	defer exporter1.Close()

	exporter2, exp2 := testexporter()
	defer exporter2.Close()

	rep, err := imp2.FinalizeCallChain(CountStringKey.Functions[0].Name, paramap, ONE2ONE, cc)

	if err != nil {
		T.Error("Finalize should not have thrown err now")
	}
	req := <-exp2.Request
	var SayHelloReq SayHelloReq
	req.Decode(&SayHelloReq)
	log.Println(SayHelloReq)

	if SayHelloReq.Greeting != "AHOI" {
		T.Error("got wrong request parameter")
	}

	err = exp2.Reply(&req, SayHelloRes{"MOINSEN"})
	if err != nil {
		T.Error(err)
	}
	var csr CountStringReq
	req = <-exp1.Request

	req.Decode(&csr)
	log.Println(csr)

	if csr.String != "MOINSEN" {
		T.Error("got wrong request parameter")
	}

	err = exp2.Reply(&req, CountStringRes{int64(len([]byte(csr.String)))})
	if err != nil {
		T.Error(err)
	}

	rply := <-rep
	var csrep CountStringRes
	rply.Decode(&csrep)

	if csrep.Size != int64(len([]byte(csr.String))) {
		T.Error("got wrong result parameter")
	}
}

func TestPersitenceLogging(T *testing.T) {
	importer, imp := testimporter()
	defer importer.Close()
	exporter, exp := testexporter()
	defer exporter.Close()
	exp.SetLogging("SayHello")
	res, _ := imp.CallFunction(Testkey.Functions[0].Name, SayHelloReq{"AHOI"}, ONE2ONE)
	req := <-exp.Request

	exp.Reply(&req, SayHelloRes{"MOINSEN"})
	<-res
}

func testexporter() (*AurSirInterface, *ExportedAppKey) {
	iface := NewInterface("testex")

	exp := iface.AddExport(Testkey, nil)
	return iface, exp

}
func testimporter() (*AurSirInterface, *ImportedAppKey) {
	iface := NewInterface("testimp")
	imp := iface.AddImport(Testkey, nil)
	return iface, imp
}

func testexporterctrstr() (*AurSirInterface, *ExportedAppKey) {
	iface := NewInterface("testex")

	exp := iface.AddExport(CountStringKey, nil)
	return iface, exp

}
