import http from 'k6/http';
import { sleep } from 'k6';
export let options = {
  vus: 10,
  duration: '30s',
};

const requests = [
	`{"block": "0x5bad55","txdetails": false}`,
	`{"block": "0x5bad55","txdetails": true}`,
	`{"block": "earliest","txdetails": true}`,
	`{"block": "earliest","txdetails": false}`,
	`{"block": "0xc6dad0","txdetails": true}`,
	`{"block": "0xc6dace","txdetails": false}`,
	`{"block": "latest","txdetails": true}`,
	`{"block": "latest","txdetails": false}`,
	// Some bad requests too
//	`{"block": "bad-block","txdetails": false}`,
//	`{"ck": "latest","txdetails": false}`,
//	`{"block": "latest","txdetails": fae}`,
//	`{"block": "0xc6dace","txdetas": false}`,
	];

let url = "SED-URL"
export default function () {
  var request = requests[Math.floor(Math.random()*requests.length)];
  let formData = request
  let headers = { 'Content-Type': 'application/x-www-form-urlencoded' };
  http.post(url+='/blockbynumber', formData, { headers: headers });
}

