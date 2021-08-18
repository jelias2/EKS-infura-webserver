import http from 'k6/http';

let url = "SED-URL"
export default function () {
  var formData = `{"block": "0xc68e80","index": "0x11"}`;
  var headers = { 'Content-Type': 'application/x-www-form-urlencoded' };
  http.post(url+='/txbyblockandindex', formData, { headers: headers });
}

