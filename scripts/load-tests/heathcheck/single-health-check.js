import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
  insecureSkipTLSVerify: true,
};

export default function () {
  http.get('http://localhost:8000/health');
  sleep(1);
}
