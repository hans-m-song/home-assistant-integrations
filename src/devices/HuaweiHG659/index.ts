import { config } from "../../config";
import { Device, DeviceContext } from "../../lib/device";
import { EntityConfiguration } from "../../lib/hass";
import { asyncInterval } from "../../lib/utils";
import { HuaweiHG659API } from "./api";
import { configuration, topics } from "./config";
import { selfTest } from "./selfTest";

export class HuaweiHG659 extends Device {
  readonly configuration: Record<string, EntityConfiguration> = configuration;
  private api: HuaweiHG659API;
  private poll: Promise<void>;
  private stopPoll: () => void;

  constructor(context: DeviceContext) {
    super();

    this.api = new HuaweiHG659API(config.huawei.hg659.endpoint);

    const [startPoll, stopPoll] = asyncInterval(async () => {
      const [{ Internet, Device }, { success, message }] = await Promise.all([
        this.api.summary(),
        selfTest(),
      ]);

      const connected =
        (Internet.Uptime ?? 0) > 0 &&
        Internet.ConnectionStatus === "Connected" &&
        success;

      await context.mqtt.push(
        topics.ROUTER_STATE,
        JSON.stringify({
          internet_connected: connected ? "ON" : "OFF",
          self_test_message: message ?? "None",
          internet_connection_status: Internet.ConnectionStatus,
          internet_err_reason: Internet.ErrReason,
          internet_uptime: Internet.Uptime,
          device_uptime: Device.UpTime,
        })
      );
    }, config.huawei.hg659.pollRate);

    this.poll = startPoll();
    this.stopPoll = stopPoll;
    process.on("exit", stopPoll);
  }

  async shutdown(): Promise<void> {
    this.stopPoll();
    await this.poll;
  }
}
