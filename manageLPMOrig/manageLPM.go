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

// ManageLPM example simple Chaincode implementation
type ManageLPM struct {
}

var CustomerIndexStr = "_Customerindex"				//name for the key/value that will store a list of all known Customer
var TransactionIndexStr = "_Transactionindex"		//name for the key/value that will store a list of all known Transaction
var MerchantIndexStr = "_Merchantindex"				//name for the key/value that will store a list of all known Merchant
var OwnerIndexStr = "_Ownerindex"					//name for the key/value that will store a list of all known Owner

var MerchantInitialBalance = "100000.00"
var StartingBalance = "100.00"

type Customer struct{							// Attributes of a Customer 
	CustomerID string `json:"customerId"`					
	UserName string `json:"userName"`
	CustomerName string `json:"customerName"`
	WalletWorth string `json:"walletWorth"`
	MerchantIDs string `json:"merchantIDs"`
	MerchantNames string `json:"merchantNames"`
	MerchantColors string `json:"merchantColors"`
	MerchantCurrencies string `json:"merchantCurrencies"`
	MerchantsPointsCount string `json:"merchantsPointsCount"`
	MerchantsPointsWorth string `json:"merchantsPointsWorth"`
}

type Transaction struct{							// Attributes of a Transaction 
	TransactionID string `json:"transactionId"`					
	TransactionDateTime string `json:"transactionDateTime"`
	TransactionType string `json:"transactionType"`				// Values are Purchase, Transfer, Accumulation (Add Points), CustomerOnBoarding
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
	IndustryColor string `json:"industryColor"`					
	PointsPerDollarSpent string `json:"pointsPerDollarSpent"`
	ExchangeRate string `json:"exchangeRate"`
	PurchaseBalance string `json:"purchaseBalance"`
	MerchantCurrency string `json:"merchantCurrency"`
	MerchantCU_date string `json:"merchantCU_date"`
	MerchantInitialBalance string `json:"merchantInitialBalance"`
}

type Owner struct{							// Attributes of a Owner
	OwnerID string `json:"ownerId"`					
	OwnerUserName string `json:"ownerUserName"`
	OwnerName string `json:"ownerName"`
}

