import mqtt from "mqtt";
import { config, entityConfiguration, topics } from "./config";
import { DataPoint } from "./poll";
import { log } from "./utils";

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

      resolve(packet);
    });
  });
};

export const pushDataPoint = (point: DataPoint) =>
  push(
    topics.SOLAR_STATE,
    JSON.stringify({
      last_updated: point.dateTime,
      power_ac: point.powerAc,
      energy_today: point.energyToday,
      status: point.status === "OK" ? "ON" : "OFF",
    })
  );

export const announce = () =>
  Promise.all(
    Object.entries(entityConfiguration).map(([topic, config]) =>
      push(topic, JSON.stringify(config))
    )
  );

export const denounce = () =>
  Promise.all(
    Object.entries(entityConfiguration).map(([topic]) => push(topic, ""))
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
    case topics.HASS_STATUS:
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

client.subscribe(topics.HASS_STATUS, (error, granted) => {
  if (error) {
    log("mqtt.client.subscribe.failure", error);
    return;
  }

  log("mqtt.client.subscribe", { topics: granted.map(({ topic }) => topic) });
});
