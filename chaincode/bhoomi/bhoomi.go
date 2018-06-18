/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


package main

// ============================================================================================================================
//  Imports
// 3 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
// 2 specific Hyperledger Fabric specific libraries for Smart Contracts
// ============================================================================================================================
import (
	"encoding/json"
	"fmt"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

//
// ============================================================================================================================
//  Define the landRecord structure, with  properties.  Structure tags are used by encoding/json library
// ============================================================================================================================
type GeoData struct {
	Latitude string `json:"latitude"`
        Longitude string `json:"longitude"`
        Length string `json:"length"`
        Width string `json:"width"`
        TotalArea string `json:"totalArea"`
        Address string `json:"address"`

}

type Owner struct{
      	OwnerName string `json:"ownerName"`
        Gender string `json:"gender"`
        AadharNo string `json:"aadharNo"`
        MobileNo string `json:"mobileNo"`
        EmailID string `json:"emailID"` 
        Address string `json:"address"`

}


type LandRecord struct {
	Pid   string `json:"pid"`
	WardNo string `json:"wardNo"`
	AreaCode string `json:"areaCode"`
	SiteNo  string `json:"siteNo"`
    GeoData GeoData `json:"geoData"`
	Owner Owner `json:"owner"`	
	Allotee Owner `json:"owner"`

}

// ============================================================================================================================
// Get LandRecord - get a land record asset from ledger
// ============================================================================================================================
func get_landRecord(stub shim.ChaincodeStubInterface, id string) (LandRecord, error) {
	var landRecord LandRecord
	landRecordAsBytes, err := stub.GetState(id)                  //getState retreives a key/value from the ledger
	if err != nil {                                          
		return landRecord, errors.New("Failed to find the landRecord - " + id)
	}
	json.Unmarshal(landRecordAsBytes, &landRecord)                   //un stringify it aka JSON.parse()

	if landRecord.Pid != id {                                     //test if landrecord is actually here or just nil
		return landRecord, errors.New("Land Record does not exist - " + id)
	}

	return landRecord, nil
}


// ============================================================================================================================
//  The Init method is called when the Smart Contract "bhoomi" is instantiated by the blockchain network
// Best practice is to have any Ledger initialization in separate function -- see initLedger()
// ============================================================================================================================
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	s.initLedger(stub)
	return shim.Success(nil)
}


 // ============================================================================================================================
//  The Invoke method is called as a result of an application request to run the Smart Contract "bhoomi"
//  The calling application program has also specified the particular smart contract function to be called, with arguments
// ============================================================================================================================
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryLandRecord" {
		return s.queryLandRecord(stub, args)
	} else if function == "initLedger" {
		return s.initLedger(stub)
	} else if function == "createLandRecord" {
		return s.createLandRecord(stub, args)
	} else if function == "transferLandRecord" {
		return s.transferLandRecord(stub, args)
	} else if function == "allotNewLandRecord" {
		return s.allotNewLandRecord(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}
 // ============================================================================================================================
//  Query Ledger for land record with PID
// ============================================================================================================================
func (s *SmartContract) queryLandRecord(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	} 
        
        // args[0] = PID
	landAsBytes, _ := stub.GetState(args[0])  
	return shim.Success(landAsBytes)
}

 // ============================================================================================================================
//  Initialize Ledger for land record with PID
// ============================================================================================================================
func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface) sc.Response {

        var owner Owner
        owner.OwnerName = "CTS"
        owner.Gender = "NA"
        owner.AadharNo = "123456789101"
        owner.MobileNo = "9999999999"
        owner.EmailID = "cts@cts.com"
        owner.Address = "cts,Bangalore"

        var geoDetails GeoData
        geoDetails.Latitude = "99.9999"
        geoDetails.Longitude = "99.9999"
        geoDetails.Length = "99"
        geoDetails.Width = "99"
        geoDetails.TotalArea = "9999"
        geoDetails.Address = "CTS, Bangalore-560043"
          
        //Initialize ledger with one sample record for unit testing
 	landRecords := []LandRecord{LandRecord{Pid: "999999999", WardNo:"999", AreaCode: "999", SiteNo: "999", GeoData:geoDetails, Owner : owner}}
	i := 0 
	for i < len(landRecords) {
		fmt.Println("i is ", i)
		landAsBytes, _ := json.Marshal(landRecords[i]) //convert to bytes = VALUE
		stub.PutState(landRecords[i].Pid, landAsBytes) //Pid is unique KEY
		fmt.Println("Added", landRecords[i])
		i = i + 1
	}

	return shim.Success(nil)
}

 // ============================================================================================================================
