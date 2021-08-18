# Infura Web Server
* Welcome to my verion of the take home project. I've created a webserver written in go to serve Infura api data over 3 possible data protocol paths. My inspiration for this came from watching a short infura demo on youtube and listing to the speaker mention a push to websockets [here](https://youtu.be/OPtt2SiQ4qk?t=430). I had never used websockets before, and I thought it would be fun task to loadtest the efficiency of websockets to standard HTTP rest using the following setups.
  * ```client <-HTTP-> jelias_infura_server <-HTTP-> Infura ```
  * ```client <-HTTP-> jelias_infura_server <-Websocket-> Infura ```
  * ```client <-Websocket-> jelias_infura_server <-Websocket-> Infura ```
 *  All basic HTTP routes can be prexfixed with "/ws" to access the backend websocket version. The /socket2socket endpoint uses the websocket proctocl and can accept all possible Infura websocket calls"

# Table of contents
- [Infura Web Server](#infura-web-server)
- [Table of contents](#table-of-contents)
  - [How to Run <a name="howtrun"></a>](#how-to-run-)
  - [How to Load Test Locally <a name="howtoloadtest"></a>](#how-to-load-test-)
  - [Repository Setup <a name="repositorysetup"></a>](#repository-setup-)
  - [Makefile Commands: <a name="makefilecommands"></a>](#makefile-commands-)
  - [Endpoint Documentation <a name="endpointdocumentation"></a>](#endpoint-documentation-)
  - [Troubleshooting <a name="troubleshooting"></a>](#troubleshooting-)
  - [Finale: Cranking Up the Loadtests <a name="crankload"></a>](#crankload-)

## How to Run <a name="howtrun"></a>
  #### Locally 
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
  #### Cloud Context
  1. Follow the steps to configure aws and EKS with terraform https://learn.hashicorp.com/tutorials/terraform/eks
  1. Build and push the image to ECR: https://www.stacksimplify.com/aws-eks/aws-ecr-eks/learn-to-use-docker-images-built-and-pushed-to-aws-ecr-and-use-in-aws-eks/
  1. Edit the k8s secret in the deployment.yaml with the appropriate values
  1. kubectl apply deployment.yaml
  1. Get the service endpoint by examining the external IP in the output of ```kubectl get service -n infura infura-webserver-loadbalancer```


## How to Load Test <a name="howtoloadtest"></a>
  * Fill out and source .envrc
  * Install k6 package for your machine [k6 installation instructions](https://k6.io/docs/getting-started/installation/)
  * Sync submodules and init 
    * ```git submodules sync &&  git submodule update --init --recursive```
  * Spin up Grafana and influx db for visualization 
    * ```$WORKSPACE/load-tests/k6 && docker-compose up -d influxdb grafana```
    * Navigate to localhost:3000 in your browser for the grafana interface
    * Import via grafana.com: Enter 2587. Click load
    * TODO add screenshot
    * On the next page find the k6 dropdown, select "myinfluxdb (default)"
    * import 
  * Run localized load-tests (dockerized or in cloud requires indivudal  configuration at this time)
    * back to local-load-tests ```cd $WORKSPACE/load-tests/local-load-tests```
  * Load test command execution ```./run-load-test.sh <endpoint>/<test-name> <# of users> <time-duration>```
  * Example ```./run-load-test.sh gasprice/gasprice-loadtest.js 10 30s```
  * Example ``` ./run-load-test.sh txbyblockandindex/single-request-body.js 100 30s```
  * Watch the results over grafana and the k6 output at end of test
  * Note: socket2socket tests results will not show up in grafana default dashboard as the websocket metrics are different fields than http -rest versions

## Repository Setup <a name="repositorysetup"></a>
  * /src contains all of the source code
  * /cmd contains the main.go which executes setup and initalization of the webserver
  * /handlers: the files which perform the logic for each endpoint
    * handlers.go: contains the pure HTTP REST endpoints and logic
    * infura-websocket-client.go: Contains the hybrid HTTP REST Websocket endpoints and logic
    * socket2socket.go: contains a websocket implementation server
  * /apis
    * apis.go contains the basic kinds and json mashalling structure for the webserver
    * block.go: structs relating to block request and responses\
    * transactions.go structs releating to transaction request and responses
  * /load-tests contains scripts for load-testing: see more info in the load-testing section
  * /build contains build artifacts
  * /deploy contains k8s deployment code and EKS terraform code

## Makefile Commands: <a name="makefilecommands"></a>
  * ```make bin``` -> Will build the binary with the artifact sent to /build
  * ```make binrun``` -> Will build a the binary and run it
  * ```make clean``` -> Will remove the build directory and perform a go mod tidy
  * ```make docker``` -> Will build the binary and the dockerized version of the web server
  * ```make docker-run``` -> Will build the binary, make the docker version, and run it on port 8000 in detached mode
  * ```make dockerk8s``` -> Will build the binary, make the docker version but without the environment variables due to injection via k8s secret


## Endpoint Documentation <a name="endpointdocumentation"></a>
* ```GET /health or /ws/health``` 
    * Will return a short message with a timestamp to display that the server is alive and running
    * ```{"status": 202, "message": "Healthcheck response", "datetime": "2021-08-15 19:03:00 607301 -0500 CDT m=+32283.828596254"}```
    * Note: using the /ws route does not actually use websockets as this does not reach out to infura
* ```GET /blocknumber or /ws/blocknumber```
    * Will return a 200 the current block of the Ethereum main chain in hex representation 
    *Example Response:  ```{"jsonrpc": "2.0","id": 1,"result": "0xc6dad0"}```
* ```GET /gasprice or /ws/gasprice```
    * Will return the current gas price of the Ethereum mainchain in hex representation
    * Example Resonse:  ```{"jsonrpc":"2.0","id":1,"result":"0xa23b835ca2"}```
* ```POST /blockbynumber or /ws/blockbynumber```
    * Takes in two parameters block of type string [required], and txdetails [required] of type bool and will return the block information and details of the included transactions txdetails is true
    * block can be an integer block number, or the string "latest", "earliest" or "pending"
    * Example Body: ```{"block": "latest", "txdetails": false }```
    * Example Reponse: ``` {"jsonrpc":"2.0","id":1,"result":{"difficulty":"0x1bab98f5272273","extraData":"0x65746865726d696e652d617369612d6561737432","gasLimit":"0x1c9c380","gasUsed":"0x1c9918a","hash":"0x5954aa6d2abdd9a354fa7ff294a7c82675db3c1116b3c1c779ddd19191167745","logsBloom":"0x35a3f18793cbf99793db6d7bc9dd7fa5ed4ad81b0ecd9674ea59e97382f4f7a2b5b64753a3ebdeb8ccec7bd68bbbc77dcf65dfd78fbdfffd0bbebff6637e7c363e3df9bd4bbf7f4ffe9f7ffe12145af12dcf6e5eb6edf79d6fc6dfc7c2e9d7925b99f9a75eef4dced4ec99a62d94fd7b865f577dddfdcf6ba7fefcffa82a5b5efbabdf2faf76bf7ed37f79793ffd7fbfed7badd5ff3b967bf3afdfce68b381b7d7fea56bee9f65e7de83f792b5e986cdd6fb9305db8379eaf3a7e1633eefdbf7e027a0676d0a19dab96eb32f91de9dbf2d9df276bf8b14da6f5b9fd3dfce7bb6fd3d732a7dac7bf9e2cdcfe3f26afda430e98cf260fcff73bab7bffd1adfbbff","miner":"0xea674fdde714fd979de3edf0f56aa9716b898ec8","mixHash":"0xe2909e7d5264275b04b1a2fa40138531bdf5c056146cf230ba104e67fdac7770","nonce":"0x7eef1d5ee98e13f3","number":"0xc6af55","parentHash":"0x24ca43f9bb904d1b6f2474ade0f3476a7491f9786c5f7a42665f61cdbbdf376f","receiptsRoot":"0xd56b5606ab26df620f6d55772a9d6fa51258e14150de1534bb6c8c82621f2040","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x215fc","stateRoot":"0x1f8db8cb8db83f2673696670aa8ae41667e1f0b27af8050df5ed1d46a6971ead","timestamp":"0x61174010","totalDifficulty":"0x61ff17f039d5622be9e","transactions":[{"blockHash":"0x5954aa6d2abdd9a354fa7ff294a7c82675db3c1116b3c1c779ddd19191167745","blockNumber":"0xc6af55","from":"0xea674fdde714fd979de3edf0f56aa9716b898ec8","gas":"0x3d090","gasPrice":"0x6c3bcfc25","hash":"0xbe2cda833ff41fab2dda3c51266c06814723825cc5b7553a949a17a7e2c0dd2a","input":"0x","nonce":"0x2306af0","r":"0x83d11143feb4fb2e80106329d8f2c7ceb1769c5e567b2430498318ce56bdd963","s":"0x45e83063f7febbe4624f7445d73264b5531704d27c9d0b9010c486836d3668f4","to":"0x78a85e5baa0a02da50cfeebd573555668cdda36d","transactionIndex":"0x0","v":"0x0","value":"0x16215043e88dff3"}, ``` 
* ```POST /txbyblockandindex or /ws/txbyblockandindex```
    * Takes in two parameters block of type string [required] and index of type string [required], and will return the specific transaction located at the block and index
    * Example Body: ```{"block": "0xc68e80","index": "0x11"}```
    * Example Response: TODO
* ```WS /socket2socket```
    * socket2socket endpoint will open a websocket connection to the server, and will allow for websocket commuication to infura websocket server. All requests from the infura websocket api documentation are valid. 
    * Example Request: ```{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["latest",false],"id":1}```
    * Example Response:  TODO
    * Example Request: ```{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0x5BAD55","0x0"],"id":1}```
    * Example Response: TODO

## Finale: Cranking Up the Loadtests <a name="crankload"></a>]
* My Loadtest proceedure was as follows
* Using EKS and terraform I deployment my webserver and exposed it via a load-balancer. The deployment requested 3Gi Memory and 1 vpcu on a AWS t3.medium with up to 5GB network performance
  * I chose a cloud enviroment as opposed to local for two reasons
  1. To remove low latency advantages that come with local deployments as this would interfere with test results
  1. To scale my server as beyond the resources of my local machine
* Once the service was deployed I used the k6 tool to perform load testing. The load-test tested a single api endpoint and slowly scaled up and down the traffic. The following coniguration was used. 
```    
Stage 1: { duration: "10s", target: 10 } Every second for 10 second, 10 users requesting would make a request
Stage 2: { duration: "30s", target: 200 } Every second for 30 second, 200 users requesting would make a request
Stage 3: { duration: "2m", target: 500 }, Every second for 2m second, 500 users requesting would make a request
Stage 4: { duration: "1m", target: 200 }, Every second for 1m second, 200 users requesting would make a request
Stage 5: { duration: "10s", target: 50 }, Every second for 10s second, 50 users requesting would make a request
Stage 6: { duration: "10s", target: 10 },  Every second for 10s second, 10 users requesting would make a request
```
* **HTTP-HTTP Results**
  * The fastest RTT to infura is a 4 way tie between all the endpoints coming in around 40 ms (Excluding Healthcheck )
    * The fastest overall endpoint was the healthcheck  21 ms this is obvious because the health check does not reach out to Infura    
  * The second fastest endpoint was getBlockNumber with 95% of requests served in under 70.82 ms
    * I reckon this is because getBlockNumber is a relatively simple to cache request and has a small payload
  * The slowest endpoint(s) are close getGasPrice and blockByNumber come in at 95% of requests being served in ~107ms
    * Worth noting that getGasPrice has much lower mean duration of 67 ms, while blockByNumber had a mean of 87 ms
  * The slowest response time of all was a 6 second call to blockByNumber
* **Websocket to Websocket Connection Results**
  * The k6 websocket load test uses different metrics so it will be difficult to compare directly to HTTP-HTTP
  * Inital Connection Time is expensive: 95% of connecting took under 86ms, average was 74.2 ms. Comparing this to HTTP RTT request in under 70ms. 
  * Approx 10.1 total messages are send per 1 second. 5.1 sending and 5.1 recieving messages.
    * This could imply the total RTT from local load test -> infura is approx <100ms = 10.1 messages / 10 seconds which above our average HTTP RTT
  * TODO: 


* **HTTP to Websocket Connection Results**
  * This approach died pretty quickly once it got to the load testing. Errors such below occured almost instantly
```
websocket: unexpected reserved bits 0x 

{"statuscode":400,"message":"websocket: close 1006 (abnormal closure): unexpected EOF"}
```
  * I would attribute these bugs to one infura client is writing and reading many HTTP requests from the websocket with zero concern for queue or concurrency
  * I think websocket client implementation which could handle connecurrent users would solve this issue. This could probably be done with RW mutex and channels
  * I think there still might be merit to this approach, a lower HTTP request time combined with an already established websocket connection could be worthwhile
  * 


## Troubleshooting <a name="troubleshooting"></a>
* ```ERRO[0010] Couldn't write stats                          error="{\"error\":\"Request Entity Too Large\"}\n" output=InfluxDBv1``` 
  * This error occurs when the loadtest data payload is too large for the influx DB. To adjust the payload limit, set the following line to the environment section of the influxdb service. Note this may crash influxdb if the test data is too much 
```
environment:
      - INFLUXDB_DB=k6
      - INFLUXDB_HTTP_MAX_BODY_SIZE=0 <--- Add this line
``` 
in the k6/docker-compose.yaml and restart the service via ```docker compose down &&  docker compose up -d influxdb grafana```

  *  Load testing error``` ERRO[0069] dial tcp 18.190.129.128:8000: socket: too many open files ```
    * This occurs when your machine has hit its limit maxmimum connection determine by ulimit -n
    * It can be overridden with the ulimit command it a user limit


## What I learned <a name="whatilearned"></a>

* Websockets, I had never implemented websockets before this was a chance to try out new technology.

* One websocket for shared across threads is a bad idea, hindsight is 20/20. Errors which showed up in the hybrid approach. Probably would need more time to engineer an elegant solution
  * One websocket per client worked suprisingly well 
```
websocket: unexpected reserved bits 0x 
{"statuscode":400,"message":"websocket: close 1006 (abnormal closure): unexpected EOF"}
```


* Too many open file descriptors when load testing: each server connection is a file descriptor and there is a soft limit in per use in linux space
```❯ ulimit -n
256
ulimit -n 2048
```

* Load testing framework k6, plus hands on with influx and grafana


