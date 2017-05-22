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
	"time"
	"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
 //       "time" 
	"strings"
	"log"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/dellwoo2/ulcontract/newulc/shared"
	"net/http" 
 //   	"encoding/binary"
  	"bytes"
)	

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
//*****************************************
//* Contract scheduler 
var scheduler string
//*****************************************
//* Contract Types
/***
type Fund struct{
 FundId string
 Units  string
}
****/
type Account struct{
  Fnds map[string]string
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
 ContID string
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
type GLtran struct{
 TranID string
 Ref string
 Dbacc string
 Db    string
 Cracc string
 Cr string
 Stat string
}
type History struct{
 Methd string
 Funct string
 Cont Contract
 Args []string
 Tranid string 
 Dte string;
 EndValue string
}
type Ods struct{
 Cont Contract
 Tranid string 
 Posted string
}
type Policy struct{
	Cont Contract
	Hist map[string]History
}

//var history map[string]History
var gltran map[string]GLtran
var policies map[string]string
var lock map[string]string
//*****************************************


var count int
var   xx = shared.Args{1, 2}
var invokeTran string
var url string
var manager string
var commsmanager string
var RFC3339    string = "2006-01-02T15:04:05Z07:00"
func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 " )
	}
        count=0;
	l := log.New(os.Stderr, "", 0)
	l.Println("*************INIT CHAINCODE Unit Linked ****************")
	manager=args[0]
	commsmanager=args[1]
	scheduler=args[2]
	url=args[3]
	err := 	stub.PutState("scheduler",[]byte(scheduler) )
	err = 	stub.PutState("url",[]byte(url) )
//	err = 	stub.PutState("glmanager",[]byte(glmanager) )
//	err = 	stub.PutState("odsmanager",[]byte(odsmanager) )
	err = 	stub.PutState("commsmanager",[]byte(commsmanager) )
//	err = 	stub.PutState("fundmanager",[]byte(fundmanager) )
	err = 	stub.PutState("manager",[]byte(manager) )
        //fmt.Println( xx.A )
	if err != nil {
		return nil, err
	}
	return nil, err
}


//****************************************
//* get the transaction time
func (t *SimpleChaincode) TransactionTime( stub shim.ChaincodeStubInterface , tranid string)( string){
	txtime, _ := stub.GetState("TIME_"+tranid)
	if txtime == nil { 
	  txtime=[]byte(time.Now().Format(RFC3339))
	}
        stub.PutState( "TIME_"+tranid ,txtime )
   return ( string( txtime) ) 

}


