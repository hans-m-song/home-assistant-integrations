import "dotenv/config";
import { log } from "./utils";

export const assertEnv = (key: string): string => {
  const value = process.env[key];
  if (!value) {
    throw new Error(`environment variable is not set: "${key}"`);
  }

  return value;
};

export const config = Object.freeze({
  sourceEndpoint: process.env.SOURCE_ENDPOINT ?? "http://localhost/home.cgi",
  destinationEndpoint: assertEnv("DESTINATION_ENDPOINT"),
  mqttNodeId: process.env.MQTT_TOPIC ?? "zeversolar",
  mqttUser: process.env.USERNAME,
  mqttPass: process.env.PASSWORD,
  serverPort: Number(process.env.SERVER_PORT) ?? 8000,
  serverHost: process.env.SERVER_HOST ?? "localhost",
  pullRate: Number(process.env.PULL_RATE) || 5000,
  pushRate: Number(process.env.PUSH_RATE) || 5000,
});

log("config", config);
