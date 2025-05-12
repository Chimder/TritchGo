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
    "http://localhost:8080/stream/stats?stream_id=318071585529",
    "http://localhost:8080/stream/stats?stream_id=319893400185",
    "http://localhost:8080/stream/stats?stream_id=319885216121",
    "http://localhost:8080/stream/stats?stream_id=323197471869",
    // "http://localhost:8080/stream/stats?stream_id=323274256893",
    "http://localhost:8080/user/stats?user_id=23161357",
    "http://localhost:8080/user/stats?user_id=545050196",
    "http://localhost:8080/user/stats?user_id=97828400",
    "http://localhost:8080/user/stats?user_id=906801651",
    "http://localhost:8080/user/stats?user_id=918839982",
    "http://localhost:8080/user/stats?user_id=29546745",
  ];

  urls.forEach((url) => {
    const res = http.get(url);
    check(res, { "status was 200": (r) => r.status === 200 });
  });

  sleep(1);
}