//***********************************************
//* Create a newpolicy
func (t *SimpleChaincode) NewPolicy(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
	fmt.Println("************* Unit Linked  New Policy****************")
	if len(args) != 11 {
	 	fmt.Println("Incorrect number of arguments. Expecting 11")
		return nil, errors.New("Incorrect number of arguments. Expecting 11")
	}
        count=0;
	fmt.Println("Policy ID="+stub.GetTxID())
  	var policy Policy 
	if policies == nil {
	  policies=make(map[string]string)
	  //*****************************************
	  // get the mal listing all policies  
	  vb, _ := stub.GetState("policies")
    	  json.Unmarshal(vb , &policies)
         }
	//*****************************************************
	//* set contract number to transaction id for now
	policy.Cont.ContID=stub.GetTxID()

  	policy.Cont.Lf.Gender=args[0]
  	policy.Cont.Lf.Dob=args[1]
  	policy.Cont.Lf.Smoker=args[2]
  	policy.Cont.Product=args[3]
  	policy.Cont.StartDate=args[4]
  	policy.Cont.Term = args[5]
  	policy.Cont.PaymentFrequency=args[6]
  	policy.Cont.Owner=args[7]
  	policy.Cont.Lf.Name=args[8]
  	policy.Cont.Email=args[9]
  	policy.Cont.SumAssured=args[10]
  	policy.Cont.Acct.Valuation="0"
  	policy.Cont.Status="PR"

	// set to ready for now till UW contract is implemented
  	policy.Cont.UWstatus="Ready"
	fmt.Println("Creating New Policy for :"+ policy.Cont.ContID )
	if _ ,ok:=policies[policy.Cont.ContID] ; ok {
		fmt.Println("Contract Already Exist ")
		return nil, errors.New("Contract exists already")
        }

       //**************************************************
       // save the history
	//year, _ , day := time.Now().Date()
        //month:=time.Now().Month()
        //hour:=time.Now().Hour()
        //min:=time.Now().Minute()
        //sec:=time.Now().Second()
	//dte:=strconv.Itoa(day)+"/"+strconv.Itoa(int(month))+"/"+strconv.Itoa(year)+":"+strconv.Itoa(hour)+":"+strconv.Itoa(min)+":"+strconv.Itoa(sec)
        dte:=  t.TransactionTime( stub, stub.GetTxID() )
	policy.Hist=make(map[string]History)
	var h History
	h.Methd="deploy"
	h.Funct="init"
	h.Tranid=stub.GetTxID()  //time.Now().String()
	h.Cont=policy.Cont
	h.Args=args
	h.Args[0]="0"
        h.Dte=dte
	policy.Hist[dte]=h

	//************************************************
	//* Funds 
	policy.Cont.Acct.Fnds=make(map[string]string)

	policies[policy.Cont.ContID]=policy.Cont.Status
        b, err := json.Marshal(policy)
	err = 	stub.PutState(policy.Cont.ContID, b)
	//*****************************************
	//* Save the stateof  the policies map
        b, err = json.Marshal(policies)
	err = 	stub.PutState("policies", b)

	if err != nil {
		return nil, err
	}
	t.welcome(stub, policy)
	return []byte("Policy Added"), err
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	invokeTran=stub.GetTxID()
	fmt.Println("DE************* Invoke Function "+ function )

	//***********************************************
	// process contract independant functions 
	if function == "init" {
		return  t.Init(stub, "init", args)
	} else if function == "schedule" {
		return t.monthlyProcessing(stub, args )
	} else if function == "activate" {
		return t.activate(stub, args  )
	} else if function == "deactivate" {
		return t.deactivate(stub, args )
	} else if function == " setJournalDone" {
		return t. setJournalDone(stub, args)
	} else if function == "NewPolicy" {
		return t. NewPolicy(stub, args)
	}
	//**************************************
	// all remaining functions require a contract number in args[0]
	fmt.Println("invoke for policy " + args[0])
	var policy Policy
	policy.Hist=make(map[string]History)
	//*****************************************
	// get Contract state 
	valAsbytes, _ := stub.GetState(args[0])
    	json.Unmarshal(valAsbytes , &policy)

        dte:=  t.TransactionTime(stub, stub.GetTxID() )
	var h History
	h.Methd="invoke"
	h.Funct=function
	h.Tranid=invokeTran
	h.Cont=policy.Cont
	h.Args=args
        h.Dte=dte


        var err error

        //xx = shared.Args{1, 2} 
	// Handle different functions

	 if function == "fundAllocation" {
                h.Args[0]="0"
		policy , err = t.fundAllocation(stub, args, policy)
	} else if function == "applyPremium" {
                h.Args[0]=args[1]
		policy,err = t.applyPremium(stub, args, policy)
	} else if function == "surrender" {
		h.Args[0]=args[1]
		policy,err = t.surrender(stub, args , policy )
	}else{ 
		fmt.Println("invoke did not find func: " + function)
		err=errors.New("Received unknown function invocation: " + function)
        }
	if  policy.Cont.ContID=="" {
		return nil,nil
        }
	h.EndValue=policy.Cont.Acct.Valuation
	policy.Hist[dte]=h
	//*************************
	//* if we are here then we need to update the ODS
	var ods Ods
	ods.Cont=policy.Cont
	ods.Tranid=invokeTran
	ods.Posted="N"
	Odsupdate(stub , ods , invokeTran ) 

        //*****************************
        // save policy state 
        b, err := json.Marshal(policy)
	err = 	stub.PutState(policy.Cont.ContID, b)

	//**********************************
	//* Save Policies map with new policy status
        policies=make(map[string]string )
	valAsbytes, _ = stub.GetState("policies")
    	json.Unmarshal(valAsbytes , &policies)
	
        policies[policy.Cont.ContID]=policy.Cont.Status

        b, err = json.Marshal(policies)
	err = 	stub.PutState("policies", b)

	if err != nil {
		return nil, err
	}

	return nil , err 
}

