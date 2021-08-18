import http from 'k6/http';
import { sleep } from 'k6';
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

// Use only infura endpoints to base against socket2socket connections

const txIndexReqs = [
	`{"block": "0xc68e80","index": "0x11"}`,
	`{"block": "0xc68e80","index": "0x10"}`,
	`{"block": "0x5bad55","index": "0xF"}`,
	`{"block": "0x5bad55","index": "0xD"}`,
	`{"block": "0xc6dace","index": "0x11"}`,
	`{"block": "0xc6dace","index": "0x11"}`,
	`{"block": "latest","index": "0x11"}`,
	`{"block": "latest","index": "0x12"}`,
	`{"block": "0xc6fd2f","index": "0x02"}`,
	`{"block": "0xc6fd2f","index": "0x12"}`,
	`{"block": "0xc6fd31","index": "0xF1"}`,
	`{"block": "0xc6fd31","index": "0x20"}`,
	`{"block": "0xc6fd38","index": "0x01"}`,
	`{"block": "0xc6fd38","index": "0x21"}`,
	];

const blockNumberReqs = [
	`{"block": "0x5bad55","txdetails": "false"}`,
	`{"block": "0x5bad55","txdetails": "true"}`,
	`{"block": "earliest","txdetails": "true"}`,
	`{"block": "earliest","txdetails": "false"}`,
	`{"block": "0xc6dad0","txdetails": "true"}`,
	`{"block": "0xc6dace","txdetails": "false"}`,
	`{"block": "latest","txdetails": "true"}`,
	`{"block": "latest","txdetails": "false"}`,
	`{"block": "0xc6fd2f","txdetails": "false"}`,
	`{"block": "0xc6fd2f","txdetails": "true"}`,
	`{"block": "0xc6fd31","txdetails": "false"}`,
	`{"block": "0xc6fd31","txdetails": "true"}`,
	`{"block": "0xc6fd38","txdetails": "false"}`,
	`{"block": "0xc6fd38","txdetails": "true"}`,
	];

let url = "SED-URL"
export default function () {
	let headers = { 'Content-Type': 'application/x-www-form-urlencoded' };

	var blockNuReq = blockNumberReqs[Math.floor(Math.random()*blockNumberReqs.length)];
	//http.post(url+='/blockbynumber', blockNuReq, { headers: headers });

	var txIReq = txIndexReqs[Math.floor(Math.random()*txIndexReqs.length)];
	//http.post(url+='/txbyblockandindex', txIReq, { headers: headers });

	//http.get(url+='/blocknumber', { headers: headers });

      let req1 = {
         method: 'GET',
	 url: url+="/blocknumber",
         params: {
           headers: { 'Content-Type': 'application/json"' },
         },
      };
      let req2 = {
         method: 'GET',
	 url: url+="/blocknumber",
         params: {
           headers: { 'Content-Type': 'application/json"' },
         },
      };

      let req3 = {
        method: 'POST',
        url: url+="/txbyblockandindex",
        body: {
          block: "0xc6fd38",
	  txdetails: "0x10"
        },
        params: {
          headers: { 'Content-Type': 'application/json"' },
        },
      };


      let req4 = {
        method: 'POST',
        url: url+="/blockbynumber",
        body: {
          block: "0xc6fd38",
	  txdetails: "true"
        },
        params: {
          headers: { 'Content-Type': 'application/json"' },
        },
      };

      let responses = http.batch([req1, req2, req3, req4]);
}

