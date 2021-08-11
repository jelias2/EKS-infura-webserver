package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"github.com/jelias2/infra-test/apis"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type Handler struct {
	Log                        *zap.Logger
	Resty                      *resty.Client
	Mainnet_http_endpoint      string
	Mainnet_websocket_endpoint string
}

// Healthcheck will display test response to make sure the server is running
func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {

	log.Info("Healthcheck ")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apis.Healthcheck{
		Status:   200,
		Message:  "Healthcheck response",
		Datetime: time.Now().String(),
	})
}

// Get all books
func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	// Hardcoded data - @todo: add database
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apis.Book{})
}

// Get ethblock number
func (h *Handler) GetBlockNumber(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Entered GetBlockNumber")
	body := `{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`
	h.Log.Info("GetBlockNumber body", zap.String("Body", body))
	result := &apis.GetBlockNumberResponse{}
	resp, err := h.Resty.R().SetBody(body).
		SetResult(result).
		Post(h.Mainnet_http_endpoint)
	h.handleResponse("GetBlockNumber", resp, err)
	json.NewEncoder(w).Encode(result)
}

// Get GetGasPrice number
func (h *Handler) GetGasPrice(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Entered GetGasPrice")
	body := `{"jsonrpc":"2.0","method":"eth_gasPrice","params": [],"id":1}`
	result := &apis.GetBlockNumberResponse{}
	resp, err := h.Resty.R().SetBody(body).
		SetResult(result).
		Post(h.Mainnet_http_endpoint)
	h.handleResponse("GetBlockNumber", resp, err)
	json.NewEncoder(w).Encode(result)
}

// GetBlockByNumber
func (h *Handler) GetBlockByNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.Log.Info("Entered GetBlockByNumber")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var GetBlockByNumberRequest apis.GetBlockByNumberRequest
	if err := json.Unmarshal(reqBody, &GetBlockByNumberRequest); err != nil {
		h.Log.Error("Error unmarshalling GetBlockByNumberRequest", zap.Error(err))
	}

	var block string
	var txdetails bool
	var err error
	switch strings.ToLower(GetBlockByNumberRequest.Block) {
	case "latest":
		block = "latest"
	case "earliest":
		block = "earliest"
	case "pending":
		block = "pending"
	default:
		block = GetBlockByNumberRequest.Block
	}

	if txdetails, err = strconv.ParseBool(GetBlockByNumberRequest.TxDetails); err != nil {
		h.Log.Error("Error Converting GetBlockByNumberRequest.TxDetails to boolean", zap.Error(err))
	}

	h.Log.Info("GetBlockByNumber params", zap.String("block", block), zap.Bool("tx_details", txdetails))
	body := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s",%s],"id":1}`, block, GetBlockByNumberRequest.TxDetails)
	h.Log.Info("GetBlockByNumber body", zap.String("Body", body))

	var resp *resty.Response
	if txdetails {
		result := &apis.GetBlockByNumberTxDetailsResponse{}
		resp, err = h.Resty.R().SetBody(body).
			SetResult(result).
			Post(h.Mainnet_http_endpoint)
		if err != nil {
			h.Log.Error("Error", zap.Error(err))
		}
		h.handleResponse("GetBlockByNumber", resp, err)

		json.NewEncoder(w).Encode(result)

	} else {
		result := &apis.GetBlockByNumberNoTxDetailsResponse{}
		resp, err = h.Resty.R().SetBody(body).
			SetResult(result).
			Post(h.Mainnet_http_endpoint)
		if err != nil {
			h.Log.Error("Error", zap.Error(err))
		}
		h.handleResponse("GetBlockByNumber", resp, err)

		json.NewEncoder(w).Encode(result)
	}

}

// func HexStringToInt(hexaString string) (intString string, err error) {
// 	// replace 0x or 0X with empty String
// 	numberStr := strings.Replace(hexaString, "0x", "", -1)
// 	numberStr = strings.Replace(numberStr, "0X", "", -1)
// 	output, err := strconv.ParseInt(numberStr, 16, 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	intString = string(output)
// 	return
// }

// func (h *Handler) PostRequest(endpoint string, result interface{}, body string) (interface{}, error) {
// 	h.Log.Info("Entered PostRequest")
// 	resp, err := h.Resty.R().SetBody(body).
// 		SetResult(result).
// 		Post(endpoint)
// 	h.Log.Info("Recieved response from PostRequest")
// 	h.handleResponse(resp, err)
// 	return resp.Body, err
// }

func (h *Handler) handleResponse(caller string, resp *resty.Response, err error) {
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

// Get single book
func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	var books []apis.Book
	// Loop through books and find one with the id from the params
	books = append(books, apis.Book{ID: "1", Isbn: "438227", Title: "apis.Book One", Author: &apis.Author{Firstname: "John", Lastname: "Doe"}})
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&apis.Book{})
}

// Add new book
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book apis.Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe

	var books []apis.Book
	// Loop through books and find one with the id from the params
	books = append(books, apis.Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &apis.Author{Firstname: "John", Lastname: "Doe"}})
	json.NewEncoder(w).Encode(book)
}

// Update book
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var books []apis.Book
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book apis.Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

// Delete book
func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var books []apis.Book
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}
