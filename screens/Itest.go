package main
 
import (
"net/http"
"fmt"
"io/ioutil"
"go-ini/ini"
 "bytes"
"time"
"strings"
)

var url string
var rurl string
var user string
var secret string
var schedule_interval string 



 
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



func startScheduler( timerccid string , ccid string )(string){
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
                 "`+ schedule_interval+`",
		 "`+ccid+`"
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

func pingScheduler( timerccid string )(string){
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
             "function": "ping",
             "args": [
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

    fmt.Println("Deploy Scheduler Ping Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Set Scheduler Body:", string(body))
    if strings.Index(string(body),"Failed to launch chaincode") == -1 {
	e="T"
    }
 }

    return timerccid
}
func Itest2()(string){
//  var timerccid string
//  var manager string
//  var commsmanager string
//  var ccid string
//  var policyregister string
 jstr:=`{
     "jsonrpc": "2.0",
     "method": "deploy",
     "params": {
         "type": 1,
         "chaincodeID": {
             "path": "github.com/dellwoo2/ulcontract/manager"
         },
         "ctorMsg": {
             "function": "init",
             "args": [
             ]
         },
         "secureContext": "admin"
     },
     "id": 1
  }
 `
    Invoke( jstr)
    return "Done"
}

func Itest()(string){
  var timerccid string
  var manager string
  var commsmanager string
  
  var policyregister string
 jstr:=`{
     "jsonrpc": "2.0",
     "method": "deploy",
     "params": {
         "type": 1,
         "chaincodeID": {
             "path": "github.com/dellwoo2/ulcontract/manager"
         },
         "ctorMsg": {
             "function": "init",
             "args": [
             ]
         },
         "secureContext": "admin"
     },
     "id": 1
  }
 `
    rsp:=Invoke( jstr)
    //*******************************
    // get the CCID for comms
    i:= strings.LastIndex( string(rsp) , "message\":\"" )
    fmt.Println("Combined GL ODS Fund Manager CCID="+ string(rsp)[i+10:i+138])
    manager=rsp[i+10:i+138]

//*********************************
  jstr=`{
     "jsonrpc": "2.0",
     "method": "deploy",
     "params": {
         "type": 1,
         "chaincodeID": {
             "path": "github.com/dellwoo2/ulcontract/Comms"
         },
         "ctorMsg": {
             "function": "init",
             "args": [
             ]
         },
         "secureContext": "admin"
     },
     "id": 1
  }
 `
    rsp=Invoke( jstr)
    //*******************************
    // get the CCID for comms
    i= strings.LastIndex( string(rsp) , "message\":\"" )
    fmt.Println("COMMS CCID="+ string(rsp)[i+10:i+138])
    commsmanager=rsp[i+10:i+138]
//*********************************
  jstr=`{
     "jsonrpc": "2.0",
     "method": "deploy",
     "params": {
         "type": 1,
         "chaincodeID": {
             "path": "github.com/dellwoo2/ulcontract/register"
         },
         "ctorMsg": {
             "function": "init",
             "args": [
             ]
         },
         "secureContext": "admin"
     },
     "id": 1
  }
 `
    rsp=Invoke( jstr)
    //*******************************
    // get the CCID for policy register
    i= strings.LastIndex( string(rsp) , "message\":\"" )
    policyregister=rsp[i+10:i+138]
    fmt.Println("Register CCID="+ policyregister)

//************************************************
 jstr=`{
     "jsonrpc": "2.0",
     "method": "deploy",
     "params": {
         "type": 1,
         "chaincodeID": {
             "path": "github.com/dellwoo2/ulcontract/timer"
         },
         "ctorMsg": {
             "function": "init",
             "args": [
		"`+url+`" 
             ]
         },
         "secureContext": "admin"
     },
     "id": 1
  }
 `
   rsp=Invoke( jstr)
   //*******************************
   // get the CCID for Timer Manager
   i= strings.LastIndex( string(rsp) , "message\":\"" )
    fmt.Println("Timer CCID="+ string(rsp)[i+10:i+138])
    timerccid=rsp[i+10:i+138]

 //*********************************************

  //**********************************************
  //* Now deploy the Smart Contract for Policies 
 jstr=`{
     "jsonrpc": "2.0",
     "method": "deploy",
     "params": {
         "type": 1,
         "chaincodeID": {
             "path": "github.com/dellwoo2/ulcontract/newulc"
         },
         "ctorMsg": {
             "function": "init",
             "args": [
		"`+manager +`",
		"`+ commsmanager +`",
 		"`+ timerccid +`", 
		"`+ url +`"
             ]
         },
         "secureContext": "admin"
     },
     "id": 1
   }
 `
 rsp=Invoke( jstr)
   //*******************************
   // get the CCID for GL Manager
   i= strings.LastIndex( string(rsp) , "message\":\"" )
    fmt.Println("UL CCID="+ string(rsp)[i+10:i+138])
    return "Done"
}
func Invoke( json string  )(string){


    var jsonStr =[]byte(json)
    url="http://192.168.99.100:7050/chaincode"
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

    fmt.Println("Invoking Chain Code:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Response:", string(body))
    return  string(body)
}



func main() {

cfg, _ := ini.InsensitiveLoad("Contract.ini")
url = cfg.Section("").Key("URL1").String()
rurl = cfg.Section("").Key("REGISTER").String()
user = cfg.Section("").Key("USER").String()
secret = cfg.Section("").Key("SECRET").String()
t:=time.Now()
fmt.Println(t.String())
fmt.Println("**************** Starting Initialising***********************")
Itest()
t=time.Now()
fmt.Println("**************** End Initialising***********************")
fmt.Println(t.String())
}

