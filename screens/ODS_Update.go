package main
 
import (
"net/http"
"fmt"
"io/ioutil"
"go-ini/ini"
 "bytes"
"time"
"strings"
"encoding/json"
"os"
)

var url string
var rurl string
var user string
var secret string
var schedule_interval string 
var manager string

type Ods struct{
 Cont Contract
 Tranid string
 Posted string
}

type Res struct{
	Status string
	Message string
 }

type Ret struct{
 Jsonrpc string
 Result Res
 Id string
 }

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





func OdsUpdate()(string){

  //***********************************************
  //* Invoke manager to create the update journal

jstr :=`{
   	  "jsonrpc": "2.0",
    	 "method": "invoke",
    	 "params": {
      	   "type": 1,
     	    "chaincodeID": {
      	       "name":"`+manager+`"
         },
         "ctorMsg": {
             "function": "CrtOdsjournal",
             "args": [ 
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` 

    rsp:=Invoke( jstr)
 //*******************************
 // get the Journal ID for the ODS Update 
	i:= strings.LastIndex( string(rsp) , "message\":\"" )
 	journalId:=string(rsp)[i+10:i+46]
	fmt.Println("Journal ID="+ journalId)
  //***********************************************
  //* Query manager to get the update
  //* but sleep for a few moments to ensure the invoke has completed
	fmt.Println("Waiting for the Journal to be propageted across the nodes")
       time.Sleep( 5 * time.Second  )

   jstr =`{
   	  "jsonrpc": "2.0",
    	 "method": "query",
    	 "params": {
      	   "type": 1,
     	    "chaincodeID": {
      	       "name":"`+manager+`"
         },
         "ctorMsg": {
             "function": "GetOdsjournal",
             "args": [
		"`+journalId+`" 
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` 

    rsp=Invoke( jstr)
    var ret Ret
    json.Unmarshal( []byte(rsp) , &ret)
    // get the payload into a map object
    transactions:=make(map[string]Ods)
    json.Unmarshal( []byte(ret.Result.Message), &transactions)	
    fmt.Println("********************************************************")
    for k , v :=range transactions {
	b,_:=json.Marshal(v.Cont)
	fmt.Println(k+"-"+string(b))
    }

 return "Done"
}

func Invoke( json string  )(string){


    var jsonStr =[]byte(json)

    fmt.Println("Update:", string(jsonStr) )
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

    fmt.Println("Invoking Chain Code:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Response:", string(body))
    return  string(body)
}

func GetJournal(journalId string )(string){

  //***********************************************
  //* Query manager to get the update

   jstr :=`{
   	  "jsonrpc": "2.0",
    	 "method": "query",
    	 "params": {
      	   "type": 1,
     	    "chaincodeID": {
      	       "name":"`+manager+`"
         },
         "ctorMsg": {
             "function": "GetOdsjournal",
             "args": [
		"`+journalId+`" 
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` 

    rsp:=Invoke( jstr)
    var ret Ret
    json.Unmarshal( []byte(rsp) , &ret)
    // get the payload into a map object
    transactions:=make(map[string]Ods)
	fmt.Println("********************************************************")
    json.Unmarshal( []byte(ret.Result.Message), &transactions)	
    for k , v :=range transactions {
	b,_:=json.Marshal(v.Cont)
	fmt.Println(k+"-"+string(b))
    }

 return "Done"
}




func main() {

cfg, _ := ini.InsensitiveLoad("Contract.ini")
url = cfg.Section("").Key("URL1").String()
user = cfg.Section("").Key("USER").String()
secret = cfg.Section("").Key("SECRET").String()
manager=cfg.Section("").Key("MANAGER").String()
rurl = cfg.Section("").Key("REGISTER").String()


signIn()
t:=time.Now()
fmt.Println(t.String())
fmt.Println("**************** Starting ODS Update***********************")

if len(os.Args) == 1 {
   OdsUpdate()
}else{
GetJournal( os.Args[1] )
}
t=time.Now()
fmt.Println("**************** End ODS Update***********************")
fmt.Println(t.String())
}

