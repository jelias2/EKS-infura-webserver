import http from 'k6/http';


const requests = [
	`{"block": "0xc68e80","index": "0x11"}`,
	`{"block": "0xc68e80","index": "0x10"}`,
	`{"block": "0x5bad55","index": "0xF"}`,
	`{"block": "0x5bad55","index": "0xD"}`,
	`{"block": "0xc68e80","index": "0xE"}`,
	`{"block": "0xc68e80","index": "0xF1"}`,
	`{"block": "0xc6dace","index": "0x0"}`,
	`{"block": "0xc6dace","index": "0x11"}`,
	`{"block": "0xc6dace","index": "0x11"}`,
	`{"block": "latest","index": "0x11"}`,
	`{"block": "latest","index": "0x12"}`,
	 // Some bad requests too
	`{"bck": "0xc6dace","index": "0x11"}`,
	`{"block": "0ace","index": "0x11"}`,
	`{"block": "0xc6dace","dex": "0x11"}`,
	`{"block": "0xc6dace","index": "true"}`,
	];

let url = "SED-URL"
export default function () {
  var request = requests[Math.floor(Math.random()*requests.length)];
  let formData = request
  let headers = { 'Content-Type': 'application/x-www-form-urlencoded' };
  http.post(url+='/txbyblockandindex', formData, { headers: headers });
}

