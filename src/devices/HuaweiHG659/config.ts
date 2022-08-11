import path from "path";
import { DeviceInformation, EntityConfiguration } from "../../lib/hass";
import { slug } from "../../lib/utils";

export const topics = {
  ROUTER_STATE: path.join("homeassistant", "huawei", "router_state"),

  CONFIG_INTERNET_CONNECTED: path.join(
    "homeassistant",
    "binary_sensor",
    "huawei",
    "router_internet_connected",
    "config"
  ),

  CONFIG_INTERNET_SELF_TEST_MESSAGE: path.join(
    "homeassistant",
    "sensor",
    "huawei",
    "router_internet_self_test_message",
    "config"
  ),

  CONFIG_INTERNET_CONNECTION_STATUS: path.join(
    "homeassistant",
    "sensor",
    "huawei",
    "router_internet_connection_status",
    "config"
  ),

  CONFIG_INTERNET_ERR_REASON: path.join(
    "homeassistant",
    "sensor",
    "huawei",
    "router_internet_err_reason",
    "config"
  ),

  CONFIG_INTERNET_UPTIME: path.join(
    "homeassistant",
    "sensor",
    "huawei",
    "router_internet_uptime",
    "config"
  ),

  CONFIG_DEVICE_UPTIME: path.join(
    "homeassistant",
    "sensor",
    "huawei",
    "router_device_uptime",
    "config"
  ),
};

const device: DeviceInformation = {
  identifiers: "huawei_router_hg659",
  name: "Huawei HG659",
  manufacturer: "Huawei",
  model: "HG659",
};

export const configuration: Record<string, EntityConfiguration> = {
  [topics.CONFIG_INTERNET_CONNECTED]: {
    name: "Internet Connected",
    unique_id: slug("_", "huawei", "internet_connected"),
    value_template: "{{ value_json.internet_connected }}",
    state_topic: topics.ROUTER_STATE,
    device_class: "power",
    device,
  },
  [topics.CONFIG_INTERNET_SELF_TEST_MESSAGE]: {
    name: "Self-test Message",
    unique_id: slug("_", "huawei", "self_test_message"),
    value_template: "{{ value_json.self_test_message }}",
    state_topic: topics.ROUTER_STATE,
    device,
  },
  [topics.CONFIG_INTERNET_CONNECTION_STATUS]: {
    name: "Internet Connection Status",
    unique_id: slug("_", "huawei", "internet_connection_status"),
    value_template: "{{ value_json.internet_connection_status }}",
    state_topic: topics.ROUTER_STATE,
    device,
  },
  [topics.CONFIG_INTERNET_ERR_REASON]: {
    name: "Internet Err Reason",
    unique_id: slug("_", "huawei", "internet_err_reason"),
    value_template: "{{ value_json.internet_err_reason }}",
    state_topic: topics.ROUTER_STATE,
    device,
  },
  [topics.CONFIG_INTERNET_UPTIME]: {
    name: "Internet Uptime",
    unique_id: slug("_", "huawei", "internet_uptime"),
    value_template: "{{ value_json.internet_uptime }}",
    state_topic: topics.ROUTER_STATE,
    state_class: "total_increasing",
    device_class: "duration",
    unit_of_measurement: "ms",
    device,
  },
  [topics.CONFIG_DEVICE_UPTIME]: {
    name: "Device Uptime",
    unique_id: slug("_", "huawei", "device_uptime"),
    value_template: "{{ value_json.device_uptime }}",
    state_topic: topics.ROUTER_STATE,
    state_class: "total_increasing",
    device_class: "duration",
    unit_of_measurement: "ms",
    device,
  },
};
