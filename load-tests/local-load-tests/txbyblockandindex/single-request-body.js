import http from 'k6/http';
import { sleep } from 'k6';
export let options = {
  vus: 10,
  duration: '30s',
};

let url = "SED-URL"

export default function () {
  var formData = `{"block": "0xc68e80","index": "0x11"}`;
  var headers = { 'Content-Type': 'application/x-www-form-urlencoded' };
  http.post(url+='/txbyblockandindex', formData, { headers: headers });
}

