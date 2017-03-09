/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"io/ioutil"
	"errors"
	"fmt"
        "time"
//	"strings"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
 	"net/http" 
	"bytes"
//	"net/smtp"
)	
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
var state string
var count int
var ccid string
var mailsent map[string]string
func main() {
 //    sstr:= "b93f36b5cdf0cc16f7e2f5a30c05431547ec049215dff9cfd6f4d8ef6b20cbdbffefd59b11fe538872e87a41a1471637ccc3c4c9ff4cbccfbafdf3ebc83f075a"
	
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}



func  xCC(sstr string ) {
    url :="https://e9aeb13602254217bdb0e8b425c82732-vp0.us.blockchain.ibm.com:5003/chaincode"
    //valAsbytes, err := stub.GetState("CCID")
    //ccstr:=string(valAsbytes)
    jsonStr := []byte( `
  {
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "` +sstr+ `"
         },
         "ctorMsg": {
             "function": "schedule",
             "args": [
                 "dtest",
                 "I am Here XXXX"
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    fmt.Println(err)
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    fmt.Println(err)
    if err != nil {
        panic(err)
    }
    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
    defer resp.Body.Close()

}



// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("***** Init Scheduler ********* " )
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	mailsent=make(map[string]string)
        count=0;
	state="inactive"
	err := stub.PutState("CCID", []byte(args[0]) )
	err = stub.PutState("state", []byte(state) )
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
        //xx = shared.Args{1, 2} 
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "activate" {
		return t.activate(stub, args)
	}else if function == "deactivate" {
		return t.deactivate(stub, args)
	}else if function == "kill" {
		return t.kill(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}
func (t *SimpleChaincode) activate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	state="active"
	err := stub.PutState("state", []byte(state))
	return []byte("Activated"), err
}
func (t *SimpleChaincode) deactivate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	state="inactive"
	err := stub.PutState("state", []byte(state))
	return []byte("De Activated"), err
}
func (t *SimpleChaincode) kill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	state="kill"
	err := stub.PutState("state", []byte(state))
	return []byte("scheduler killed"), err
}
// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "schedule" {
		go t.schedule(stub, args)
		return nil, nil 

	} 
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// Schedule - query function to call invoke methods on contract
func (t *SimpleChaincode) schedule(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var key, value string
	var err error
	var tx int64
	fmt.Println("running set scheduler")
	if len(args) != 3 {
		tx=300
	}else{
		tx, _ = strconv.ParseInt( args[2] , 10 , 64);
        }
	sbytes, err := stub.GetState("state")
	state= string(sbytes[:]) 
        for i := 0; i < 1000000 && state!="kill"; i++ {
		time.Sleep(time.Duration( tx )*time.Second  )
    		fmt.Print("Timer Iteration=")
    		fmt.Println(i)
		if state=="active" { 
			t.callCC(stub , args )
			t.callDD(stub , args )
		}
	}
	return []byte(state) , err
}

func (t *SimpleChaincode) callCC(stub shim.ChaincodeStubInterface , args []string) {
    //url :="https://e9aeb13602254217bdb0e8b425c82732-vp0.us.blockchain.ibm.com:5003/chaincode"
    url :=args[1]
    //valAsbytes, err := stub.GetState("CCID")
    //ccstr:=string(valAsbytes)
    jsonStr := []byte( `
  {
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "` + args[0]+ `"
         },
         "ctorMsg": {
             "function": "schedule",
             "args": [
                 "dtest",
                 "I am Here XXXX"
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }`)
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
    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
    defer resp.Body.Close()

}



func (t *SimpleChaincode) callDD(stub shim.ChaincodeStubInterface , args []string) {
    //url :="https://e9aeb13602254217bdb0e8b425c82732-vp0.us.blockchain.ibm.com:5003/chaincode"
    url :=args[1]
    //valAsbytes, err := stub.GetState("CCID")
    //ccstr:=string(valAsbytes)
    jsonStr := []byte( `
{
     "jsonrpc": "2.0",
     "method": "query",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "` + args[0]+ `"
         },
         "ctorMsg": {
             "function": "statement",
             "args": [
                 "Contract"
             ]
         },
         "secureContext":"admin"
     },
     "id": 2
 }
 `)
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
    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
    defer resp.Body.Close()

}





// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var key, jsonResp string
	//var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	//key = args[0]
	//valAsbytes, err := stub.GetState(key)
//	x:=int64(count)
//        str1:=strconv.FormatInt(x,10)
//	valAsbytes:=[]byte(str1)
	resp, _ := http.Get("http://www.bbc.com")
  	bb, _ := ioutil.ReadAll(resp.Body)
	valAsbytes:=bb[0:50]
	//if err != nil {
	//	jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
	//	return nil, errors.New(jsonResp)
	//}


	return valAsbytes, nil
}
