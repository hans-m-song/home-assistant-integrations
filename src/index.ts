import "source-map-support/register";
import { config } from "./config";
import { asyncInterval } from "./utils";
import { pull } from "./pull";
import { listen } from "./server";

(async () => {
  let data: any[] = [];

  listen();

  const [startPull, stopPull] = asyncInterval(async () => {
    const chunk = await pull(config.sourceEndpoint);
    if (!chunk) {
      return;
    }
    data.push(chunk);
  }, config.pullRate);

  // const [startPush, stopPush] = asyncInterval(async () => {
  //   if (data.length < 1) {
  //     return;
  //   }
  //   push(client, "", JSON.stringify(data));
  //   // reset data
  //   data = [];
  // }, config.pushRate);

  startPull();
  // startPush();

  process.on("exit", () => {
    console.log("stopping runtime");
    stopPull();
    // stopPush();
  });
})();
