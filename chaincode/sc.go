package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//PUBLICKEY ...
const PUBLICKEY = "PUBLICKEY"

//KEYRING ...
const KEYRING = "KEYRING"

//TOPIC ...
var TOPIC = struct {
	stage map[string]string
}{
	stage: map[string]string{
		"prepare": "prepare",
		"start":   "start",
		"finish":  "finish",
	},
}

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

//{topic:string, stage:start}
//stage: prepare(one can init public key), start(one can submit tx for topic) or finish (no more action is allowed)
func (t *sc) configureTopic(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}
	recBytes := args[0]
	var params map[string]interface{}
	err := json.Unmarshal([]byte(recBytes), &params)
	if err != nil {
		return shim.Error("Failed to unmarshal topic.")
	}
	if len(params) != 2 {
		return shim.Error("Incorrect number of fields in params. Expecting 2.")
	}
	topic := getSafeString(params["topic"])
	stage := getSafeString(params["stage"])
	logger.Info("configureTopic: topic=", topic)
	logger.Info("configureTopic: stage=", stage)
	//Store the topic record
	err = stub.PutState(topic, []byte(recBytes))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *sc) checkTopic(stub shim.ChaincodeStubInterface, topic string, status string) error {
	//check the status of current topic
	topicDetails, err := stub.GetState(topic)
	if err != nil {
		return errors.New("Failed to get state")
	} else if topicDetails == nil {
		return errors.New("The topic does not exist:" + topic)
	}
	var topicMap map[string]interface{}
	err = json.Unmarshal([]byte(topicDetails), &topicMap)
	if err != nil {
		return errors.New("failed to unmarshal topic")
	}
	actualStatus := getSafeString(topicMap["stage"])
	if status != actualStatus {
		return errors.New("The actual status " + actualStatus + " is different from expected status " + status)
	}
	return nil
}

//{topic:string ,uid:string,x:big.Int,y:big.Int}
//only available in prepare stage
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
	if len(params) != 4 {
		return shim.Error("Incorrect number of fields in params. Expecting 4.")
	}

	topic := getSafeString(params["topic"])
	err = t.checkTopic(stub, topic, TOPIC.stage["prepare"])
	if err != nil {
		return shim.Error(err.Error())
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
	err = stub.PutState(topic+"_"+PUBLICKEY+":"+uid, []byte(recBytes))
	if err != nil {
		return shim.Error(err.Error())
	}
	//Store the key ring key id
	keyRingDetails, err := stub.GetState(topic + "_" + KEYRING)
	var keyRingBytes []byte
	var keyRing []map[string]interface{}
	if err != nil {
		return shim.Error("Failed to get uid: " + uid)
	} else if keyRingDetails == nil {
		keyRing = make([]map[string]interface{}, 0)
	} else {
		err = json.Unmarshal([]byte(keyRingDetails), &keyRing)
		if err != nil {
			return shim.Error("failed to unmarshal keyRing")
		}
	}
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	keyRing = append(keyRing, uidMap)
	keyRingBytes, _ = json.Marshal(keyRing)
	err = stub.PutState(topic+"_"+KEYRING, []byte(keyRingBytes))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("SUCCESS"))
}

//get details of a particular public key by uid
//parameters: {topic:string, uid:string}
func (t *sc) getPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//checking the number of argument
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}
	var err error
	var params map[string]interface{}
	recBytes := args[0]
	err = json.Unmarshal([]byte(recBytes), &params)
	if err != nil {
		return shim.Error("Failed to unmarshal ")
	}
	if len(params) != 2 {
		return shim.Error("Incorrect number of fields. Expecting 2.")
	}
	topic := getSafeString(params["topic"])
	uid := getSafeString(params["uid"])

	keyDetails, err := stub.GetState(topic + "_" + PUBLICKEY + ":" + uid)
	if err != nil {
		return shim.Error("Failed to get uid: " + uid)
	} else if keyDetails == nil {
		return shim.Error("This uid does not exist: " + uid)
	}
	return shim.Success(keyDetails)
}

//{topic:string,uid:string}
func (t *sc) getKeyRing(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var err error
	var params map[string]interface{}
	recBytes := args[0]
	err = json.Unmarshal([]byte(recBytes), &params)
	if err != nil {
		return shim.Error("failed to unmarshal")
	}
	if len(params) != 1 {
		return shim.Error("Incorrect number of fields. Expecting 2.")
	}
	topic := getSafeString(params["topic"])

	keyRingDetails, err := stub.GetState(topic + "_" + KEYRING)
	if err != nil {
		return shim.Error("Failed to get keyring for topic:" + topic)
	} else if keyRingDetails == nil {
		return shim.Error("This topic does not exist: " + topic)
	}
	return shim.Success(keyRingDetails)
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
	} else if function == "getPublicKey" {
		return t.getPublicKey(stub, args)
	} else if function == "getKeyRing" {
		return t.getKeyRing(stub, args)
	} else if function == "configureTopic" {
		return t.configureTopic(stub, args)
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
