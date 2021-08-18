import http from 'k6/http';
import { sleep } from 'k6';

let url = "SED-URL"

export default function () {
  var formData = `{"block": "latest","txdetails": "false"}`;
  var headers = { 'Content-Type': 'application/json' };
  http.post(url+='/blockbynumber', formData, { headers: headers });
  sleep(1);
}