func (t *SimpleChaincode) setJournalDone(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  valAsbytes, err := stub.GetState("gltran")
  gltran=make(map[string]GLtran)
  json.Unmarshal(valAsbytes , &gltran)
  for key, value := range gltran {
    fmt.Println("Key:", key, "Value:", value)
    if value.Stat=="N" {
	value.Stat="Y"
    }
  }
  b, err := json.Marshal(gltran)
  stub.PutState("gltran", b)
  return b, err

}
func (t *SimpleChaincode) journal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var journal map[string]GLtran
  journal=make(map[string]GLtran)
  valAsbytes, err := stub.GetState("gltran")
  gltran=make(map[string]GLtran)
  json.Unmarshal(valAsbytes , &gltran)
  for key, value := range gltran {
    fmt.Println("Key:", key, "Value:", value)
    if value.Stat=="N" {
	journal[key]=value
	value.Stat="Y"
    }
  }
  byt, _ := json.Marshal(journal)
  return byt, err
}



func (t *SimpleChaincode) activate(stub shim.ChaincodeStubInterface, args []string ) ([]byte, error) {
	valAsbytes, err := stub.GetState("scheduler")
	scheduler=string(valAsbytes)

	valAsbytes, err = stub.GetState("url")
	url=string(valAsbytes)
	stub.PutState("ccid" , []byte(args[0]))

	 var jsonStr = []byte( `{
   	  "jsonrpc": "2.0",
    	 "method": "invoke",
    	 "params": {
      	   "type": 1,
     	    "chaincodeID": {
      	       "name":"`+scheduler+`"
         },
         "ctorMsg": {
             "function": "activate",
             "args": [
                 "`+args[0]+`"
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` )
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err2 := client.Do(req)
    err=err2
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Set Activate Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Set Scheduler Body:", string(body))
    return []byte("Scheduler Activated"),err
}

func (t *SimpleChaincode) deactivate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	valAsbytes, err := stub.GetState("scheduler")
	scheduler=string(valAsbytes)
	return []byte("Scheduler De-activated"),err
}

func (t *SimpleChaincode) surrender(stub shim.ChaincodeStubInterface, args []string , policy Policy ) ( Policy, error) {
	var contract Contract=policy.Cont
	// surrender amount is args[1]
	surr, _ := strconv.ParseFloat( args[1] , 10);
	surrchg:=20
 	fmt.Print("DE***** Contract value="+contract.Acct.Valuation)
	i, err := strconv.ParseFloat( contract.Acct.Valuation , 10);
	i = i - float64(surrchg)
	if i > float64(surr) {
		i = i - float64(surr)
        }else{
         surr=i
         i=0
        }
        contract.Acct.Valuation= strconv.FormatFloat(i,  'f' , 2,  64)
 	log.Print("DE***** Contract value="+contract.Acct.Valuation)
        
        if i==0 {
		contract.Status="SR"
        }
	//********************************************
	//gl posting 
	//*************************
	// GL Posting Surrender Payment
	var glt GLtran
	glt.Ref="Policy Surrender Payment"
 	glt.Dbacc="PSUSP"
 	glt.Db= strconv.FormatFloat(float64(surr),  'f' , 2,  64)
 	glt.Cracc="BK001"
 	glt.Cr=strconv.FormatFloat(-float64(surr),  'f' , 2,  64)
	glPost(stub, glt , "SUR")
	//************************************
	// GL Posting Surrender Charge
	glt.Ref="Policy Surrender Charge"
 	glt.Dbacc="PSUSP"
 	glt.Db= strconv.FormatFloat(float64(surrchg),  'f' , 2,  64)
 	glt.Cracc="MGCHG"
 	glt.Cr=strconv.FormatFloat(-float64(surrchg),  'f' , 2,  64)
	glPost(stub, glt , "SCG")
    
	//*****************************************************
	// email
        var subject string
        var body string
	if i == 0 {
	  subject="Your Policy has been Surrendered"
	  body=`Dear Mr `+ contract.Lf.Name+ `#N         Your request to surrender your policy has been accepted #N and payment of $`+strconv.FormatFloat(surr,  'f' , 2,  64)+` has been made directly to your bank account #N Many thanks`
        }else{
	  subject="Partial Surrender of your policy"
	  body=`Dear Mr `+ contract.Lf.Name+ `#N         Your request to partial surrender of your policy has been accepted #N and payment of $`+strconv.FormatFloat(surr,  'f' , 2,  64)+` has been made directly to your bank account #N       The value remaining in your policy is `+ contract.Acct.Valuation + ` #N Many thanks`

        }
 	t.mailto(stub, subject , body, policy)
	policy.Cont=contract
	return policy , err
}

