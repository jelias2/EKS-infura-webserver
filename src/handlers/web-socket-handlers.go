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
	err := h.WebSocket.WriteMessage(websocket.TextMessage, getBlockBody)
	if err != nil {
		h.Log.Info("Error writing WebSocketGetBlockNumber websocket message", zap.Error(err))
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: 400,
			Message:    "Error writing websocket message",
		})
		return
	}
	_, message, err := h.WebSocket.ReadMessage()
	if err != nil {
		h.Log.Info("Error reading WebSocketGetBlockNumber websocket message", zap.Error(err))
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: 400,
			Message:    "Error reading websocket message",
		})
		return
	}
	wsGetBlockNumberResponse := &apis.GetBlockNumberResponse{}
	json.Unmarshal(message, wsGetBlockNumberResponse)
	json.NewEncoder(w).Encode(wsGetBlockNumberResponse)
}

// WebSocketGetGasPrice
func (h *Handler) WebSocketGetGasPrice(w http.ResponseWriter, r *http.Request) {

	getBlockBody, _ := json.Marshal(createRequestBody(apis.GetGasPrice, []string{}))
	err := h.WebSocket.WriteMessage(websocket.TextMessage, getBlockBody)
	if err != nil {
		h.Log.Info("Error writing WebSocketGetGasPrice websocket message", zap.Error(err))
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: 400,
			Message:    "Error writing websocket message",
		})
		return
	}
	_, message, err := h.WebSocket.ReadMessage()
	if err != nil {
		h.Log.Info("Error reading WebSocketGetGasPrice websocket message", zap.Error(err))
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: 400,
			Message:    "Error writing websocket message",
		})
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
		return
	}
	json.NewEncoder(w).Encode(h.WebSocketGetBlockByNumberHandler(formmattedRequest, apis.GetBlockByNumberNoTxDetailsResponse{}))

}

func (h *Handler) WebSocketGetBlockByNumberHandler(body []byte, umarshallStruct interface{}) interface{} {
	var err error
	err = h.WebSocket.WriteMessage(websocket.TextMessage, []byte(body))
	if err != nil {
		h.Log.Info("Error writing WebSocketGetGasPrice websocket message", zap.Error(err))
		return &apis.ErrorResponse{StatusCode: 400, Message: err.Error()}
	}

	_, message, err := h.WebSocket.ReadMessage()
	if err != nil {
		h.Log.Info("Error reading WebSocketGetGasPrice websocket message", zap.Error(err))
		return apis.ErrorResponse{StatusCode: 400, Message: err.Error()}
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
		return &apis.ErrorResponse{StatusCode: 400, Message: err.Error()}
	}
}
