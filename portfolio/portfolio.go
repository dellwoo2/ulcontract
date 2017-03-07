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


type Fund struct{
 FundId string
 Units  string
}
var contracts [1000]string

var fundarray [100]Fund

var positionmap map[string]int64
var requirementmap map[string]int64
var buysell [100]Fund

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
    	var fnd [20]Fund
	var newpositionmap map[string]int64
	var  contractfunds  string
	for i:=0 ; i < 1000 ; i++ {
	   //**************************
	   //* call contract for fund position
	   contractfunds =getContractFunds(contracts[i]) [
           json.Unmarshal([]byte(contractfunds), &fnd)	   
	   for( j :=0 ; j , 20 ; j++ {
	    // add to fund list
	    newpositionmap[fnd[j].FundId]=newpositionmap[fnd[j].FundId]+ newpositionmap[fnd[j].Units  	
	   }
	}
	positionmap=newpositionmap 

	return nil, err
}
// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "buysell" { //read a variable
		return t.fund(stub, args)
	} else if function == "requiredPosition" { //read a variable
		return t.fund(stub, args)
	} else if function == "position" { //read a variable
		return t.fund(stub, args)
	} 


	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) buysell(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
        b, err := json.Marshal(buysell)
	return b , nil
}

func (t *SimpleChaincode) position(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return []byte(val) , nil
}
func (t *SimpleChaincode) requiredPosition(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	return []byte(val) , nil
}


// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	return nil, nil
}
