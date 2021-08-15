// script.js
import http from "k6/http";
import { sleep } from "k6";

let url = "SED-URL"

export default function () {
  http.get(url);
  sleep(1);
}
