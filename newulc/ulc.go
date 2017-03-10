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
var glmanager string
var odsmanager string
var commsmanager string
func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 " )
	}
        count=0;
	l := log.New(os.Stderr, "", 0)
	l.Println("*************INIT CHAINCODE Unit Linked ****************")
	glmanager=args[0]
	odsmanager=args[1]
	commsmanager=args[2]
	scheduler=args[3]
	url=args[4]
	err := 	stub.PutState("scheduler",[]byte(scheduler) )
	err = 	stub.PutState("url",[]byte(url) )
	err = 	stub.PutState("glmanager",[]byte(glmanager) )
	err = 	stub.PutState("odsmanager",[]byte(odsmanager) )
	err = 	stub.PutState("commsmanager",[]byte(commsmanager) )
        //fmt.Println( xx.A )
	if err != nil {
		return nil, err
	}
	return nil, err
}

//***********************************************
//* Create a newpolicy
func (t *SimpleChaincode) NewPolicy(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
	if len(args) != 12 {
		return nil, errors.New("Incorrect number of arguments. Expecting 12")
	}
        count=0;
	l := log.New(os.Stderr, "", 0)
	l.Println("************* Unit Linked  New Policy****************")
	l.Println("Policy ID="+stub.GetTxID())
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

	if _ ,ok:=policies[policy.Cont.ContID] ; ok {
		fmt.Println("Contract Already Exist ")
		return nil, errors.New("Contract exists already")
        }

       //**************************************************
       // save the history
	policy.Hist=make(map[string]History)
	var h History
	h.Methd="deploy"
	h.Funct="init"
	h.Tranid=stub.GetTxID()  //time.Now().String()
	h.Cont=policy.Cont
	h.Args=args
	policy.Hist[h.Tranid]=h

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
	fmt.Println("invoke is running " + function)
	l := log.New(os.Stderr, "", 0)
	fmt.Println("invoke for policy " + args[0])
	var policy Policy
	policy.Hist=make(map[string]History)
	//*****************************************
	// get Contract state 
	valAsbytes, _ := stub.GetState(args[0])
    	json.Unmarshal(valAsbytes , &policy)

	var h History
	h.Methd="invoke"
	h.Funct=function
	h.Tranid=invokeTran
	h.Cont=policy.Cont
	h.Args=args
	policy.Hist[h.Tranid]=h

        var err error
	l.Println("DE************* Invoke Function "+ function )
        //xx = shared.Args{1, 2} 
	// Handle different functions
	if function == "init" {
		return  t.Init(stub, "init", args)
	} else if function == "fundAllocation" {
		policy , err = t.fundAllocation(stub, args, policy)
	} else if function == "applyPremium" {
		policy,err = t.applyPremium(stub, args, policy)
	} else if function == "schedule" {
		_ ,err = t.monthlyProcessing(stub, args )
	} else if function == "surrender" {
		policy,err = t.surrender(stub, args , policy )
	} else if function == "activate" {
		_ ,err = t.activate(stub, args  )
	} else if function == "deactivate" {
		_ ,err = t.deactivate(stub, args )
	} else if function == " setJournalDone" {
		return t. setJournalDone(stub, args)
	} else if function == " NewPolicy" {
		return t. NewPolicy(stub, args)
	}else{ 
		fmt.Println("invoke did not find func: " + function)
		err=errors.New("Received unknown function invocation: " + function)
        }
	if  policy.Cont.ContID=="" {
		return nil,nil
        }
	//*************************
	//* if we are here then we need to update the ODS
	var ods Ods
	ods.Cont=policy.Cont
	ods.Tranid=invokeTran
	ods.Posted="N"
	Odsupdate(stub , ods , invokeTran ) 

        //*****************************
        // save policy sate 
        b, err := json.Marshal(policy)
	err = 	stub.PutState(policy.Cont.ContID, b)

	//**********************************
	//* Save Ploicies map with new policy status
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
	surr, _ := strconv.ParseFloat( args[0] , 10);
	surrchg:=20
 	fmt.Print("DE***** Contract value="+contract.Acct.Valuation)
	i, _ := strconv.ParseFloat( contract.Acct.Valuation , 10);
	i = i - float64(surr)
	i = i - float64(surrchg)
        contract.Acct.Valuation= strconv.FormatFloat(i,  'f' , 2,  64)
 	log.Print("DE***** Contract value="+contract.Acct.Valuation)
        b, err := json.Marshal(contract)
	err = 	stub.PutState("Contract", b)
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
	subject:="Your Policy has been Surrendered"
	body:=`Dear Mr `+ contract.Lf.Name+ `#N	Your request to surrender your policy has been accepted #N and payment of $`+strconv.FormatFloat(surr,  'f' , 2,  64)+` has been made directly to your bank account
	Many thanks`
 	t.mailto(stub, subject , body, policy)
	policy.Cont=contract
	return policy , err
}

