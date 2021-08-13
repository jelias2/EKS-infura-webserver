package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jelias2/infra-test/src/apis"
	"go.uber.org/zap"
)

// WebSocketGetGasPrice
func (h *Handler) WebSocketGetBlockNumber(w http.ResponseWriter, r *http.Request) {

	getBlockBody, _ := json.Marshal(createRequestBody(apis.GetBlockNumber, []string{}))
	message, ErrorResponse := h.WebSocketWriteAndRead(getBlockBody)
	if ErrorResponse.Message != "" && ErrorResponse.StatusCode != 0 {
		json.NewEncoder(w).Encode(ErrorResponse)
		return
	}
	wsGetBlockNumberResponse := &apis.GetBlockNumberResponse{}
	json.Unmarshal(message, wsGetBlockNumberResponse)
	json.NewEncoder(w).Encode(wsGetBlockNumberResponse)
}

// WebSocketGetGasPrice
func (h *Handler) WebSocketGetGasPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	getBlockBody, _ := json.Marshal(createRequestBody(apis.GetGasPrice, []string{}))
	message, ErrorResponse := h.WebSocketWriteAndRead(getBlockBody)
	if ErrorResponse.Message != "" && ErrorResponse.StatusCode != 0 {
		json.NewEncoder(w).Encode(ErrorResponse)
		return
	}
	wsGetGasResponse := &apis.GetGasPriceResponse{}
	json.Unmarshal(message, wsGetGasResponse)
	json.NewEncoder(w).Encode(wsGetGasResponse)
}

// WebSocketGetGasPrice
// GetBlockByNumber
func (h *Handler) WebSocketGetBlockByNumber(w http.ResponseWriter, r *http.Request) {

	var txdetails bool
	w.Header().Set("Content-Type", "application/json")
	formmattedRequest, validRequest, txdetails := h.ParseGetBlockByNumberRequest(r)
	if !validRequest {
		wsError := &apis.ErrorResponse{}
		json.Unmarshal(formmattedRequest, wsError)
		json.NewEncoder(w).Encode(wsError)
		return
	}

	if txdetails {
		json.NewEncoder(w).Encode(h.WebSocketGetBlockByNumberHandler(formmattedRequest, apis.GetBlockByNumberTxDetailsResponse{}))
	}
	json.NewEncoder(w).Encode(h.WebSocketGetBlockByNumberHandler(formmattedRequest, apis.GetBlockByNumberNoTxDetailsResponse{}))
}

func (h *Handler) WebSocketGetBlockByNumberHandler(body []byte, umarshallStruct interface{}) interface{} {
	var message []byte
	var errorResponse apis.ErrorResponse
	message, errorResponse = h.WebSocketWriteAndRead(body)
	if errorResponse.Message != "" && errorResponse.StatusCode != 0 {
		return errorResponse
	}

	// Umarshall response into the type of the umarshallStruct
	switch umarshallStruct.(type) {
	case apis.GetBlockByNumberTxDetailsResponse:
		wsResult := &apis.GetBlockByNumberTxDetailsResponse{}
		json.Unmarshal(message, wsResult)
		return wsResult
	case apis.GetBlockByNumberNoTxDetailsResponse:
		wsResult := &apis.GetBlockByNumberNoTxDetailsResponse{}
		json.Unmarshal(message, wsResult)
		return wsResult
	default:
		h.Log.Error("Improper Type")
		return &apis.ErrorResponse{StatusCode: http.StatusBadRequest, Message: "Error Unmarshalling GetBlockResponse"}
	}
}

// WebSocketGetTransactionByBlockNumberAndIndex
func (h *Handler) WebSocketGetTransactionByBlockNumberAndIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var getTxReq apis.GetTransactionByBlockNumberAndIndexRequest
	if err := json.Unmarshal(reqBody, &getTxReq); err != nil {
		h.Log.Error("Error unmarshalling GetBlockByNumberRequest", zap.Error(err))
	}

	if getTxReq.Block == "" || getTxReq.Index == "" {
		json.NewEncoder(w).Encode(apis.MalformedRequestError)
		return
	}

	getBlockTxIndex, _ := json.Marshal(createRequestBody(apis.GetTransactionByBlockNumberAndIndex, []string{getTxReq.Block, getTxReq.Index}))
	message, errorResponse := h.WebSocketWriteAndRead(getBlockTxIndex)
	if errorResponse.Message != "" && errorResponse.StatusCode != 0 {
		json.NewEncoder(w).Encode(errorResponse)
	}

	wsGetTxByBlockAndIndexResp := &apis.GetTransactionByBlockNumberAndIndexResponse{}
	json.Unmarshal(message, wsGetTxByBlockAndIndexResp)
	json.NewEncoder(w).Encode(wsGetTxByBlockAndIndexResp)

}

func (h *Handler) WebSocketWriteAndRead(body []byte) ([]byte, apis.ErrorResponse) {
	errorMessage := ""
	statusCode := 0
	err := h.WebSocket.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		h.Log.Info("Error writing WebSocketGetGasPrice websocket message", zap.Error(err))
		errorMessage = err.Error()
		statusCode = http.StatusBadRequest
	}
	_, message, err := h.WebSocket.ReadMessage()
	if err != nil {
		h.Log.Info("Error reading WebSocketGetGasPrice websocket message", zap.Error(err))
		errorMessage = err.Error()
		statusCode = http.StatusBadRequest
	}
	errorResponse := apis.ErrorResponse{
		StatusCode: statusCode,
		Message:    errorMessage,
	}
	return message, errorResponse
}
