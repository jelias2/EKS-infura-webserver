import { randomString, randomIntBetween } from "https://jslib.k6.io/k6-utils/1.1.0/index.js";
import ws from 'k6/ws';
import { check, sleep } from 'k6';

let sessionDuration = 30000 // user session between 30s

export let options = {
  vus: 10,
  iterations: 10, 
  duration: '30s',
};

export default function () {
 //let url = ('ws://host.docker.internal:8000/socket2socket')
  let url = ('ws://localhost:8000/socket2socket')
  let params = { tags: { my_tag: 'my ws session' } };

  let res = ws.connect(url, params, function (socket) {
    socket.on('open', function open() {
      console.log(`VU ${__VU}: connected`);

      socket.setInterval(function timeout() {		
          socket.send(`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0x5BAD55","0x0"],"id":1}`);
      }, randomIntBetween(2000, 8000));

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
