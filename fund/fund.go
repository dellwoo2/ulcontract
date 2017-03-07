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
//	"io/ioutil"
	"errors"
	"fmt"
//        "time"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
// 	"net/http" 
//	"bytes"
 	"encoding/json"
)	
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type Fnd struct{
  FundID string
  Unitvalue float64
}

var fundarray [100]Fnd
var fundmap map[string]float64

var lastupdated string

func main() {
	
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}


// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "update" {
		return t.update(stub, "init", args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}


func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to Update")
	}
        err:=json.Unmarshal([]byte(args[0]), &fundarray)
        lastupdated =args[1]
	err = stub.PutState("lastupdated" , []byte(lastupdated ))
	err = stub.PutState("fundarray" , []byte(args[0]) )
        for i:=0 ; i < 100 ; i++{
           fundmap[fundarray[i].FundID]=fundarray[i].Unitvalue
	} 
	return nil, err
}
// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	fbyt, _ := stub.GetState("fundarray" )
        json.Unmarshal(fbyt, &fundarray)
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "fund" { //read a variable
		return t.fund(stub, args)
	} 
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) fund(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var key, jsonResp string
	//var err error
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}
	val:=strconv.FormatFloat(fundmap[args[0]], 'F', -1, 64)

	return []byte(val) , nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var key, jsonResp string
	//var err error
        //fbyt, err := json.Marshal(fundarray)
	fbyt, err := stub.GetState("fundarray" )
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	return fbyt, err
}
