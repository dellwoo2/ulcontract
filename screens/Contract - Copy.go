package main
 
import (
"fmt"
"net/http"
"html/template"
"io/ioutil"
"go-ini/ini"
 "bytes"
"time"
"strings"
"strconv"
)

//*****************************************
//* Contract Types

type Fund struct{
 FundId string
 Units  string
}
type Account struct{
  Fnds [20]Fund
  LastvaluationDate string
  Valuation string
}
type Life struct{
 Name string
 Gender string
 Dob    string
 Smoker string
}
type Contract struct{
 Acct Account
 Product string
 StartDate string
 SumAssured string
 Term  string
 PaymentFrequency string
 Owner  string
 Beneficiary string
 Lf  Life
 Status string
 Email string
}

var cm map[string]string
var url string
var rurl string
var user string
var secret string
var schedule_interval string
var count int
var ccid string
var timerccid string
var startInit string
var odsmanager string
var glmanager string
var commsmanager string
func load(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/load/"):]
fmt.Print("title="+title)
str:=`<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>Go Web Programming</title>
</head>
<body>
<form action="http://127.0.0.1:8080/process?hello=world&thread=123" method="post" enctype="application/x-www-form-urlencoded">
<input type="text" name="hello" value="sau sheong"/>
<input type="text" name="post" value="456"/>
<input type="submit"/>
</form>
</body>
</html>`
fmt.Fprintln(w, str)
}
type Page struct {
    Title string
    Body  []byte
}
type Test struct{
  Title string
  Body string
  AA string
  BB string
  CC string
  DD string
  EE string
}
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    fmt.Print("Page="+title)
    fmt.Print("Page Body="+string(p.Body))
    if err != nil {
        fmt.Print(err)
        p = &Page{Title: title}
    }
    t, _ := template.ParseFiles("Edit.html")
    var x Test
    x.AA="THIS IS AA"
    x.BB="THIS IS BB"
    x.CC="THIS IS CC"
    x.DD="THIS IS DD"
    x.EE="THIS IS EE"
    x.Title="THIS IS THE TITLE"
    x.Body="Some Text"
    t.Execute(w, x)
}

func loadPage(title string) (*Page ,error){
    filename := title
    body, err := ioutil.ReadFile(filename)
    fmt.Print(err)
    return &Page{Title: title, Body: body},nil
}
 
func process(w http.ResponseWriter, r *http.Request) {
var cont Contract
r.ParseForm()
//fmt.Fprintln(w, r.Form)
//fmt.Print(r.Form)
fmt.Print("calling process\n")

    title := r.URL.Path[len("/process/"):]
	  fmt.Println("Title="+title+":")
    if title=="doit" {
	//call the CC API
	if r.FormValue("cid") == "" {
		fmt.Println("Insert new contract")
		// fmt.Println(r.FormValue("product"))
		 //fmt.Println(r.FormValue("owner"))
		// fmt.Println(r.FormValue("name"))
		// fmt.Println(r.FormValue("dob"))
		 cont.Product=r.FormValue("product")
		 cont.Term=r.FormValue("term")
		 cont.SumAssured=r.FormValue("si")
		 cont.PaymentFrequency=r.FormValue("freq")
		 cont.Term=r.FormValue("term")
		 cont.Owner=r.FormValue("owner")
		 cont.Lf.Name=r.FormValue("name")
		 cont.Lf.Dob=r.FormValue("dob")
		 cont.Lf.Gender=r.FormValue("gender")
		 cont.Lf.Smoker=r.FormValue("smoker")
		 cont.StartDate=r.FormValue("start")
		 cont.Beneficiary=r.FormValue("beneficiary")
		 cont.Email=r.FormValue("email")
		 ccid=createContract(cont)
	 	 //var cm map[string]string
		 x:=int64(count)
		 cm["000" + strconv.FormatInt( x ,10)]=ccid
		 count++
		 fmt.Println(ccid)
		 timerccid=createScheduler( ccid )
		 time.Sleep(time.Second * 3 )
		 setScheduler()
		 if startInit =="Y"{
		 	go startScheduler()
                 }
        	title="EnterContract.html"
           } else {
		fmt.Print("Update existing contract")
	  }
        }
	if title=="pay" {
		fmt.Print("Process Pay")
		payment( r.FormValue("pay") )
        	title="Payment.html"
 	}



    p, err := loadPage("colour_orange/"+title)
    //t, err := template.ParseFiles("colour_orange/EnterContract.html")
	fmt.Print(err)
	fmt.Print("\n")
    	//t.Execute(w, x)
	//fmt.Fprintln(w, string(p.Body))
	w.Write(p.Body)
}
 
