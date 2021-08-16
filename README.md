# Infura Web Server
* Welcome to my verion of the take home project

## How to Run
1. In the root directory fill out the .envrc file with some Infura Project Credentials
```export PROJECT_ID=<Infura-ProjectID>
export PROJECT_SECRET=<Infura-Project-Secret>
export MAINNET_HTTP_ENDPOINT=<Infura HTTP Endpoint>
export MAINNET_WEBSOCKET_ENDPOINT=<Infura-WS-Endpoint>
```
1. Source the .envrc file 
    * ```$ source .envrc```
1. Make the docker image
    * ```make docker```
1. Deploy the webserver locally via docker or via binary 
    * ```make docker-run``` -> via docker on localhost:8000
    * ```make binrun``` -> run plain binary localhost:8000
1. Begin using the endpoints via the Endpoint Documentation section below
    * If you are familar with postman you can download and import the postman api collection from `/load-tests/jelias-infura-rest.postman_collection.json```

## How to Load Test
  * Change directories to the load-testing area
     * ```cd load-tests```
  * Spin up Grafana and influx db for visualization
    * ```docker compose up```
    * 
  * For localized load-tests (dockerized or in cloud requires specialized configuration at this time)
    * ```cd local-load-tests```
  * 
    
## Repository Setup
  * /src contains all of the source code
  * /cmd contains the main.go which executes setup and initalization of the webserver
  * /handlers: the files which perform the logic for each endpoint
    * handlers.go: contains the HTTP REST endpoints and logic
    * socket2socket.go: contains a websocket implementation server
  * /apis
    * apis.go contains the basic kinds and json mashalling structure for the webserver
    * block.go: structs relating to block request and responses\
    * transactions.go structs releating to transaction request and responses
  * /load-tests contains scripts for load-testing: see more info in the load-testing section
  * /build contains build artifacts
* Makefile Commands:
  * ```make bin``` -> Will build the binary with the artifact sent to /build
  * ```make binrun``` -> Will build a the binary and run it
  * ```make clean``` -> Will remove the build directory and perform a go mod tidy
  * ```make docker``` -> Will build the binary and the dockerized version of the web server
  * ```make docker-run``` -> Will build the binary, make the docker version, and run it on port 8000 in detached mode
## Endpoint Documentation
* ```GET /health``` 
    * Will return a short message with a timestamp to display that the server is alive and running
    * ```{"status": 202, "message": "Healthcheck response", "datetime": "2021-08-15 19:03:00 607301 -0500 CDT m=+32283.828596254"}```
* ```GET /blocknumber```
    * Will return a 200 the current block of the Ethereum main chain in hex representation 
    *Example Response:  ```{"jsonrpc": "2.0","id": 1,"result": "0xc6dad0"}```
* ```GET /gasprice```
    * Will return the current gas price of the Ethereum mainchain in hex representation
    * Example Resonse:  ```{"jsonrpc":"2.0","id":1,"result":"0xa23b835ca2"}```
* ```POST /blockbynumber```
    * Takes in two parameters block of type string [required], and txdetails [required] of type bool and will return the block information and details of the included transactions txdetails is true
    * block can be an integer block number, or the string "latest", "earliest" or "pending"
    * Example Body: ```{"block": "latest", "txdetails": false }```
    * Example Reponse: ``` {"jsonrpc":"2.0","id":1,"result":{"difficulty":"0x1bab98f5272273","extraData":"0x65746865726d696e652d617369612d6561737432","gasLimit":"0x1c9c380","gasUsed":"0x1c9918a","hash":"0x5954aa6d2abdd9a354fa7ff294a7c82675db3c1116b3c1c779ddd19191167745","logsBloom":"0x35a3f18793cbf99793db6d7bc9dd7fa5ed4ad81b0ecd9674ea59e97382f4f7a2b5b64753a3ebdeb8ccec7bd68bbbc77dcf65dfd78fbdfffd0bbebff6637e7c363e3df9bd4bbf7f4ffe9f7ffe12145af12dcf6e5eb6edf79d6fc6dfc7c2e9d7925b99f9a75eef4dced4ec99a62d94fd7b865f577dddfdcf6ba7fefcffa82a5b5efbabdf2faf76bf7ed37f79793ffd7fbfed7badd5ff3b967bf3afdfce68b381b7d7fea56bee9f65e7de83f792b5e986cdd6fb9305db8379eaf3a7e1633eefdbf7e027a0676d0a19dab96eb32f91de9dbf2d9df276bf8b14da6f5b9fd3dfce7bb6fd3d732a7dac7bf9e2cdcfe3f26afda430e98cf260fcff73bab7bffd1adfbbff","miner":"0xea674fdde714fd979de3edf0f56aa9716b898ec8","mixHash":"0xe2909e7d5264275b04b1a2fa40138531bdf5c056146cf230ba104e67fdac7770","nonce":"0x7eef1d5ee98e13f3","number":"0xc6af55","parentHash":"0x24ca43f9bb904d1b6f2474ade0f3476a7491f9786c5f7a42665f61cdbbdf376f","receiptsRoot":"0xd56b5606ab26df620f6d55772a9d6fa51258e14150de1534bb6c8c82621f2040","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x215fc","stateRoot":"0x1f8db8cb8db83f2673696670aa8ae41667e1f0b27af8050df5ed1d46a6971ead","timestamp":"0x61174010","totalDifficulty":"0x61ff17f039d5622be9e","transactions":[{"blockHash":"0x5954aa6d2abdd9a354fa7ff294a7c82675db3c1116b3c1c779ddd19191167745","blockNumber":"0xc6af55","from":"0xea674fdde714fd979de3edf0f56aa9716b898ec8","gas":"0x3d090","gasPrice":"0x6c3bcfc25","hash":"0xbe2cda833ff41fab2dda3c51266c06814723825cc5b7553a949a17a7e2c0dd2a","input":"0x","nonce":"0x2306af0","r":"0x83d11143feb4fb2e80106329d8f2c7ceb1769c5e567b2430498318ce56bdd963","s":"0x45e83063f7febbe4624f7445d73264b5531704d27c9d0b9010c486836d3668f4","to":"0x78a85e5baa0a02da50cfeebd573555668cdda36d","transactionIndex":"0x0","v":"0x0","value":"0x16215043e88dff3"}, ``` 
* ```POST /txbyblockandindex```
    * Takes in two parameters block of type string [required] and index of type string [required], and will return the specific transaction located at the block and index
    * Example Body: ```{"block": "0xc68e80","index": "0x11"}```
    * Example Response: TODO
* ```WS /socket2socket```
    * socket2socket endpoint will open a websocket connection to the server, and will allow for websocket commuication to infura websocket server. All requests from the infura websocket api documentation are valid. 
    * Example Request: ```{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["latest",false],"id":1}```
    * Example Response:  TODO
    * Example Request: ```{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0x5BAD55","0x0"],"id":1}```
    * Example Response: TODO

## Load-Testing
 How to run: 
1. Install k6 package [here](https://k6.io/docs/getting-started/installation/)
2. suor