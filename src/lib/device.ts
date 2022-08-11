import { EntityConfiguration } from "./hass";
import { HTTP } from "./http";
import { MQTT } from "./mqtt";

export interface DeviceContext {
  mqtt: typeof MQTT;
  http: typeof HTTP;
}

export abstract class Device {
  abstract readonly configuration: Record<string, EntityConfiguration>;

  async announce(context: DeviceContext): Promise<void> {
    await Promise.all(
      Object.entries(this.configuration).map(([topic, entity]) =>
        context.mqtt.push(topic, JSON.stringify(entity))
      )
    );
  }

  async denounce(context: DeviceContext): Promise<void> {
    await Promise.all(
      Object.keys(this.configuration).map((topic) =>
        context.mqtt.push(topic, "")
      )
    );
  }

  abstract shutdown(context: DeviceContext): Promise<void>;
}

export class DeviceManager {
  readonly context: DeviceContext;
  readonly devices: Device[];

  constructor(context: DeviceContext) {
    this.context = context;
    this.devices = [];
  }

  add(device: Device) {
    this.devices.push(device);
    return this;
  }

  all<T = void>(
    fn: (value: Device, index: number, array: Device[]) => Promise<T>
  ) {
    return Promise.all(this.devices.map(fn));
  }

  async announce(): Promise<void> {
    await this.all((device) => device.announce(this.context));
  }

  async denounce(): Promise<void> {
    await this.all((device) => device.denounce(this.context));
  }

  async shutdown(): Promise<void> {
    await this.all((device) => device.shutdown(this.context));
  }
}
