package handlers

import (
	"encoding/json"
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