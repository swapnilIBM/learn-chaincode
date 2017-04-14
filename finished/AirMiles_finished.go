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


func (t *AirMilesChaincode) addtrip(tripJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.adduser start ")
	
	//var usr UserDetails
	//var users UserDetails
	var mdet MilesDetails
	trip := &TripDetails{}
	
	
	err := json.Unmarshal([]byte(tripJSON), trip)
	if err != nil {
		fmt.Println("Failed to unmarshal trip ")
	}	
	
		
	fmt.Println("AirMilesID ID : ", trip.AirMilesID)
	
	trip.TripID = trip.DepartureTime[:len(trip.DepartureTime)-4]
	
	mdBytes, err := stub.GetState(trip.AirMilesID)
	err = json.Unmarshal(mdBytes, &mdet)
	
	//var PointBalanceI int
	//var PointsRewardedI int
	//var PointsConsumedI int
	
	PointBalanceI,_ := strconv.Atoi(mdet.PointBalance)
	PointsRewardedI,_ := strconv.Atoi(trip.PointsRewarded)
	PointsConsumedI,_ := strconv.Atoi(trip.PointsConsumed)
	mdet.PointBalance= strconv.Itoa(PointBalanceI + PointsRewardedI - PointsConsumedI)
	mdet.UpdatedDate = trip.DepartureTime
	
	 	
	body, err := json.Marshal(mdet)
	if err != nil {
        panic(err)
    }
    fmt.Println(string(body))	
	err = stub.PutState(mdet.AirMilesID, []byte(string(body)))
	if err != nil {
		fmt.Println("Failed to update miles balance ")
	}
	body1, err := json.Marshal(trip)
	
	if err != nil {
        panic(err)
    }
    fmt.Println(string(body1))	
	err = stub.PutState(trip.AirMilesID+"_"+trip.TripID, []byte(string(body1)))
	if err != nil {
		fmt.Println("Failed to create miles details ")
	}
			
	fmt.Println("Created trip  with Key : "+ trip.TripID)
	fmt.Println("In initialize.adduser end ")
	return nil,nil	
	
}

func (t *AirMilesChaincode) gettripdetails(userID string, traveldate string, stub shim.ChaincodeStubInterface)([]TripDetails, error) {
	fmt.Println("In query.gettripdetails start ")

	key := userID
	tdate := traveldate
	var users UserDetails
	var mdet MilesDetails
	var trip TripDetails
	var triparray []TripDetails
	var milesid string
	var bytemilesid []byte
	var bytetrip []byte
	var hours []string
	hours := []string{"00","01", "02", "03","04","05", "06", "07","08","09", "10", "11","12","13", "14", "15","16","17", "18", "19","20","21", "22", "23"}
	bytemilesid,err := t.getmilesid(key,stub);
	if err != nil {
		fmt.Println("Error retrieving trip details for user" , key)
		return nil, errors.New("Error retrieving trip details for user" + key)
	}
	milesid = string(bytemilesid);
	
	j := 0
	for i := 0; i < 24; i++ {
		
		bytetrip,err := stub.GetState(milesid + "_"+tdate+ hours[i]);
		err = json.Unmarshal(bytetrip, &trip)
		if err != nil {
			fmt.Printf("Error while unmarshaling the trip : %s", err)
		} else {
			TripDetails[j] = trip
			j += 1
		}
		
	}
	
	return TripDetails, nil
}

func (t *AirMilesChaincode) adduser(userJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.adduser start ")
	
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
	fmt.Println("AirMilesID is : ", strconv.FormatInt(int64(secs), 10))
	 
	md.AirMilesID = strconv.FormatInt(int64(secs), 10)
	md.PointBalance	= usr.AirMilesPoint
	md.CreatedDate = now.String()
	md.UpdatedDate = now.String()
	
	usr.AirMilesID = md.AirMilesID 	
	body, err := json.Marshal(usr)
	if err != nil {
        panic(err)
    }
    fmt.Println(string(body))	
	err = stub.PutState(usr.UserID + "_" + usr.PhoneNumber, []byte(string(body)))
	if err != nil {
		fmt.Println("Failed to create User ")
	}
	//Storing miles id
	err = stub.PutState(usr.UserID, []byte(string(usr.AirMilesID)))
	if err != nil {
		fmt.Println("Failed to set Miles ID")
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
	fmt.Println("In initialize.adduser end ")
	return nil,nil	
	
}


func (t *AirMilesChaincode) getmilesid(userID string, stub shim.ChaincodeStubInterface)([]byte, error) {
	fmt.Println("In query.GetUsers start ")

	key := userID
//	var users UserDetails
	//var mdet MilesDetails
	//var balance string
	
	bytemilesid, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving milesid for " , userID)
		return nil, errors.New("Error retrieving milesid for" + userID)
	}
	
	return bytemilesid, nil
}


func (t *AirMilesChaincode) getbalance(userID string, stub shim.ChaincodeStubInterface)([]byte, error) {
	fmt.Println("In query.GetUsers start ")

	key := userID
	var users UserDetails
	var mdet MilesDetails
	var balance string
	var byteBalance []byte
	userBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving Users" , userID)
		return nil, errors.New("Error retrieving Users" + userID)
	}
	err = json.Unmarshal(userBytes, &users)
	fmt.Println("Users   : " , users);
	fmt.Println("In query.GetUsers end ")
	
	mdBytes, err := stub.GetState(users.AirMilesID)
	err = json.Unmarshal(mdBytes, &mdet)
	
	balance = mdet.PointBalance
	
	byteBalance = []byte(balance)
	
	return byteBalance, nil
}


// Invoke isur entry point to invoke a chaincode function
func (t *AirMilesChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "adduser" {
		var testBytes []byte
		fmt.Println("invoking adduser " + function)
		testBytes,err := t.adduser(args[0],stub)
		if err != nil {
			fmt.Println("Error performing adduser ")
			return nil, err
		}
		fmt.Println("Processed adduser successfully. ")
		return testBytes, nil
	} else if function == "addtrip" {
		var testBytes []byte
		fmt.Println("invoking addtrip " + function)
		testBytes,err := t.addtrip(args[0],stub)
		if err != nil {
			fmt.Println("Error performing addtrip ")
			return nil, err
		}
		fmt.Println("Processed addtrip successfully. ")
		return testBytes, nil
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
	} else if function == "getbalance" { //Get a miles point balance
		return t.getbalance(args[0] + "_" + args[1],stub)
	} else if function == "getmilesid" { //Get a miles id 
		return t.getmilesid(args[0],stub)
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
