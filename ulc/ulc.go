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
//	"encoding/json"
	"errors"
	"fmt"
//        "time"
	"log"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/dellwoo2/ulcontract/ulc/shared"
 //	"net/http" 
 //   	"encoding/binary"
 //  	"bytes"
)	
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//*****************************************
//* Contract Types

type Fund struct{
 fundId string
 units  string
}
type Account struct{
  funds [100]Fund
  lastvaluationDate string
  valuation string
}
type Life struct{
 gender string
 dob    string
 smoker string
}
type Contract struct{
 account Account
 product string
 startDate string
 term  string
 paymentFrequency string
 owner  string
 beneficiary string
 life  Life
}

//*****************************************

var contract Contract


var count int
var   xx = shared.Args{1, 2}
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}
        count=0;
 	//xx = &shared.Args{1, 2} 
/****
 gender 0
 dob    1
 smoker 2
 product 3
 startDate 4
 Term  int64 5
 PaymentFrequency 
 Owner  7
 Beneficiary 8
**********/
  //contract=new(Contract)
  //contract.account=new(Account)
  //contract.life=new(Life)
  contract.life.gender=args[0]
  contract.life.dob=args[1]
  contract.life.smoker=args[2]
  contract.product=args[3]
  contract.startDate=args[4]
  contract.term = args[5]
  contract.paymentFrequency=args[6]
  contract.owner=args[7]
  contract.account.valuation="0"
  //var bin_buf bytes.Buffer
  //binary.Write( &bin_buf, binary.BigEndian, contract )
	err := stub.PutState("owner", []byte(contract.owner))
	err = stub.PutState("paymentFrequency", []byte(contract.paymentFrequency))
	err = stub.PutState("startDate",  []byte(contract.startDate) )
	err = stub.PutState("product", []byte(contract.product))
	err = stub.PutState("life.smoker", []byte(contract.life.smoker))
	err = stub.PutState("life.dob", []byte(contract.life.dob))
	err = stub.PutState("life.gender",  []byte(contract.life.gender))
	err = stub.PutState("account.valuation",  []byte(contract.account.valuation))

        //fmt.Println( xx.A )

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
        log.Print("DE***** Invoke Function")
        //xx = shared.Args{1, 2} 
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args )
	} else if function == "fundAllocation" {
		return t.fundAllocation(stub, args)
	} else if function == "applyPremium" {
		return t.applyPremium(stub, args)
	} else if function == "schedule" {
		return t.monthlyProcessing(stub, args)
	} else if function == "valuation" {
		return t.valuation(stub, args)
	}


	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}
func (t *SimpleChaincode) fundAllocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	return  []byte("allocated"), nil
}
func (t *SimpleChaincode) applyPremium(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return  []byte("applied"), nil
}
func (t *SimpleChaincode) monthlyProcessing(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	bonus:=121
	valAsbytes, _ := stub.GetState("account.valuation")
	contract.account.valuation= string(valAsbytes[:]) 
	i, _ := strconv.ParseInt( contract.account.valuation , 10, 64);
	i = i + int64(bonus)
        contract.account.valuation= strconv.FormatInt(i, 64)
 	log.Print("DE***** Contract value="+contract.account.valuation)
	stub.PutState("account.valuation",  []byte(contract.account.valuation))
	return  []byte("processed"), nil
}
func (t *SimpleChaincode) valuation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	valAsbytes, _ := stub.GetState("account.valuation")
	contract.account.valuation= "Valuation="+string(valAsbytes[:]) 


	return  []byte(contract.account.valuation), nil
}


// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "valuation" {
		return t.valuation(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	xx = shared.Args{1, 2} 
	fmt.Println("Writing in Invoke DE********************")
	fmt.Println(xx.A)
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
/****************
        for i := 0; i < 1000; i++ {
		time.Sleep(time.Second * 5)
		count++
	}
******************/
	return  []byte(value), nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key , jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]

	valAsbytes, err := stub.GetState("owner")
	contract.owner= string(valAsbytes[:]) 

	valAsbytes, err = stub.GetState("paymentFrequency")
	contract.paymentFrequency= string(valAsbytes[:]) 

	valAsbytes, err = stub.GetState("startDate")
	contract.startDate= string(valAsbytes[:]) 

	valAsbytes, err = stub.GetState("product")
	contract.product= string(valAsbytes[:]) 

	valAsbytes, err = stub.GetState("life.smoker")
	contract.life.smoker= string(valAsbytes[:]) 

	valAsbytes, err = stub.GetState("life.dob")
	contract.life.dob= string(valAsbytes[:]) 

	valAsbytes, err = stub.GetState("life.gender")
	contract.life.gender= string(valAsbytes[:]) 

//	 b, err := json.Marshal(contract)
//	str1:=string(b)
 //       buff:=bytes.NewBuffer(valAsbytes)
 //       var cont Contract
//	binary.Read(buff, binary.BigEndian, &cont)
        str1:= "Product:"+contract.product+","+
               "startDate:"+contract.startDate+","+
		"owner:"+contract.owner+","+
		"DOB:"+contract.life.dob

//	x:=int64(count)
//        str1:=strconv.FormatInt(x,10)
//	valAsbytes:=[]byte(str1)
//	resp, _ := http.Get("http://www.bbc.com")
//  	bb, _ := ioutil.ReadAll(resp.Body)
//	valAsbytes:=bb[0:50]
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}


	return []byte(str1), nil
}
