import http from "k6/http";
import { check, sleep } from "k6";
import { Counter } from "k6/metrics";

export const requests = new Counter("http_reqs");

export const options = {
  stages: [
    // { duration: "1m", target: 500 },
    { duration: "1m", target: 1200 },
    { duration: "1m", target: 3000 },
    // { duration: "1m", target: 3000 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],
  },
};

// export default function () {
//   const res = http.get("http://localhost:8080/stream/stats?stream_id=316005331705");
//   check(res, { "status was 200": (r) => r.status === 200 });
//   sleep(1);
// }

export default function () {
  const urls = [
    "http://localhost:8080/stream/stats?stream_id=317402421241",
    "http://localhost:8080/user/stats?user_id=28354765",
  ];

  urls.forEach((url) => {
    const res = http.get(url);
    check(res, { "status was 200": (r) => r.status === 200 });
  });

  sleep(1);
}
