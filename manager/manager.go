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
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
 	"net/http" 
//	"bytes"
//	"net/smtp"
)	
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//*****************************************
// Fund Manager types
type ContractMovement struct{
 Fundunits map[string]string
 ContID string
 Dte string
}
type Prices struct{
  FundID string
  UnitPrice string
  PriceDate string
}
type Position struct{
 dte string
 portfolio map[string]string
}
//***********************************
//* holds all movements not yet applied
var pendingmovements map[string]string

//***********************************
//*  holds a dated history of the fund portfolio position 
var movementhistory  map[string]Position
var contractmovements map[string]ContractMovement 
//***********************************
//* Portfolio holds the existing position
var portfoliomap map[string]string
//*************************************
//* Holds the latest prices for the funds 
var pricemap map[string]Prices

var lastupdated string

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
 Name string
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
 Status string
 Email string
 UWstatus string
}
type Ods struct{
 Cont Contract
 Tranid string
 Posted string
}



type GLtran struct{
 TranID string
 Ref string
 Dbacc string
 Db    string
 Cracc string
 Cr string
 Stat string
}

var state string
var count int
var ccid string
var gltran map[string]GLtran
var locked bool
var transactions map[string]Ods

func main() {
 //    sstr:= "b93f36b5cdf0cc16f7e2f5a30c05431547ec049215dff9cfd6f4d8ef6b20cbdbffefd59b11fe538872e87a41a1471637ccc3c4c9ff4cbccfbafdf3ebc83f075a"
	
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}


// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("***** Init Combined ODS GL & FUND Manager ********* ")
	gltran=make(map[string]GLtran)
 	byt, _ := json.Marshal(gltran)
	err := stub.PutState("gltran", byt)
        count=0;
	locked=false

	transactions=make(map[string]Ods)
 	byt, _ = json.Marshal(transactions)
	err = stub.PutState("transactions", byt)

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
	}else if function == "updateJ" {
		return t.updateJ(stub, args)
	}else if function == "updateT" {
		return t.updateT(stub, args)
	}else if function == "CrtGljournal" {
		return t.CrtGljournal(stub, args)
	}else if function == "updateOds" {
		return t.updateOds(stub, args)
	}else if function == "CrtOdsjournal" {
		return t.CrtOdsjournal(stub, args)
	}else if function == "updateFunds" {
		return t.updateFunds(stub, "init", args)
	}else if function == "crtFndjournal" {
		return t.crtFndjournal(stub, "init", args)
	}else if function == "setnewFndbalance" {
		return t.setnewFndbalance(stub, "init", args)
	}
	fmt.Println("invoke did not find func: " + function)
	locked=false
	return nil, errors.New("Received unknown function invocation: " + function)
}

//*******************************************************************************
//*    	Fund Manager Invoke methods              
func (t *SimpleChaincode) setnewFndbalance (stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) updateFunds(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	for i:=0 ; i<=100 && locked== true; i++ {
	    time.Sleep(time.Millisecond * 20 )	
	}
	locked=true
	var cf ContractMovement
	contractmovements=make( map[string]ContractMovement) 
	cf.ContID=args[0]

	fu:=make( map[string]string )
    	json.Unmarshal([]byte(args[1]) , &fu)
	cf.Fundunits= fu
 	valAsbytes, _ := stub.GetState("contractmovements")
	json.Unmarshal(valAsbytes , &contractmovements)
	contractmovements[stub.GetTxID()]=cf
        b, err := json.Marshal(contractmovements)
	err = 	stub.PutState("contractmovements", b)
        locked=false
	return nil, err
}


//**************************************************
//* get the latest movements 
func (t *SimpleChaincode) crtFndjournal(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	for i:=0 ; i<=100 && locked== true; i++ {
	    time.Sleep(time.Millisecond * 20 )	
	}
	locked=true
	valAsbytes, _ := stub.GetState("contractmovements")
	var cm map[string]ContractMovement 

	//**************************************
	//* Save the curent movements as an intermediate
	//* for safe keeping in case a cockup somehow
	cm=make(map[string]ContractMovement )
 	json.Unmarshal(valAsbytes , &cm)
	err := 	stub.PutState("cmJournal", valAsbytes)

	var mvmap map[string]string
	mvmap= make( map[string]string)

	//******************************
	//* Zero off the contract movements 
	contractmovements= make( map[string]ContractMovement)	
        b, err := json.Marshal(contractmovements)
	err = 	stub.PutState("contractmovements", b)
	//********************************
	// we can safely unlock it now
 	locked=false
        //create an aggregated movement map
	for k , v := range cm {
		for  mk ,fu := range v.Fundunits {
		  if vv, ok := mvmap[mk] ; ok {
		  	x, _ := strconv.ParseFloat( fu , 10);
			y, _ := strconv.ParseFloat( vv , 10);
			nmv:= strconv.FormatFloat( x+y ,  'f' , 2,  64)  
			mvmap[k]=nmv
		   }else{
		  	x, _ := strconv.ParseFloat( fu , 10);
			nmv:= strconv.FormatFloat( x ,  'f' , 2,  64)  	
		   	mvmap[k]=nmv	
		   }
		  //************************
		  //* update the movements map

                }		

	}
	//*************************************
	// Save the Mvmap keyed on the transaction ID
        valAsbytes, err = json.Marshal(mvmap)
	journalid:=stub.GetTxID()
	err = 	stub.PutState( journalid, valAsbytes)

  return []byte(journalid), err
}

