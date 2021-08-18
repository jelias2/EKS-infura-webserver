import http from 'k6/http';
import { check, sleep } from "k6";

export let options = {
  stages: [
   { duration: "10s", target: 10 },
   { duration: "30s", target: 200 },
   { duration: "2m", target: 500 },
   { duration: "1m", target: 200 },
   { duration: "10s", target: 50 },
   { duration: "10s", target: 10 },
  ],
}
let url = "http://ada22e80555734f068d24df886149515-237930780.us-east-2.elb.amazonaws.com:8000/ws/blocknumber"

export default function () {
  	//var formData = `{"block": "latest","txdetails": "false"}`;
  	var headers = { 'Content-Type': 'application/json' };
    let res = http.get(url, { headers: headers });
    check(res, {
        "is status 200": (r) => r.status === 200
    });
  sleep(1);
}

