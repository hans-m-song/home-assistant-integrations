import path from "path";
import { config } from "./config";

/**
 * @param entityType https://www.home-assistant.io/docs/configuration/customizing-devices/#device-class
 */
export const sensorTopic = (entityType: string, topic: string) =>
  path.join("homeassistant", entityType, `${config.mqttNodeId}_${topic}`);

export const Topic = Object.freeze({
  HassStatus: "homeassistant/status",

  LastUpdatedConfig: sensorTopic("sensor", "last_updated/config"),
  LastUpdatedState: sensorTopic("sensor", "last_updated/state"),

  PowerAcConfig: sensorTopic("sensor", "power_ac/config"),
  PowerAcState: sensorTopic("sensor", "power_ac/state"),

  EnergyTodayConfig: sensorTopic("sensor", "energy_today/config"),
  EnergyTodayState: sensorTopic("sensor", "energy_today/state"),

  StatusConfig: sensorTopic("binary_sensor", "status/config"),
  StatusState: sensorTopic("binary_sensor", "status/state"),
});

export type Device = {
  name: string;
  identifiers?: string;
  manufacturer: string;
  model: string;
  // TODO
  sw_version?: string;
  hw_version?: string;
};

export type DeviceConfiguration = {
  name: string;
  unique_id: string;
  state_topic: string;
  /**
   * https://developers.home-assistant.io/docs/core/entity/sensor#available-state-classes
   */
  state_class?: string;
  last_reset?: string;
  /**
   * https://developers.home-assistant.io/docs/core/entity/sensor#available-device-classes
   */
  device_class: string;
  unit_of_measurement?: string;
  device?: Device;
};
