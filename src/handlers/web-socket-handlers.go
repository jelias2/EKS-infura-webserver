package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jelias2/infra-test/src/apis"
	"go.uber.org/zap"
)

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
func (h *Handler) WebSocketGetBlockByNumber(w http.ResponseWriter, r *http.Request) {

	getBlockBody, _ := json.Marshal(createRequestBody(apis.GetBlockByNumber, []string{}))
	err := h.WebSocket.WriteMessage(websocket.TextMessage, getBlockBody)
	if err != nil {
		h.Log.Info("Error writing WebSocketGetBlockByNumber websocket message", zap.Error(err))
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: 400,
			Message:    "Error writing websocket message",
		})
		return
	}
	_, message, err := h.WebSocket.ReadMessage()
	if err != nil {
		h.Log.Info("Error reading WebSocketGetBlockByNumber websocket message", zap.Error(err))
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: 400,
			Message:    "Error writing websocket message",
		})
		return
	}
	wsGetBlockNumberResponse := &apis.GetBlockNumberResponse{}
	json.Unmarshal(message, wsGetBlockNumberResponse)
	json.NewEncoder(w).Encode(wsGetBlockNumberResponse)
}