func (t *SimpleChaincode) fundAllocation(stub shim.ChaincodeStubInterface, args []string , policy Policy) ( Policy, error ) {
	//var contract Contract=policy.Cont
    	var f map[string]string
	f= make(map[string]string)
    	var mv map[string]string
	mv= make(map[string]string)	
        //******************************
	// arg[1] is the map of funds & allocations strig, string 
        json.Unmarshal([]byte(args[1]), &f)
	//******************************
	// Now loop round the array creating a map of unit changes
	for k, v := range f {
		if vv, ok :=policy.Cont.Acct.Fnds[k]; ok {
			x, _ := strconv.ParseFloat( v , 10);
			y, _ := strconv.ParseFloat( vv , 10);
			mv[k]= strconv.FormatFloat( x  - y ,  'f' , 2,  64) 
		}else{
			x, _ := strconv.ParseFloat( v , 10);
			mv[k]=strconv.FormatFloat( x ,  'f' , 2,  64)  
		}
		policy.Cont.Acct.Fnds[k]=v	
	}
        b, err := json.Marshal(mv)
	postFundUpdate(stub, policy, b )
	return  policy, err
}

func 	postFundUpdate(stub shim.ChaincodeStubInterface, policy Policy , mv  []byte )( Policy, error ) {
	_ , errx := stub.GetState("FM"+ invokeTran)
	if errx == nil {
		fmt.Println("Fund Manager update already done")
		stub.PutState("FM"+ invokeTran , []byte("Y") )
	}
	stub.PutState("FM"+ invokeTran , []byte("Y") )


	s:=strings.Replace(string(mv), "\"", "\\\"", -1)
	 var jsonStr = []byte( `{
   	  "jsonrpc": "2.0",
    	 "method": "invoke",
    	 "params": {
      	   "type": 1,
     	    "chaincodeID": {
      	       "name":"`+manager+`"
         },
         "ctorMsg": {
             "function": "updateFunds",
             "args": [
		 "`+policy.Cont.ContID+`",
                 "`+ s+`" 
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
    }` )
    fmt.Println(string(jsonStr))
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Fund update Satus:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Fund Update Body:", string(body))
    return  policy, err
}

