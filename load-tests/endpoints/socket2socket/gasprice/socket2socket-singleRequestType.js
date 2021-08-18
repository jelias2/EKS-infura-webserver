import ws from 'k6/ws';
import { check, sleep } from 'k6';

let url="SED-URL"
export default function () {
  let params = { tags: { my_tag: 'my ws session' } };

  let res = ws.connect(url+="/socket2socket", params, function (socket) {
    socket.on('open', function open() {
      console.log(`VU ${__VU}: connected`);
          socket.send(`{"jsonrpc":"2.0","method":"eth_gasPrice","params": [],"id":1}`);
    });
  });

  check(res, { 'Connected successfully': (r) => r && r.status === 101 });
}

