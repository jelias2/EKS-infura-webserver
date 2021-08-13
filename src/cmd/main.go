package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jelias2/infra-test/src/handlers"
	"go.uber.org/zap"
)

//TODO:  Make all logs consistent format!!!

// Main function
func main() {

	log, _ := zap.NewProduction()
	defer log.Sync()
	log.Info("Beginning Webserver main.go...")
	// Init router
	log.Info("Creating mux router and initalizing mux router")
	r := mux.NewRouter()

	projectID := os.Getenv("PROJECT_ID")
	projectSecret := os.Getenv("PROJECT_SECRET")
	mainnetHTTPEndpoint := os.Getenv("MAINNET_HTTP_ENDPOINT")
	mainnetWebsocketEndpoint := os.Getenv("MAINNET_WEBSOCKET_ENDPOINT")

	log.Info("Config vars",
		zap.String("Project_id", projectID),
		zap.String("projectSecret", projectSecret),
		zap.String("mainnetHTTPEndpoint", mainnetHTTPEndpoint),
		zap.String("mainnetWebsocketEndpoint", mainnetWebsocketEndpoint),
	)

	interrupt := make(chan os.Signal, 1) // Channel to listen for interrupt signal to terminate gracefully
	signal.Notify(interrupt, os.Interrupt)

	log.Info("Websocket connecting to", zap.String("Url", mainnetWebsocketEndpoint))
	ws_client, _, err := websocket.DefaultDialer.Dial(mainnetWebsocketEndpoint, nil)
	if err != nil {
		log.Fatal("Fatal Dial Error:", zap.Error(err))
	}
	defer ws_client.Close()

	go func() {
		for {
			sig := <-interrupt
			if sig != nil {
				log.Info("Websocket Recieved Interrupt, closing channel")
				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := ws_client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Info("Error Closing Socket:", zap.Error(err))
					return
				}
				// Without this ctrl-c will kills the websocket, and leave the webserver hanging
				os.Exit(1)
			}
		}
	}()

	handler := &handlers.Handler{
		Log:                        log,
		Resty:                      resty.New(),
		Mainnet_websocket_endpoint: mainnetWebsocketEndpoint,
		Mainnet_http_endpoint:      mainnetHTTPEndpoint,
		WebSocket:                  ws_client,
	}

	// Route handles & endpoints
	r.HandleFunc("/health", handler.Healthcheck).Methods("GET")
	r.HandleFunc("/", handler.Healthcheck).Methods("GET")

	r.HandleFunc("/blocknumber", handler.GetBlockNumber).Methods("GET")
	r.HandleFunc("/gasprice", handler.GetGasPrice).Methods("GET")
	r.HandleFunc("/blockbynumber", handler.GetBlockByNumber).Methods("POST")
	r.HandleFunc("/txbyblockandindex", handler.GetTransactionByBlockNumberAndIndex).Methods("POST")

	r.HandleFunc("/ws/blocknumber", handler.WebSocketGetBlockNumber).Methods("GET")
	r.HandleFunc("/ws/gasprice", handler.WebSocketGetGasPrice).Methods("GET")
	r.HandleFunc("/ws/blockbynumber", handler.WebSocketGetBlockByNumber).Methods("POST")
	r.HandleFunc("/ws/txbyblockandindex", handler.WebSocketGetTransactionByBlockNumberAndIndex).Methods("POST")

	// Start server
	log.Info("Beginning to server traffic on port")
	http.ListenAndServe(":8000", r)
}
