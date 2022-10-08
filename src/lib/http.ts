import express from "express";
import { log } from "./utils";

const client = express();

client.use("*", (req, res) => {
  const { headers, url, method, params, query, body } = req;
  log("server.receive", { headers, url, method, params, query, body });
  res.sendStatus(200);
});

const listen = (port: number) =>
  client.listen(port, () => {
    log("server.listen", "started server", { port: port });
  });

export const HTTP = { client, listen };
