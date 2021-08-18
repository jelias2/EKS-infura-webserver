import ws from 'k6/ws';
import { check, sleep } from 'k6';


export let options = {
  stages: [
    { duration: "10s", target: 10 },
    { duration: "30s", target: 200 },
    { duration: "2m", target: 500 },
    { duration: "1m", target: 200 },
    { duration: "10s", target: 50 },
    { duration: "10s", target: 10 },
  ],
}

const requests = [
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["0x5bad55",false],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["0x5bad55",true],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["earliest",false],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["earliest",true],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["0xc6dad0",true],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["0xc6dad0",false],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["latest",true],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["latest",false],"id":1}`,
	// Some bad requests too
	`{"sonrpc":"2.0","method":"eth_getckByNumber","params": ["latest",true],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","pams": ["latest",true],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["last",tue],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["latest",true],"i":1}`,
	];

let url="SED-URL"
export default function () {
  let url = (url+='/socket2socket')
  let params = { tags: { my_tag: 'my ws session' } };

  let res = ws.connect(url, params, function (socket) {
    socket.on('open', function open() {
      console.log(`VU ${__VU}: connected`);
        socket.send(`{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`);
    });
  });

  check(res, { 'Connected successfully': (r) => r && r.status === 101 });
}
