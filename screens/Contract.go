/*
 * File: Contract.go
 * Date: 01 June 2017
 * Author: Ellwood Technology Consulting
 * 
 * Copyright (2017) Ellwood Technology Consultimg Sdn Bhd, all rights reserved. 
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Ellwood Technology Consultimg Sdn Bhd,
 * The intellectual and technical concepts contained
 * herein are proprietary to Ellwood Technology Consultimg Sdn Bhd
 * and its suppliers  and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Ellwood Technology Consultimg Sdn Bhd.
 */

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
//"strconv"
"encoding/json"
"sort"
)
type Res struct{
	Status string
	Message string
 }

type Ret struct{
 Jsonrpc string
 Result Res
 Id string
 }

type Wallet struct{
 User string
 Policies map[string]string
 Message string
}

var wa Wallet



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
 ContID string
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
var polid string
var cont Contract
var regid string
 
type Page struct {
    Title string
    Body  []byte
}
type Txn struct{
  id string
  Txnstr string
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

type History struct{
 Methd string
 Funct string
 Cont Contract
 Args []string
 Tranid string 
 Dte string
 EndValue string
}

var hist map[string]History

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

r.ParseForm()
//fmt.Fprintln(w, r.Form)
//fmt.Print(r.Form)
fmt.Print("calling process\n")
                wa.Message=""
    var txn Txn
    title := r.URL.Path[len("/process/"):]
	  fmt.Println("Title="+title+":")

    if title=="Txn.html" {
	txn=getTransactions(cont.ContID)

    }
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
		 cont.ContID=createContract(cont)
		 register(user,cont.ContID)
	 	 //var cm map[string]string
		  //x:=int64(count)
		 //cm["000" + strconv.FormatInt( x ,10)]=ccid
		 count++
		 fmt.Println(ccid)
        	 title="EnterContract.html"
 		 //policyList("admin" )
		 wa.Policies[cont.ContID]=cont.ContID
           } else {
		fmt.Print("Update existing contract")
	  }
        }
	if title=="pay" {
		fmt.Print("Process Pay")
		payment( r.FormValue("pay"), cont )
        	title="Payment.html"
                wa.Message="Payment submitted"
 	} 
	if title=="surr" {
		fmt.Print("Process Surrender")
		surrender( r.FormValue("surramount"), cont )
        	title="Surrender.html"
                wa.Message="Surrender request submitted"
 	}
	if title=="fund" {
    		var f map[string]string
		f= make(map[string]string)
		f["A"]=r.FormValue("funda")
		f["B"]=r.FormValue("fundb")
		f["C"]=r.FormValue("fundc")
		f["D"]=r.FormValue("fundd")
	 	cont.Acct.Fnds[0].FundId="A"
	 	cont.Acct.Fnds[0].Units=f["A"]
	 	cont.Acct.Fnds[1].FundId="B"
	 	cont.Acct.Fnds[1].Units=f["B"]
	 	cont.Acct.Fnds[2].FundId="C"
	 	cont.Acct.Fnds[2].Units=f["C"]
	 	cont.Acct.Fnds[3].FundId="D"
	 	cont.Acct.Fnds[3].Units=f["D"]


		fundupdate(cont.ContID, f)
        	title="fund_Switch.html"
	}

    //**p, err := loadPage("colour_orange/"+title)
    t, err := template.ParseFiles("colour_orange/"+title)
	fmt.Print(err)
	fmt.Print("\n")
	if title == "Payment.html" || title == "fund_Switch.html" || title == "Surrender.html" {
		fmt.Print("t.Execute(w, wa.Policies )")
    	   t.Execute(w, wa.Policies )	
	}else 	if title == "Txn.html" {
    	   t.Execute(w, template.HTML(txn.Txnstr) )
	}else{
    	   t.Execute(w, cont )
        }
	//fmt.Fprintln(w, string(p.Body))
	//**w.Write(p.Body)
}
func getTransactions( contid string) Txn {
  var jsonStr = []byte( `
  {
     "jsonrpc": "2.0",
     "method": "query",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"`+ccid+`"
         },
         "ctorMsg": {
             "function": "transactions",
             "args": [
		"`+contid+`"
             ]
         },
         "secureContext":"admin"
     },
     "id": 2
 }` )
    fmt.Println("Transaction History :"+string( jsonStr) )	
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
    fmt.Println("Transaction History:", string(body))
    var ret Ret
    err=json.Unmarshal(body , &ret)
    str:=`      <table>
		<tr>
		 <td>Date</td>
		 <td>transaction</td>
		 <td>Transaction amt</td>
		 <td>Contract Value</td>
		 <td>Sel</td>
		</tr>`
    json.Unmarshal([]byte(ret.Result.Message), &hist)
    var a []string
    for kk , _ := range hist {
         a= append(a,kk)
     }	
    sort.Strings(a)
    for _ , k := range a {

	str=str+`<TR><TD>`+hist[k].Dte+`</TD><TD>`+hist[k].Funct+`</TD><TD>`+hist[k].Args[0]+`</TD><TD>`+hist[k].EndValue+`</TD><TD></TD></TR>`
     }	
    str=str+`</TABLE>`
  fmt.Println("TRANS="+str)
  var txn Txn;
  txn.Txnstr=str
  return txn 
}
func fundupdate(contid string,  f map[string]string)(){
  b, err := json.Marshal(f )
  var jsonStr = []byte( `
  {
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"`+ccid+`"
         },
         "ctorMsg": {
             "function": "fundAllocation",
             "args": [
		"`+contid+`",
		"`+strings.Replace(string(b),"\"","\\\"",-1) +`"
             ]
         },
         "secureContext":"admin"
     },
     "id": 2
 }` )
    fmt.Println("Fund request :"+string( jsonStr) )	
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
    fmt.Println("Fund Update:", string(body))
    fmt.Println(err)
}

