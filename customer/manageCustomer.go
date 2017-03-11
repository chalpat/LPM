/*/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
"errors"
"fmt"
"strconv"
"encoding/json"

"github.com/hyperledger/fabric/core/chaincode/shim"	
)

// ManageCustomer example simple Chaincode implementation
type ManageCustomer struct {
}

var CustomerIndexStr = "_Customerindex"				// name for the key/value that will store a list of all known Customer
var TransactionIndexStr = "_Transactionindex"		// name for the key/value that will store a list of all known Transaction

type Customer struct{							// Attributes of a Customer 
	CustomerID string `json:"customerId"`					
	CustomerName string `json:"customerName"`
	WalletWorth string `json:"walletWorth"`
	MerchantIDs string `json:"merchantIDs"`
	MerchantNames string `json:"merchantNames"`
	MerchantCurrencies string `json:"merchantCurrencies"`
	MerchantsPointsCount string `json:"merchantsPointsCount"`
	MerchantsPointsWorth string `json:"merchantsPointsWorth"`
}

type Transaction struct{							// Attributes of a Transaction 
	TransactionID string `json:"transactionId"`					
	TransactionDateTime string `json:"transactionDateTime"`
	TransactionType string `json:"transactionType"`				// Values are Purchase, Transfer, Accumulation (Add Points)
	TransactionFrom string `json:"transactionFrom"`
	TransactionTo string `json:"transactionTo"`
	Credit string `json:"credit"`
	Debit string `json:"debit"`
	CustomerID string `json:"customerId"`
}

// ============================================================================================================================
// Main - start the chaincode for Customer management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageCustomer))
	if err != nil {
		fmt.Printf("Error starting Customer management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageCustomer) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting ' ' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	// Initialize the chaincode
	msg = args[0]
	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))		//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(CustomerIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(TransactionIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	tosend := "{ \"message\" : \"ManageCustomer chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *ManageCustomer) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}
// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *ManageCustomer) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} 
	fmt.Println("invoke did not find func: " + function)
	errMsg := "{ \"message\" : \"Received unknown function invocation\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil			//error
}
// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *ManageCustomer) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getCustomerByID" {													//Read a Customer by transId
		return t.getCustomerByID(stub, args)
	} else if function == "getActivityHistory" {													//Read a Customer by Buyer's name
		return t.getActivityHistory(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error
	errMsg := "{ \"message\" : \"Received unknown function query\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// getCustomerByID - get Customer details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageCustomer) getCustomerByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var customerId string
	var err error
	fmt.Println("start getCustomerByID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'customerId' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set customerId
	customerId = args[0]
	valAsbytes, err := stub.GetState(customerId)									//get the customerId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \""+ customerId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Print("valAsbytes : ")
	fmt.Println(valAsbytes)
	fmt.Println("end getCustomerByID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getActivityHistory - get Customer Transaction Activity details from chaincode state
// ============================================================================================================================
func (t *ManageCustomer) getActivityHistory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var transactionIndex []string
	var valIndex Transaction
	var customerId string
	var err error
	fmt.Println("start getActivityHistory")
	
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'customerId' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set customerId
	customerId = args[0]
	fmt.Println("customerId in getActivityHistory::" + customerId)

	transactionAsBytes, err := stub.GetState(TransactionIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Transaction index string")
	}
	json.Unmarshal(transactionAsBytes, &transactionIndex)								//un stringify it aka JSON.parse()
	jsonResp = "{"
	for i,val := range transactionIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getActivityHistory")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.CustomerID == customerId{
			fmt.Println("Customer found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(transactionIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} else{
			errMsg := "{ \"message\" : \""+ customerId+ " Not Found.\", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
		
	}
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getActivityHistory")
	return []byte(jsonResp), nil											//send it onward
}