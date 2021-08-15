import http from 'k6/http';
import { sleep } from 'k6';
export let options = {
  vus: 20,
  duration: '30s',
};
export default function () {
  http.get('http://host.docker.internal:8000/health');
  sleep(1);
}

