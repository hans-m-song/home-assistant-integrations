import "source-map-support/register";
import { config } from "./config";
import { HTTP } from "./lib/http";
import { MQTT } from "./lib/mqtt";
import { DeviceContext, DeviceManager } from "./lib/device";
import { ZeversolarTLC5000 } from "./devices/ZeversolarTLC5000";
import { HuaweiHG659 } from "./devices/HuaweiHG659";
import { HomeAssistant } from "./devices/HomeAssistant";

(async () => {
  const context: DeviceContext = { mqtt: MQTT, http: HTTP };
  context.http.listen(config.http.port);

  const manager = new DeviceManager(context);
  manager.add(new HomeAssistant(context, manager));

  if (config.huawei.hg659.endpoint) {
    manager.add(new HuaweiHG659(context));
  }

  if (config.zeversolar.tlc5000.endpoint) {
    manager.add(new ZeversolarTLC5000(context));
  }

  await manager.announce();
  process.on("exit", manager.shutdown);
})();
