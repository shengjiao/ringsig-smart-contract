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
	checkInvoke(t, stub, [][]byte{[]byte("setStage"), []byte(topicStr)})
	//uid001's key
	keyStr := "{\"topic\":\"USpresident\",\"uid\":\"u01\",\"x\":\"6398110988917490402227161619364033929651640452825541292775484387736563454761\",\"y\":\"44075383121363997802359041633178364043324461069553219394880481874330231188370\"}"
	checkInvoke(t, stub, [][]byte{[]byte("initPublicKey"), []byte(keyStr)})
	//getPublicKey
	var valueMap map[string]interface{}
	json.Unmarshal([]byte(keyStr), &valueMap)
	checkQuery(t, stub, valueMap, [][]byte{[]byte("getPublicKey"), []byte("{\"topic\":\"USpresident\",\"uid\":\"u01\"}")})

	//uid002's key
	keyStr = "{\"topic\":\"USpresident\",\"uid\":\"u02\",\"x\":\"67498550438000454545972363987014904694407640704607446338179461020493945081534\",\"y\":\"67179481346289296301599017224010916382463550533037488721315774963103041194559\"}"
	checkInvoke(t, stub, [][]byte{[]byte("initPublicKey"), []byte(keyStr)})

	//get keyRing
	keyRing := "[{\"uid\":\"u01\"},{\"uid\":\"u02\"}]"
	var valueMapArray []map[string]interface{}
	json.Unmarshal([]byte(keyRing), &valueMapArray)
	checkQueryArray(t, stub, valueMapArray, [][]byte{[]byte("getKeyRing"), []byte("{\"topic\":\"USpresident\"}")})

	//configure topic
	topicStr = "{\"topic\":\"USpresident\",\"stage\":\"start\"}"
	checkInvoke(t, stub, [][]byte{[]byte("setStage"), []byte(topicStr)})

	//submit transaction
	transaction := "{\"topic\":\"USpresident\",\"uid\":\"u01\",\"msg\":\"Hello, world.\",\"sig\":{\"hsx\":\"97769613060326423184466863334433925251257344864036598209839835582457743858796\",\"hsy\":\"35573932117220247294212087135482053550981809836517735404307116170034931105481\",\"c\":[\"57509808921387526446863189170044818150780961915097858159263146585833014118016\",\"2840124079242068908122670502835382084740910938587574934331953180838035099633\"],\"t\":[\"40711078032794012211083313132698384679743458610309089999267032328825394655227\",\"51649693558643593918242089311909556086167419639407443247363869832628152428431\"]},\"keyIndex\":[{\"uid\":\"u01\"},{\"uid\":\"u02\"}]}"
	checkInvoke(t, stub, [][]byte{[]byte("submit"), []byte(transaction)})
}
