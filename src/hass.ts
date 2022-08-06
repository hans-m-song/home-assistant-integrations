import path from "path";
import { config } from "./config";

/**
 * @param entityType https://www.home-assistant.io/docs/configuration/customizing-devices/#device-class
 */
export const sensorTopic = (entityType: string, topic: string) =>
  path.join("home-assistant", entityType, config.mqttNodeId, topic);

export const Topic = Object.freeze({
  HassStatus: "home-assistant/status",

  LastUpdatedConfig: sensorTopic("sensor", "last_updated/config"),
  LastUpdatedState: sensorTopic("sensor", "last_updated/state"),

  PowerAcConfig: sensorTopic("sensor", "power_ac/config"),
  PowerAcState: sensorTopic("sensor", "power_ac/state"),

  EnergyTodayConfig: sensorTopic("sensor", "energy_today/config"),
  EnergyTodayState: sensorTopic("sensor", "energy_today/state"),

  StatusConfig: sensorTopic("binary_sensor", "status/config"),
  StatusState: sensorTopic("binary_sensor", "status/state"),
});

export type DeviceConfiguration = {
  name: string;
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
};
