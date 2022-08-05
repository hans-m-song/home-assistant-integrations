import "source-map-support/register";
import { config } from "./config";
import { asyncInterval } from "./utils";
import { pull } from "./pull";

let data = [];

const [startPull, stopPull] = asyncInterval(async () => {
  const chunk = await pull();
  if (!chunk) {
    return;
  }

  console.log("received data:", chunk);
  data.push(chunk);
}, config.pullRate);

const [startPush, stopPush] = asyncInterval(async () => {
  if (data.length < 1) {
    return;
  }

  console.log("pushing chunks", { count: data.length });

  // reset data
  data = [];
}, config.pushRate);

startPull();
startPush();

process.on("exit", () => {
  console.log("stopping runtime");
  stopPull();
  stopPush();
});
