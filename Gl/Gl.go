/*
 * File: Gl.go
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

func main() {
 //    sstr:= "b93f36b5cdf0cc16f7e2f5a30c05431547ec049215dff9cfd6f4d8ef6b20cbdbffefd59b11fe538872e87a41a1471637ccc3c4c9ff4cbccfbafdf3ebc83f075a"
	
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}


// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("***** Init GL Manager ********* ")
	gltran=make(map[string]GLtran)
 	byt, _ := json.Marshal(gltran)
	err := stub.PutState("gltran", byt)
        count=0;
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
	}else if function == "updateJ" {
		return t.updateJ(stub, args)
	}else if function == "updateT" {
		return t.updateT(stub, args)
	}else if function == "Crtjournal" {
		return t.Crtjournal(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	locked=false
	return nil, errors.New("Received unknown function invocation: " + function)
}

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

func (t *SimpleChaincode) Crtjournal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	}else if function == "Getjournal" {
		return t.Getjournal(stub, args)
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
