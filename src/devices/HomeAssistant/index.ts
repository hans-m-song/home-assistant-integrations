import { config } from "../../config";
import { Device, DeviceContext, DeviceManager } from "../../lib/device";
import { EntityConfiguration } from "../../lib/hass";
import { asyncInterval, log } from "../../lib/utils";

export class HomeAssistant extends Device {
  configuration: Record<string, EntityConfiguration> = {};
  private poll: Promise<void>;
  private stopPoll: () => void;

  constructor(context: DeviceContext, manager: DeviceManager) {
    super();

    context.mqtt.subscribe("homeassistant/status", async (payload) => {
      log("homeassistant.status_message", { payload });

      if (payload === "online") {
        await manager.announce();
      }
    });

    const [startPoll, stopPoll] = asyncInterval(
      () => manager.announce(),
      config.mqtt.announceRate
    );

    this.poll = startPoll();
    this.stopPoll = stopPoll;
    process.on("exit", stopPoll);
  }

  async shutdown(): Promise<void> {
    this.stopPoll();
    await this.poll;
  }
}
