/*
 * File: Portfolio.go
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
