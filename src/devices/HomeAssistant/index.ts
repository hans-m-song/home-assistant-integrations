import { Device, DeviceContext, DeviceManager } from "../../lib/device";
import { EntityConfiguration } from "../../lib/hass";
import { log } from "../../lib/utils";

export class HomeAssistant extends Device {
  configuration: Record<string, EntityConfiguration> = {};

  constructor(context: DeviceContext, manager: DeviceManager) {
    super();

    context.mqtt.subscribe("homeassistant/status", async (payload) => {
      log("homeassistant.status_message", { payload });

      if (payload === "online") {
        await manager.announce();
      }
    });
  }

  async shutdown(): Promise<void> {}
}
