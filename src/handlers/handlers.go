package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"time"

	"github.com/gorilla/websocket"
	"github.com/jelias2/infra-test/src/apis"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Handler struct {
	Log                        *zap.Logger
	Resty                      *resty.Client
	WebSocket                  *websocket.Conn
	Mainnet_http_endpoint      string
	Mainnet_websocket_endpoint string
}

// Healthcheck will display test response to make sure the server is running
func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apis.Healthcheck{
		Status:   http.StatusAccepted,
		Message:  "Healthcheck response",
		Datetime: time.Now().String(),
	})
}

// Get ethblock number
func (h *Handler) GetBlockNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	getBlockBody := createRequestBody(apis.GetBlockNumber, []string{})
	result := &apis.GetBlockNumberResponse{}
	resp, err := h.Resty.R().SetBody(getBlockBody).
		SetResult(result).
		Post(h.Mainnet_http_endpoint)
	h.debugResponse("GetBlockNumber", resp, err)
	json.NewEncoder(w).Encode(result)
}

// Get GetGasPrice number
func (h *Handler) GetGasPrice(w http.ResponseWriter, r *http.Request) {

	h.Log.Info("Entered GetGasPrice")
	getGasBody := createRequestBody(apis.GetGasPrice, []string{})
	result := &apis.GetGasPriceResponse{}
	resp, err := h.Resty.R().SetBody(getGasBody).
		SetResult(result).
		Post(h.Mainnet_http_endpoint)
	h.debugResponse("GetBlockNumber", resp, err)
	json.NewEncoder(w).Encode(result)
}

// GetBlockByNumber
func (h *Handler) GetBlockByNumber(w http.ResponseWriter, r *http.Request) {

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
		json.NewEncoder(w).Encode(h.TxDetailsResponse(formmattedRequest, apis.GetBlockByNumberTxDetailsResponse{}))
		return
	} else {
		json.NewEncoder(w).Encode(h.TxDetailsResponse(formmattedRequest, apis.GetBlockByNumberNoTxDetailsResponse{}))
	}
}

/*
* ParseGetBlockByNumber Request will take in an http request
* validate that block and txdetails exist and are valid
* it will then return either an Error message body, or the
* request body along with the value of txDetails
 */
func (h *Handler) ParseGetBlockByNumberRequest(r *http.Request) ([]byte, bool, bool) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var getBlockByNumberRequest apis.GetBlockByNumberRequest
	if err := json.Unmarshal(reqBody, &getBlockByNumberRequest); err != nil {
		h.Log.Error("Error unmarshalling GetBlockByNumberRequest", zap.Error(err))
		errorBody, _ := json.Marshal(apis.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error()})
		return errorBody, false, false
	}

	txdetails, err := strconv.ParseBool(getBlockByNumberRequest.TxDetails)
	if getBlockByNumberRequest.Block == "" || err != nil {
		errorBody, _ := json.Marshal(apis.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    apis.MalformedRequestMessage})
		return errorBody, false, false
	}
	// Can't use create RequestBody because 2nd param is bool with no quotes
	body := []byte(fmt.Sprintf(apis.GetBlockByNumberRequestBodyTemplate, getBlockByNumberRequest.Block, getBlockByNumberRequest.TxDetails))
	h.Log.Info("GetBlockByNumber body", zap.String("Body", string(body)))
	return body, true, txdetails
}

func (h *Handler) TxDetailsResponse(body []byte, unmashallStruct interface{}) interface{} {
	var err error
	var resp *resty.Response

	resp, err = h.Resty.R().SetBody(body).
		Post(h.Mainnet_http_endpoint)
	if err != nil {
		h.Log.Error("Error", zap.Error(err))
	}
	h.debugResponse("GetBlockByNumber", resp, err)
	switch unmashallStruct.(type) {
	case apis.GetBlockByNumberTxDetailsResponse:
		result := &apis.GetBlockByNumberTxDetailsResponse{}
		json.Unmarshal(resp.Body(), result)
		return result
	case apis.GetBlockByNumberNoTxDetailsResponse:
		result := &apis.GetBlockByNumberNoTxDetailsResponse{}
		json.Unmarshal(resp.Body(), result)
		return result
	default:
		return &apis.ErrorResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}
}

// GetBlockByNumber
func (h *Handler) GetTransactionByBlockNumberAndIndex(w http.ResponseWriter, r *http.Request) {

	var err error
	var resp *resty.Response
	w.Header().Set("Content-Type", "application/json")
	h.Log.Info("Entered GetTransactionByBlockNumberAndIndex")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var getTxReq apis.GetTransactionByBlockNumberAndIndexRequest
	if err := json.Unmarshal(reqBody, &getTxReq); err != nil {
		h.Log.Error("Error unmarshalling GetBlockByNumberRequest", zap.Error(err))
	}

	if getTxReq.Block == "" || getTxReq.Index == "" {
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    apis.MalformedRequestMessage,
		})
		return
	}

	getBlockNumberAndTxBody := createRequestBody(apis.GetTransactionByBlockNumberAndIndex, []string{getTxReq.Block, getTxReq.Index})
	result := &apis.GetTransactionByBlockNumberAndIndexResponse{}
	resp, err = h.Resty.R().SetBody(getBlockNumberAndTxBody).
		SetResult(result).
		Post(h.Mainnet_http_endpoint)
	if err != nil {
		h.Log.Error("Error", zap.Error(err))
	}
	h.debugResponse("GetTransactionByBlockNumberAndIndex", resp, err)
	json.NewEncoder(w).Encode(result)

}

func (h *Handler) debugResponse(caller string, resp *resty.Response, err error) {
	h.Log.Info("Handling response from", zap.String("caller", caller))
	// Explore response object
	h.Log.Info("Response Info:",
		zap.Error(err),
		zap.Int("Status Code:", resp.StatusCode()),
		zap.String("Status     :", resp.Status()),
		zap.String("Proto      :", resp.Proto()),
		zap.Time("Received At:", resp.ReceivedAt()))
	// zap.String("Body :\n", string(resp.Body())))
}

func createRequestBody(method apis.RPCCall, params []string) *apis.InfuraRequestBody {
	return &apis.InfuraRequestBody{
		JsonRPC: apis.RPCVersion2,
		Method:  method,
		Params:  params,
		ID:      apis.RequestID,
	}
}
