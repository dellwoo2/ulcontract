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
 //       "time"
	"strings"
//	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
 	"net/http" 
//	"bytes"
	"net/smtp"
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





// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("***** Init Comms  ********* " )
	mailsent=make(map[string]string)
        count=0;
	state="inactive"
	err := stub.PutState("state", []byte(state) )
	return nil, err
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("***** Invoke ********* "+ function )
	fmt.Println("invoke is running " + function)
        //xx = shared.Args{1, 2} 
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "kill" {
		return t.kill(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
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
	} else if function == "mailto" {
		return t.mailto(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


func (t *SimpleChaincode) mailto(stub shim.ChaincodeStubInterface,  args []string ) ([]byte, error) {
    if( mailsent[args[0]]=="Y"){
	return  []byte("Mail not sent"), nil
    }

    mailsent[args[0]]="Y"
    // Set up authentication information.
    auth := smtp.PlainAuth(
        "",
        "dannyellwood",
        "Fr@nkly51",
        "smtp.gmail.com",
    )
    // Connect to the server, authenticate, set the sender and recipient,
    // and send the email all in one step.

str1:=`From:dannyellwood@gmail.com.org;
To: `+args[3]+ ` 
Subject: `+ args[1]+ "\n" + strings.Replace( args[2] , "#N","\n",-1) 
fmt.Println("MAIL STR="+str1)
    err := smtp.SendMail(
        "smtp.gmail.com:587",
        auth,
        "dannyellwood@gmail.com.org",
        []string{ args[3] },
        []byte(str1),
    )
    if err != nil {
    fmt.Println("Emailing Error")
     fmt.Print(err)
    }
    return  []byte("Mail sent"), err
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
