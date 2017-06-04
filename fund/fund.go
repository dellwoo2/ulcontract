/*
 * File: Funds.go
 * Date: 1 June 2017
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
       "time"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
// 	"net/http" 
//	"bytes"
 	"encoding/json"
)	
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

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
var locked bool
func main() {
	
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}


// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("***** Init Fund Manager ********* ")
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "updateFunds" {
		return t.update(stub, "init", args)
	}else if function == "crtjournal" {
		return t.crtjournal(stub, "init", args)
	}else if function == "setnewbalance" {
		return t.setnewbalance(stub, "init", args)
	}

	fmt.Println("invoke did not find func: " + function)


	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) setnewbalance (stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	for i:=0 ; i<=100 && locked== true; i++ {
	    time.Sleep(time.Millisecond * 20 )	
	}
	locked=true
	var cf ContractMovement
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
func (t *SimpleChaincode) crtjournal(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
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

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "fundprice" { //read a variable
		return t.fundprice(stub, args)
	} else if function == "getjournal" { //pass back a journal
		return t.getjournal(stub, args)
	} 
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) getjournal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

/******************************************************************
func (t *SimpleChaincode) movements(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
   //var totalmap map[string]ContractFunds
   //***********************************
   //* Portfolio holds the existing position
   //var portfoliomap map[string]ContractFunds

	valAsbytes, _ := stub.GetState("Totalmap")
    	json.Unmarshal(valAsbytes , &totalmap)
	valAsbytes, _ = stub.GetState("Newmap")
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
*************************************************/