//  Create a new Land Record on the ledger
// ============================================================================================================================
func (s *SmartContract) createLandRecord(stub shim.ChaincodeStubInterface, args []string) sc.Response {
        var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

    
         id := args[0]     //args[0] = Pid = KEY -- Declare and assign using :=

	//check if land record with this Pid already exists
	landRec, err := get_landRecord(stub, id)
	if err == nil {
		fmt.Println("This land record already exists - " + id)
		fmt.Println(landRec)
		return shim.Error("This land record already exists - " + id)  // Halt as land record already exists
	}
      //Initialize values for other objects as empty string in the ledger
        var owner Owner
        owner.OwnerName = ""
        owner.Gender = ""
        owner.AadharNo = ""
        owner.MobileNo = ""
        owner.EmailID = ""
        owner.Address = ""

        var geoDetails GeoData
        geoDetails.Latitude = ""
        geoDetails.Longitude = ""
        geoDetails.Length = ""
        geoDetails.Width = ""
        geoDetails.TotalArea = ""
        geoDetails.Address = ""
	
	var landRecord = LandRecord{Pid: id, WardNo: args[1], AreaCode: args[2], SiteNo: args[3], GeoData: geoDetails , Owner: owner }
	landAsBytes, _ := json.Marshal(landRecord) //convert to bytes = VALUE
	stub.PutState(id, landAsBytes)  //update to ledger

	return shim.Success(nil)
}

 // ============================================================================================================================
//  Update Geo Details on the ledger
// ============================================================================================================================
func (s *SmartContract) updateGeoDetails(stub shim.ChaincodeStubInterface, args []string) sc.Response {
        var err error
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	id := args[0]     //args[0] = Pid = KEY

	//check if land record with this Pid already exists
	landRecord, err := get_landRecord(stub, id)
	if err != nil {         //halt if error
		return shim.Error("Failed to find the landRecord - " + id)
	}
        
        //assign the values for the geoDetails object
        landRecord.GeoData.Latitude = args[1]
        landRecord.GeoData.Longitude = args[2]
        landRecord.GeoData.Length = args[3]
        landRecord.GeoData.Width = args[4]
        landRecord.GeoData.TotalArea = args[5]
        landRecord.GeoData.Address = args[6]


	landAsBytes, _ := json.Marshal(landRecord) // convert to bytes aka Json.stringify()
	stub.PutState(id, landAsBytes) //update ledger

	return shim.Success(nil)
}

 // ============================================================================================================================
//  Change ownership of the land record on the ledger
// ============================================================================================================================
func (s *SmartContract) transferLandRecord(stub shim.ChaincodeStubInterface, args []string) sc.Response {
        var err error
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	id := args[0]     //args[0] = Pid = KEY

	//check if land record with this Pid already exists
	landRecord, err := get_landRecord(stub, id)
	if err != nil {         //halt if error
		return shim.Error("Failed to find the landRecord - " + id)
	}
        
        //assign the values for the new owner object
      	landRecord.Owner.OwnerName = args[1]
        landRecord.Owner.Gender = args[2]
        landRecord.Owner.AadharNo = args[3]
        landRecord.Owner.MobileNo = args[4]
        landRecord.Owner.EmailID = args[5]
        landRecord.Owner.Address = args[6]

	landAsBytes, _ := json.Marshal(landRecord) // convert to bytes aka Json.stringify()
	stub.PutState(id, landAsBytes) //update ledger

	return shim.Success(nil)
}

 // ============================================================================================================================
//  Allot new land record to a owner
// ============================================================================================================================
func (s *SmartContract) allotNewLandRecord(stub shim.ChaincodeStubInterface, args []string) sc.Response {
        var err error
	if len(args) != 13 {
		return shim.Error("Incorrect number of arguments. Expecting 13")
	}

	id := args[0] //args[0] = Pid = KEY

	//check if land record with this Pid already exists
	landRecord, err := get_landRecord(stub, id)
	if err != nil {         //halt if error
		return shim.Error("Failed to find the landRecord - " + id)
	}
        
        //assign the values for the new owner object
      	landRecord.Allotee.OwnerName = args[1]
        landRecord.Allotee.Gender = args[2]
        landRecord.Allotee.AadharNo = args[3]
        landRecord.Allotee.MobileNo = args[4]
        landRecord.Allotee.EmailID = args[5]
        landRecord.Allotee.Address = args[6]

        landRecord.GeoData.Latitude = args[7]
        landRecord.GeoData.Longitude = args[8]
        landRecord.GeoData.Length = args[9]
        landRecord.GeoData.Width = args[10]
        landRecord.GeoData.TotalArea = args[11]
        landRecord.GeoData.Address = args[12]

	landAsBytes, _ := json.Marshal(landRecord) // convert to bytes aka Json.stringify()  -- _ is to skip reading the first return val
	stub.PutState(id, landAsBytes) //update ledger

	return shim.Success(nil)
}

// ============================================================================================================================
//  Main function for unit testing
// ============================================================================================================================
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

