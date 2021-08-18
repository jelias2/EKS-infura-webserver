import ws from 'k6/ws';
import { check} from 'k6';

let url="SED-URL"
export default function () {


  url+='/socket2socket'
  let params = { tags: { my_tag: 'my ws session' } };

  let res = ws.connect(url, params, function (socket) {
    socket.on('open', function open() {
      console.log(`VU ${__VU}: connected`);
      socket.send(`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params": ["0x5BAD55","0x0"],"id":1}`);
    });
  });

  check(res, { 'Connected successfully': (r) => r && r.status === 101 });
}
