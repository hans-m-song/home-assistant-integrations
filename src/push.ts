import mqtt from "mqtt";
import { config } from "./config";
import { DeviceConfiguration, Topic } from "./hass";
import { DataPoint } from "./pull";
import { log, midnight } from "./utils";

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
    push(Topic.LastUpdatedState, JSON.stringify({ value: point.dateTime })),
    push(Topic.PowerAcState, JSON.stringify({ value: point.pacW })),
    push(Topic.EnergyTodayState, JSON.stringify({ value: point.eTodayKWh })),
    push(Topic.StatusState, JSON.stringify({ value: point.status })),
  ]);

const discoveryConfig: DeviceConfiguration[] = [
  {
    name: "last_updated",
    state_topic: Topic.LastUpdatedState,
    device_class: "timestamp",
  },
  {
    name: "power_ac",
    state_topic: Topic.PowerAcState,
    state_class: "measurement",
    device_class: "power",
    unit_of_measurement: "W",
  },
  {
    name: "energy_today",
    state_topic: Topic.EnergyTodayState,
    state_class: "total_increasing",
    last_reset: midnight(),
    device_class: "energy",
    unit_of_measurement: "kWh",
  },
  {
    name: "status",
    state_topic: Topic.StatusState,
    device_class: "power",
  },
];

export const announce = () =>
  Promise.all(
    discoveryConfig.map((config) =>
      push(config.state_topic, JSON.stringify(config))
    )
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
