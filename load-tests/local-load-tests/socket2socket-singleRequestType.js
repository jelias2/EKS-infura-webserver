import { randomString, randomIntBetween } from "https://jslib.k6.io/k6-utils/1.1.0/index.js";
import ws from 'k6/ws';
import { check, sleep } from 'k6';

let sessionDuration = 30000 // user session between 30s

export let options = {
  stages: [
    { duration: "10s", target: 10 },
    { duration: "30s", target: 200 },
    { duration: "2m", target: 500 },
    { duration: "1m", target: 200 },
    { duration: "10s", target: 50 },
    { duration: "10s", target: 10 },

  //  { duration: "10s", target: 10 },
  //  { duration: "30s", target: 100 },
  //  { duration: "10s", target: 50 },
  //  { duration: "10s", target: 10 },
//
  ],
}

export default function () {
 //let url = ('ws://host.docker.internal:8000/socket2socket')
  let url = ('ws://a68db72516d7f4c8ba06db566ce9604c-631413878.us-east-2.elb.amazonaws.com:8000/socket2socket')
  let params = { tags: { my_tag: 'my ws session' } };

  let res = ws.connect(url, params, function (socket) {
    socket.on('open', function open() {
      console.log(`VU ${__VU}: connected`);
	
      socket.send(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`);

    });
  });

  check(res, { 'Connected successfully': (r) => r && r.status === 101 });
}
