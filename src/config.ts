import "dotenv/config";
import path from "path";
import { Device, EntityConfiguration } from "./hass";
import { clamp, log, midnight, slug } from "./utils";

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

export const config = Object.freeze({
  sourceEndpoint: process.env.SOURCE_ENDPOINT ?? "http://localhost/home.cgi",
  destinationEndpoint: assertEnv("DESTINATION_ENDPOINT"),
  mqttNodeId: process.env.MQTT_TOPIC ?? "zeversolar",
  mqttUser: process.env.MQTT_USER,
  mqttPass: process.env.MQTT_PASS,
  serverPort: clamp(numberEnv("SERVER_PORT", 8000), 65536),
  pollRate: clamp(numberEnv("POLL_RATE", 5000), Infinity, 1000),
});

log("config", config);

const device: Device = {
  identifiers: "zeversolar_inverter_TLC5000",
  name: "Solar Inverter",
  manufacturer: "Zeversolar",
  model: "TLC5000",
};

log("device", device);

/**
 * @param entityType https://www.home-assistant.io/docs/configuration/customizing-devices/#device-class
 */
export const topics = Object.freeze({
  HASS_STATUS: "homeassistant/status",

  SOLAR_STATE: path.join("homeassistant", config.mqttNodeId, "solar_state"),

  CONFIG_LAST_UPDATED: path.join(
    "homeassistant",
    "sensor",
    config.mqttNodeId,
    `solar_last_updated`,
    "config"
  ),

  CONFIG_POWER_AC: path.join(
    "homeassistant",
    "sensor",
    config.mqttNodeId,
    `solar_power_ac`,
    "config"
  ),

  CONFIG_ENERGY_TODAY: path.join(
    "homeassistant",
    "sensor",
    config.mqttNodeId,
    `solar_energy_today`,
    "config"
  ),

  CONFIG_STATUS: path.join(
    "homeassistant",
    "binary_sensor",
    config.mqttNodeId,
    `solar_status`,
    "config"
  ),
});

export const entityConfiguration: Record<string, EntityConfiguration> = {
  [topics.CONFIG_LAST_UPDATED]: {
    name: "Solar Last Updated",
    unique_id: slug("_", config.mqttNodeId, "last_updated"),
    value_template: "{{ value_json.last_updated }}",
    state_topic: topics.SOLAR_STATE,
    device_class: "timestamp",
    device,
  },
  [topics.CONFIG_POWER_AC]: {
    name: "Solar Power AC (W)",
    unique_id: slug("_", config.mqttNodeId, "power_ac"),
    value_template: "{{ value_json.power_ac }}",
    state_topic: topics.SOLAR_STATE,
    state_class: "measurement",
    device_class: "power",
    unit_of_measurement: "W",
    device,
  },
  [topics.CONFIG_ENERGY_TODAY]: {
    name: "Solar Energy Today (kWh)",
    unique_id: slug("_", config.mqttNodeId, "energy_today"),
    value_template: "{{ value_json.energy_today }}",
    state_topic: topics.SOLAR_STATE,
    state_class: "total_increasing",
    last_reset: midnight(),
    device_class: "energy",
    unit_of_measurement: "kWh",
    device,
  },
  [topics.CONFIG_STATUS]: {
    name: "Solar Status",
    unique_id: slug("_", config.mqttNodeId, "status"),
    value_template: "{{ value_json.status }}",
    state_topic: topics.SOLAR_STATE,
    device_class: "power",
    device: device,
  },
};

log("topics", topics);
