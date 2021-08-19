import http from 'k6/http';
import { check, sleep} from 'k6';

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

let	blockNumberReqs = [
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
let gaspriceUrl=(url+'/gasprice')
let blockNumUrl=(url+'/blocknumber')
let blockByNumUrl=(url+'/blockbynumber')
let txIndexUrl=(url+'/txbyblockandindex')
export default function () {
	let headers = { 'Content-Type': 'application/json' };

	let res = http.get(gaspriceUrl, { headers: headers });

	res = http.get(blockNumUrl, { headers: headers });

	let blockNuReq = blockNumberReqs[Math.floor(Math.random()*blockNumberReqs.length)];
	res = http.post(blockByNumUrl, blockNuReq, { headers: headers });


	var txIReq = txIndexReqs[Math.floor(Math.random()*txIndexReqs.length)];
	http.post(txIndexUrl, txIReq, { headers: headers });
}

