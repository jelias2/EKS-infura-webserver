package main

import (
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"github.com/jelias2/infra-test/handlers"
	"go.uber.org/zap"
)

// Main function
func main() {

	log, _ := zap.NewProduction()
	defer log.Sync()
	log.Info("Beginning Webserver main.go...")
	// Init router
	log.Info("Creating mux router and initalizing mux router")
	r := mux.NewRouter()

	project_id := os.Getenv("PROJECT_ID")
	project_secret := os.Getenv("PROJECT_SECRET")
	mainnet_http_endpoint := os.Getenv("MAINNET_HTTP_ENDPOINT")
	mainnet_websocket_endpoint := os.Getenv("MAINNET_WEBSOCKET_ENDPOINT")

	log.Info("Creating go-resty client",
		zap.String("Project_id", project_id),
		zap.String("project_secret", project_secret),
		zap.String("mainnet_http_endpoint", mainnet_http_endpoint),
		zap.String("mainnet_websocket_endpoint", mainnet_websocket_endpoint),
	)

	handler := &handlers.Handler{
		Log:                        log,
		Resty:                      resty.New(),
		Mainnet_websocket_endpoint: mainnet_websocket_endpoint,
		Mainnet_http_endpoint:      mainnet_http_endpoint,
	}

	// Route handles & endpoints
	r.HandleFunc("/health", handler.Healthcheck).Methods("GET")
	r.HandleFunc("/", handler.Healthcheck).Methods("GET")
	r.HandleFunc("/blocknumber", handler.GetBlockNumber).Methods("GET")
	r.HandleFunc("/gasprice", handler.GetGasPrice).Methods("GET")
	r.HandleFunc("/getblockbynumber", handler.GetBlockByNumber).Methods("POST")
	r.HandleFunc("/books/{id}", handler.GetBlockByNumber).Methods("POST")
	// r.HandleFunc("/books", handler.CreateBook).Methods("POST")
	// r.HandleFunc("/books/{id}", handler.UpdateBook).Methods("PUT")
	// r.HandleFunc("/books/{id}", handler.DeleteBook).Methods("DELETE")

	// Start server
	log.Info("Beginning to server traffic on port")
	http.ListenAndServe(":8000", r)
}