func (t *SimpleChaincode) applyPremium(stub shim.ChaincodeStubInterface, args []string, policy Policy) (Policy, error) {

	// payment is arg[1]
	premium, _ := strconv.ParseFloat( args[1] , 10);

 	log.Print("DE***** Contract value="+policy.Cont.Acct.Valuation + "Payment="+args[1])


	i, _ := strconv.ParseFloat( policy.Cont.Acct.Valuation , 10);
	i = i + float64(premium)
        policy.Cont.Acct.Valuation= strconv.FormatFloat(i,  'f' , 2,  64)
 	log.Print("DE***** Contract value now="+policy.Cont.Acct.Valuation)

	//*************************
	// GL Posting
	var glt GLtran
	glt.Ref="Premium Payment"
 	glt.Dbacc="BK001"
 	glt.Db= strconv.FormatFloat(float64(premium),  'f' , 2,  64)
 	glt.Cracc="CSUSP"
 	glt.Cr=strconv.FormatFloat(-float64(premium),  'f' , 2,  64)
	glPost(stub, glt , "PPY" )

	//*************************************
	//* Now set the policy in force 

	if 	policy.Cont.Status=="IF" {
		//*****************************************************
		// email
		subject:="Thank you for your Payment"
		body:=`Dear Mr `+ policy.Cont.Lf.Name + `#N Thank you for your payment of $` +strconv.FormatFloat(premium,  'f' , 2,  64)+ ` for your policy `+policy.Cont.ContID+` #N Many thanks`
 		t.mailto(stub, subject , body, policy )
	} else{

	  if  policy.Cont.UWstatus=="Ready"{
  		policy.Cont.Status="IF"
		subject:="Your Policy is now in Force"
		body:=`Dear Mr `+policy.Cont.Lf.Name+ ` #N Thank you for your payment of $`+strconv.FormatFloat(premium,  'f' , 2,  64)+ ` for your new policy #N we are pleased to inform you that your policy is now in force #N Many thanks`
		t.mailto(stub, subject , body , policy )
	  }else{
		//*****************************************************
		// email
		subject:="Thank you for your Payment"
		body:=`Dear Mr `+ policy.Cont.Lf.Name + `#N Thank you for your payment of $` +strconv.FormatFloat(premium,  'f' , 2,  64)+ ` for your policy #N Many thanks`
 		t.mailto(stub, subject , body, policy )
	  }
       }
	policy.Cont.Status="IF"
 	//t.activate(stub, args )
	return  policy, nil
}

