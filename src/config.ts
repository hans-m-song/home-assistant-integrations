import "dotenv/config";
import { clamp, log } from "./lib/utils";

const assertEnv = (key: string): string => {
  const value = process.env[key];
  if (!value) {
    throw new Error(`environment variable is not set: "${key}"`);
  }

  return value;
};

const numberEnv = (key: string, defaultValue: number) => {
  const value = Number(process.env[key]);
  if (isNaN(value)) {
    return defaultValue;
  }

  return value;
};

const clampPollValue = (value: number) => clamp(value, Infinity, 1000);

export const config = Object.freeze({
  debug: process.env.DEBUG === "true",
  verbose: process.env.VERBOSE === "true",

  zeversolar: {
    tlc5000: {
      endpoint: process.env.ZEVERSOLAR_TLC5000_ENDPOINT ?? "",
      pollRate: clampPollValue(numberEnv("ZEVERSOLAR_POLL_RATE", 5000)),
    },
  },

  huawei: {
    hg659: {
      endpoint: process.env.HUAWEI_HG659_ENDPOINT ?? "",
      pollRate: clampPollValue(numberEnv("HUAWEI_HG659_POLL_RATE", 5000)),
    },
  },

  mqtt: {
    endpoint: assertEnv("MQTT_ENDPOINT"),
    user: process.env.MQTT_USER,
    pass: process.env.MQTT_PASS,
    announceRate: clampPollValue(numberEnv("MQTT_ANNOUNCE_RATE", 60000)),
  },

  http: {
    port: clamp(numberEnv("HTTP_PORT", 8000), 65536),
  },
});

log("config", config);
