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
	"strconv"
	"time"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

// AirMilesChaincode example simple Chaincode implementation
type AirMilesChaincode struct {
}

type UserDetails struct{
	UserID string `json:"UserID"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	PhoneNumber string `json:"PhoneNumber"`
	AirMilesID string `json:"AirMilesID"`
	UserType string `json:"UserType"`
	AirMilesPoint string `json:"AirMilesPoint"`
}


type MilesDetails struct{
	UserID string `json:"UserID"`
	AirMilesID string `json:"AirMilesID"`
	PointBalance string `json:"PointBalance"`
	CreatedDate string `json:"CreatedDate"`
	UpdatedDate string `json:"UpdatedDate"`
}

type TripDetails struct{
	TripID string `json:"TripID"`
	AirMilesID string `json:"AirMilesID"`
	Airlines string `json:"Airlines"`
	DepartureLocation string `json:"DepartureLocation"`
	ArrivalLocation string `json:"ArrivalLocation"`
	DepartureTime string `json:"DepartureTime"`
	ArrivalTime string `json:"ArrivalTime"`
	IsPartnerAirlines string `json:"IsPartnerAirlines"`
	PointsConsumed string `json:"PointsConsumed"`
	PointsRewarded string `json:"PointsRewarded"`
}


func main() {
	err := shim.Start(new(AirMilesChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *AirMilesChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func AddUser(userJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.AddUser start ")
	
	//var usr UserDetails
	usr := &UserDetails{}
	md := &MilesDetails{}
	err := json.Unmarshal([]byte(userJSON), usr)
	if err != nil {
		fmt.Println("Failed to unmarshal user ")
	}	
	
		
	fmt.Println("User ID : ",usr.UserID)
	
	usr.AirMilesPoint = "100"
	
	md.UserID = usr.UserID;
	now := time.Now()
    secs := now.Unix()
	fmt.Println("AirMilesID is : ", strconv.Itoa(secs))
	 
	md.AirMilesID = strconv.Itoa(secs)
	md.PointBalance	= usr.AirMilesPoint
	md.CreatedDate = now
	md.UpdatedDate = now
	
	usr.AirMilesID = md.AirMilesID 	
	body, err := json.Marshal(usr)
	if err != nil {
        panic(err)
    }
    fmt.Println(string(body))	
	err = stub.PutState(usr.UserID + "_" + usr.UserType, []byte(string(body)))
	if err != nil {
		fmt.Println("Failed to create User ")
	}
	body1, err := json.Marshal(md)
	
	if err != nil {
        panic(err)
    }
    fmt.Println(string(body1))	
	err = stub.PutState(md.AirMilesID, []byte(string(body1)))
	if err != nil {
		fmt.Println("Failed to create miles details ")
	}
		
	
	fmt.Println("Created User  with Key : "+ usr.UserID)
	fmt.Println("In initialize.AddUser end ")
	return nil,nil	
	
}

func GetBalance(userID string, stub shim.ChaincodeStubInterface)([]byte, error) {
	fmt.Println("In query.GetUsers start ")

	key := userID
	var users UserDetails
	var mdet MilesDetails
	userBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving Users" , userID)
		return users, errors.New("Error retrieving Users" + userID)
	}
	err = json.Unmarshal(userBytes, &users)
	fmt.Println("Users   : " , users);
	fmt.Println("In query.GetUsers end ")
	
	mdBytes, err := stub.GetState(users.AirMilesID)
	err = json.Unmarshal(mdBytes, &mdet)
	
	balance = mdet.PointBalance
	
	return []byte(balance), nil
}


// Invoke isur entry point to invoke a chaincode function
func (t *AirMilesChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "AddUser" {
		fmt.Println("invoking AddUser " + function)
		testBytes,err := AddUser(args[0],stub)
		if err != nil {
			fmt.Println("Error performing AddUser ")
			return nil, err
		}
		fmt.Println("Processed AddUser successfully. ")
		return testBytes, nil
	}
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *AirMilesChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "GetBalance" { //read a variable
		return t.GetBalance(args,stub)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *AirMilesChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
func (t *AirMilesChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
