import "source-map-support/register";
import { config } from "./config";
import { asyncInterval } from "./utils";
import { poll } from "./poll";
import { listen } from "./server";
import { announce, pushDataPoint } from "./push";

(async () => {
  await announce();

  const [startPoll, stopPoll] = asyncInterval(async () => {
    const chunk = await poll(config.sourceEndpoint);
    if (!chunk) {
      return;
    }

    pushDataPoint(chunk);
  }, config.pollRate);

  process.on("exit", stopPoll);

  startPoll();
  listen();
})();
