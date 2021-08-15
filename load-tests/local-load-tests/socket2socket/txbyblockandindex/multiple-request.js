import { randomString, randomIntBetween } from "https://jslib.k6.io/k6-utils/1.1.0/index.js";
import ws from 'k6/ws';
import { check, sleep } from 'k6';

let sessionDuration = 30000 // user session between 30s

export let options = {
  vus: 10,
  iterations: 10, 
  duration: '30s',
};

const requests = [
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0xc68e80","0x11"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0xc68e80","0x10"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0x5bad55","0xF"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0x5bad55","0xD"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0xc68e80","0xE"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0xc68e80","0xF1"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0xc6dace","0x0"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0xc6dace","0x11"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["latest","0x11"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["latest","0x12"],"id":1}`,
	// Some bad requests too
	`{"jsopc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["latest","0x12"],"id":1}`,
	`{"jsonrpc":"2.0","thod":"eth_getTransactionByBlockNumberAndIndex","params": ["latest","0x12"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTrationByBlockNumberAndIndex","params": ["latest","0x12"],"id":1}`,
	`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","pars": ["latest","0x12"],"id":1}`,
	];


let url="SED-URL"

export default function () {
  let params = { tags: { my_tag: 'my ws session' } };
  let res = ws.connect(url+="/socket2socket", params, function (socket) {
    socket.on('open', function open() {
      console.log(`VU ${__VU}: connected`);

      var request = requests[Math.floor(Math.random()*requests.length)];
      socket.setInterval(function timeout() {		
          socket.send(request);
      });

    });

    socket.setTimeout(function () {
      console.log(`VU ${__VU}: ${sessionDuration}ms passed, leaving the chat`);
      socket.send(JSON.stringify({'event': 'LEAVE'}));

    }, sessionDuration-3000);

    socket.setTimeout(function () {
      console.log(`Closing the socket forcefully 3s after graceful LEAVE`);
      socket.close();
    }, sessionDuration);
  });

  check(res, { 'Connected successfully': (r) => r && r.status === 101 });
}
