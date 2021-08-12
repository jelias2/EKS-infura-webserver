package apis

// Default Request Fields
const RPCVersion2 = "2.0"
const RequestID = 1

// RPC Methods
type RPCCall string

const GetBlockNumber RPCCall = "eth_blockNumber"
const GetGasPrice RPCCall = "eth_gasPrice"
const GetBlockByNumber RPCCall = "eth_getBlockByNumber"
const GetTransactionByBlockNumberAndIndex RPCCall = "eth_getTransactionByBlockNumberAndIndex"

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

type InfuraRequestBody struct {
	JsonRPC string   `json:"jsonrpc"`
	Method  RPCCall  `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}