func glPost( stub shim.ChaincodeStubInterface, glt GLtran, pid string)( error){
	var err error
        var valAsbytes, b []byte
  	gltran=make(map[string]GLtran)
	valAsbytes, err = stub.GetState("gltran")
    	json.Unmarshal(valAsbytes , &gltran)
	glt.TranID=invokeTran
 	glt.Stat="N"
	gltran[invokeTran+pid]=glt
        b, err = json.Marshal(gltran)
	stub.PutState("gltran", b)
 	err=Glupdate(stub, glt , pid)
	return   err
}
func Glupdate(stub shim.ChaincodeStubInterface, glt GLtran, pid string ) ( error) {
	valAsbytes, err := stub.GetState("manager")
	manager=string(valAsbytes)
	_ ,errx := stub.GetState("GL"+ invokeTran) 
        if errx == nil {
		fmt.Println("GL Posting already done")
		stub.PutState("GL"+ invokeTran , []byte("Y") )
	}
	stub.PutState("GL"+ invokeTran , []byte("Y") )
	valAsbytes, err = stub.GetState("url")
	url=string(valAsbytes)
        b, err := json.Marshal(glt)
	
	s:=strings.Replace(string(b), "\"", "\\\"", -1)
	 var jsonStr = []byte( `{
   	  "jsonrpc": "2.0",
    	 "method": "invoke",
    	 "params": {
      	   "type": 1,
     	    "chaincodeID": {
      	       "name":"`+manager+`"
         },
         "ctorMsg": {
             "function": "updateT",
             "args": [
                 "`+ s +`" 
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` )
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err2 := client.Do(req)
    err=err2
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("GL Post Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("GL Post Body:", string(body))
    return err
}
func Odsupdate(stub shim.ChaincodeStubInterface, ods Ods, pid string ) ( error) {
	 _ , errx := stub.GetState("ODS"+ invokeTran) 
	if errx == nil {
		fmt.Println("ODS Posting already done")
		stub.PutState("ODS"+ invokeTran , []byte("Y") )
	}
	stub.PutState("ODS"+ invokeTran , []byte("Y") )


	valAsbytes, err := stub.GetState("manager")
	manager=string(valAsbytes)
	valAsbytes, err = stub.GetState("url")
	url=string(valAsbytes)
        b, err := json.Marshal(ods)
	s:=strings.Replace(string(b), "\"", "\\\"", -1)
	 var jsonStr = []byte( `{
   	  "jsonrpc": "2.0",
    	 "method": "invoke",
    	 "params": {
      	   "type": 1,
     	    "chaincodeID": {
      	       "name":"`+manager+`"
         },
         "ctorMsg": {
             "function": "updateOds",
             "args": [
                 "`+ s +`" 
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` )
    fmt.Println("ODS REQ="+string(jsonStr))
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err2 := client.Do(req)
    err=err2
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("ODS Update Status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("ODS Update Body:", string(body))
    return err
}
func (t *SimpleChaincode) monthlyProcessing(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
        fmt.Println("Starting Scheduled Processing")
	var err error
	  policies=make(map[string]string)
	  //*****************************************
	  // get the mal listing all policies  
	  vb, _ := stub.GetState("policies")
    	  json.Unmarshal(vb , &policies)

   //Iterate over the contracts & process all that are in force
        for key , value := range policies {
        fmt.Println("Sheduler Looking at Policy" + key +" Status="+value)
		if value == "IF" {
		  var policy Policy
		  policy.Hist=make(map[string]History)
		  //*****************************************
		  // get Contract state 
		  valAsbytes, _ := stub.GetState(key)
    		  json.Unmarshal(valAsbytes , &policy)

		  policy , err = t.ProcessPolicy(stub, args , policy)

		  policies[policy.Cont.ContID]=policy.Cont.Status

        	  b, _ := json.Marshal(policy)
		  err =	stub.PutState(key , b)
	        }
	}
        fmt.Println("completed Scheduled Processing")
	return nil , err
}
func (t *SimpleChaincode) ProcessPolicy(stub shim.ChaincodeStubInterface, args []string, policy Policy) (Policy, error) {
		  var err error
		  policy , err = t.ProcessCharges(stub, args , policy)
		  if policy.Cont.Status == "LS"{
                     t.lapseNotification(stub , args , policy )		
		  }else{
		    t.statement(stub , args , policy )
                  } 
	return policy, err	
}


type Res struct{
  COI string
  FMC string
  AMC string
}
type Calc struct{
  Service string
  DOB string
  CalcDate string
  Smoker string
  Gender string
  Suminsured string
}


func (t *SimpleChaincode) ProcessCharges(stub shim.ChaincodeStubInterface, args []string, policy Policy) (Policy, error) {
	var contract Contract=policy.Cont
        fmt.Println("Scheduled Processing for contract" + policy.Cont.ContID)


	//***************************************
	// Cal Calc Engine for Charges
        var x Calc

	year, _ , day := time.Now().Date()
        month:=time.Now().Month();
        //hour:=time.Now().Hour();
        //min:=time.Now().Minute();
        //sec:=time.Now().Second()
	dte:=strconv.Itoa(day)+"/"+strconv.Itoa(int(month))+"/"+strconv.Itoa(year)
	//dtex:=strconv.Itoa(day)+"/"+strconv.Itoa(int(month))+"/"+strconv.Itoa(year)+":"+strconv.Itoa(hour)+":"+strconv.Itoa(min)+":"+strconv.Itoa(sec)
        dtex:=  t.TransactionTime( stub, stub.GetTxID() )
	x.CalcDate=dte
        x.DOB=contract.Lf.Dob
        x.Gender=contract.Lf.Gender
        x.Smoker=contract.Lf.Smoker
        x.Suminsured=contract.SumAssured
        x.Service="DemoCharges"
        b := new(bytes.Buffer)
        json.NewEncoder(b).Encode(x)
        fmt.Println(b) 
        //res, errx := http.Post("http://203.106.175.109:8080/test", "application/json; charset=utf-8",  b )
        var resx Res;
	coi:=33.00
        fmc:=10.00
	adc:=12.00
        totalcharges:=55.00
        /*if errx == nil {
          body, _ := ioutil.ReadAll(res.Body)

          json.Unmarshal(body , &resx)
         coi, _ =strconv.ParseFloat(resx.COI,10)
	 fmc, _ =strconv.ParseFloat(resx.FMC,10)
         adc, _ =strconv.ParseFloat(resx.AMC,10)
         totalcharges=coi+fmc+adc
         fmt.Println( "COI="+ resx.COI +" FMC=" + resx.FMC +" AMC="+ resx.AMC )
       } else{   ******************/
          //fmt.Println(errx)
	  resx.COI="33"
	  resx.FMC="10"
	  resx.AMC="12"
       // }


	//****************************
	// Write history record first
        charges:=[]string{resx.COI,resx.FMC, resx.AMC}
	var h History
	h.Methd="Scheduled Processing"
	h.Funct="Deduct Charges"
	h.Tranid=invokeTran
	h.Cont=policy.Cont
	h.Args=charges
        h.Args[0]= strconv.FormatFloat( totalcharges ,  'f' , 2,  64)
        h.Dte=dtex

        //*****************************
        // save policy sate 
        b1, _ := json.Marshal(policy)
	stub.PutState(policy.Cont.ContID, b1)

	//*******************************************
	//* Do the valuation
 	fmt.Print("DE***** Contract value="+contract.Acct.Valuation)

	i, _ := strconv.ParseFloat( contract.Acct.Valuation , 10);
	i = i - float64(coi+fmc+adc)
        if i< 0 {
          i = 0 
        }
        contract.Acct.Valuation= strconv.FormatFloat(i,  'f' , 2,  64)
 	log.Print("DE***** Contract value="+contract.Acct.Valuation)
        h.EndValue=policy.Cont.Acct.Valuation
	policy.Hist[dtex]=h
        //********************************************************
        //* Check lapsing rules
        if i < ( coi+fmc+adc ) {
         contract.Status="LS"
	}
	//*************************
	// GL Posting COI
	var glt GLtran
	glt.Ref="COI Deduction"
 	glt.Dbacc="PSUSP"
 	glt.Db= strconv.FormatFloat(float64(coi),  'f' , 2,  64)
 	glt.Cracc="PRESV"
 	glt.Cr=strconv.FormatFloat(-float64(coi),  'f' , 2,  64)
	glPost(stub, glt , "COI")
	//*************************
	// GL Posting Fund Management Charge
	glt.Ref="Fund Management Charge"
 	glt.Dbacc="PSUSP"
 	glt.Db= strconv.FormatFloat(float64(fmc),  'f' , 2,  64)
 	glt.Cracc="FDEXP"
 	glt.Cr=strconv.FormatFloat(-float64(fmc),  'f' , 2,  64)
	glPost(stub, glt , "FMC")
	//*************************
	// GL Posting Admmin Charge
	glt.Ref="Admin Charge"
 	glt.Dbacc="PSUSP"
 	glt.Db= strconv.FormatFloat(float64(adc),  'f' , 2,  64)
 	glt.Cracc="MGEXP"
 	glt.Cr=strconv.FormatFloat(-float64(adc),  'f' , 2,  64)
	glPost(stub, glt , "ADC" )
	policy.Cont=contract
	return policy, nil

}


func (t *SimpleChaincode) statement(stub shim.ChaincodeStubInterface, args []string , policy Policy) ([]byte, error) {
	var contract Contract=policy.Cont
	//subject:="Your Monthly Statement"
	//body:= `Your Statement of account: #N Account Holder:`+contract.Owner+`#NPolicy No: `+contract.ContID +` #N Value:`+contract.Acct.Valuation+` #N Yours Sincerely, Danny`
	subject:="Your monthly statement"
	body:=`Dear Mr `+ contract.Lf.Name + `#N Policy Number=`+policy.Cont.ContID+`#N Value of Policy=`+contract.Acct.Valuation+` #N `
	t.mailto(stub, subject , body, policy )
 return nil,nil
}

func (t *SimpleChaincode) lapseNotification(stub shim.ChaincodeStubInterface, args []string , policy Policy) ([]byte, error) {
	var contract Contract=policy.Cont
	//subject:="Your Policy Has Lapsed"
	//body:= `: #N Account Holder:`+contract.Owner+`#NPolicy No: `+contract.ContID +` #N Value:`+contract.Acct.Valuation+` #N Yours Sincerely, Danny`
	subject:="Your monthly statement"
	body:=`Dear Mr `+ contract.Lf.Name + `#N Policy Number=`+policy.Cont.ContID+`#N The remaining value of your policy is below the minimum to support it#N`+
        `Value of Policy=`+contract.Acct.Valuation+` #N The policy has now lapsed. If payment is made within one month the policy will be reinstated.`+
	`#N If no payment is made the remaining value will be paid out within 14 days`
	t.mailto(stub, subject , body, policy )
 return nil,nil
}

func(t *SimpleChaincode) welcome(stub shim.ChaincodeStubInterface , policy Policy) ([]byte, error) {
	var contract Contract=policy.Cont
	subject:="Thank you for your application"
	body:=`Dear Mr `+ contract.Lf.Name + `#N Policy Number=`+policy.Cont.ContID+`#N Thank you for your application, which has now been accepted #N We will activate your new Policy as soon as payment is received`

 t.mailto(stub, subject , body, policy )
 return nil,nil
}



func (t *SimpleChaincode) mailto(stub shim.ChaincodeStubInterface, subject string, body string , policy Policy ) ([]byte, error) {
	var contract Contract=policy.Cont
 	valAsbytes, err := stub.GetState("scheduler")
	scheduler=string(valAsbytes)

	valAsbytes, err = stub.GetState("url")
	url=string(valAsbytes)


	 var jsonStr = []byte( `{
   	  "jsonrpc": "2.0",
    	 "method": "query",
    	 "params": {
      	   "type": 1,
     	    "chaincodeID": {
      	       "name":"`+commsmanager+`"
         },
         "ctorMsg": {
             "function": "mailto",
             "args": [
                 "`+stub.GetTxID()+`",
		 "`+subject+`",
		 "`+body+`",
		 "`+contract.Email+`"
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }` )


    fmt.Println("Send Email:", string(jsonStr) )
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    //req.Header.Set("Postman-Token", "")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("accept", "application/json")
    client := &http.Client{}
    resp, err2 := client.Do(req)
    err=err2
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Email To Status:", resp.Status)
   
   
     return  []byte("Mail sent"), err
}



func (t *SimpleChaincode) valuation(stub shim.ChaincodeStubInterface, args []string , policy Policy) ([]byte, error) {
	var contract Contract=policy.Cont
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)

	return  []byte("Valuation="+contract.Acct.Valuation), err
}



// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	// first those that dont take policy number
	if function == "journal" {
		return t.journal(stub, args )
	}

	//**************************************
	// Now the functions that require a policy number
	fmt.Println("invoke for policy " + args[0])
	var policy Policy
	policy.Hist=make(map[string]History)
	//*****************************************
	// get Contract state 
	valAsbytes, _ := stub.GetState(args[0])
    	json.Unmarshal(valAsbytes , &policy)
	// Handle different functions
	if function == "statement" {
		return t.statement(stub, args, policy)
	} else if function == "valuation" {
		return t.valuation(stub, args, policy)
	} else if function == "transactions" {
		return t.transactions(stub, args, policy)
	} else if function == "dump" {
		return t.dump(stub, args, policy)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


func (t *SimpleChaincode) transactions(stub shim.ChaincodeStubInterface, args []string , policy Policy) ([]byte, error) {
	var valAsbytes []byte
    	valAsbytes, err :=json.Marshal( policy.Hist )
	fmt.Println( "TRANSACTION HISTORY="+ string(valAsbytes))
	fmt.Println( err)
  return valAsbytes, err
}


// read - query function to read key/value pair
func (t *SimpleChaincode) dump(stub shim.ChaincodeStubInterface, args []string , policy Policy ) ([]byte, error) {
    	valAsbytes, err :=json.Marshal( policy)
	fmt.Println( "POLICY DUMP="+ string(valAsbytes))
	return valAsbytes, err
}
