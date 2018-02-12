package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, value map[string]interface{}, args [][]byte) {
	res := stub.MockInvoke("1", args)

	if res.Status != shim.OK {
		fmt.Println("Query failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query failed to get value")
		t.FailNow()
	}

	str := string(res.Payload)
	var valueMap map[string]interface{}
	json.Unmarshal([]byte(str), &valueMap)

	fmt.Println("After query: ", string(res.Payload))
	fmt.Println("=======")
	valueString, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Parameter:", string(valueString))

	eq := reflect.DeepEqual(value, valueMap)
	if eq {
		fmt.Println("They're equal.")
	} else {
		fmt.Println("Query failed")
		t.FailNow()
	}
	fmt.Println("---------------------")
}

func checkQueryArray(t *testing.T, stub *shim.MockStub, value []map[string]interface{}, args [][]byte) {
	res := stub.MockInvoke("1", args)

	if res.Status != shim.OK {
		fmt.Println("Query failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query failed to get value")
		t.FailNow()
	}

	str := string(res.Payload)
	var valueMap []map[string]interface{}
	json.Unmarshal([]byte(str), &valueMap)

	fmt.Println("After query: ", string(res.Payload))
	fmt.Println("=======")
	valueString, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Parameter:", string(valueString))

	for i := range value {
		eq := reflect.DeepEqual(value[i], valueMap[i])
		if !eq {
			fmt.Println("Query failed")
			t.FailNow()
		}
	}

	fmt.Println("They're equal.")
	fmt.Println("---------------------")
}

func checkQueryString(t *testing.T, stub *shim.MockStub, value string, args [][]byte) {
	res := stub.MockInvoke("1", args)

	if res.Status != shim.OK {
		fmt.Println("Query failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query failed to get value")
		t.FailNow()
	}

	str := string(res.Payload)

	fmt.Println("After query: ", string(res.Payload))
	fmt.Println("=======")
	fmt.Println("Parameter:", value)

	if str == value {
		fmt.Println("They're equal.")
	} else {
		fmt.Println("Query failed")
		t.FailNow()
	}
	fmt.Println("---------------------")
}

func TestEndToEndWorkflow(t *testing.T) {
	scc := new(sc)
	stub := shim.NewMockStub("sc", scc)

	// Init without any argument
	checkInit(t, stub, [][]byte{[]byte("init")})

	//configure topic
	topicStr := "{\"topic\":\"USpresident\",\"stage\":\"prepare\"}"
	checkInvoke(t, stub, [][]byte{[]byte("configureTopic"), []byte(topicStr)})
	//uid001's key
	keyStr := "{\"topic\":\"USpresident\",\"uid\":\"u01\",\"x\":\"79325335377659446719061365985594928216557351703018226449107942108421649247394\",\"y\":\"32519086782170187642508289520183198737894713738318360286903070346702754579434\"}"
	checkInvoke(t, stub, [][]byte{[]byte("initPublicKey"), []byte(keyStr)})
	//getPublicKey
	var valueMap map[string]interface{}
	json.Unmarshal([]byte(keyStr), &valueMap)
	checkQuery(t, stub, valueMap, [][]byte{[]byte("getPublicKey"), []byte("{\"topic\":\"USpresident\",\"uid\":\"u01\"}")})

	//uid002's key
	keyStr = "{\"topic\":\"USpresident\",\"uid\":\"u02\",\"x\":\"110855819708902045581138277350236606000241599806345538157278800268663949421927\",\"y\":\"66713262917835133565738665462638578244558791813636461880138560176459282784472\"}"
	checkInvoke(t, stub, [][]byte{[]byte("initPublicKey"), []byte(keyStr)})

	keyRing := "[{\"uid\":\"u01\"},{\"uid\":\"u02\"}]"
	var valueMapArray []map[string]interface{}
	json.Unmarshal([]byte(keyRing), &valueMapArray)
	checkQueryArray(t, stub, valueMapArray, [][]byte{[]byte("getKeyRing"), []byte("{\"topic\":\"USpresident\"}")})
}
