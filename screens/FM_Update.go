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


type Res struct{
	Status string
	Message string
 }

type Ret struct{
 Jsonrpc string
 Result Res
 Id string
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





func FMUpdate()(string){

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
             "function": "crtFndjournal",
             "args": [ 
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` 

    rsp:=Invoke( jstr)
 //*******************************
 // get the Journal ID for the GL Update 
	i:= strings.LastIndex( string(rsp) , "message\":\"" )
 	journalId:=string(rsp)[i+10:i+46]
	fmt.Println("GL Journal ID="+ journalId)
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
             "function": "getFndjournal",
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
    transactions:=  make(map[string]string)
    json.Unmarshal( []byte(ret.Result.Message), &transactions)	
    for k , v :=range transactions {
	fmt.Println(k+"-"+ v )
    }

 return "Done"
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
             "function": "getFndjournal",
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
    transactions:=  make(map[string]string)
    json.Unmarshal( []byte(ret.Result.Message), &transactions)	
    for k , v :=range transactions {
	fmt.Println(k+"-"+ v )
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
fmt.Println("**************** Starting Fund Manager Update***********************")
if len(os.Args) == 1 {
  FMUpdate()
}else{
  GetJournal( os.Args[1] )
}

t=time.Now()
fmt.Println("**************** End Fund Manager Update***********************")
fmt.Println(t.String())
}

