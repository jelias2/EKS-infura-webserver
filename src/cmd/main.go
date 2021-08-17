package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jelias2/infra-test/src/apis"
	"github.com/jelias2/infra-test/src/handlers"
	"go.uber.org/zap"
)

var (
	projectID                string
	projectSecret            string
	mainnetHTTPEndpoint      string
	mainnetWebsocketEndpoint string
	err                      error
)

// Main function
//TODO:  Make all logs consistent format!!!
func main() {
	flag.Parse()
	log, _ := zap.NewProduction()
	defer log.Sync()
	log.Info("Beginning Webserver main.go...")
	// Init router
	log.Info("Creating mux router and initalizing mux router")
	r := mux.NewRouter()

	projectID = os.Getenv("PROJECT_ID")
	projectSecret = os.Getenv("PROJECT_SECRET")
	mainnetHTTPEndpoint = os.Getenv("MAINNET_HTTP_ENDPOINT")
	mainnetWebsocketEndpoint = os.Getenv("MAINNET_WEBSOCKET_ENDPOINT")

	log.Info("Config vars",
		zap.String("Project_id", projectID),
		zap.String("projectSecret", projectSecret),
		zap.String("mainnetHTTPEndpoint", mainnetHTTPEndpoint),
		zap.String("mainnetWebsocketEndpoint", mainnetWebsocketEndpoint),
	)

	log.Info("Creating websockets for endpoint connecting to", zap.String("Url", mainnetWebsocketEndpoint))
	var wsClients = make(map[apis.ClientName]*websocket.Conn)
	for _, endpoint := range apis.AllWsClients {
		ws_client, _, err := websocket.DefaultDialer.Dial(mainnetWebsocketEndpoint, nil)
		if err != nil {
			log.Fatal("Error creating websocket clients")
		}
		wsClients[endpoint] = ws_client
	}

	handler := &handlers.Handler{
		Log:                        log,
		Resty:                      resty.New(),
		Mainnet_websocket_endpoint: mainnetWebsocketEndpoint,
		Mainnet_http_endpoint:      mainnetHTTPEndpoint,
		WsClients:                  wsClients,
	}

	defer handler.WsClients[apis.WsBlockNumber].Close()
	defer handler.WsClients[apis.WsBlockByNumber].Close()
	defer handler.WsClients[apis.WsGasPrice].Close()
	defer handler.WsClients[apis.WsTxByBlockNumberAndIndex].Close()

	interrupt := make(chan os.Signal, 1) // Channel to listen for interrupt signal to terminate gracefully
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		for {
			// On system interrupt shut all sockets down
			sig := <-interrupt
			if sig != nil {
				log.Info("Websocket Recieved Interrupt, closing channel")
				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				for _, endpoint := range apis.AllWsClients {
					handler.WsClients[endpoint].WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Fatal("Error Closing Websocket", zap.String("Websocket", string(endpoint)), zap.Error(err))
					}
				}
				// Without this ctrl-c will kills the websocket, and leave the webserver hanging
				os.Exit(1)
			}
		}
	}()

	// Route handles & endpoints
	r.HandleFunc("/health", handler.Healthcheck).Methods("GET")
	r.HandleFunc("/", handler.Healthcheck).Methods("GET")

	r.HandleFunc("/blocknumber", handler.GetBlockNumber).Methods("GET")
	r.HandleFunc("/gasprice", handler.GetGasPrice).Methods("GET")
	r.HandleFunc("/blockbynumber", handler.GetBlockByNumber).Methods("POST")
	r.HandleFunc("/txbyblockandindex", handler.GetTransactionByBlockNumberAndIndex).Methods("POST")

	r.HandleFunc("/ws/health", handler.Healthcheck).Methods("GET")
	r.HandleFunc("/ws/blocknumber", handler.WebSocketGetBlockNumber).Methods("GET")
	r.HandleFunc("/ws/gasprice", handler.WebSocketGetGasPrice).Methods("GET")
	r.HandleFunc("/ws/blockbynumber", handler.WebSocketGetBlockByNumber).Methods("POST")
	r.HandleFunc("/ws/txbyblockandindex", handler.WebSocketGetTransactionByBlockNumberAndIndex).Methods("POST")

	r.HandleFunc("/socket2socket", handler.Socket2socket)

	// Start server
	log.Info("Beginning to server traffic on port")
	log.Fatal("Error Serving traffic ", zap.Error(http.ListenAndServe(":8000", r)))
}
