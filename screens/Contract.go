package main
 
import (
"fmt"
"net/http"
"html/template"
"io/ioutil"
"go-ini/ini"
 "bytes"
"time"
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
var user string
var secret string


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
	  fmt.Println("Title="+title)
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
		 ccid:=createContract(cont)
		fmt.Println(ccid)

        }else{
		fmt.Print("Update existing contract")

	}
        title="EnterContract.html"
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
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Login:", string(body))
   return  string(body)
}
func createContract( cont Contract)(string){
  var ccid string=""

  args:="\""+ cont.Lf.Gender+"\",\n" +
	"\""+ cont.Lf.Dob +"\",\n" +
	"\""+ cont.Lf.Smoker +"\",\n" +
	"\""+ cont.Product +"\",\n" +
	"\""+ cont.StartDate +"\",\n" +
	"\""+ cont.Term +"\",\n" +
	"\""+ cont.PaymentFrequency +"\",\n" +
	"\""+ cont.Owner +"\",\n" +
	"\""+ cont.Lf.Name +"\",\n"+
	"\""+ cont.Email +"\",\n" +
	"\""+ cont.SumAssured +"\"\n" 
  var jsonStr = []byte( `{
     "jsonrpc": "2.0",
     "method": "deploy",

     "params": {
         "type": 1,
         "chaincodeID": {
             "path": "https://github.com/dellwoo2/ulcontract/ulc"
         },
         "ctorMsg": {
             "function": "init",
             "args": [
                 `+ args + `
             ]
         },
         "secureContext": "admin"
     },
     "id": 1
 }`) 
   fmt.Println(string(jsonStr))
/***********************
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
***************************/


 return ccid
}
func main() {
server := http.Server{
Addr: "127.0.0.1:8080",
}
cfg, err := ini.InsensitiveLoad("Contract.ini")
url = cfg.Section("").Key("URL1").String()
user = cfg.Section("").Key("USER").String()
secret = cfg.Section("").Key("SECRET").String()
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

