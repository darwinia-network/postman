package main

import (
	"fmt"
	"io/ioutil"

	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	// "github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/websocket"
)

const (
	ENCODE_KEY = "0xcdacb51c37fcd27f3b87230d9a1c26509f7d076895629ddec219b5e71b9bc2ad"
	SCALE_TYPE = "Vec<(BlockNumber, u64, TcHeaderThing)>"
)

// Register metadata
func register() {
	if coded, err := rpc.GetMetadataByHash(nil); err == nil {
		metadata.Latest(&metadata.RuntimeRaw{Spec: 1, Raw: util.TrimHex(coded)})
		return
	}
	register()
}

func init() {
	websocket.SetEndpoint("ws://localhost:9944")
	register()
	config, err := ioutil.ReadFile("../crab.json")
	if err != nil {
		panic("Could not find crab.json in ../crab.json")
	}

	substrate.RegCustomTypes(config)
}

func main() {
	v := rpc.JsonRpcResult{}
	err := websocket.SendWsRequest(nil, &v, rpc.StateGetStorage(0, ENCODE_KEY, ""))
	if err != nil {
		fmt.Println(err)
	}

	dataHex, err := v.ToString()
	if err != nil {
		fmt.Println(err)
	}

	r, err := storage.Decode(dataHex, SCALE_TYPE, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r)
}
