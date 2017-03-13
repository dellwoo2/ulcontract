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
	"encoding/json"
//	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
 	"net/http" 
//	"bytes"
//	"net/smtp"
)	
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Wallet struct{
 User string
 Policies map[string]string
}

var locked bool
var register map[string]Wallet

func main() {
	
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}


// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("***** Init Policy Register  ********* ")
	register=make(map[string]Wallet)
 	byt, _ := json.Marshal(register)
	err := stub.PutState("register", byt)
	locked=false
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("***** Invoke ********* "+ function )
	for i:=0 ; i<=100 && locked== true; i++ {
	    time.Sleep(time.Millisecond * 20 )	
	}
	locked=true
	fmt.Println("invoke is running " + function)
        //xx = shared.Args{1, 2} 
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "update" {
		return t.update(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	locked=false
	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	locked=true
	valAsbytes, err := stub.GetState("register")
  	register=make(map[string]Wallet)
  	json.Unmarshal(valAsbytes , &register)

    	fmt.Println("Update Register: user=" + args[0] +" Policy="+args[1] )
	if w, ok :=register[args[0]] ; ok {
		w.Policies[args[1]]="Y"
	}else{
	  var wa Wallet
	  wa.User=args[0]
	  wa.Policies[args[1]]="Y"
	  register[args[0]]=wa
	}

 	byt, _ := json.Marshal(register)
	err = stub.PutState("register", byt)
	fmt.Println(err)
	locked=false
	return []byte("Updated"), err
}




// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}else if function == "get" {
		return t.get(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}



func (t *SimpleChaincode) get(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	valAsbytes, err := stub.GetState("register")
  	register=make(map[string]Wallet)
  	json.Unmarshal(valAsbytes , &register)
	var w1 Wallet
	w1.Policies=make(map[string]string)
	var byt []byte
	if w, ok :=register[args[0]] ; ok {
	    	byt, err = json.Marshal( w )

        } else{
		byt, err = json.Marshal( w1 )
	}

	fmt.Println("Policy List for user "+args[0]+"="+string(byt))
	
	return  byt, err 
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
