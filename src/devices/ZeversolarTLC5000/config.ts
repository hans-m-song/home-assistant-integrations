import path from "path";
import { DeviceInformation, EntityConfiguration } from "../../lib/hass";
import { midnight, slug } from "../../lib/utils";

export const topics = {
  SOLAR_STATE: path.join("homeassistant", "zeversolar", "solar_state"),

  CONFIG_LAST_UPDATED: path.join(
    "homeassistant",
    "sensor",
    "zeversolar",
    "solar_last_updated",
    "config"
  ),

  CONFIG_POWER_AC: path.join(
    "homeassistant",
    "sensor",
    "zeversolar",
    "solar_power_ac",
    "config"
  ),

  CONFIG_ENERGY_TODAY: path.join(
    "homeassistant",
    "sensor",
    "zeversolar",
    "solar_energy_today",
    "config"
  ),

  CONFIG_STATUS: path.join(
    "homeassistant",
    "binary_sensor",
    "zeversolar",
    "solar_status",
    "config"
  ),
};

const device: DeviceInformation = {
  identifiers: "zeversolar_inverter_TLC5000",
  name: "Solar Inverter",
  manufacturer: "Zeversolar",
  model: "TLC5000",
};

export const configuration: Record<string, EntityConfiguration> = {
  [topics.CONFIG_LAST_UPDATED]: {
    name: "Solar Last Updated",
    unique_id: slug("_", "zeversolar", "last_updated"),
    value_template: "{{ value_json.last_updated }}",
    state_topic: topics.SOLAR_STATE,
    device_class: "timestamp",
    device,
  },
  [topics.CONFIG_POWER_AC]: {
    name: "Solar Power AC (W)",
    unique_id: slug("_", "zeversolar", "power_ac"),
    value_template: "{{ value_json.power_ac }}",
    state_topic: topics.SOLAR_STATE,
    state_class: "measurement",
    device_class: "power",
    unit_of_measurement: "W",
    device,
  },
  [topics.CONFIG_ENERGY_TODAY]: {
    name: "Solar Energy Today (kWh)",
    unique_id: slug("_", "zeversolar", "energy_today"),
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
    unique_id: slug("_", "zeversolar", "status"),
    value_template: "{{ value_json.status }}",
    state_topic: topics.SOLAR_STATE,
    device_class: "power",
    device: device,
  },
};