func signIn()(string){
   var jsonStr = []byte( `{
  	"enrollId": "`+ user+`",
  	"enrollSecret": "`+secret+`"
	}` )
    req, err := http.NewRequest("POST", rurl, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Login:", string(body))
    fmt.Println(err)
   return  string(body)
}
func createContract( cont Contract)(string){

  args:="\""+ cont.Lf.Gender+"\"," +
	"\""+ cont.Lf.Dob +"\"," +
	"\""+ cont.Lf.Smoker +"\"," +
	"\""+ cont.Product +"\"," +
	"\""+ cont.StartDate +"\"," +
	"\""+ cont.Term +"\"," +
	"\""+ cont.PaymentFrequency +"\"," +
	"\""+ cont.Owner +"\"," +
	"\""+ cont.Lf.Name +"\","+
	"\""+ cont.Email +"\"," +
	"\""+ cont.SumAssured +"\"" 
  var jsonStr = []byte( `{"jsonrpc":"2.0","method":"deploy","params":{"type":1,"chaincodeID":{"path": "https://github.com/dellwoo2/ulcontract/ulc"},"ctorMsg":{"function":"init","args":[`+args +`]},"secureContext":"admin"},"id":1}` ) 
  fmt.Println(string(jsonStr))

   req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
/***************************
 var jsonRes = []byte( `{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "05140f2cb0e30fd59b7da4ca1c27ed6ac37329ed7b7d242c3788f72be64207337d73ed159db89993450ac852062e676cda357d5e4681d3d3c8d5ead155bf8f6d"
  },
  "id": 1
}`)
*********************/
i:= strings.LastIndex( string(body) , "message\":\"" )
fmt.Println("CCID="+ string(body)[i+10:i+138])
 ccid=string(body)[i+10:i+138]
 return ccid
}

func createScheduler( ccid string )(string){
  var jsonStr = []byte( `{"jsonrpc": "2.0","method": "deploy",
     "params": {
         "type": 1,
         "chaincodeID": {
             "path": "https://github.com/dellwoo2/ulcontract/timer"
         },
         "ctorMsg": {
             "function": "init",
             "args": [
                  "`+ccid+`"] },"secureContext":"admin"},"id": 1}`)

    fmt.Println("Timer DEPLOY:", string(jsonStr) )
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Timer Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("TimerBody:", string(body))

    i:= strings.LastIndex( string(body) , "message\":\"" )
    fmt.Println("TIMER CCID="+ string(body)[i+10:i+138])
    timerccid=string(body)[i+10:i+138]
    return timerccid
}

func setScheduler(  )(string){
  var jsonStr = []byte( `{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"`+ccid+`"
         },
         "ctorMsg": {
             "function": "setscheduler",
             "args": [
                 "`+timerccid+`",
                 "`+ccid+`",
                 "`+url+`",
                 "`+glmanager+`",
                 "`+odsmanager+`",
		 "`+commsmanager+`"
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` )

    fmt.Println("Set Timer:", string(jsonStr) )
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Set Scheduler Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("SEt Scheduler Body:", string(body))
    return timerccid
}

func startScheduler( )(string){
    e:="F"
  for e=="F" {
    time.Sleep(time.Second * 10 )
    fmt.Println("****STARTING SCHEDULER****")
    var jsonStr = []byte( `{
     "jsonrpc": "2.0",
     "method": "query",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"` +timerccid+ `"
         },
         "ctorMsg": {
             "function": "schedule",
             "args": [
                 "`+ccid+`","`+url+`","`+ schedule_interval+`" 
             ]
         },
         "secureContext": "admin"
     }, "id": 3 }` )

    fmt.Println("Start Scheduler Request:"+string(jsonStr))
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("STart Scheduler Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Set Scheduler Body:", string(body))
    if strings.Index(string(body),"Failed to launch chaincode") == -1 {
	e="T"
    }
 }

    return timerccid
}

func payment(  payment string  )(string){
  var jsonStr = []byte( `{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"`+ccid+`"
         },
         "ctorMsg": {
             "function": "applyPremium",
             "args": [
                 "`+payment+`"
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` )

    fmt.Println("Payment:", string(jsonStr) )
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Payment Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Payment Response:", string(body))
    return timerccid
}



func main() {
cm=make(map[string]string)
server := http.Server{
Addr: "127.0.0.1:8080",
}
count=1
cfg, err := ini.InsensitiveLoad("Contract.ini")
url = cfg.Section("").Key("URL1").String()
rurl = cfg.Section("").Key("REGISTER").String()
user = cfg.Section("").Key("USER").String()
secret = cfg.Section("").Key("SECRET").String()
schedule_interval = cfg.Section("").Key("SCHEDULE").String()
startInit= cfg.Section("").Key("START_SCHEDULE").String()
glmanager= cfg.Section("").Key("GL_MANAGER").String()
odsmanager= cfg.Section("").Key("ODS_MANAGER").String()
commsmanager= cfg.Section("").Key("COMMS_MANAGER").String()

fmt.Print(err)
fmt.Print(url)
fmt.Print(user)
fmt.Print(secret)
signIn()
t:=time.Now()
fmt.Println(t.String())
http.HandleFunc("/process/", process)
http.Handle("/process/style/", http.StripPrefix("/process/style/", http.FileServer(http.Dir("/Go/src/github.com/dellwoo2/ulcontract/screens/style"))))
http.HandleFunc("/load/", load)
http.HandleFunc("/edit/", editHandler)
server.ListenAndServe()
}

