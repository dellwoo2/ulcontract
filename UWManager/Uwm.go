/*
 * File: Uwm.go
 * Date: 01 June 2017
 * Author: Ellwood Technology Consulting
 * 
 * Copyright (2017) Ellwood Technology Consultimg Sdn Bhd, all rights reserved. 
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Ellwood Technology Consultimg Sdn Bhd,
 * The intellectual and technical concepts contained
 * herein are proprietary to Ellwood Technology Consultimg Sdn Bhd
 * and its suppliers  and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Ellwood Technology Consultimg Sdn Bhd.
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

type Follow_up struct{
  Typ string
  LastUpdated string
  Stat string

}

type UWcase struct{
 ContID string
 CCID string
 Ref string
 Followups map[string]Follow_up
 UWStat string
}
var state string
var UWcases map[string]UWcase
var locked bool

func main() {
 //    sstr:= "b93f36b5cdf0cc16f7e2f5a30c05431547ec049215dff9cfd6f4d8ef6b20cbdbffefd59b11fe538872e87a41a1471637ccc3c4c9ff4cbccfbafdf3ebc83f075a"
	
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}


// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("***** Init Underwriting Manager ********* ")
	UWcases=make(map[string]UWcase)
 	byt, _ := json.Marshal(UWcases)
	err := stub.PutState("UWcases", byt)
	state="active"
	err = stub.PutState("state",[]byte(state))
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
	}else if function == "StartCase" {
		return t.StartCase(stub, args)
	}else if function == "UpdateCase" {
		return t.UpdateCase(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)
	locked=false
	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) StartCase(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) UpdateCase(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}




// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}else if function == "GetStatus" {
		return t.GetStatus(stub, args)
	}else if function == "GetCase" {
		return t.GetCase(stub, args)
	}else if function == "Schedule" {
		return t.schedule(stub, args)
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
		//***************************************
		//* Scheduler payload here 
		//* Iterate over cases - Process new ones and ones with 
		//* a timed out status 
	}
	return []byte(state) , err
}


func (t *SimpleChaincode) GetCase(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	valAsbytes, err := stub.GetState("UWcases")
	UWcases=make(map[string]UWcase)
 	json.Unmarshal(valAsbytes , &UWcases)
	if val, ok :=UWcases[args[0]] ; ok {
	  valAsbytes, _ = json.Marshal(val)
	}	
	return  valAsbytes, err 
}
func (t *SimpleChaincode) GetStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

  return nil , nil 
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
