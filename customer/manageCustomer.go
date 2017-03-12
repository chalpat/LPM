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
	} else if function == "createCustomer" {											//create a new Customer
		return t.createCustomer(stub, args)
	}else if function == "deleteCustomer" {									// delete a Customer
		return t.deleteCustomer(stub, args)
	}else if function == "updateCustomer" {									//update a Customer
		return t.updateCustomer(stub, args)
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
	fmt.Print("customerId in getCustomerByID : "+customerId)
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
// ============================================================================================================================
// Delete - remove a Customer from chain
// ============================================================================================================================
func (t *ManageCustomer) deleteCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'customerId' as an argument\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set customerId
	customerId := args[0]
	err := stub.DelState(customerId)						//remove the Customer from chaincode
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to delete state\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//get the Customer index
	customerAsBytes, err := stub.GetState(CustomerIndexStr)
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get Customer index\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	var customerIndex []string
	json.Unmarshal(customerAsBytes, &customerIndex)								//un stringify it aka JSON.parse()
	//remove marble from index
	for i,val := range customerIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + customerId)
		if val == customerId{															//find the correct Customer
			fmt.Println("found Customer with matching customerId")
			customerIndex = append(customerIndex[:i], customerIndex[i+1:]...)			//remove it
			for x:= range customerIndex{											//debug prints...
				fmt.Println(string(x) + " - " + customerIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(customerIndex)									//save new index
	err = stub.PutState(CustomerIndexStr, jsonAsBytes)

	tosend := "{ \"customerID\" : \""+customerId+"\", \"message\" : \"Customer deleted succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	// TODO:: All transactions related to customer should be deleted.

	fmt.Println("Customer deleted succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Write - update customer into chaincode state
// ============================================================================================================================
func (t *ManageCustomer) updateCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Customer")
	if len(args) != 15 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 15\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set customerId
	customerId := args[0]
	customerAsBytes, err := stub.GetState(customerId)					//get the Customer for the specified customerId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get state for " + customerId + "\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	res := Customer{}
	res_trans := Transaction{}
	transactionId := args[8]
	json.Unmarshal(customerAsBytes, &res)
	if res.CustomerID == customerId{
		fmt.Println("Customer found with customerId : " + customerId)
		fmt.Println(res);
		res.CustomerName = args[1]
		res.WalletWorth = args[2]
		res.MerchantIDs = args[3]
		res.MerchantNames = args[4]
		res.MerchantCurrencies = args[5]
		res.MerchantsPointsCount = args[6]
		res.MerchantsPointsWorth = args[7]
		res_trans.TransactionID = transactionId
		res_trans.TransactionDateTime = args[9]
		res_trans.TransactionType = args[10]
		res_trans.TransactionFrom = args[11]
		res_trans.TransactionTo = args[12]
		res_trans.Credit = args[13]
		res_trans.Debit = args[14]
		res_trans.CustomerID = customerId
	}else{
		errMsg := "{ \"message\" : \""+ customerId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Customer json string manually
	customer_json := 	`{`+
		`"customerId": "` + res.CustomerID + `" , `+
		`"customerName": "` + res.CustomerName + `" , `+
		`"walletWorth": "` + res.WalletWorth + `" , `+
		`"merchantIDs": "` + res.MerchantIDs + `" , `+ 
		`"merchantNames": "` + res.MerchantNames + `" , `+ 
		`"merchantCurrencies": "` + res.MerchantCurrencies + `" , `+ 
		`"merchantsPointsCount": "` + res.MerchantsPointsCount + `" , `+ 
		`"merchantsPointsWorth": "` +  res.MerchantsPointsWorth + `" `+ 
	`}`
	err = stub.PutState(customerId, []byte(customer_json))									//store Customer with id as key
	if err != nil {
		return nil, err
	}

	//build the Transaction json string manually
	transaction_json := `{`+
		`"transactionId": "` + transactionId + `" , `+
		`"transactionDateTime": "` + res_trans.TransactionDateTime + `" , `+
		`"transactionType": "` + res_trans.TransactionType + `" , `+
		`"transactionFrom": "` + res_trans.TransactionFrom + `" , `+ 
		`"transactionTo": "` + res_trans.TransactionTo + `" , `+ 
		`"credit": "` + res_trans.Credit + `" , `+ 
		`"debit": "` + res_trans.Debit + `" , `+ 
		`"customerId": "` +  res_trans.CustomerID + `" `+ 
	`}`
	err = stub.PutState(transactionId, []byte(transaction_json))									//store Transaction with id as key
	if err != nil {
		return nil, err
	}

	//get the Transaction index
	transactionAsBytes, err := stub.GetState(TransactionIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Transaction index")
	}
	var transactionIndex []string	
	json.Unmarshal(transactionAsBytes, &transactionIndex)							//un stringify it aka JSON.parse()
	
	//append
	transactionIndex = append(transactionIndex, transactionId)									//add Transaction transactionId to index list
	
	jsonAsBytes, _ := json.Marshal(transactionIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(TransactionIndexStr, jsonAsBytes)						//store name of Transaction
	if err != nil {
		return nil, err
	}

	tosend := "{ \"customerID\" : \""+customerId+"\", \"message\" : \"Customer details updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Customer details updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// create Customer - create a new Customer, store into chaincode state
// ============================================================================================================================
func (t *ManageCustomer) createCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 8 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 8\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start createCustomer")
	customerId := args[0]
	customerName := args[1]
	walletWorth := args[2]
	merchantIDs := args[3]
	merchantNames := args[4]
	merchantCurrencies := args[5]
	merchantsPointsCount := args[6]
	merchantsPointsWorth := args[7]
	
	customerAsBytes, err := stub.GetState(customerId)
	if err != nil {
		return nil, errors.New("Failed to get Customer customerID")
	}
	res := Customer{}
	json.Unmarshal(customerAsBytes, &res)
	if res.CustomerID == customerId{
		errMsg := "{ \"message\" : \"This Customer arleady exists\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil				//all stop a Customer by this name exists
	}
	
	//build the Customer json string manually
	customer_json := 	`{`+
		`"customerId": "` + customerId + `" , `+
		`"customerName": "` + customerName + `" , `+
		`"walletWorth": "` + walletWorth + `" , `+
		`"merchantIDs": "` + merchantIDs + `" , `+ 
		`"merchantNames": "` + merchantNames + `" , `+ 
		`"merchantCurrencies": "` + merchantCurrencies + `" , `+ 
		`"merchantsPointsCount": "` + merchantsPointsCount + `" , `+ 
		`"merchantsPointsWorth": "` +  merchantsPointsWorth + `" `+ 
	`}`
	fmt.Println("customer_json: " + customer_json)
	fmt.Print("customer_json in bytes array: ")
	fmt.Println([]byte(customer_json))
	err = stub.PutState(customerId, []byte(customer_json))									//store Customer with customerId as key
	if err != nil {
		return nil, err
	}
	//get the Customer index
	customerIndexAsBytes, err := stub.GetState(CustomerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Customer index")
	}
	var customerIndex []string	
	json.Unmarshal(customerIndexAsBytes, &customerIndex)							//un stringify it aka JSON.parse()
	
	//append
	customerIndex = append(customerIndex, customerId)									//add Customer customerID to index list
	
	jsonAsBytes, _ := json.Marshal(customerIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(CustomerIndexStr, jsonAsBytes)						//store name of Customer
	if err != nil {
		return nil, err
	}

	tosend := "{ \"customerID\" : \""+customerId+"\", \"message\" : \"Customer created succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("end createCustomer")
	return nil, nil
}