//*******************************************************************************
//*                  ODS Methods
func (t *SimpleChaincode) updateOds(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	locked=true
	valAsbytes, err := stub.GetState("transactions")
  	transactions=make(map[string]Ods)
  	json.Unmarshal(valAsbytes , &transactions)

	var ods Ods
  	json.Unmarshal([]byte(args[0]) , &ods)
    	fmt.Println("Update ODS:" + string(args[0]) )
	transactions[ods.Tranid]=ods
 	byt, _ := json.Marshal(transactions)
	err = stub.PutState("transactions", byt)
	fmt.Println(err)
	locked=false
	return []byte("ODS Updated"), err
}

func (t *SimpleChaincode) CrtOdsjournal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	locked=true
	valAsbytes, err := stub.GetState("transactions")
  	transactions=make(map[string]Ods)
  	json.Unmarshal(valAsbytes , &transactions)

	var journal map[string]Ods
	journal=make(map[string]Ods)
 	for key, value := range transactions {
		if value.Posted=="N" {
	    		fmt.Println("ODS JournalKey:", key, "ODS Journal Value:", value)
			value.Posted="Y"
			transactions[key]=value
			journal[key]=value
		}
    	}
 	byt, _ := json.Marshal(transactions)
	err = stub.PutState("transactions", byt)
 	jb, _ := json.Marshal(journal)
	locked=false
	invokeTran:=stub.GetTxID()
	err = stub.PutState(invokeTran, jb)

	return jb , err
}



//********************************************************************************

func (t *SimpleChaincode) updateT(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	locked=true
	valAsbytes, err := stub.GetState("gltran")
  	gltran=make(map[string]GLtran)
  	json.Unmarshal(valAsbytes , &gltran)

	var gt GLtran
  	json.Unmarshal([]byte(args[0]) , &gt)
    	fmt.Println("Update Tran:" + string(args[0]) )
	gltran[gt.TranID]=gt
 	byt, _ := json.Marshal(gltran)
	err = stub.PutState("gltran", byt)
	fmt.Println(err)
	locked=false
	return []byte("Updated"), err
}

func (t *SimpleChaincode) updateJ(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	locked=true
	valAsbytes, err := stub.GetState("gltran")
  	gltran=make(map[string]GLtran)
  	json.Unmarshal(valAsbytes , &gltran)

	var journal map[string]GLtran
	journal=make(map[string]GLtran)
  	json.Unmarshal([]byte(args[0]) , &journal)
 	for key, value := range journal {
	    fmt.Println("UpdateKey:", key, "Update Value:", value)
		value.Stat="N"
		gltran[key]=value
    	}
 	byt, _ := json.Marshal(gltran)
	err = stub.PutState("gltran", byt)
	fmt.Println(err)
	locked=false
	return []byte("Updated"), err
}

func (t *SimpleChaincode) CrtGljournal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	locked=true
	valAsbytes, err := stub.GetState("gltran")
  	gltran=make(map[string]GLtran)
  	json.Unmarshal(valAsbytes , &gltran)

	var journal map[string]GLtran
	journal=make(map[string]GLtran)
 	for key, value := range gltran {
		if value.Stat=="N" {
	    		fmt.Println("JournalKey:", key, "Journal Value:", value)
			value.Stat="Y"
			gltran[key]=value
			journal[key]=value
		}
    	}
 	byt, _ := json.Marshal(gltran)
	err = stub.PutState("gltran", byt)
 	jb, _ := json.Marshal(journal)
	invokeTran:=stub.GetTxID()
	err = stub.PutState(invokeTran, jb)
	locked=false
	return jb , err
}


// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}else if function == "GetGljournal" {
		return t.Getjournal(stub, args)
	}else if function == "GetOdsjournal" {
		return t.GetOdsjournal(stub, args)
	} else if function == "fundprice" { //read a variable
		return t.fundprice(stub, args)
	} else if function == "getFndjournal" { //pass back a journal
		return t.getFndjournal(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


//*********************************************************************************
//* Fund Manager Query Methods 
func (t *SimpleChaincode) getFndjournal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	valAsbytes, err := stub.GetState(args[0])
	return valAsbytes , err
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



//*********************************************************************************
//* ODS Query Methods 

func (t *SimpleChaincode) GetOdsjournal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	valAsbytes, err := stub.GetState(args[0])
	
	return  valAsbytes, err 
}


//**********************************************************************************
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
		}
	}
	return []byte(state) , err
}

func (t *SimpleChaincode) Getjournal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	valAsbytes, err := stub.GetState(args[0])
	
	return  valAsbytes, err 
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