// ============================================================================================================================
// Main - start the chaincode for LPM management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageLPM))
	if err != nil {
		fmt.Printf("Error starting LPM management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageLPM) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
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
	err = stub.PutState(MerchantIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	tosend := "{ \"message\" : \"ManageLPM chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *ManageLPM) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}
// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *ManageLPM) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {									//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "createCustomer" {				//create a new Customer
		return t.createCustomer(stub, args)
	}else if function == "updateCustomerAccumulationSC" {		//update a Customer - Add points
		return t.updateCustomerAccumulationSC(stub, args)
	}else if function == "updateCustomerPurchaseSC" {			//update a Customer - Purchase
		return t.updateCustomerPurchaseSC(stub, args)
	}else if function == "updateCustomerTransferSC" {			//update a Customer - Transfer
		return t.updateCustomerTransferSC(stub, args)
	}else if function == "updateCustomerAccumulation" {		//update a Customer - Add points
		return t.updateCustomerAccumulation(stub, args)
	}else if function == "updateCustomerPurchase" {			//update a Customer - Purchase
		return t.updateCustomerPurchase(stub, args)
	}else if function == "updateCustomerTransfer" {			//update a Customer - Transfer
		return t.updateCustomerTransfer(stub, args)
	}else if function == "deleteCustomer" {					//delete a Customer
		return t.deleteCustomer(stub, args)
	}else if function == "createMerchant" {					//create a new Merchant
		return t.createMerchant(stub, args)
	}else if function == "updateMerchant" {					//update a Merchant
		return t.updateMerchant(stub, args)
	}else if function == "deleteMerchant" {					//delete a Merchant
		return t.deleteMerchant(stub, args)
	}else if function == "createOwner" {					//create a owner
		return t.createOwner(stub, args)
	}else if function == "updateMerchantsPPDS" {			//update a Merchant's PPDS
		return t.updateMerchantsPPDS(stub, args)
	}else if function == "associateCustomer" {				//associate a customer to Merchant
		return t.associateCustomer(stub, args)
	}else if function == "updateMerchantsExchangeRate" {	//update a Merchant's Exchange Rate
		return t.updateMerchantsExchangeRate(stub, args)
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
func (t *ManageLPM) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getCustomerByID" {						//Read a Customer by Id
		return t.getCustomerByID(stub, args)
	}else if function == "getCustomerDetailsByID" {			//Read Customer Details by Id 
		return t.getCustomerDetailsByID(stub, args)
	}else if function == "getActivityHistory" {				//Read all transactions 
		return t.getActivityHistory(stub, args)
	}else if function == "getActivityHistoryForMerchant" {	//Read all transactions 
		return t.getActivityHistoryForMerchant(stub, args)
	}else if function == "getAllCustomers" {				//Read all Customers
		return t.getAllCustomers(stub, args)
	}else if function == "getCustomersByMerchantID" {		//Read a Customer by transId
		return t.getCustomersByMerchantID(stub, args)
	}else if function == "getMerchantByName" {				//Read Merchant by Name
		return t.getMerchantByName(stub, args)
	}else if function == "getMerchantByID" {				//Read Merchant by Id
		return t.getMerchantByID(stub, args)
	}else if function == "getMerchantDetailsByID" {			//Read Merchant details by Id
		return t.getMerchantDetailsByID(stub, args)
	}else if function == "getMerchantsByIndustry" {			//Read all Merchants by Industry
		return t.getMerchantsByIndustry(stub, args)
	}else if function == "getAllMerchants" {				//Read all Merchants
		return t.getAllMerchants(stub, args)
	}else if function == "getMerchantsAccountBalance" {		//Read Merchant Account Balance
		return t.getMerchantsAccountBalance(stub, args)
	}else if function == "getMerchantsUserCount" {			//Read Merchant's User Count
		return t.getMerchantsUserCount(stub, args)
	}else if function == "getOwnersMerchantUserCount" {		//Read Owner's Merchant and User Count
		return t.getOwnersMerchantUserCount(stub, args)
	}else if function == "getOwnerByID" {					//Read Owner by Id
		return t.getOwnerByID(stub, args)
	}
	fmt.Println("query did not find func: " + function)		//error
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
func (t *ManageLPM) getCustomerByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	fmt.Print("customerId in getCustomerByID: "+customerId)
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
// getCustomerDetailsByID - get Customer details for a specific ID from chaincode state POST Implementation
// ============================================================================================================================
func (t *ManageLPM) getCustomerDetailsByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var customerId string
	var err error
	fmt.Println("start getCustomerDetailsByID")
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
	fmt.Print("customerId in getCustomerDetailsByID : "+customerId)
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
	fmt.Println("end getCustomerDetailsByID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getActivityHistory - get Customer Transaction Activity details for a given customer from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getActivityHistory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var transactionIndex []string
	var valIndex Transaction
	var customerId string
	var err error
	var transactionTypeCustomerOnBoarding string
	fmt.Println("start getActivityHistory")
	
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'customerId' as argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	// set customerId
	customerId = args[0]
	fmt.Println("customerId in getActivityHistory::" + customerId)
	customerAsBytes, err := stub.GetState(customerId)
	if err != nil {
		return nil, errors.New("Failed to get Customer customerID")
	}
	res_Customer := Customer{}
	json.Unmarshal(customerAsBytes, &res_Customer)
		
	transactionAsBytes, err := stub.GetState(TransactionIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Transaction index string")
	}
	json.Unmarshal(transactionAsBytes, &transactionIndex)								//un stringify it aka JSON.parse()
	jsonResp = "{"
	transactionTypeCustomerOnBoarding = "CustomerOnBoarding"
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
		if valIndex.TransactionType == transactionTypeCustomerOnBoarding{
			if valIndex.CustomerID == customerId{
				fmt.Println("Customer found for CustomerOnBoarding")
				jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
				fmt.Println("jsonResp inside if for CustomerOnBoarding")
				fmt.Println(jsonResp)
				fmt.Println("transactionIndex for CustomerOnBoarding::")
				fmt.Println(transactionIndex)
				fmt.Println("length for CustomerOnBoarding::")
				fmt.Println(len(transactionIndex))
				if i < len(transactionIndex)-1 {
					fmt.Println("i for CustomerOnBoarding::")
					fmt.Println(i)
					jsonResp = jsonResp + ","
					fmt.Println("jsonResp inside if if for CustomerOnBoarding")
					fmt.Println(jsonResp)
				}
			}	
		} else if valIndex.TransactionFrom == res_Customer.UserName{
			fmt.Println("Customer found other than CustomerOnBoarding")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			fmt.Println("transactionIndex::")
			fmt.Println(transactionIndex)
			fmt.Println("length::")
			fmt.Println(len(transactionIndex))
			if i < len(transactionIndex)-1 {
				fmt.Println("i::")
				fmt.Println(i)
				jsonResp = jsonResp + ","
				fmt.Println("jsonResp inside if if")
				fmt.Println(jsonResp)
			}
        } 
	}
	jsonResp = jsonResp + "}"
	if strings.Contains(jsonResp, "},}"){
		fmt.Println("in if for jsonResp contains wrong json")	
		jsonResp = strings.Replace(jsonResp, "},}", "}}", -1)
	}
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getActivityHistory")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  getActivityHistoryForMerchant - get Customer Transaction Activity details for a given merchant from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getActivityHistoryForMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var transactionIndex []string
	var valIndex Transaction
	var merchantName string
	var transactionTypeCustomerOnBoarding string
	var err error
	fmt.Println("start getActivityHistoryForMerchant")
	
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantName' as argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantName
	merchantName = args[0]
	fmt.Println("merchantName in getActivityHistoryForMerchant::" + merchantName)

	transactionAsBytes, err := stub.GetState(TransactionIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Transaction index string")
	}
	json.Unmarshal(transactionAsBytes, &transactionIndex)								//un stringify it aka JSON.parse()
	jsonResp = "{"
	transactionTypeCustomerOnBoarding = "CustomerOnBoarding"
	for i,val := range transactionIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getActivityHistoryForMerchant")
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

		if valIndex.TransactionType == transactionTypeCustomerOnBoarding{
			if valIndex.TransactionFrom == merchantName{
				fmt.Println("Customer's merchant found for CustomerOnBoarding")
				jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
				fmt.Println("jsonResp inside if for CustomerOnBoarding")
				fmt.Println(jsonResp)
				fmt.Println("transactionIndex for CustomerOnBoarding::")
				fmt.Println(transactionIndex)
				fmt.Println("length for CustomerOnBoarding::")
				fmt.Println(len(transactionIndex))
				if i < len(transactionIndex)-1 {
					fmt.Println("i for CustomerOnBoarding::")
					fmt.Println(i)
					jsonResp = jsonResp + ","
					fmt.Println("jsonResp inside if if for CustomerOnBoarding")
					fmt.Println(jsonResp)
				}
			}
		} else if valIndex.TransactionTo == merchantName{
			fmt.Println("Customer's merchant found for other transactions")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			fmt.Println("transactionIndex::")
			fmt.Println(transactionIndex)
			fmt.Println("length::")
			fmt.Println(len(transactionIndex))
			if i < len(transactionIndex)-1 {
				fmt.Println("i::")
				fmt.Println(i)
				jsonResp = jsonResp + ","
				fmt.Println("jsonResp inside if if")
				fmt.Println(jsonResp)
			}
		}
	}
	jsonResp = jsonResp + "}"
	if strings.Contains(jsonResp, "},}"){
		fmt.Println("in if for jsonResp contains wrong json")	
		jsonResp = strings.Replace(jsonResp, "},}", "}}", -1)
	}
	fmt.Println("final jsonResp in getActivityHistoryForMerchant: " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getActivityHistoryForMerchant")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  getAllCustomers- get details of all Merchants from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getAllCustomers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var customerIndex []string
	var err error
	fmt.Println("start getAllCustomers")
		
	customerAsBytes, err := stub.GetState(CustomerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Customer index")
	}
	json.Unmarshal(customerAsBytes, &customerIndex)			//un stringify it aka JSON.parse()
	jsonResp = "{"
	for i,val := range customerIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Customer")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(customerIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	jsonResp = jsonResp + "}"
	if strings.Contains(jsonResp, "},}"){
		fmt.Println("in if for jsonResp contains wrong json")	
		jsonResp = strings.Replace(jsonResp, "},}", "}}", -1)
	}
	fmt.Println("jsonResp in getAllCustomers::")
	fmt.Println(jsonResp)

	fmt.Println("end getAllCustomers")
	return []byte(jsonResp), nil			//send it onward
}
// ============================================================================================================================
// getCustomersByMerchantID - get Customers for a specific Merchant ID from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getCustomersByMerchantID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
			/*if i < len(customerIndex)-1 {
				jsonResp = jsonResp + ","
			}*/
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
func (t *ManageLPM) getMerchantByName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	fmt.Println("merchantName" + merchantName)
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
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.MerchantName == merchantName{
			fmt.Println("Merchant found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			/*if i < len(merchantIndex)-1 {
				jsonResp = jsonResp + ","
			}*/
		} 
	}
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp in getMerchantByName: " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getMerchantByName")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
// getMerchantByID - get Merchant details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getMerchantByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
func (t *ManageLPM) getMerchantDetailsByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
func (t *ManageLPM) getMerchantsByIndustry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
		} 
	}
	jsonResp = jsonResp + "}"
	if strings.Contains(jsonResp, "},}"){
		fmt.Println("in if for jsonResp contains wrong json")	
		jsonResp = strings.Replace(jsonResp, "},}", "}}", -1)
	}
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getMerchantsByIndustry")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  getAllMerchants- get details of all Merchants from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getAllMerchants(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	if strings.Contains(jsonResp, "},}"){
		fmt.Println("in if for jsonResp contains wrong json")	
		jsonResp = strings.Replace(jsonResp, "},}", "}}", -1)
	}
	fmt.Println("jsonResp in getAllMerchants::")
	fmt.Println(jsonResp)

	fmt.Println("end getAllMerchants")
	return []byte(jsonResp), nil			//send it onward
}
// ============================================================================================================================
// getMerchantsAccountBalance - get merchants account balance from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getMerchantsAccountBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, merchantId, errResp string
	var err error
	var customerIndex []string
	accountBalance := float64(0.0)
	var valIndex Customer
	var merchantIndex Merchant
	var merchantIndexForPointsWorth int
	fmt.Println("start getMerchantsAccountBalance")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantId' as argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantName
	merchantId = args[0]
	
	// Get Merchants Balance from Merchant Struct START
	merchantAsbytes, err := stub.GetState(merchantId)
	if err != nil {
		errMsg := "{ \"message\" : \""+ merchantId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	json.Unmarshal(merchantAsbytes, &merchantIndex)
	fmt.Print("purchaseBalance for merchant : ")
	fmt.Print(merchantIndex.PurchaseBalance)
	accountBalanceMerchant, _ := strconv.ParseFloat(merchantIndex.PurchaseBalance, 64)
	merchantInitialBalance, _ := strconv.ParseFloat(merchantIndex.MerchantInitialBalance, 64)
	// Get Merchants Balance from Merchant Struct END

	customerAsBytes, err := stub.GetState(CustomerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Customer index string")
	}
	json.Unmarshal(customerAsBytes, &customerIndex)			//un stringify it aka JSON.parse()
	fmt.Print("customerIndex : ")
	fmt.Println(customerIndex)
	jsonResp = "{"
	for i,val := range customerIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getMerchantsAccountBalance")
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
			fmt.Println("Merchant found for Customer::"+valIndex.CustomerID)

			// find the index of the merchant to take the merchantPointsWorth in that index only
			fmt.Println("valIndex.MerchantIDs::"+valIndex.MerchantIDs)
			stringSliceMerchantIDs := strings.Split(valIndex.MerchantIDs, ",")
			for j,val := range stringSliceMerchantIDs{
				if val == merchantId{
					fmt.Println(strconv.Itoa(j) + " - looking at " + val + " for index")
					merchantIndexForPointsWorth = j
				}
			}
			// Sum all pointsWorth at the merchantIndex
			fmt.Println("valIndex.MerchantsPointsWorth::"+valIndex.MerchantsPointsWorth)
			stringSlice := strings.Split(valIndex.MerchantsPointsWorth, ",")
			for k,val := range stringSlice{
				if merchantIndexForPointsWorth == k{
					fmt.Println(strconv.Itoa(k) + " - looking at " + val + " for balance")
					valToBeAdded, _ := strconv.ParseFloat(val, 64)
					fmt.Println("accountBalance1::")
					fmt.Println(accountBalance)			
					fmt.Println("valToBeAdded::")
					fmt.Println(valToBeAdded)
	     		    		accountBalance = accountBalance + valToBeAdded
					fmt.Println("accountBalance2::")
					fmt.Println(accountBalance)			
	     			}
    			}
			/*fmt.Println("accountBalance3::")
			fmt.Println(accountBalance)
			fmt.Println("accountBalanceMerchant::")
			fmt.Println(accountBalanceMerchant)
    			accountBalance = accountBalance + accountBalanceMerchant
			fmt.Println("accountBalance4::")
			fmt.Println(accountBalance)*/
		} 
	}
	fmt.Println("accountBalance3::")
	fmt.Println(accountBalance)
	accountBalance = accountBalance + merchantInitialBalance + accountBalanceMerchant
	fmt.Println("accountBalance4::")
	fmt.Println(accountBalance)
	merchantInitialBalanceVar, _ := strconv.ParseFloat(MerchantInitialBalance, 64)
	amountFromCustomerOnBoarding :=  merchantInitialBalanceVar - merchantInitialBalance
	fmt.Println("amountFromCustomerOnBoarding::")
	fmt.Println(amountFromCustomerOnBoarding)
	accountBalance = accountBalance - amountFromCustomerOnBoarding
	fmt.Println("accountBalance5::")
	fmt.Println(accountBalance)
	jsonResp = jsonResp + "\"merchantAccountBalance\":" + strconv.FormatFloat(accountBalance, 'f', 2, 64)
	jsonResp = jsonResp + "}"
	if strings.Contains(jsonResp, "},}"){
		fmt.Println("in if for jsonResp contains wrong json")	
		jsonResp = strings.Replace(jsonResp, "},}", "}}", -1)
	}
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getMerchantsAccountBalance")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
// getMerchantsUserCount - get merchants user count from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getMerchantsUserCount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, merchantId, errResp string
	var err error
	var customerIndex []string
	var userCount = 0;
	var valIndex Customer
	fmt.Println("start getMerchantsUserCount")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantId' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantName
	merchantId = args[0]

	customerAsBytes, err := stub.GetState(CustomerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Customer index string")
	}
	json.Unmarshal(customerAsBytes, &customerIndex)			//un stringify it aka JSON.parse()
	fmt.Print("customerIndex : ")
	fmt.Println(customerIndex)
	jsonResp = "{"
	for i,val := range customerIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getMerchantsUserCount")
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
			fmt.Println("Merchant found")
			userCount++;
		} 
	}
	jsonResp = jsonResp + "\"merchantUsersCount\":" + strconv.Itoa(userCount)
	jsonResp = jsonResp + "}"
	
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getMerchantsUserCount")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
// getOwnerByID - get Owner details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getOwnerByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var ownerId string
	var err error
	fmt.Println("start getOwnerByID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'ownerId' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set ownerId
	ownerId = args[0]
	fmt.Print("ownerId in getOwnerByID : "+ownerId)
	valAsbytes, err := stub.GetState(ownerId)									//get the ownerId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \""+ ownerId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Print("valAsbytes : ")
	fmt.Println(valAsbytes)
	fmt.Println("end getOwnerByID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
// getOwnersMerchantUserCount - get owners merchants and users count from chaincode state
// ============================================================================================================================
func (t *ManageLPM) getOwnersMerchantUserCount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	var merchantIndex []string
	var customerIndex []string
	var merchantCount = 0;
	var userCount = 0;
	fmt.Println("start getOwnersMerchantUserCount")
	
	merchantAsBytes, err := stub.GetState(MerchantIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant index string")
	}
	json.Unmarshal(merchantAsBytes, &merchantIndex)			//un stringify it aka JSON.parse()
	fmt.Print("merchantIndex : ")
	fmt.Println(merchantIndex)
	jsonResp = "{"
	for i,val := range merchantIndex{
		fmt.Println("Merchant found")
		fmt.Println(strconv.Itoa(i) + " - looking at " + val)
		merchantCount++;
	}
	customerAsBytes, err := stub.GetState(CustomerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Customer index string")
	}
	json.Unmarshal(customerAsBytes, &customerIndex)			//un stringify it aka JSON.parse()
	fmt.Print("customerIndex : ")
	fmt.Println(customerIndex)
	jsonResp = "{"
	for j,valCustomer := range customerIndex{
		fmt.Println("Customer found")
		fmt.Println(strconv.Itoa(j) + " - looking at " + valCustomer)
		userCount++;
	}

	jsonResp = jsonResp + "\"merchantCount\":" + strconv.Itoa(merchantCount) + "," + "\"userCount\":" + strconv.Itoa(userCount)
	jsonResp = jsonResp + "}"
	
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getOwnersMerchantUserCount")
	return []byte(jsonResp), nil
}
// ============================================================================================================================
// create Customer - create a new Customer, store into chaincode state
// ============================================================================================================================
func (t *ManageLPM) createCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 13 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 13\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start createCustomer")
	customerId := args[0]
	userName := args[1]
	customerName := args[2]
	walletWorth := args[3]
	merchantID := args[4]
	merchantName := args[5]
	merchantColor := args[6]
	merchantCurrency := args[7]
	merchantsPointsCount := args[8]
	merchantsPointsWorth := args[9]
	transactionID := args[10]
 	transactionDateTime := args[11]
	transactionType := args[12]

	merchantAsBytes, err := stub.GetState(merchantID)
	if err != nil {
		return nil, errors.New("Failed to get Merchant merchantID")
	}
	res_Merchant := Merchant{}
	json.Unmarshal(merchantAsBytes, &res_Merchant)
	if res_Merchant.MerchantID == merchantID{
		fmt.Println("Merchant found with merchantID in createCustomer: " + merchantID)
		fmt.Println(res_Merchant);
	}else{
		errMsg := "{ \"message\" : \""+ merchantID+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	floatStartingBalance, _ := strconv.ParseFloat(StartingBalance, 64)
	floatInitialBalance, _ := strconv.ParseFloat(res_Merchant.MerchantInitialBalance, 64) 
	_merchantInitialBalance := floatInitialBalance - floatStartingBalance
	//build the Merchant json string manually
	merchant_json := 	`{`+
		`"merchantId": "` + res_Merchant.MerchantID + `" , `+
		`"merchantUserName": "` + res_Merchant.MerchantUserName + `" , `+
		`"merchantName": "` + res_Merchant.MerchantName + `" , `+
		`"merchantIndustry": "` + res_Merchant.MerchantIndustry + `" , `+
		`"industryColor": "` + res_Merchant.IndustryColor + `" , `+
		`"pointsPerDollarSpent": "` + res_Merchant.PointsPerDollarSpent + `" , `+ 
		`"exchangeRate": "` + res_Merchant.ExchangeRate + `" , `+ 
		`"purchaseBalance": "` + res_Merchant.PurchaseBalance + `" , `+
		`"merchantCurrency": "` + res_Merchant.MerchantCurrency + `" , `+ 
		`"merchantCU_date": "` + res_Merchant.MerchantCU_date + `" , `+ 
		`"merchantInitialBalance": "` + strconv.FormatFloat(_merchantInitialBalance, 'f', 2, 64) + `" `+ 
	`}`
	fmt.Println("merchant_json:::::::::::::::::::::::::::::::::::::::::::::::::::: " + merchant_json)
	fmt.Print("merchant_json in bytes array: ")
	fmt.Println([]byte(merchant_json))
	err = stub.PutState(merchantID, []byte(merchant_json))									//store Merchant with merchantId as key
	if err != nil {
		return nil, err
	}

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
		`"userName": "` + userName + `" , `+
		`"walletWorth": "` + walletWorth + `" , `+
		`"merchantIDs": "` + merchantID + `" , `+ 
		`"merchantNames": "` + merchantName + `" , `+ 
		`"merchantColors": "` + merchantColor + `" , `+
		`"merchantCurrencies": "` + merchantCurrency + `" , `+ 
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

	// build the Transaction json string manually
	transaction_json := `{`+
		`"transactionId": "` + transactionID + `" , `+
		`"transactionDateTime": "` + transactionDateTime + `" , `+
		`"transactionType": "` + transactionType + `" , `+
		`"transactionFrom": "` + merchantName + `" , `+ 
		`"transactionTo": "` + userName + `" , `+ 
		`"credit": "` + merchantsPointsWorth + `" , `+ 
		`"debit": "` + "0.00" + `" , `+ 
		`"customerId": "` +  customerId + `" `+ 
	`}`
	err = stub.PutState(transactionID, []byte(transaction_json))					//store Transaction with id as key
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
	transactionIndex = append(transactionIndex, transactionID)			//add Transaction transactionID to index list
	
	transactioJsonAsBytes, _ := json.Marshal(transactionIndex)
	fmt.Print("update transaction jsonAsBytes: ")
	fmt.Println(transactioJsonAsBytes)
	err = stub.PutState(TransactionIndexStr, transactioJsonAsBytes)						//store name of Transaction
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
// ============================================================================================================================
// Write - update customer during accumulation into chaincode state
// ============================================================================================================================
func (t *ManageLPM) updateCustomerAccumulation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Customer - accumulation")
	if len(args) != 11 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 11\", \"code\" : \"503\"}"
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
 	transactionId := args[4]
	json.Unmarshal(customerAsBytes, &res)
	if res.CustomerID == customerId{
		fmt.Println("Customer found with customerId : " + customerId)
		fmt.Println(res);
		res.WalletWorth = args[1]
		res.MerchantsPointsCount = args[2]
		res.MerchantsPointsWorth = args[3]
		res_trans.TransactionID = transactionId
 		res_trans.TransactionDateTime = args[5]
 		res_trans.TransactionType = args[6]
 		res_trans.TransactionFrom = args[7]
 		res_trans.TransactionTo = args[8]
 		res_trans.Credit = args[9]
 		res_trans.Debit = args[10]
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
		`"userName": "` + res.UserName + `" , `+
		`"customerName": "` + res.CustomerName + `" , `+
		`"walletWorth": "` + res.WalletWorth + `" , `+
		`"merchantIDs": "` + res.MerchantIDs + `" , `+
		`"merchantNames": "` + res.MerchantNames + `" , `+
		`"merchantColors": "` + res.MerchantColors + `" , `+
		`"merchantCurrencies": "` + res.MerchantCurrencies + `" , `+
		`"merchantsPointsCount": "` + res.MerchantsPointsCount + `" , `+ 
		`"merchantsPointsWorth": "` +  res.MerchantsPointsWorth + `" `+ 
	`}`
	err = stub.PutState(customerId, []byte(customer_json))									//store Customer with id as key
	if err != nil {
		return nil, err
	}

	// build the Transaction json string manually
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
	fmt.Print("update customer jsonAsBytes: ")
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
// Write - update customer during redemption into chaincode state
// ============================================================================================================================
func (t *ManageLPM) updateCustomerPurchase(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Customer - purchase")
	if len(args) != 20 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 20\", \"code\" : \"503\"}"
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
	res_trans1 := Transaction{}
 	res_trans2 := Transaction{}
 	res_Merchant := Merchant{}

  	transactionId1 := args[4]
 	transactionId2 := args[11]
	json.Unmarshal(customerAsBytes, &res)
	if res.CustomerID == customerId{
		fmt.Println("Customer found with customerId : " + customerId)
		fmt.Println(res);
		res.WalletWorth = args[1]
		res.MerchantsPointsCount = args[2]
		res.MerchantsPointsWorth = args[3]
		res_trans1.TransactionID = transactionId1
 		res_trans1.TransactionDateTime = args[5]
 		res_trans1.TransactionType = args[6]
 		res_trans1.TransactionFrom = args[7]
 		res_trans1.TransactionTo = args[8]
 		res_trans1.Credit = args[9]
 		res_trans1.Debit = args[10]
 		res_trans1.CustomerID = customerId
 		res_trans2.TransactionID = transactionId2
 		res_trans2.TransactionDateTime = args[12]
 		res_trans2.TransactionType = args[6]
 		res_trans2.TransactionFrom = args[13]
 		res_trans2.TransactionTo = args[14]
 		res_trans2.Credit = args[15]
 		res_trans2.Debit = args[16]
 		res_trans2.CustomerID = customerId
 		res_Merchant.MerchantID = args[17]
 		res_Merchant.PurchaseBalance = args[18]
 		res_Merchant.MerchantCU_date = args[19]
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
		`"userName": "` + res.UserName + `" , `+
		`"customerName": "` + res.CustomerName + `" , `+
		`"walletWorth": "` + res.WalletWorth + `" , `+
		`"merchantIDs": "` + res.MerchantIDs + `" , `+
		`"merchantNames": "` + res.MerchantNames + `" , `+
		`"merchantColors": "` + res.MerchantColors + `" , `+
		`"merchantCurrencies": "` + res.MerchantCurrencies + `" , `+
		`"merchantsPointsCount": "` + res.MerchantsPointsCount + `" , `+ 
		`"merchantsPointsWorth": "` +  res.MerchantsPointsWorth + `" `+ 
	`}`
	fmt.Println("customer_json:::::::::::::::::::::::::::::::::::::::::::::::::::: " + customer_json)
	err = stub.PutState(customerId, []byte(customer_json))							//store Customer with id as key
	if err != nil {
		return nil, err
	}

	// build the Transaction1 json string manually
 	transaction_json1 := `{`+
 		`"transactionId": "` + transactionId1 + `" , `+
 		`"transactionDateTime": "` + res_trans1.TransactionDateTime + `" , `+
 		`"transactionType": "` + res_trans1.TransactionType + `" , `+
 		`"transactionFrom": "` + res_trans1.TransactionFrom + `" , `+ 
 		`"transactionTo": "` + res_trans1.TransactionTo + `" , `+ 
 		`"credit": "` + res_trans1.Credit + `" , `+ 
 		`"debit": "` + res_trans1.Debit + `" , `+ 
 		`"customerId": "` +  res_trans1.CustomerID + `" `+ 
    `}`
	fmt.Println("transaction_json1:::::::::::::::::::::::::::::::::::::::::::::::::::: " + transaction_json1)
	err = stub.PutState(transactionId1, []byte(transaction_json1))					//store Transaction with id as key
	if err != nil {
		return nil, err
	}

	//get the Transaction index
	transactionAsBytes1, err := stub.GetState(TransactionIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Transaction index")
	}
	var transactionIndex1 []string	
	json.Unmarshal(transactionAsBytes1, &transactionIndex1)							//un stringify it aka JSON.parse()
	
	//append
	transactionIndex1 = append(transactionIndex1, transactionId1)					//add Transaction transactionId to index list
	
	jsonAsBytes1, _ := json.Marshal(transactionIndex1)
	fmt.Print("update customer jsonAsBytes1: ")
	fmt.Println(jsonAsBytes1)
	err = stub.PutState(TransactionIndexStr, jsonAsBytes1)							//store name of Transaction
	if err != nil {
		return nil, err
	}

	// build the Transaction2 json string manually
 	transaction_json2 := `{`+
 		`"transactionId": "` + transactionId2 + `" , `+
 		`"transactionDateTime": "` + res_trans2.TransactionDateTime + `" , `+
 		`"transactionType": "` + res_trans2.TransactionType + `" , `+
 		`"transactionFrom": "` + res_trans2.TransactionFrom + `" , `+ 
 		`"transactionTo": "` + res_trans2.TransactionTo + `" , `+ 
 		`"credit": "` + res_trans2.Credit + `" , `+ 
 		`"debit": "` + res_trans2.Debit + `" , `+ 
 		`"customerId": "` +  res_trans2.CustomerID + `" `+ 
 	`}`
	fmt.Println("transaction_json2:::::::::::::::::::::::::::::::::::::::::::::::::::: " + transaction_json2)
 	err = stub.PutState(transactionId2, []byte(transaction_json2))					//store Transaction with id as key
 	if err != nil {
 		return nil, err
 	}
 
 	//get the Transaction index
 	transactionAsBytes2, err := stub.GetState(TransactionIndexStr)
 	if err != nil {
 		return nil, errors.New("Failed to get Transaction index")
 	}
 	var transactionIndex2 []string	
 	json.Unmarshal(transactionAsBytes2, &transactionIndex2)							//un stringify it aka JSON.parse()
 	
 	//append
 	transactionIndex2 = append(transactionIndex2, transactionId2)					//add Transaction transactionId to index list
 	jsonAsBytes2, _ := json.Marshal(transactionIndex2)
 	fmt.Println(jsonAsBytes2)
 	err = stub.PutState(TransactionIndexStr, jsonAsBytes2)							//store name of Transaction
  	if err != nil {
  		return nil, err
  	}

	// update the Merchant START
	fmt.Println("res_Merchant.MerchantID in updateCustomerPurchase::"+res_Merchant.MerchantID)
	fmt.Println("res_Merchant.PurchaseBalance in updateCustomerPurchase::"+res_Merchant.PurchaseBalance)
	merchant_args := []string{res_Merchant.MerchantID, res_Merchant.PurchaseBalance, res_Merchant.MerchantCU_date}
	t.updateMerchantsPurchaseBal(stub, merchant_args)	// Call to Internal Function
 	// update the Merchant END

	tosend := "{ \"customerID\" : \""+customerId+"\", \"message\" : \"Customer details updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Customer details updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Write - update customer during transfer into chaincode state
// ============================================================================================================================
func (t *ManageLPM) updateCustomerTransfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Customer - transfer")
	if len(args) != 21 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 21\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set customerIds
	customerId1 := args[0]
	customerId2 := args[17]
	res1 := Customer{}
	res2 := Customer{}
	res_trans1 := Transaction{}
 	res_trans2 := Transaction{}
  	transactionId1 := args[4]
 	transactionId2 := args[11]
	
	customer1AsBytes, err := stub.GetState(customerId1)					//get the Customer for the specified customerId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get state for " + customerId1 + "\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	json.Unmarshal(customer1AsBytes, &res1)
	if res1.CustomerID == customerId1{
		fmt.Println("Customer found with customerId1 : " + customerId1)
		fmt.Println(res1);
		res1.WalletWorth = args[1]
		res1.MerchantsPointsCount = args[2]
		res1.MerchantsPointsWorth = args[3]
	}else{
		errMsg := "{ \"message\" : \""+ customerId1+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	res_trans1.TransactionID = transactionId1
	res_trans1.TransactionDateTime = args[5]
	res_trans1.TransactionType = args[6]
	res_trans1.TransactionFrom = args[7]
	res_trans1.TransactionTo = args[8]
	res_trans1.Credit = args[9]
	res_trans1.Debit = args[10]
	res_trans1.CustomerID = customerId1
	res_trans2.TransactionID = transactionId2
	res_trans2.TransactionDateTime = args[12]
	res_trans2.TransactionType = args[6]
	res_trans2.TransactionFrom = args[13]
	res_trans2.TransactionTo = args[14]
	res_trans2.Credit = args[15]
	res_trans2.Debit = args[16]
	res_trans2.CustomerID = customerId2

	customer2AsBytes, err := stub.GetState(customerId2)					//get the Customer for the specified customerId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get state for " + customerId2 + "\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	json.Unmarshal(customer2AsBytes, &res2)
	if res2.CustomerID == customerId2{
		fmt.Println("Customer found with customerId2 : " + customerId2)
		fmt.Println(res2);
		res2.WalletWorth = args[18]
		res2.MerchantsPointsCount = args[19]
		res2.MerchantsPointsWorth = args[20]
	}else{
		errMsg := "{ \"message\" : \""+ customerId2+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//build the Customer1 json string manually
	customer1_json := 	`{`+
		`"customerId": "` + res1.CustomerID + `" , `+
		`"userName": "` + res1.UserName + `" , `+
		`"customerName": "` + res1.CustomerName + `" , `+
		`"walletWorth": "` + res1.WalletWorth + `" , `+
		`"merchantIDs": "` + res1.MerchantIDs + `" , `+
		`"merchantNames": "` + res1.MerchantNames + `" , `+
		`"merchantColors": "` + res1.MerchantColors + `" , `+
		`"merchantCurrencies": "` + res1.MerchantCurrencies + `" , `+
		`"merchantsPointsCount": "` + res1.MerchantsPointsCount + `" , `+ 
		`"merchantsPointsWorth": "` +  res1.MerchantsPointsWorth + `" `+ 
	`}`
	fmt.Println("customer1_json:::::::::::::::::::::::::::::::::::::::::::::::::::: " + customer1_json)
	err = stub.PutState(customerId1, []byte(customer1_json))							//store Customer with id as key
	if err != nil {
		return nil, err
	}

	//build the Customer2 json string manually
	customer2_json := 	`{`+
		`"customerId": "` + res2.CustomerID + `" , `+
		`"userName": "` + res2.UserName + `" , `+
		`"customerName": "` + res2.CustomerName + `" , `+
		`"walletWorth": "` + res2.WalletWorth + `" , `+
		`"merchantIDs": "` + res2.MerchantIDs + `" , `+
		`"merchantNames": "` + res2.MerchantNames + `" , `+
		`"merchantColors": "` + res2.MerchantColors + `" , `+
		`"merchantCurrencies": "` + res2.MerchantCurrencies + `" , `+
		`"merchantsPointsCount": "` + res2.MerchantsPointsCount + `" , `+ 
		`"merchantsPointsWorth": "` +  res2.MerchantsPointsWorth + `" `+ 
	`}`
	fmt.Println("customer2_json:::::::::::::::::::::::::::::::::::::::::::::::::::: " + customer2_json)
	err = stub.PutState(customerId2, []byte(customer2_json))							//store Customer with id as key
	if err != nil {
		return nil, err
	}

	// build the Transaction1 json string manually
 	transaction_json1 := `{`+
 		`"transactionId": "` + transactionId1 + `" , `+
 		`"transactionDateTime": "` + res_trans1.TransactionDateTime + `" , `+
 		`"transactionType": "` + res_trans1.TransactionType + `" , `+
 		`"transactionFrom": "` + res_trans1.TransactionFrom + `" , `+ 
 		`"transactionTo": "` + res_trans1.TransactionTo + `" , `+ 
 		`"credit": "` + res_trans1.Credit + `" , `+ 
 		`"debit": "` + res_trans1.Debit + `" , `+ 
 		`"customerId": "` +  res_trans1.CustomerID + `" `+ 
    `}`
	fmt.Println("transaction_json1:::::::::::::::::::::::::::::::::::::::::::::::::::: " + transaction_json1)
	err = stub.PutState(transactionId1, []byte(transaction_json1))					//store Transaction with id as key
	if err != nil {
		return nil, err
	}

	//get the Transaction index
	transactionAsBytes1, err := stub.GetState(TransactionIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Transaction index")
	}
	var transactionIndex1 []string	
	json.Unmarshal(transactionAsBytes1, &transactionIndex1)							//un stringify it aka JSON.parse()
	
	//append
	transactionIndex1 = append(transactionIndex1, transactionId1)					//add Transaction transactionId to index list
	
	jsonAsBytes1, _ := json.Marshal(transactionIndex1)
	fmt.Print("update customer jsonAsBytes1: ")
	fmt.Println(jsonAsBytes1)
	err = stub.PutState(TransactionIndexStr, jsonAsBytes1)							//store name of Transaction
	if err != nil {
		return nil, err
	}

	// build the Transaction2 json string manually
 	transaction_json2 := `{`+
 		`"transactionId": "` + transactionId2 + `" , `+
 		`"transactionDateTime": "` + res_trans2.TransactionDateTime + `" , `+
 		`"transactionType": "` + res_trans2.TransactionType + `" , `+
 		`"transactionFrom": "` + res_trans2.TransactionFrom + `" , `+ 
 		`"transactionTo": "` + res_trans2.TransactionTo + `" , `+ 
 		`"credit": "` + res_trans2.Credit + `" , `+ 
 		`"debit": "` + res_trans2.Debit + `" , `+ 
 		`"customerId": "` +  res_trans2.CustomerID + `" `+ 
 	`}`
	fmt.Println("transaction_json2:::::::::::::::::::::::::::::::::::::::::::::::::::: " + transaction_json2)
 	err = stub.PutState(transactionId2, []byte(transaction_json2))					//store Transaction with id as key
 	if err != nil {
 		return nil, err
 	}
 
 	//get the Transaction index
 	transactionAsBytes2, err := stub.GetState(TransactionIndexStr)
 	if err != nil {
 		return nil, errors.New("Failed to get Transaction index")
 	}
 	var transactionIndex2 []string	
 	json.Unmarshal(transactionAsBytes2, &transactionIndex2)							//un stringify it aka JSON.parse()
 	
 	//append
 	transactionIndex2 = append(transactionIndex2, transactionId2)					//add Transaction transactionId to index list
 	jsonAsBytes2, _ := json.Marshal(transactionIndex2)
 	fmt.Println(jsonAsBytes2)
 	err = stub.PutState(TransactionIndexStr, jsonAsBytes2)							//store name of Transaction
  	if err != nil {
  		return nil, err
  	}

	tosend := "{ \"customerID\" : \""+customerId1+"\", \"message\" : \"Customer details updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Customer details updated succcessfully for transfer")
	return nil, nil
}
// ============================================================================================================================
// Write - update customer during accumulation into chaincode state - SmartContracts
// ============================================================================================================================
func (t *ManageLPM) updateCustomerAccumulationSC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Customer details updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Write - update customer during redemption into chaincode state - SmartContracts
// ============================================================================================================================
func (t *ManageLPM) updateCustomerPurchaseSC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Customer details updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Write - update customer during transfer into chaincode state - SmartContracts
// ============================================================================================================================
func (t *ManageLPM) updateCustomerTransferSC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Customer details updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Delete - remove a Customer and all his transactions from chain
// ============================================================================================================================
func (t *ManageLPM) deleteCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
// create Merchant - create a new Merchant, store into chaincode state
// ============================================================================================================================
func (t *ManageLPM) createMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 10 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 10\", \"code\" : \"503\"}"
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
	industryColor := args[4]
	pointsPerDollarSpent := args[5]
	exchangeRate := args[6]
	purchaseBalance := args[7]
	merchantCurrency := args[8]
	merchantCU_date := args[9]
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
	fmt.Println("MerchantInitialBalance::"+MerchantInitialBalance)
	//build the Merchant json string manually
	merchant_json := 	`{`+
		`"merchantId": "` + merchantID + `" , `+
		`"merchantUserName": "` + merchantUserName + `" , `+
		`"merchantName": "` + merchantName + `" , `+
		`"merchantIndustry": "` + merchantIndustry + `" , `+ 
		`"industryColor": "` + industryColor + `" , `+
		`"pointsPerDollarSpent": "` + pointsPerDollarSpent + `" , `+ 
		`"exchangeRate": "` + exchangeRate + `" , `+ 
		`"purchaseBalance": "` + purchaseBalance + `" , `+
		`"merchantCurrency": "` + merchantCurrency + `" , `+
		`"merchantCU_date": "` + merchantCU_date + `" , `+
		`"merchantInitialBalance": "` + MerchantInitialBalance + `" `+ 
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
// ============================================================================================================================
// Write - update merchant into chaincode state
// ============================================================================================================================
func (t *ManageLPM) updateMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Merchant")
	if len(args) != 10 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 10\", \"code\" : \"503\"}"
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
		res.IndustryColor = args[4]
		res.PointsPerDollarSpent = args[5]
		res.ExchangeRate = args[6]
		res.PurchaseBalance = args[7]
		res.MerchantCurrency = args[8]
		res.MerchantCU_date = args[9]
		res.MerchantInitialBalance = args[10]
	}else{
		errMsg := "{ \"message\" : \""+ merchantId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Merchant json string manually
	merchant := 	`{`+
		`"merchantId": "` + res.MerchantID + `" , `+
		`"merchantUserName": "` + res.MerchantUserName + `" , `+
		`"merchantName": "` + res.MerchantName + `" , `+
		`"merchantIndustry": "` + res.MerchantIndustry + `" , `+ 
		`"industryColor": "` + res.IndustryColor + `" , `+
		`"pointsPerDollarSpent": "` + res.PointsPerDollarSpent + `" , `+
		`"exchangeRate": "` + res.ExchangeRate + `" , `+ 
		`"purchaseBalance": "` + res.PurchaseBalance + `" , `+ 
		`"merchantCurrency": "` + res.MerchantCurrency + `" , `+
		`"merchantCU_date": "` + res.MerchantCU_date + `" , `+
		`"merchantInitialBalance": "` +  res.MerchantInitialBalance + `" `+ 
	`}`
	err = stub.PutState(merchantId, []byte(merchant))						//store Merchant with id as key
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
// Write - update merchant's purchase balance into chaincode state -- Internal Function
// ============================================================================================================================
func (t *ManageLPM) updateMerchantsPurchaseBal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Merchant - Purchase Balance")
	if len(args) != 3 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 3\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId := args[0]
	fmt.Println("merchantId in updateMerchantsPurchaseBal: " + merchantId)
	newPurchaseBal := args[1]
	fmt.Println("newPurchaseBal in updateMerchantsPurchaseBal : " + newPurchaseBal)
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
		fmt.Println("Merchants old purchaseBalance : " + res.PurchaseBalance)
		fmt.Println("Merchants new purchaseBalance : " + newPurchaseBal)

		currentPurchaseBalanceFloat, _ := strconv.ParseFloat(res.PurchaseBalance, 64)
		newPurchaseBalFloat, _ := strconv.ParseFloat(newPurchaseBal, 64) 
		purchaseBalanceCalculated := currentPurchaseBalanceFloat + newPurchaseBalFloat

		res.PurchaseBalance = strconv.FormatFloat(purchaseBalanceCalculated, 'f', 2, 64)
		res.MerchantCU_date = args[2]
	}else{
		errMsg := "{ \"message\" : \""+ merchantId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Merchant json string manually
	merchant_json := 	`{`+
		`"merchantId": "` + res.MerchantID + `" , `+
		`"merchantUserName": "` + res.MerchantUserName + `" , `+
		`"merchantName": "` + res.MerchantName + `" , `+
		`"merchantIndustry": "` + res.MerchantIndustry + `" , `+ 
		`"industryColor": "` + res.IndustryColor + `" , `+
		`"pointsPerDollarSpent": "` + res.PointsPerDollarSpent + `" , `+
		`"exchangeRate": "` + res.ExchangeRate + `" , `+ 
		`"purchaseBalance": "` + res.PurchaseBalance + `" , `+ 
		`"merchantCurrency": "` + res.MerchantCurrency + `" , `+
		`"merchantCU_date": "` + res.MerchantCU_date + `" , `+
		`"merchantInitialBalance": "` +  res.MerchantInitialBalance + `" `+ 
	`}`

	fmt.Println("merchant_json in updateMerchantsPurchaseBal::" + merchant_json)
		
	err = stub.PutState(merchantId, []byte(merchant_json))						//store Merchant with id as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"merchantId\" : \""+merchantId+"\", \"message\" : \"Merchant purchase balance details updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Merchant purchase balance updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Write - update merchant's PPDS into chaincode state
// ============================================================================================================================
func (t *ManageLPM) updateMerchantsPPDS(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Merchant - Points Per Dolllar Spent")
	if len(args) != 3 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 3\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId := args[0]
	newPPDS := args[1]
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
		fmt.Println("Merchants old pointsPerDollarSpent : " + res.PointsPerDollarSpent)
		fmt.Println("Merchants new pointsPerDollarSpent : " + newPPDS)
		res.PointsPerDollarSpent = newPPDS 
		res.MerchantCU_date = args[2]
	}else{
		errMsg := "{ \"message\" : \""+ merchantId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Merchant json string manually
	merchant_json := 	`{`+
		`"merchantId": "` + res.MerchantID + `" , `+
		`"merchantUserName": "` + res.MerchantUserName + `" , `+
		`"merchantName": "` + res.MerchantName + `" , `+
		`"merchantIndustry": "` + res.MerchantIndustry + `" , `+ 
		`"industryColor": "` + res.IndustryColor + `" , `+
		`"pointsPerDollarSpent": "` + res.PointsPerDollarSpent + `" , `+
		`"exchangeRate": "` + res.ExchangeRate + `" , `+ 
		`"purchaseBalance": "` + res.PurchaseBalance + `" , `+ 
		`"merchantCurrency": "` + res.MerchantCurrency + `" , `+
		`"merchantCU_date": "` + res.MerchantCU_date + `" , `+
		`"merchantInitialBalance": "` +  res.MerchantInitialBalance + `" `+ 
	`}`

	fmt.Println("merchant_json in updateMerchantsPPDS::" + merchant_json)
		
	err = stub.PutState(merchantId, []byte(merchant_json))						//store Merchant with id as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"merchantId\" : \""+merchantId+"\", \"message\" : \"Merchant points per dollar spent updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Merchant points per dollar spent updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Write - update merchant's Exchange Rate into chaincode state
// ============================================================================================================================
func (t *ManageLPM) updateMerchantsExchangeRate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Merchant - ExchangeRate")
	if len(args) != 3 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 3\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId := args[0]
	newExchangeRate := args[1]
	merchantAsBytes, err := stub.GetState(merchantId)				//get the Merchant for the specified merchant from chaincode state
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
		fmt.Println("Merchants old exchangeRate : " + res.ExchangeRate)
		fmt.Println("Merchants new exchangeRate : " + newExchangeRate)
		res.ExchangeRate = newExchangeRate
		res.MerchantCU_date = args[2]
	}else{
		errMsg := "{ \"message\" : \""+ merchantId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Merchant json string manually
	merchant_json := 	`{`+
		`"merchantId": "` + res.MerchantID + `" , `+
		`"merchantUserName": "` + res.MerchantUserName + `" , `+
		`"merchantName": "` + res.MerchantName + `" , `+
		`"merchantIndustry": "` + res.MerchantIndustry + `" , `+ 
		`"industryColor": "` + res.IndustryColor + `" , `+
		`"pointsPerDollarSpent": "` + res.PointsPerDollarSpent + `" , `+
		`"exchangeRate": "` + res.ExchangeRate + `" , `+ 
		`"purchaseBalance": "` + res.PurchaseBalance + `" , `+ 
		`"merchantCurrency": "` + res.MerchantCurrency + `" , `+
		`"merchantCU_date": "` + res.MerchantCU_date + `" , `+
		`"merchantInitialBalance": "` +  res.MerchantInitialBalance + `" `+ 
	`}`

	fmt.Println("merchant_json in updateMerchantsExchangeRate::" + merchant_json)
		
	err = stub.PutState(merchantId, []byte(merchant_json))						//store Merchant with id as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"merchantId\" : \""+merchantId+"\", \"message\" : \"Merchant exchange rate updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Merchant exchange rate updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Delete - remove a merchant from chain
// ============================================================================================================================
func (t *ManageLPM) deleteMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
// create Owner - create a Owner, store into chaincode state
// ============================================================================================================================
func (t *ManageLPM) createOwner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 3\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start createOwner")
	ownerId := args[0]
	ownerUserName := args[1]
	ownerName := args[2]
	
	ownerAsBytes, err := stub.GetState(ownerId)
	if err != nil {
		return nil, errors.New("Failed to get Owner ownerID")
	}
	res := Owner{}
	json.Unmarshal(ownerAsBytes, &res)
	if res.OwnerID == ownerId{
		errMsg := "{ \"message\" : \"This Owner arleady exists\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil				//all stop a Owner by this name exists
	}
	
	//build the Owner json string manually
	owner_json := 	`{`+
		`"ownerId": "` + ownerId + `" , `+
		`"ownerName": "` + ownerName + `" , `+
		`"ownerUserName": "` +  ownerUserName + `" `+
	`}`
	fmt.Println("owner_json: " + owner_json)
	fmt.Print("owner_json in bytes array: ")
	fmt.Println([]byte(owner_json))
	err = stub.PutState(ownerId, []byte(owner_json))									//store Owner with ownerId as key
	if err != nil {
		return nil, err
	}
	//get the Owner index
	ownerIndexAsBytes, err := stub.GetState(OwnerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Owner index")
	}
	var ownerIndex []string
	json.Unmarshal(ownerIndexAsBytes, &ownerIndex)							//un stringify it aka JSON.parse()
	
	//append
	ownerIndex = append(ownerIndex, ownerId)									//add Owner ownerID to index list
	
	jsonAsBytes, _ := json.Marshal(ownerIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(OwnerIndexStr, jsonAsBytes)						//store name of Owner
	if err != nil {
		return nil, err
	}

	tosend := "{ \"ownerID\" : \""+ownerId+"\", \"message\" : \"Owner created succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("end createOwner")
	return nil, nil
}
// ============================================================================================================================
// associate Customer - associate a customer to Merchant, store into chaincode state
// ============================================================================================================================
func (t *ManageLPM) associateCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var walletWorth, merchantIDs, merchantNames, merchantColors, merchantCurrencies, merchantsPointsCount, merchantsPointsWorth string
	if len(args) != 5 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 5\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start associateCustomer")
	customerId := args[0]
	merchantId := args[1]
		
	customerAsBytes, err := stub.GetState(customerId)
	if err != nil {
		return nil, errors.New("Failed to get Customer customerID")
	}
	
	merchantAsBytes, err := stub.GetState(merchantId)
	if err != nil {
		return nil, errors.New("Failed to get Merchant merchantID")
	}
	
	res := Customer{}
	res_Merchant := Merchant{}
	res_trans := Transaction{}
	
	json.Unmarshal(merchantAsBytes, &res_Merchant)
	
	if res_Merchant.MerchantID == merchantId{
		fmt.Println("Merchant found with merchantId in associateCustomer: " + merchantId)
		fmt.Println(res_Merchant);
	}else{
		errMsg := "{ \"message\" : \""+ merchantId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	// Calculation	
	floatStartingBalance, _ := strconv.ParseFloat(StartingBalance, 64)
	floatExchangeRate, _ := strconv.ParseFloat(res_Merchant.ExchangeRate, 64)
	pointsToBeCredited := floatStartingBalance / floatExchangeRate
	floatInitialBalance, _ := strconv.ParseFloat(res_Merchant.MerchantInitialBalance, 64) 
	_merchantInitialBalance := floatInitialBalance - floatStartingBalance
	fmt.Println("_merchantInitialBalance in associateCustomer::"+strconv.FormatFloat(_merchantInitialBalance, 'f', 2, 64))
	json.Unmarshal(customerAsBytes, &res)
	floatWalletWorth, _ := strconv.ParseFloat(res.WalletWorth, 64)
	newWalletWorth := floatWalletWorth + floatStartingBalance
	//fmt.Println("newWalletWorth in associateCustomer: " + strconv.FormatFloat(newWalletWorth, 'f', 2, 64))
	if res.CustomerID == customerId{
		fmt.Println("Customer found with customerId in associateCustomer: " + customerId)
		fmt.Println(res);
		walletWorth = strconv.FormatFloat(newWalletWorth, 'f', 2, 64)
		merchantIDs = res.MerchantIDs + "," + res_Merchant.MerchantID
		merchantNames = res.MerchantNames + "," + res_Merchant.MerchantName
		merchantColors = res.MerchantColors + "," + res_Merchant.IndustryColor
		merchantCurrencies = res.MerchantCurrencies + "," + res_Merchant.MerchantCurrency
		merchantsPointsCount = res.MerchantsPointsCount + "," + strconv.FormatFloat(pointsToBeCredited, 'f', 2, 64)
		merchantsPointsWorth = res.MerchantsPointsWorth + "," + StartingBalance
	
		res_trans.TransactionID = args[2]
 		res_trans.TransactionDateTime = args[3]
 		res_trans.TransactionType = args[4]
 		res_trans.TransactionFrom = res_Merchant.MerchantName
 		res_trans.TransactionTo = res.UserName
 		//res_trans.Credit = strconv.FormatFloat(pointsToBeCredited, 'f', 2, 64)
 		res_trans.Credit = strconv.FormatFloat(floatStartingBalance, 'f', 2, 64)
 		res_trans.Debit = "0.00"
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
		`"customerId": "` + customerId + `" , `+
		`"customerName": "` + res.CustomerName + `" , `+
		`"userName": "` + res.UserName + `" , `+
		`"walletWorth": "` + walletWorth + `" , `+
		`"merchantIDs": "` + merchantIDs + `" , `+ 
		`"merchantNames": "` + merchantNames + `" , `+ 
		`"merchantColors": "` + merchantColors + `" , `+
		`"merchantCurrencies": "` + merchantCurrencies + `" , `+ 
		`"merchantsPointsCount": "` + merchantsPointsCount + `" , `+ 
		`"merchantsPointsWorth": "` +  merchantsPointsWorth + `" `+ 
	`}`
	fmt.Println("customer_json:::::::::::::::::::::::::::::::::::::::::::::::::::: " + customer_json)
	fmt.Print("customer_json in bytes array: ")
	fmt.Println([]byte(customer_json))
	err = stub.PutState(customerId, []byte(customer_json))									//store Customer with customerId as key
	if err != nil {
		return nil, err
	}

	//build the Merchant json string manually
	merchant_json := 	`{`+
		`"merchantId": "` + res_Merchant.MerchantID + `" , `+
		`"merchantUserName": "` + res_Merchant.MerchantUserName + `" , `+
		`"merchantName": "` + res_Merchant.MerchantName + `" , `+
		`"merchantIndustry": "` + res_Merchant.MerchantIndustry + `" , `+
		`"industryColor": "` + res_Merchant.IndustryColor + `" , `+
		`"pointsPerDollarSpent": "` + res_Merchant.PointsPerDollarSpent + `" , `+ 
		`"exchangeRate": "` + res_Merchant.ExchangeRate + `" , `+ 
		`"purchaseBalance": "` + res_Merchant.PurchaseBalance + `" , `+
		`"merchantCurrency": "` + res_Merchant.MerchantCurrency + `" , `+ 
		`"merchantCU_date": "` + res_Merchant.MerchantCU_date + `" , `+ 
		`"merchantInitialBalance": "` + strconv.FormatFloat(_merchantInitialBalance, 'f', 2, 64) + `" `+ 
	`}`
	fmt.Println("merchant_json:::::::::::::::::::::::::::::::::::::::::::::::::::: " + merchant_json)
	fmt.Print("merchant_json in bytes array: ")
	fmt.Println([]byte(merchant_json))
	err = stub.PutState(merchantId, []byte(merchant_json))									//store Merchant with merchantId as key
	if err != nil {
		return nil, err
	}

	//build the Transaction json string manually
	transaction_json := `{`+
		`"transactionId": "` + res_trans.TransactionID + `" , `+
		`"transactionDateTime": "` + res_trans.TransactionDateTime + `" , `+
		`"transactionType": "` + res_trans.TransactionType + `" , `+
		`"transactionFrom": "` + res_trans.TransactionFrom + `" , `+ 
		`"transactionTo": "` + res_trans.TransactionTo + `" , `+ 
		`"credit": "` + res_trans.Credit + `" , `+ 
		`"debit": "` + res_trans.Debit + `" , `+ 
		`"customerId": "` +  res_trans.CustomerID + `" `+ 
	`}`
	fmt.Println("transaction_json:::::::::::::::::::::::::::::::::::::::::::::::::::: " + transaction_json)
	err = stub.PutState(res_trans.TransactionID, []byte(transaction_json))					//store Transaction with id as key
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
	transactionIndex = append(transactionIndex, res_trans.TransactionID)			//add Transaction res_trans.TransactionID to index list
	
	jsonAsBytes, _ := json.Marshal(transactionIndex)
	fmt.Print("update transaction jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(TransactionIndexStr, jsonAsBytes)						//store name of Transaction
	if err != nil {
		return nil, err
	}

	tosend := "{ \"customerID\" : \""+customerId+"\", \"message\" : \"Customer associated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("end associateCustomer")
	return nil, nil
}
