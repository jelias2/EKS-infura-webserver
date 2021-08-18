import http from 'k6/http';

let url = "SED-URL"
export default function () {
  http.get(url+='/blocknumber');
}

