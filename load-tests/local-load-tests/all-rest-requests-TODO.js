import http from 'k6/http';
import { sleep } from 'k6';
export let options = {
  vus: 10,
  duration: '30s',
};

// Use only infura endpoints to base against socket2socket connections
const endpoints = ["/blockbynumber", "/gasprice", "/blocknumber","/txbyblockandindex"]

const blockIndexReq = [
	`{"block": "0xc68e80","index": "0x11"}`,
	`{"block": "0xc68e80","index": "0x10"}`,
	`{"block": "0x5bad55","index": "0xF"}`,
	`{"block": "0x5bad55","index": "0xD"}`,
	`{"block": "0xc6dace","index": "0x11"}`,
	`{"block": "0xc6dace","index": "0x11"}`,
	`{"block": "latest","index": "0x11"}`,
	`{"block": "latest","index": "0x12"}`,
	];

const blockNumberReq = [
	`{"block": "0x5bad55","txdetails": false}`,
	`{"block": "0x5bad55","txdetails": true}`,
	`{"block": "earliest","txdetails": true}`,
	`{"block": "earliest","txdetails": false}`,
	`{"block": "0xc6dad0","txdetails": true}`,
	`{"block": "0xc6dace","txdetails": false}`,
	`{"block": "latest","txdetails": true}`,
	`{"block": "latest","txdetails": false}`,
	];

let url = "SED-URL"
export default function () {

	var request = blockIndexReq[Math.floor(Math.random()*blockIndexReq.length)];
	var request = blockNumberReq[Math.floor(Math.random()*blockNumberReq.length)];

	var endpoint = endpoint[Math.floor(Math.random()*endpoint.length)];
  
	let formData = request
	let headers = { 'Content-Type': 'application/x-www-form-urlencoded' };
	http.post(url+='/blockbynumber', formData, { headers: headers });
}

