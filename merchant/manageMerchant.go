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
"strings"

"github.com/hyperledger/fabric/core/chaincode/shim"	
)

// ManageMerchant example simple Chaincode implementation
type ManageMerchant struct {
}

var CustomerIndexStr = "_Customerindex"				// name for the key/value that will store a list of all known Customer
var TransactionIndexStr = "_Transactionindex"		// name for the key/value that will store a list of all known Transaction
var MerchantIndexStr = "_Merchantindex"				//name for the key/value that will store a list of all known Merchants

type Customer struct{							// Attributes of a Customer 
	CustomerID string `json:"customerId"`					
	UserName string `json:"userName"`
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

type Merchant struct{							// Attributes of a Merchant
	MerchantID string `json:"merchantId"`					
	MerchantUserName string `json:"merchantUserName"`
	MerchantName string `json:"merchantName"`
	MerchantIndustry string `json:"merchantIndustry"`					
	PointsPerDollarSpent string `json:"pointsPerDollarSpent"`
	MerchantCurrency string `json:"merchantCurrency"`
	MerchantCU_date string `json:"merchantCU_date"`
}

// ============================================================================================================================
// Main - start the chaincode for Customer management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageMerchant))
	if err != nil {
		fmt.Printf("Error starting Customer management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageMerchant) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
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
	tosend := "{ \"message\" : \"ManageMerchant chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *ManageMerchant) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}
// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *ManageMerchant) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}else if function == "createMerchant" {											//create a new Merchant
		return t.createMerchant(stub, args)
	}else if function == "deleteMerchant" {									// delete a Merchant
		return t.deleteMerchant(stub, args)
	}else if function == "updateMerchant" {									//update a Merchant
		return t.updateMerchant(stub, args)
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
func (t *ManageMerchant) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getCustomersByMerchantID" {													//Read a Customer by transId
		return t.getCustomersByMerchantID(stub, args)
	}else if function == "getMerchantByName" {													//Read all Merchants
		return t.getMerchantByName(stub, args)
	}else if function == "getMerchantByID" {													//Read all Merchants
		return t.getMerchantByID(stub, args)
	}else if function == "getMerchantDetailsByID" {													//Read all Merchants
		return t.getMerchantDetailsByID(stub, args)
	}else if function == "getMerchantsByIndustry" {													//Read all Merchants
		return t.getMerchantsByIndustry(stub, args)
	}else if function == "getAllMerchants" {													//Read all Merchants
		return t.getAllMerchants(stub, args)
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
// getCustomersByMerchantID - get Customers for a specific Merchant ID from chaincode state
// ============================================================================================================================
func (t *ManageMerchant) getCustomersByMerchantID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, merchantId, errResp string
	var err error
	var customerIndex []string
	var valIndex Customer
	fmt.Println("start getCustomersByMerchantID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantId' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId = args[0]

	customerAsBytes, err := stub.GetState(CustomerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Customer index string")
	}
	json.Unmarshal(customerAsBytes, &customerIndex)								//un stringify it aka JSON.parse()
	fmt.Print("customerIndex : ")
	fmt.Println(customerIndex)
	
	jsonResp = "{"
	for i,val := range customerIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getCustomersByMerchantID")
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
		if strings.Contains(valIndex.MerchantIDs, merchantId){
			fmt.Println("Customer found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(customerIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} else{
			errMsg := "{ \"message\" : \""+ merchantId+ " Not Found.\", \"code\" : \"503\"}"
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
	fmt.Println("end getCustomersByMerchantID")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  getMerchantByName - get Merchant details by name from chaincode state
// ============================================================================================================================
func (t *ManageMerchant) getMerchantByName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, merchantName, errResp string
	var merchantIndex []string
	var valIndex Merchant
	fmt.Println("start getMerchantByName")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantName' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchant's name
	merchantName = args[0]
	//fmt.Println("merchantName" + merchantName)
	merchantAsBytes, err := stub.GetState(MerchantIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant index string")
	}
	json.Unmarshal(merchantAsBytes, &merchantIndex)				//un stringify it aka JSON.parse()
	fmt.Print("merchantIndex : ")
	fmt.Println(merchantIndex)
	jsonResp = "{"
	for i,val := range merchantIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getMerchantByName")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.MerchantName == merchantName{
			fmt.Println("Merchant found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(merchantIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} 
		
	}
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getMerchantByName")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
// getMerchantByID - get Merchant details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageMerchant) getMerchantByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var merchantId string
	var err error
	fmt.Println("start getMerchantByID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantId' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId = args[0]
	valAsbytes, err := stub.GetState(merchantId)									//get the merchantId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \""+ merchantId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("end getMerchantByID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
// getMerchantDetailsByID - get Merchant details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageMerchant) getMerchantDetailsByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var merchantId string
	var err error
	fmt.Println("start getMerchantDetailsByID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantId' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId = args[0]
	valAsbytes, err := stub.GetState(merchantId)									//get the merchantId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \""+ merchantId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("end getMerchantDetailsByID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
// getMerchantsByIndustry - get Merchants for a given Industry from chaincode state
// ============================================================================================================================
func (t *ManageMerchant) getMerchantsByIndustry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, industryName, errResp string
	var err error
	var merchantIndex []string
	var valIndex Merchant
	fmt.Println("start getMerchantsByIndustry")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'industryName' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	industryName = args[0]

	merchantAsBytes, err := stub.GetState(MerchantIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant index string")
	}
	json.Unmarshal(merchantAsBytes, &merchantIndex)			//un stringify it aka JSON.parse()
	fmt.Print("merchantIndex : ")
	fmt.Println(merchantIndex)
	
	jsonResp = "{"
	for i,val := range merchantIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getMerchantsByIndustry")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.MerchantIndustry == industryName{
			fmt.Println("Merchant found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(merchantIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} else{
			errMsg := "{ \"message\" : \""+ industryName+ " Not Found.\", \"code\" : \"503\"}"
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
	fmt.Println("end getMerchantsByIndustry")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  getAllMerchants- get details of all Merchants from chaincode state
// ============================================================================================================================
func (t *ManageMerchant) getAllMerchants(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var merchantIndex []string
	var err error
	fmt.Println("start getAllMerchants")
		
	merchantAsBytes, err := stub.GetState(MerchantIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant index")
	}
	json.Unmarshal(merchantAsBytes, &merchantIndex)								//un stringify it aka JSON.parse()
	jsonResp = "{"
	for i,val := range merchantIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Merchant")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(merchantIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp in getAllMerchants::")
	fmt.Println(jsonResp)

	fmt.Println("end getAllMerchants")
	return []byte(jsonResp), nil			//send it onward
}
// ============================================================================================================================
// Delete - remove a merchant from chain
// ============================================================================================================================
func (t *ManageMerchant) deleteMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantId' as an argument\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId := args[0]
	err := stub.DelState(merchantId)													//remove the Merchant from chaincode
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to delete state\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//get the Merchant index
	merchantAsBytes, err := stub.GetState(MerchantIndexStr)
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get Merchant index\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	var merchantIndex []string
	json.Unmarshal(merchantAsBytes, &merchantIndex)								//un stringify it aka JSON.parse()
	//remove marble from index
	for i,val := range merchantIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + merchantId)
		if val == merchantId{															//find the correct Merchant
			fmt.Println("found Merchant with matching merchantId")
			merchantIndex = append(merchantIndex[:i], merchantIndex[i+1:]...)			//remove it
			for x:= range merchantIndex{											//debug prints...
				fmt.Println(string(x) + " - " + merchantIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(merchantIndex)									//save new index
	err = stub.PutState(MerchantIndexStr, jsonAsBytes)

	tosend := "{ \"merchantID\" : \""+merchantId+"\", \"message\" : \"Merchant deleted succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Merchant deleted succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Write - update merchant into chaincode state
// ============================================================================================================================
func (t *ManageMerchant) updateMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Merchant")
	if len(args) != 7 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 7\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId := args[0]
	merchantAsBytes, err := stub.GetState(merchantId)									//get the Merchant for the specified merchant from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get state for " + merchantId + "\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	res := Merchant{}
	json.Unmarshal(merchantAsBytes, &res)
	if res.MerchantID == merchantId{
		fmt.Println("Merchant found with merchantId : " + merchantId)
		//fmt.Println(res);
		res.MerchantUserName = args[1]
		res.MerchantName = args[2]
		res.MerchantIndustry = args[3]
		res.PointsPerDollarSpent = args[4]
		res.MerchantCurrency = args[5]
		res.MerchantCU_date = args[6]
	}else{
		errMsg := "{ \"message\" : \""+ merchantId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Merchant json string manually
	order := 	`{`+
		`"merchantID": "` + res.MerchantID + `" , `+
		`"merchantUserName": "` + res.MerchantUserName + `" , `+
		`"merchantName": "` + res.MerchantName + `" , `+
		`"merchantIndustry": "` + res.MerchantIndustry + `" , `+ 
		`"pointsPerDollarSpent": "` + res.PointsPerDollarSpent + `" , `+ 
		`"merchantCurrency": "` + res.MerchantCurrency + `" , `+
		`"merchantCU_date": "` +  res.MerchantCU_date + `" `+ 
		`}`
	err = stub.PutState(merchantId, []byte(order))						//store Merchant with id as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"merchantId\" : \""+merchantId+"\", \"message\" : \"Merchant details updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Merchant details updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// create Merchant - create a new Merchant, store into chaincode state
// ============================================================================================================================
func (t *ManageMerchant) createMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 7 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 8\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start createMerchant")
	merchantID := args[0]
	merchantUserName := args[1]
	merchantName := args[2]
	merchantIndustry := args[3]
	pointsPerDollarSpent := args[4]
	merchantCurrency := args[5]
	merchantCU_date := args[6]
	merchantAsBytes, err := stub.GetState(merchantID)
	if err != nil {
		return nil, errors.New("Failed to get Merchant merchantID")
	}
	res := Merchant{}
	json.Unmarshal(merchantAsBytes, &res)
	fmt.Print("res: ")
	fmt.Println(res)
	if res.MerchantID == merchantID{
		fmt.Println("This Merchant arleady exists: " + merchantID)
		fmt.Println(res);
		errMsg := "{ \"message\" : \"This Merchant arleady exists\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil				//all stop a Merchant by this name exists
	}
	//build the Merchant json string manually
	merchant_json := 	`{`+
		`"merchantId": "` + merchantID + `" , `+
		`"merchantUserName": "` + merchantUserName + `" , `+
		`"merchantName": "` + merchantName + `" , `+
		`"merchantIndustry": "` + merchantIndustry + `" , `+ 
		`"pointsPerDollarSpent": "` + pointsPerDollarSpent + `" , `+ 
		`"merchantCurrency": "` + merchantCurrency + `" , `+
		`"merchantCU_date": "` + merchantCU_date + `" `+ 
	`}`
	fmt.Println("merchant_json: " + merchant_json)
	fmt.Print("merchant_json in bytes array: ")
	fmt.Println([]byte(merchant_json))
	err = stub.PutState(merchantID, []byte(merchant_json))		//store Merchant with merchantId as key
	if err != nil {
		return nil, err
	}
	//get the Merchant index
	merchantIndexAsBytes, err := stub.GetState(MerchantIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant index")
	}
	var merchantIndex []string
	json.Unmarshal(merchantIndexAsBytes, &merchantIndex)							//un stringify it aka JSON.parse()
	
	//append
	merchantIndex = append(merchantIndex, merchantID)									//add Merchant merchantID to index list
	fmt.Println("! Merchant index after appending merchantID: ", merchantIndex)
	jsonAsBytes, _ := json.Marshal(merchantIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(MerchantIndexStr, jsonAsBytes)						//store name of Merchant
	if err != nil {
		return nil, err
	}

	tosend := "{ \"merchantID\" : \""+merchantID+"\", \"message\" : \"Merchant created succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("end createMerchant")
	return nil, nil
}
