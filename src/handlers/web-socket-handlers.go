package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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
	var err error

	w.Header().Set("Content-Type", "application/json")
	h.Log.Info("Entered WebsocketGetBlockByNumber")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var getBlockByNumberRequest apis.GetBlockByNumberRequest
	if err := json.Unmarshal(reqBody, &getBlockByNumberRequest); err != nil {
		h.Log.Error("Error unmarshalling GetBlockByNumberRequest", zap.Error(err))
	}

	txdetails, err = strconv.ParseBool(getBlockByNumberRequest.TxDetails)
	if getBlockByNumberRequest.Block == "" || err != nil {
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: 400,
			Message:    "Malformed Request",
		})
		return
	}

	// Can't use create RequestBody because 2nd param is bool with no quotes
	body := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s",%s],"id":1}`, getBlockByNumberRequest.Block, getBlockByNumberRequest.TxDetails))
	h.Log.Info("GetBlockByNumber body", zap.String("Body", string(body)))

	//TODO: Maybe use a function pointer in here to make it simpler, pass in an unmarshal function which is hardcode to type?
	if txdetails {
		json.NewEncoder(w).Encode(h.WebSocketTxDetailsResponse(body))
	} else {
		json.NewEncoder(w).Encode(h.WebSocketNoTxDetailsResponse(body))
	}

}

func (h *Handler) WebSocketTxDetailsResponse(body []byte) interface{} {
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
	wsResult := &apis.GetBlockByNumberTxDetailsResponse{}
	//TODO: check for marshal error here
	json.Unmarshal(message, wsResult)
	return wsResult
}

func (h *Handler) WebSocketNoTxDetailsResponse(body []byte) interface{} {
	err := h.WebSocket.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		h.Log.Info("Error writing WebSocketNoTxDetailsResponse websocket message", zap.Error(err))
		return &apis.ErrorResponse{StatusCode: 400, Message: err.Error()}
	}
	_, message, err := h.WebSocket.ReadMessage()
	if err != nil {
		h.Log.Info("Error writing WebSocketNoTxDetailsResponse websocket message", zap.Error(err))
		return &apis.ErrorResponse{StatusCode: 400, Message: err.Error()}
	}
	wsResult := &apis.GetBlockByNumberNoTxDetailsResponse{}
	//TODO: check for marshal error here
	json.Unmarshal(message, wsResult)
	return wsResult
}
