
import ws from 'k6/ws';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: "10s", target: 10 },
    { duration: "30s", target: 200 },
    { duration: "2m", target: 500 },
    { duration: "1m", target: 200 },
    { duration: "10s", target: 50 },
    { duration: "10s", target: 10 },
   ],
};

export default function () {
  let url = "SED-URL"
  let params = { tags: { my_tag: 'my ws session' } };

  let res = ws.connect(url+="/socket2socket", params, function (socket) {
    socket.on('open', function open() {
      console.log(`VU ${__VU}: connected`);
          socket.send(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["latest", false],"id":1}`);
    });
  });

  check(res, { 'Connected successfully': (r) => r && r.status === 101 });
}
