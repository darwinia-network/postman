package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/websocket"
)

// Environment variables
var (
	INTERVAL_SECONDS      = 10
	ALERTMANAGER_ENDPOINT = "http://127.0.0.1:9093/api/v2"
	SHADOW_ENDPOINT       = "https://testnet.shadow.darwinia.network"
	NODE_WS_ENDPOINT      = "wss://crab.darwinia.network"
)

// Init module
func init() {
	ale := os.Getenv("ALERTMANAGER_ENDPOINT")
	if ale != "" {
		ALERTMANAGER_ENDPOINT = ale
	}

	nwe := os.Getenv("NODE_WS_ENDPOINT")
	if nwe != "" {
		NODE_WS_ENDPOINT = nwe
	}

	she := os.Getenv("SHADOW_ENDPOINT")
	if she != "" {
		SHADOW_ENDPOINT = she
	}

	ins := os.Getenv("INTERVAL_SECONDS")
	if ins != "" {
		if secs, err := strconv.Atoi(ins); err == nil {
			INTERVAL_SECONDS = secs
		} else {
			panic(err)
		}
	}

	websocket.SetEndpoint(NODE_WS_ENDPOINT)
	register()
}

func main() {
	for {
		ride()
		time.Sleep(time.Duration(INTERVAL_SECONDS) * time.Second)
	}
}

func ride() {
	v := rpc.JsonRpcResult{}
	ce(websocket.SendWsRequest(nil, &v, rpc.StateGetStorage(0, ENCODE_KEY, "")))

	// Get pending headers
	dataHex, err := v.ToString()
	ce(err)

	// Decode headerthing codec
	r, err := storage.Decode(dataHex, SCALE_TYPE, nil)
	ce(err)

	// check if empty headers
	if r == "null" || r == "" {
		log.Println("No pending headers...")
		return
	}

	// Decode headerthing json
	var headers []PendingHeader
	ce(json.Unmarshal([]byte(r.ToString()), &headers))

	// Check headers
	for _, item := range headers {
		if eq, ht := checkHeader(item); !eq {
			alert := GenAlert(item, ht)
			alert.emit()
			continue
		}

		log.Printf("Ethereum block %d looks nice ~ \n", item.EthereumBlockNumber)
	}
}

func checkHeader(header PendingHeader) (bool, HeaderThing) {
	resp, err := http.Get(fmt.Sprintf("%s/eth/header/%d", SHADOW_ENDPOINT, header.EthereumBlockNumber))
	if err != nil {
		log.Println(err)
		time.Sleep(3 * time.Second)
		return checkHeader(header)
	}

	defer resp.Body.Close()
	var canonicalHT ComplexHeaderThing
	err = json.NewDecoder(resp.Body).Decode(&canonicalHT)
	if err != nil {
		log.Println(err)
		time.Sleep(3 * time.Second)
		return checkHeader(header)
	}

	// deep compare
	return reflect.DeepEqual(canonicalHT.HeaderThing, header.HeaderThing.HeaderThing()), canonicalHT.HeaderThing
}

/**
 * Type registry
 */
type PendingHeader struct {
	BlockNumber         uint64           `json:"col1"`
	EthereumBlockNumber uint64           `json:"col2"`
	HeaderThing         ScaleHeaderThing `json:"col3"`
}

func (p *PendingHeader) toString() string {
	b, err := json.Marshal(p)
	ce(err)
	return string(b)
}

type ScaleHeaderThing struct {
	Header  ScaleHeader `json:"header"`
	MmrRoot string      `json:"mmr_root"`
}

func (h *ScaleHeaderThing) HeaderThing() (header HeaderThing) {
	header.MmrRoot = h.MmrRoot
	header.Header = h.Header.Header()
	return
}

type ScaleHeader struct {
	ParentHash       string   `json:"parent_hash"`
	TimeStamp        uint64   `json:"timestamp"`
	Number           uint64   `json:"number"`
	Author           string   `json:"author"`
	TransactionsRoot string   `json:"transactions_root"`
	UnclesHash       string   `json:"uncles_hash"`
	ExtraData        string   `json:"extra_data"`
	StateRoot        string   `json:"state_root"`
	ReceiptsRoot     string   `json:"receipts_root"`
	LogBloom         string   `json:"log_bloom"`
	GasUsed          []uint64 `json:"gas_used"`
	GasLimited       []uint64 `json:"gas_limit"`
	Difficulty       []uint64 `json:"difficulty"`
	Seal             []string `json:"seal"`
	Hash             string   `json:"hash"`
}

func (h *ScaleHeader) Header() (header Header) {
	if !strings.HasPrefix(h.Seal[0], "0x") {
		h.Seal[0] = "0x" + h.Seal[0]
	}

	if !strings.HasPrefix(h.Seal[1], "0x") {
		h.Seal[1] = "0x" + h.Seal[1]
	}

	header.ParentHash = h.ParentHash
	header.TimeStamp = h.TimeStamp
	header.Number = h.Number
	header.Author = strings.ToLower(h.Author)
	header.TransactionsRoot = h.TransactionsRoot
	header.UnclesHash = h.UnclesHash
	header.ExtraData = "0x" + h.ExtraData
	header.StateRoot = h.StateRoot
	header.ReceiptsRoot = h.ReceiptsRoot
	header.LogBloom = h.LogBloom
	header.GasUsed = h.GasUsed[0]
	header.GasLimited = h.GasLimited[0]
	header.Difficulty = h.Difficulty[0]
	header.Seal = h.Seal
	header.Hash = h.Hash
	return
}

