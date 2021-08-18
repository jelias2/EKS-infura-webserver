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

let url = "SED-URL"

export default function () {
  var formData = `{"block": "latest","txdetails": "false"}`;
  var headers = { 'Content-Type': 'application/json' };
  http.post(url+='/blockbynumber', formData, { headers: headers });
  sleep(1);
}

