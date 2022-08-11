import { config } from "../../config";
import { DeviceContext, Device } from "../../lib/device";
import { EntityConfiguration } from "../../lib/hass";
import { asyncInterval } from "../../lib/utils";
import { configuration, topics } from "./config";
import { pull } from "./pull";

export class ZeversolarTLC5000 extends Device {
  configuration: Record<string, EntityConfiguration> = configuration;
  private poll: Promise<void>;
  private stopPoll: () => void;

  constructor(context: DeviceContext) {
    super();

    const [startPoll, stopPoll] = asyncInterval(async () => {
      const point = await pull(config.zeversolar.tlc5000.endpoint);
      if (!point) {
        return;
      }

      await context.mqtt.push(
        topics.SOLAR_STATE,
        JSON.stringify({
          last_updated: point.dateTime,
          power_ac: point.powerAc,
          energy_today: point.energyToday,
          status: point.status === "OK" ? "ON" : "OFF",
        })
      );
    }, config.zeversolar.tlc5000.pollRate);

    this.poll = startPoll();
    this.stopPoll = stopPoll;
    process.on("exit", stopPoll);
  }

  async shutdown(): Promise<void> {
    this.stopPoll();
    await this.poll;
  }
}
