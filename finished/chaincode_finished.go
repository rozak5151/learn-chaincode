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
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

type Customer struct{
	PhoneNumber string `json:"PhoneNumber"`
	Operator string `json:"Operator"`
	Name string `json:"Name"`
	Email string `json:"Email"`
	Code string `json:"Code"`

	//access code
}

type TheSimple struct{

}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

 fmt.Println("LALALALALLALAAALAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")


	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "makecustomer" {
		return t.makecustomer(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "getcustomerdata" {
		return t.getcustomerdata(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}


func (t *SimpleChaincode) makecustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	// var key, value string
	// var err error
	// fmt.Println("running write()")
	//
	// if len(args) != 2 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	// }
	//
	// key = args[0] //rename for funsies
	// value = args[1]
	// err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	// if err != nil {
	// 	return nil, err
	// }
	// return nil, nil

	//PhoneNumber, Owner, CusomerName, Code, Email
	var customer_name, operator, code, email, phone_number string
	var err error
	phone_number = args[0]
	customer := new(Customer)
	_, err = stub.GetState(phone_number)

	if err == nil {
		return nil, errors.New("Customer already exists")
	}

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5. name of the key and value to set")
	}

	operator = args[1]
	customer_name = args[2]
	code = args[3]
	email = args[4]

  customer.Operator = operator
	customer.Name = customer_name
	customer.Code = code
	customer.Email = email
	customer.PhoneNumber = phone_number

  customerJSONBytes, err := json.Marshal(&customer)
	if err != nil {
		return nil, errors.New("Marshal operation went wrong")
	}

	err = stub.PutState(phone_number, customerJSONBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}


func (t *SimpleChaincode) getcustomerdata(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	customer := new(Customer)
	var phone_number, jsonResp string
	phone_number = args[0]
	var err error

	customerJSONBytes, err := stub.GetState(phone_number)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + phone_number + "\"}"
		return nil, errors.New(jsonResp)
	}

	json.Unmarshal(customerJSONBytes, &customer)

	return []byte(customer.Name), nil
}
