import http from 'k6/http';
import { sleep } from 'k6';
export let options = {
  vus: 10,
  duration: '30s',
};

let url = "SED-URL"
export default function () {
  http.get(url+='/blocknumber');
}

