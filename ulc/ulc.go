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
	"os"
//	"io/ioutil"
	"encoding/json"
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
	"net/smtp"
)	

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
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
 Gender string
 Dob    string
 Smoker string
}
type Contract struct{
 Acct Account
 Product string
 StartDate string
 SumAssured string
 Term  string
 PaymentFrequency string
 Owner  string
 Beneficiary string
 Lf  Life
}

//*****************************************

var contract Contract


var count int
var   xx = shared.Args{1, 2}
func main() {
/************
	bonus:=121
	contract.account.valuation="8989.89"
	i, _ := strconv.ParseFloat( contract.account.valuation , 10);
	i = i + float64(bonus)
        contract.account.valuation= strconv.FormatFloat(i,  'f' , 2,  64)
 	fmt.Print("DE***** Contract value="+contract.account.valuation)
**************/
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
	l := log.New(os.Stderr, "", 0)
	l.Println("DE************* INIT CHAINCODE")
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
  contract.Lf.Gender=args[0]
  contract.Lf.Dob=args[1]
  contract.Lf.Smoker=args[2]
  contract.Product=args[3]
  contract.StartDate=args[4]
  contract.Term = args[5]
  contract.PaymentFrequency=args[6]
  contract.Owner=args[7]
  contract.Acct.Valuation="0"
  //var bin_buf bytes.Buffer
  //binary.Write( &bin_buf, binary.BigEndian, contract )
/*************************
	err := stub.PutState("owner", []byte(contract.owner))
	err = stub.PutState("paymentFrequency", []byte(contract.paymentFrequency))
	err = stub.PutState("startDate",  []byte(contract.startDate) )
	err = stub.PutState("product", []byte(contract.product))
	err = stub.PutState("life.smoker", []byte(contract.life.smoker))
	err = stub.PutState("life.dob", []byte(contract.life.dob))
	err = stub.PutState("life.gender",  []byte(contract.life.gender))
	err = stub.PutState("account.valuation",  []byte(contract.account.valuation))
*********************/
        b, err := json.Marshal(contract)
	err = 	stub.PutState("Contract", b)

        //fmt.Println( xx.A )

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	l := log.New(os.Stderr, "", 0)
	l.Println("DE************* Invoke Function")
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
	} else if function == "surrender" {
		return t.surrender(stub, args)
	} 
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}
func (t *SimpleChaincode) surrender(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	surr, _ := strconv.ParseFloat( args[0] , 10);
	surrchg:=20
 	fmt.Print("DE***** Contract value="+contract.Acct.Valuation)
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)

	i, _ := strconv.ParseFloat( contract.Acct.Valuation , 10);
	i = i - float64(surr)
	i = i - float64(surrchg)
        contract.Acct.Valuation= strconv.FormatFloat(i,  'f' , 2,  64)
 	log.Print("DE***** Contract value="+contract.Acct.Valuation)
        b, err := json.Marshal(contract)
	err = 	stub.PutState("Contract", b)

	return  []byte("surrendered"), err
}

func (t *SimpleChaincode) fundAllocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    	var fnd[20]Fund
        json.Unmarshal([]byte(args[0]), &fnd)
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)
        contract.Acct.Fnds=fnd
        b, err := json.Marshal(contract)
	err = 	stub.PutState("Contract", b)
	return  []byte("Funds_allocated"), err
}
func (t *SimpleChaincode) applyPremium(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	premium, _ := strconv.ParseFloat( args[0] , 10);

 	fmt.Print("DE***** Contract value="+contract.Acct.Valuation)
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)

	i, _ := strconv.ParseFloat( contract.Acct.Valuation , 10);
	i = i + float64(premium)
        contract.Acct.Valuation= strconv.FormatFloat(i,  'f' , 2,  64)
 	log.Print("DE***** Contract value="+contract.Acct.Valuation)
        b, err := json.Marshal(contract)
	err = 	stub.PutState("Contract", b)
	return  []byte("applied"), err
}
func (t *SimpleChaincode) monthlyProcessing(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


	bonus:=33

 	fmt.Print("DE***** Contract value="+contract.Acct.Valuation)
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)

	i, _ := strconv.ParseFloat( contract.Acct.Valuation , 10);
	i = i - float64(bonus)
        contract.Acct.Valuation= strconv.FormatFloat(i,  'f' , 2,  64)
 	log.Print("DE***** Contract value="+contract.Acct.Valuation)
        b, err := json.Marshal(contract)
	err = 	stub.PutState("Contract", b)
	return  []byte("processed"), err
}
func (t *SimpleChaincode) statement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    // Set up authentication information.
    auth := smtp.PlainAuth(
        "",
        "dannyellwood",
        "Fr@nkly51",
        "smtp.gmail.com",
    )
    // Connect to the server, authenticate, set the sender and recipient,
    // and send the email all in one step.
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)

str1:=`From:dannyellwood@gmail.com.org;
To: dellwoo2@csc.com
Subject: Monthly Statement;

Body: Your Statement of account:
Account Holder:`+contract.Owner+`
Value:`+contract.Acct.Valuation+`

Sincerely, Danny `


    err = smtp.SendMail(
        "smtp.gmail.com:587",
        auth,
        "dannyellwood@gmail.com.org",
        []string{"dellwoo2@csc.com" },
        []byte(str1),
    )
    if err != nil {
     fmt.Print(err)
    }
	return  []byte("Mail sent"), err
}
func (t *SimpleChaincode) valuation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)

	return  []byte("Valuation="+contract.Acct.Valuation), err
}


// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}else if function == "statement" {
		return t.statement(stub, args)
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
/*******************************
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
*****************************/
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)

//	binary.Read(buff, binary.BigEndian, &cont)
/************************
        str1:= "Product:"+contract.product+","+
               "startDate:"+contract.startDate+","+
		"owner:"+contract.owner+","+
		"DOB:"+contract.life.dob
************************/
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


	return valAsbytes, nil
}
