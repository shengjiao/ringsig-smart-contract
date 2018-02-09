package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//PUBLICKEY ...
const PUBLICKEY = "PUBLICKEY"

var (
	DefaultCurve = elliptic.P256()
)

var logger = shim.NewLogger("Voting")

//Voting chaincode implementation
type sc struct {
}

// Init method for chaincode
func (t *sc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("Init")
	logger.Info("smart contract version 1")
	_, args := stub.GetFunctionAndParameters()

	if len(args) < 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0.")
	}
	return shim.Success(nil)
}

//{x:big.Int,y:big.Int}
func (t *sc) initPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//checking the number of argument
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}
	recBytes := args[0]
	var params map[string]interface{}
	err := json.Unmarshal([]byte(recBytes), &params)
	if err != nil {
		return shim.Error("Failed to unmarshal key Records.")
	}
	if len(params) != 3 {
		return shim.Error("Incorrect number of fields in params. Expecting 3.")
	}
	uid := getSafeString(params["uid"])
	x, _ := new(big.Int).SetString(getSafeString(params["x"]), 10)
	y, _ := new(big.Int).SetString(getSafeString(params["y"]), 10)

	logger.Info("initPublicKey: uid=", uid)
	logger.Info("initPublicKey: x=", x)
	logger.Info("initPublicKey: y=", y)
	pk := ecdsa.PublicKey{DefaultCurve, x, y}
	logger.Info("initPublicKey: pk=", pk)

	//Store the user key records
	err = stub.PutState(PUBLICKEY+":"+uid, []byte(recBytes))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("SUCCESS"))
}

//check whether string has value or not
func getSafeString(input interface{}) string {
	var safeValue string
	var isOk bool

	if input == nil {
		safeValue = ""
	} else {
		safeValue, isOk = input.(string)
		if isOk == false {
			safeValue = ""
		}
	}
	return safeValue
}

//invoke
func (t *sc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//myLogger.Debug("Invoke Chaincode...")
	function, args := stub.GetFunctionAndParameters()
	if function == "initPublicKey" {
		//request a new voyage creation
		return t.initPublicKey(stub, args)
	}
	return shim.Error("Invalid invoke function name.")
}

func main() {
	logger.SetLevel(shim.LogInfo)

	err := shim.Start(new(sc))
	if err != nil {
		fmt.Printf("Error starting Voting: %s", err)
	}
}