func signIn()(string){
   var jsonStr = []byte( `{
  	"enrollId": "`+ user +`",
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

func  policyList(user string  )(string){ 
   var jsonStr = []byte( `
{
     "jsonrpc": "2.0",
     "method": "query",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"`+regid+`"
         },
         "ctorMsg": {
             "function": "get",
             "args": [
		"`+user+`"
             ]
         },
         "secureContext":"admin"
     },
     "id": 2
 }
` )
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
    fmt.Println("POLICY LIST:", string(body))
    //*************************************
    // get the list & marshal it
    var ret Ret
    err=json.Unmarshal(body , &ret)
    
    json.Unmarshal([]byte(ret.Result.Message), &wa)
          for key, value := range wa.Policies{
		fmt.Println("KEY="+key+" Value="+value)
                cont.ContID=key	
	  }	
   return  string(body)
}


func  register(user string ,id string )(string){ 
   var jsonStr = []byte( `
{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"`+regid+`"
         },
         "ctorMsg": {
             "function": "update",
             "args": [
		"`+user+`",
           "`+id+`"
             ]
         },
         "secureContext":"admin"
     },
     "id": 2
 }
` )
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
    fmt.Println("Register:", string(body))
    fmt.Println(err)
   return  string(body)
}



func createContract( cont Contract)(string){
	 cont.Acct.Fnds[0].FundId="A"
	 cont.Acct.Fnds[0].Units="20"
	 cont.Acct.Fnds[1].FundId="B"
	 cont.Acct.Fnds[1].Units="20"
	 cont.Acct.Fnds[2].FundId="C"
	 cont.Acct.Fnds[2].Units="20"
	 cont.Acct.Fnds[3].FundId="D"
	 cont.Acct.Fnds[3].Units="20"

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
  

var jsonStr = []byte( `
{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"`+ccid+`"
         },
         "ctorMsg": {
             "function": "NewPolicy",
             "args": [`+ args +` 
        ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }`)


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
fmt.Println("CONTRACT ID="+ string(body)[i+10:i+36])
 contractId:=string(body)[i+10:i+46]
 return contractId
}
/**********************************
func activate(   )(string){
  var jsonStr = []byte( `{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"`+ccid+`"
         },
         "ctorMsg": {
             "function": "activate",
             "args": [
                 "`+ccid+`"
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

    fmt.Println("activate  Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(" activate Response:", string(body))
    return ccid


}
***************************************/


func payment(  payment string , cont Contract )(string){
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
                 "`+cont.ContID+`",
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
    return ccid
}
func surrender(  surrvalue string , cont Contract )(string){
  var jsonStr = []byte( `{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name":"`+ccid+`"
         },
         "ctorMsg": {
             "function": "surrender",
             "args": [
                 "`+cont.ContID+`",
                 "`+surrvalue+`"
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` )

    fmt.Println("Surrender:", string(jsonStr) )
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

    fmt.Println("Surrender Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Surender Response:", string(body))
    return ccid
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
schedule_interval = cfg.Section("").Key("SCHEDULE_INTERVAL").String()
glmanager= cfg.Section("").Key("GL_MANAGER").String()
odsmanager= cfg.Section("").Key("ODS_MANAGER").String()
commsmanager= cfg.Section("").Key("COMMS_MANAGER").String()
ccid= cfg.Section("").Key("CCID").String()
regid=cfg.Section("").Key("POLICY_REGISTER").String()
fmt.Print(err)
fmt.Print(url)
fmt.Print(user)
fmt.Print(ccid)
signIn()
t:=time.Now()
fmt.Println(t.String())
http.HandleFunc("/process/", process)
http.Handle("/process/style/", http.StripPrefix("/process/style/", http.FileServer(http.Dir("/Go/src/github.com/dellwoo2/ulcontract/screens/style"))))
http.HandleFunc("/edit/", editHandler)
wa.Policies=make(map[string]string)
 policyList("admin" )
server.ListenAndServe()
}