type ComplexHeaderThing struct {
	HeaderThing   HeaderThing `json:"header_thing"`
	Confirmations uint64      `json:"confirmations"`
}

type HeaderThing struct {
	Header  Header `json:"header"`
	MmrRoot string `json:"mmr_root"`
}

func (h *HeaderThing) toString() string {
	b, err := json.Marshal(h)
	ce(err)
	return string(b)
}

type Header struct {
	ParentHash       string   `json:"parent_hash"`
	TimeStamp        uint64   `json:"timestamp"`
	Number           uint64   `json:"number"`
	Author           string   `json:"author"`
	TransactionsRoot string   `json:"transactions_root"`
	UnclesHash       string   `json:"uncles_hash"`
	ExtraData        string   `json:"extra_data"`
	StateRoot        string   `json:"state_root"`
	ReceiptsRoot     string   `json:"receipts_root"`
	LogBloom         string   `json:"log_bloom"`
	GasUsed          uint64   `json:"gas_used"`
	GasLimited       uint64   `json:"gas_limit"`
	Difficulty       uint64   `json:"difficulty"`
	Seal             []string `json:"seal"`
	Hash             string   `json:"hash"`
}

func register() {
	if coded, err := rpc.GetMetadataByHash(nil); err == nil {
		metadata.Latest(&metadata.RuntimeRaw{Spec: 1, Raw: util.TrimHex(coded)})
		substrate.RegCustomTypes([]byte(Registry))
		return
	}
	register()
}

/**
 * Prometheus Alert
 */
type Alert struct {
	StartsAt   string          `json:"startsAt,omitempty"`
	EndsAt     string          `json:"endsAt,omitempty"`
	Annotation AlertAnnotation `json:"annotations"`
	Label      AlertLabel      `json:"labels"`
	URL        string          `json:"generatorUrl"`
}

type AlertAnnotation struct {
	Pending   string `json:"pending"`
	Canonical string `json:"canonical"`
}

type AlertLabel struct {
	AlertName      string `json:"alertname"`
	ShadowEndpoint string `json:"shadow"`
	NodeWsEndpoint string `json:"node_ws"`
}

func GenAlert(pending PendingHeader, canonical HeaderThing) (alert Alert) {
	alert.StartsAt = time.Now().Format(time.RFC3339)
	alert.Label = AlertLabel{
		AlertName:      "PendingHeadersInconsistent",
		ShadowEndpoint: SHADOW_ENDPOINT,
		NodeWsEndpoint: NODE_WS_ENDPOINT,
	}
	alert.URL = "https://github.com/darwinia-network/postman"
	alert.Annotation = AlertAnnotation{
		Pending:   pending.toString(),
		Canonical: canonical.toString(),
	}

	return
}

func (a *Alert) emit() {
	j, err := json.Marshal([]Alert{*a})
	ce(err)

	// Post alert to Alertmanager
	_, err = http.Post(
		fmt.Sprintf(ALERTMANAGER_ENDPOINT+"/alerts"),
		"application/json",
		bytes.NewBuffer(j),
	)
	ce(err)
}

/**
 * Util
 */
func ce(err error) {
	if err != nil {
		log.Println(err)
	}
}

/**
 * Static config
 */
const (
	ENCODE_KEY = "0xcdacb51c37fcd27f3b87230d9a1c26509f7d076895629ddec219b5e71b9bc2ad"
	SCALE_TYPE = "Vec<(BlockNumber, u64, EthereumHeaderThing)>"
)

const Registry = `{
  "U256": "[u64; 4]",
  "EthAddress": "H160",
  "Bloom": "[u8; 256]",
  "EthereumHeaderThing": {
      "type": "struct",
      "type_mapping": [
        [
          "header",
          "EthereumHeader"
        ],
        [
          "mmr_root",
          "H256"
        ]
      ]
    },
  "EthereumHeader": {
      "type": "struct",
      "type_mapping": [
        [
          "parent_hash",
          "H256"
        ],
        [
          "timestamp",
          "u64"
        ],
        [
          "number",
          "u64"
        ],
        [
          "author",
          "EthereumAddress"
        ],
        [
          "transactions_root",
          "H256"
        ],
        [
          "uncles_hash",
          "H256"
        ],
        [
          "extra_data",
          "Bytes"
        ],
        [
          "state_root",
          "H256"
        ],
        [
          "receipts_root",
          "H256"
        ],
        [
          "log_bloom",
          "Bloom"
        ],
        [
          "gas_used",
          "U256"
        ],
        [
          "gas_limit",
          "U256"
        ],
        [
          "difficulty",
          "U256"
        ],
        [
          "seal",
          "Vec<Bytes>"
        ],
        [
          "hash",
          "Option<H256>"
        ]
      ]
    }
}`
