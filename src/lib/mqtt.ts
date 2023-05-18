import mqtt from "mqtt";
import { config } from "../config";
import { log, verbose } from "./utils";

const client = mqtt.connect(config.mqtt.endpoint, {
  username: config.mqtt.user,
  password: config.mqtt.pass,
});

process.on("beforeExit", async () => {
  if (client.connected) {
    log("mqtt.end", "attempting disconnection");
    await new Promise((resolve) => client.end(undefined, undefined, resolve));
  }
});

client.on("connect", (packet) => log("mqtt.client.connect", { packet }));
client.on("reconnect", () => log("mqtt.client.reconnect"));
client.on("close", () => log("mqtt.client.close"));
client.on("disconnect", () => log("mqtt.client.disconnect"));
client.on("offline", () => log("mqtt.client.offline"));
client.on("error", (error) => log("mqtt.client.error", error));
client.on("end", () => log("mqtt.client.end"));
client.on("message", (topic, _, { cmd }) => {
  log("mqtt.client.message", { topic, cmd });
});

const subscribe = async (
  topic: string,
  fn: (payload: string, packet: mqtt.IPublishPacket) => void
): Promise<boolean> => {
  const subscription = await new Promise((resolve) => {
    client.subscribe(topic, (error, granted) => {
      if (error) {
        log("mqtt.client.subscribe.failure", error);
        resolve(false);
        return;
      }

      const grantedTopics = granted.map(({ topic }) => topic);
      log("mqtt.client.subscribe", { grantedTopics });
      resolve(true);
    });
  });

  if (!subscription) {
    return false;
  }

  client.on("message", (incomingTopic, payload, packet) => {
    if (incomingTopic === topic) {
      fn(payload.toString(), packet);
    }
  });

  return true;
};

const push = async (topic: string, message: string | Buffer) => {
  verbose("mqtt.push", { topic, message });

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

export const MQTT = { client, push, subscribe };
