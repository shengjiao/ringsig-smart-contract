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
	keyStr := "{\"topic\":\"USpresident\",\"uid\":\"u01\",\"x\":\"76246489839868595250736483335544835586964803968400689679787038113601921231908\",\"y\":\"56041144687754187004494289427517687774152975165664140788239714683881245853753\"}"
	checkInvoke(t, stub, [][]byte{[]byte("initPublicKey"), []byte(keyStr)})
	//getPublicKey
	var valueMap map[string]interface{}
	json.Unmarshal([]byte(keyStr), &valueMap)
	checkQuery(t, stub, valueMap, [][]byte{[]byte("getPublicKey"), []byte("{\"topic\":\"USpresident\",\"uid\":\"u01\"}")})

	//uid002's key
	keyStr = "{\"topic\":\"USpresident\",\"uid\":\"u02\",\"x\":\"29118928702052831236093327210741276250501097675031360223624055807789787515116\",\"y\":\"38525044966898496489930698969687718996165907795558657949585270845383985923206\"}"
	checkInvoke(t, stub, [][]byte{[]byte("initPublicKey"), []byte(keyStr)})

	//get keyRing
	keyRing := "[{\"uid\":\"u01\"},{\"uid\":\"u02\"}]"
	var valueMapArray []map[string]interface{}
	json.Unmarshal([]byte(keyRing), &valueMapArray)
	checkQueryArray(t, stub, valueMapArray, [][]byte{[]byte("getKeyRing"), []byte("{\"topic\":\"USpresident\"}")})

	//configure topic
	topicStr = "{\"topic\":\"USpresident\",\"stage\":\"start\"}"
	checkInvoke(t, stub, [][]byte{[]byte("setStage"), []byte(topicStr)})

	//submit transaction, signed by u01 using ring key
	transaction := "{\"topic\":\"USpresident\",\"msg\":\"Trump\",\"sig\":{\"hsx\":\"30703293276322474432077759229812626311153619506355410404284365095286437266201\",\"hsy\":\"109783253905845432603851663791351392322055706254190241088647146192362309998077\",\"c\":[\"104027820470093965906600879384753086420996910311983518503674757699941113744916\",\"105231547088978014403579853369679548230648615247725480470623886656193328104399\"],\"t\":[\"20174120779801100673831565358206618705638099255145070882496441119328888406400\",\"48829932724179419142539218998261985156434163456715641621909206899798522029727\"]},\"keyIndex\":[{\"uid\":\"u01\"},{\"uid\":\"u02\"}]}"
	checkInvoke(t, stub, [][]byte{[]byte("submit"), []byte(transaction)})
	//submit transaction, signed by u01 using ring key
	transaction = "{\"topic\":\"USpresident\",\"msg\":\"Trump\",\"sig\":{\"hsx\":\"30703293276322474432077759229812626311153619506355410404284365095286437266201\",\"hsy\":\"109783253905845432603851663791351392322055706254190241088647146192362309998077\",\"c\":[\"59112343535051908299949916799690771981324662596786637431858167202968713739265\",\"110922536158145308210393977035514666520999183389619518492691865922278137126937\"],\"t\":[\"8550274795537384873677367681416591161155493586211715822160186529518315177324\",\"45654221062156635923665323117389128732870414043385491830657389056106751127396\"]},\"keyIndex\":[{\"uid\":\"u01\"},{\"uid\":\"u02\"}]}"
	checkInvoke(t, stub, [][]byte{[]byte("submit"), []byte(transaction)})
	//submit transaction, signed by u02 using ring key
	transaction = "{\"topic\":\"USpresident\",\"msg\":\"Trump\",\"sig\":{\"hsx\":\"58332199994191092592348420150349264464568582559535380724749575699277739971963\",\"hsy\":\"86603447641809474551814830879628653895345272791274427601439274239349981735045\",\"c\":[\"5954247435713859385819772819793847208014044760965454371014604395457780505506\",\"42708764911580717095513513061874846510523026428464105575044450392456085875650\"],\"t\":[\"46755348651082555636719026939679789609984908257232004212728608438107147717531\",\"101704588159729850901643832343634371764886146333455311761307244157188789091395\"]},\"keyIndex\":[{\"uid\":\"u01\"},{\"uid\":\"u02\"}]}"
	checkInvoke(t, stub, [][]byte{[]byte("submit"), []byte(transaction)})
}
