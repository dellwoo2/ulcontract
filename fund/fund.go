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

type ContractFunds struct
 Fmap  map[string]Funds
 ContID string
 Dte string
}
type Prices struct{
  FundID string
  UnitPrice string
  PriceDate string
}

//***********************************
//* total is the list of new funds for contracts
var totalmap map[string]ContractFunds
var newmap map[string]ContractFunds
var movementhistory  map[string]ContractFunds
//***********************************
//* Portfolio holds the existing position
var portfoliomap map[string]ContractFunds
//*************************************
//* Holds the latest prices for the funds 
var pricemap map[string]Prices

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
	newmap=make(map[string]ContractFunds)
	valAsbytes, _ := stub.GetState("Newmap")
    	json.Unmarshal(valAsbytes , &newmap)

	fmt.Println("invoke is running " + function)
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "update" {
		return t.update(stub, "init", args)
	}else if function == "setnewbalance" {
		return t.setnewbalance(stub, "init", args)
	}

	fmt.Println("invoke did not find func: " + function)
        b, err := json.Marshal(newmap)
	err = 	stub.PutState("Newmap", b)

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) setnewbalance (stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var cf ContractFund
	cf.Fmap=make(map[string]Funds)

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to Update")
	}
    	json.Unmarshal(args[1] , &cf.Fmap)
	cf.ContID=args[0]
	newmap[args[0]]=cf

	return nil, err
}
// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "fundprice" { //read a variable
		return t.fund(stub, args)
	} else if function == "movements" { //read a variable
		return t.fund(stub, args)
	} else if function == "portfolio" { //read a variable
		return t.fund(stub, args)

	} 
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) fundprice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var key, jsonResp string
	//var err error
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}
	pm:= pricemap[args[0]]

	return []byte( pm.UnitPrice ) , nil
}

func (t *SimpleChaincode) movements(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
   //var totalmap map[string]ContractFunds
   //***********************************
   //* Portfolio holds the existing position
   //var portfoliomap map[string]ContractFunds

	valAsbytes, _ := stub.GetState("Totalmap")
    	json.Unmarshal(valAsbytes , &totalmap)
	valAsbytes, _ := stub.GetState("Newmap")
    	json.Unmarshal(valAsbytes , &newmap)
	//loop round total map getting movements
	var mv Fund 
	var mvmap map[string]Fund
	mvmap=make(map[string]Fund)
          for key, newvalue := range newmap{
            fmt.Println("Key:", key, "Value:", fvalue)
	    if existingval,ok "=totalmap[key] ; ok }
                 mv.FundId=key
		 a , err :=strconv.ParseFloat(newvalue.Units , 10)
		 b , err :=strconv.ParseFloat(existingval.Units , 10)
		 mv.Units=strconv.FormatFloat(float64( a - b ),  'f' , 2,  64)
		 mv.FundId=key
		 mvmap[mv.FundId]=mv
            }else{
		 mv.Units=newvalue.Units
		 mv.FundId=key
		 mvmap[mv.FundId]=mv
	    }
	  }
        byt, err := json.Marshal(mvmap)
	err = 	stub.PutState("policies", b)
	return byt , nil
}

func (t *SimpleChaincode) portfolio(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	valAsbytes, _ := stub.GetState("Portfoliomap")
    	json.Unmarshal(valAsbytes , &portfoliomap)

	return valAsbytes , nil
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
