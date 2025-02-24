// import grpc from "k6/net/grpc";
// import { check, sleep } from "k6";
// import { Counter } from "k6/metrics";

// export const requests = new Counter("grpc_reqs");

// export const options = {
//   stages: [
//     // { duration: "1m", target: 500 },
//     { duration: "1m", target: 1200 },
//     { duration: "1m", target: 3000 },
//   ],
//   thresholds: {
//     grpc_req_duration: ["p(95)<500"],
//   },
// };

// export default function () {
//   const client = new grpc.Client();
//   client.load(["./proto"], "stream.proto");

//   client.connect("localhost:50051", { insecure: true });

//   const resUserStats = client.invoke("stream.StreamStatsService/GetUserStats", {
//     userId: "181357564",
//   });

//   check(resUserStats, {
//     "OK": (r) => r.status === 0,
//     "response": (r) => r && r.body.stats.length > 0,
//   });

//   const resStreamStats = client.invoke("stream.StreamStatsService/GetStreamStats", {
//     streamId: "316005331705",
//   });

//   check(resStreamStats, {
//     "OK": (r) => r.status === 0,
//     "response": (r) => r && r.body.stats.length > 0,
//   });

//   client.close();
//   sleep(1);
// }
