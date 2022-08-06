import mqtt from "mqtt";
import { config } from "./config";
import { DeviceConfiguration, Device, Topic } from "./hass";
import { DataPoint } from "./poll";
import { log, midnight, slug } from "./utils";

const client = mqtt.connect(config.destinationEndpoint, {
  username: config.mqttUser,
  password: config.mqttPass,
});

process.on("beforeExit", async () => {
  if (client.connected) {
    log("mqtt.end", "attempting disconnection");
    await new Promise((resolve) => client.end(undefined, undefined, resolve));
  }
});

export const push = async (topic: string, message: string | Buffer) => {
  log("mqtt.push", { topic, message });

  return new Promise<mqtt.Packet | undefined>((resolve) => {
    client.publish(topic, message, (error, packet) => {
      if (error) {
        log("mqtt.push.error", error);
        resolve(undefined);
      }

      log("mqtt.push.success", { packet });
      resolve(packet);
    });
  });
};

export const pushDataPoint = (point: DataPoint) =>
  Promise.all([
    push(Topic.LastUpdatedState, point.dateTime),
    push(Topic.PowerAcState, point.pacW),
    push(Topic.EnergyTodayState, point.eTodayKWh),
    push(Topic.StatusState, point.status === "OK" ? "ON" : "OFF"),
  ]);

const device: Device = {
  identifiers: "zeversolar_inverter_TLC5000",
  name: "Zeversolar Inverter",
  manufacturer: "Zeversolar",
  model: "TLC5000",
};

const discoveryConfig: Record<string, DeviceConfiguration> = {
  [Topic.LastUpdatedConfig]: {
    name: "Solar Last Updated",
    unique_id: slug("_", config.mqttNodeId, "last_updated"),
    state_topic: Topic.LastUpdatedState,
    device_class: "timestamp",
    device,
  },
  [Topic.PowerAcConfig]: {
    name: "Solar Power AC (W)",
    unique_id: slug("_", config.mqttNodeId, "power_ac"),
    state_topic: Topic.PowerAcState,
    state_class: "measurement",
    device_class: "power",
    unit_of_measurement: "W",
    device,
  },
  [Topic.EnergyTodayConfig]: {
    name: "Solar Energy Today (kWh)",
    unique_id: slug("_", config.mqttNodeId, "energy_today"),
    state_topic: Topic.EnergyTodayState,
    state_class: "total_increasing",
    last_reset: midnight(),
    device_class: "energy",
    unit_of_measurement: "kWh",
    device,
  },
  [Topic.StatusConfig]: {
    name: "Solar Status",
    unique_id: slug("_", config.mqttNodeId, "status"),
    state_topic: Topic.StatusState,
    device_class: "power",
    device: device,
  },
};

export const announce = () =>
  Promise.all(
    Object.entries(discoveryConfig).map(([topic, config]) =>
      push(topic, JSON.stringify(config))
    )
  );

export const denounce = () =>
  Promise.all(
    Object.entries(discoveryConfig).map(([topic]) => push(topic, ""))
  );

client.on("connect", (packet) => log("mqtt.client.connect", { packet }));
client.on("reconnect", () => log("mqtt.client.reconnect"));
client.on("close", () => log("mqtt.client.close"));
client.on("disconnect", () => log("mqtt.client.disconnect"));
client.on("offline", () => log("mqtt.client.offline"));
client.on("error", (error) => log("mqtt.client.error", error));
client.on("end", () => log("mqtt.client.end"));
client.on("message", (topic, payload, _packet) => {
  switch (topic) {
    case Topic.HassStatus:
      if (payload.toString() === "online") {
        log("mqtt.client.message", "sending configuration", { topic });
        announce();
      }
      break;

    default:
      log("mqtt.client.message", "unhandled topic received", { topic });
      break;
  }
});

client.subscribe(Topic.HassStatus, (error, granted) => {
  if (error) {
    log("mqtt.client.subscribe.failure", error);
    return;
  }

  log("mqtt.client.subscribe", { topics: granted.map(({ topic }) => topic) });
});
