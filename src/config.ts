import "dotenv/config";
import { log } from "./utils";

const assertEnv = (key: string): string => {
  const value = process.env[key];
  if (!value) {
    throw new Error(`environment variable is not set: "${key}"`);
  }

  return value;
};

const numberEnv = (key: string, defaultValue: number) => {
  const value = Number(process.env[key]);
  if (isNaN(value) || value < 0 || value > 65536) {
    console.log(
      `number environment variable was invalid: "${key}",`,
      `using default: "${defaultValue}"`
    );
    return defaultValue;
  }

  return value;
};

export const config = Object.freeze({
  sourceEndpoint: process.env.SOURCE_ENDPOINT ?? "http://localhost/home.cgi",
  destinationEndpoint: assertEnv("DESTINATION_ENDPOINT"),
  mqttNodeId: process.env.MQTT_TOPIC ?? "zeversolar",
  mqttUser: process.env.MQTT_USER,
  mqttPass: process.env.MQTT_PASS,
  serverPort: numberEnv("SERVER_PORT", 8000),
  pullRate: numberEnv("PULL_RATE", 5000),
  pushRate: numberEnv("PUSH_RATE", 5000),
});

log("config", config);
