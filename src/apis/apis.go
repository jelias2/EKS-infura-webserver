package apis

import "net/http"

// Default Request Fields
const RPCVersion2 = "2.0"
const RequestID = 1

// RPC Methods
type RPCCall string

const GetBlockNumber RPCCall = "eth_blockNumber"
const GetGasPrice RPCCall = "eth_gasPrice"
const GetBlockByNumber RPCCall = "eth_getBlockByNumber"
const GetLogs RPCCall = "eth_getLogs"
const GetStorageAt RPCCall = "eth_getStorageAt"

const GetTransactionByBlockNumberAndIndex RPCCall = "eth_getTransactionByBlockNumberAndIndex"

//const GetBlockByNumberRequestBodyTemplate string = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s",%s],"id":1}`
const BooleanRequestBodyTemplate string = `{"jsonrpc":"2.0","method":"%s","params":["%s",%s],"id":1}`

const MalformedRequestMessage = "Malformed Request"

type Healthcheck struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Datetime string `json:"datetime"`
}

type GetBlockNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

type GetGasPriceResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statuscode"`
	Message    string `json:"message"`
}

var MalformedRequestError = ErrorResponse{
	StatusCode: http.StatusBadRequest,
	Message:    MalformedRequestMessage,
}

type InfuraRequestBody struct {
	JsonRPC string   `json:"jsonrpc"`
	Method  RPCCall  `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}
