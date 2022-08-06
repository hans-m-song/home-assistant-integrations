import express from "express";
import { config } from "./config";
import { log } from "./utils";

const app = express();

app.use("*", (req, res) => {
  const { headers, url, method, baseUrl, originalUrl, params, query, body } =
    req;
  log("server.receive", {
    headers,
    url,
    method,
    params,
    query,
    body,
  });
  res.sendStatus(200);
});

export const listen = () =>
  app.listen(config.serverPort, config.serverHost, () => {
    log("server.listen", "started server", {
      host: config.serverHost,
      port: config.serverPort,
    });
  });
