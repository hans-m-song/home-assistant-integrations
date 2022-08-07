export type Device = {
  name: string;
  identifiers?: string;
  manufacturer: string;
  model: string;
  // TODO
  sw_version?: string;
  hw_version?: string;
};

export type EntityConfiguration = {
  name: string;
  unique_id: string;
  value_template?: string;
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
