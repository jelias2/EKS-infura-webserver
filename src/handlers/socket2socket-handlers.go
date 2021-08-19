package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jelias2/infra-test/src/apis"
	"go.uber.org/zap"
)

func (h *Handler) Socket2socket(w http.ResponseWriter, r *http.Request) {

	var infuraReq []byte
	var ok bool
	infuraClient, clientConn := h.UpgradeConnection(w, r)

	for {
		var err error

		_, msg, err := clientConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.Log.Error("Client closed unexpecttedly")
			} else {
				h.Log.Error("Failed to read message from websocket client", zap.Error(err))
			}
			break
		}
		h.Log.Info("Recieved Websocket Message", zap.String("Message", string(msg)))
		if bytes.Contains(msg, []byte("true")) || bytes.Contains(msg, []byte("false")) {
			if infuraReq, ok = h.formatInfuraBooeanRequestMsg(msg); !ok {
				clientConn.WriteMessage(websocket.TextMessage, infuraReq)
				continue
			}
		} else {
			if infuraReq, ok = h.formatInfuraRequestMsg(msg); !ok {
				clientConn.WriteMessage(websocket.TextMessage, []byte("Failed to unmarshalll message client"))
				continue
			}
		}

		if ok {
			clientResp := h.WriteAndReadToInfura(infuraClient, infuraReq)
			h.Log.Info("Writing Client websocket message", zap.ByteString("Response", clientResp))
			if err = clientConn.WriteMessage(websocket.TextMessage, clientResp); err != nil {
				h.Log.Info("Error wrting client message", zap.Error(err))
				break
			}
		}
	}

	infuraClient.Close()
	clientConn.Close()
}

func (h *Handler) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, *websocket.Conn) {
	upgrader := websocket.Upgrader{}
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		json.NewEncoder(w).Encode(&apis.ErrorResponse{
			StatusCode: http.StatusUpgradeRequired,
			Message:    "Error Upgrading to WebSocket connection",
		})
	}

	mainnetWebsocketEndpoint := h.Mainnet_websocket_endpoint
	infuraClient, _, err := websocket.DefaultDialer.Dial(mainnetWebsocketEndpoint, nil)
	if err != nil {
		log.Fatal("Fatal Dial Error:", zap.Error(err))
		clientConn.WriteMessage(websocket.CloseMessage, []byte("Failed to unmarshalll message client"))
		clientConn.Close()
	}
	return infuraClient, clientConn
}

func (h *Handler) WriteAndReadToInfura(infuraClient *websocket.Conn, infuraReq []byte) []byte {
	var err error
	var infuraResp []byte
	if err = infuraClient.WriteMessage(websocket.TextMessage, infuraReq); err != nil {
		h.Log.Info("Error writing Infura websocket message", zap.Error(err))
		return []byte("Error writing Infura websocket message")
	}

	if _, infuraResp, err = infuraClient.ReadMessage(); err != nil {
		h.Log.Info("Error reading Infura websocket message", zap.Error(err))
		return []byte("Error reading Infura websocket message")
	}
	return infuraResp

}

// eth_getBlockByNumber, eth_getBlockByHash contain boolean, have len 2 params array
func (h *Handler) formatInfuraBooeanRequestMsg(body []byte) ([]byte, bool) {
	body = bytes.Replace(body, []byte("true"), []byte("\"true\""), 1)
	body = bytes.Replace(body, []byte("false"), []byte("\"false\""), 1)
	reqBody := &apis.InfuraRequestBody{}
	if err := json.Unmarshal(body, reqBody); err != nil {
		h.Log.Info("Error Parsing boolean websocket message", zap.Error(err))
		return []byte("Error Parsing boolean websocket message"), false
	}
	body = []byte(fmt.Sprintf(apis.BooleanRequestBodyTemplate, reqBody.Method, reqBody.Params[0], reqBody.Params[1]))
	return body, true
}

func (h *Handler) formatInfuraRequestMsg(clientReq []byte) ([]byte, bool) {
	var reqBody = &apis.InfuraRequestBody{}
	var msg []byte
	var err error
	if err = json.Unmarshal(clientReq, reqBody); err != nil {
		return []byte("Failed to unmarshall request message for Infura"), false
	}

	if msg, err = json.Marshal(h.CreateRequestBody(reqBody.Method, reqBody.Params)); err != nil {
		return []byte("Failed to create formmated request message for Infura"), false
	}
	return msg, true
}