func (t *SimpleChaincode) fundAllocation(stub shim.ChaincodeStubInterface, args []string , policy Policy) ( Policy, error ) {
	var contract Contract=policy.Cont
    	var fnd[20]Fund
        json.Unmarshal([]byte(args[0]), &fnd)
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)
        contract.Acct.Fnds=fnd
	policy.Cont=contract
	return  policy, err
}


func (t *SimpleChaincode) applyPremium(stub shim.ChaincodeStubInterface, args []string, policy Policy) (Policy, error) {
	var contract Contract=policy.Cont
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
	if  contract.UWstatus=="Ready"{
  		contract.Status="InForce"
		t.activate(stub, args)	

		subject:="Your Policy is now in Force"
		body:=`Dear Mr `+contract.Lf.Name+ ` #N Thank you for your payment of $`+strconv.FormatFloat(premium,  'f' , 2,  64)+ ` for your new policy #N we are pleased to inform you that your policy is now in force #N Many thanks`
		t.mailto(stub, subject , body , policy )
	}else
	{
		//*****************************************************
		// email
		subject:="Thank you for your Payment"
		body:=`Dear Mr `+ contract.Lf.Name + `#N Thank you for your payment of $` +strconv.FormatFloat(premium,  'f' , 2,  64)+ ` for your policy #N Many thanks`
 		t.mailto(stub, subject , body, policy )
	}
	contract.Status="IF"
	policy.Cont=contract
 	t.activate(stub, args )
	return  policy, err
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
	valAsbytes, err := stub.GetState("glmanager")
	glmanager=string(valAsbytes)


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
      	       "name":"`+glmanager+`"
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
	valAsbytes, err := stub.GetState("odsmanager")
	odsmanager=string(valAsbytes)
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
      	       "name":"`+odsmanager+`"
         },
         "ctorMsg": {
             "function": "update",
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
	var err error
   //Iterate over the contracts & process all that are in force
        for key , value := range policies {
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
	return nil , err
}

func (t *SimpleChaincode) ProcessPolicy(stub shim.ChaincodeStubInterface, args []string, policy Policy) (Policy, error) {
	var contract Contract=policy.Cont
	coi:=33
        fmc:=10
	adc:=12
 	fmt.Print("DE***** Contract value="+contract.Acct.Valuation)
	valAsbytes, err := stub.GetState("Contract")
    	json.Unmarshal(valAsbytes , &contract)

	i, _ := strconv.ParseFloat( contract.Acct.Valuation , 10);
	i = i - float64(coi+fmc+adc)
        contract.Acct.Valuation= strconv.FormatFloat(i,  'f' , 2,  64)
 	log.Print("DE***** Contract value="+contract.Acct.Valuation)
        b, err := json.Marshal(contract)
	err = 	stub.PutState("Contract", b)
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
	return policy, err

}


func (t *SimpleChaincode) statement(stub shim.ChaincodeStubInterface, args []string , policy Policy) ([]byte, error) {
	var contract Contract=policy.Cont
	subject:="Your Monthly Statement for "+contract.ContID 
	body:= `Your Statement of account: #N Account Holder:`+contract.Owner+` #N Value:`+contract.Acct.Valuation+` #N Yours Sincerely, Danny `
	t.mailto(stub, subject , body, policy )
 return nil,nil
}

func(t *SimpleChaincode) welcome(stub shim.ChaincodeStubInterface , policy Policy) ([]byte, error) {
	var contract Contract=policy.Cont
	subject:="Thank you for your application"
	body:=`Dear Mr `+ contract.Lf.Name + `#N Thank you for your application, which has now been accepted #N We will activate your new Policy as soon as payment is received `

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
	fmt.Println("invoke for policy " + args[0])
	var policy Policy
	//*****************************************
	// get Contract state 
	valAsbytes, _ := stub.GetState(args[0])
    	json.Unmarshal(valAsbytes , &policy)
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args, policy )
	}else if function == "statement" {
		return t.statement(stub, args, policy)
	} else if function == "valuation" {
		return t.valuation(stub, args, policy)
	} else if function == "activate" {
		return t.activate(stub, args)
	} else if function == "deactivate" {
		return t.deactivate(stub, args)
	} else if function == "transactions" {
		return t.transactions(stub, args, policy)
	} else if function == "journal" {
		return t.journal(stub, args )
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


func (t *SimpleChaincode) transactions(stub shim.ChaincodeStubInterface, args []string , policy Policy) ([]byte, error) {
	var valAsbytes []byte
	err:=json.Unmarshal(valAsbytes , &policy.Hist)
  return valAsbytes, err
}


// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string , policy Policy ) ([]byte, error) {
	var key , jsonResp string
	var err error
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
    	valAsbytes, err :=json.Marshal( policy );

	return valAsbytes, nil
}
