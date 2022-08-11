import axios, { AxiosInstance } from "axios";
import { log } from "../../lib/utils";
import {
  DiagnoseInternetResponse,
  DeviceCountResponse,
  WizardWifiResponse,
  DeviceinfoResponse,
  WandetectResponse,
} from "./types";

export class HuaweiHG659API {
  instance: AxiosInstance;
  constructor(baseURL: string) {
    this.instance = axios.create({
      baseURL,
    });

    this.instance.interceptors.response.use((value) => {
      if (typeof value.data === "string") {
        try {
          value.data = JSON.parse(
            value.data.replace(/^while\(1\); \/\*/, "").replace(/\*\/$/, "")
          );
        } catch (error) {
          log("huawei.api.parse.failure", error);
        }
      }

      return value;
    });
  }

  async getDiagnoseInternet(): Promise<DiagnoseInternetResponse | null> {
    try {
      const response = await this.instance.get<DiagnoseInternetResponse>(
        "/api/system/diagnose_internet"
      );
      return response.data;
    } catch (error) {
      log("huawei.api.get_diagnose_internet.failure", error);
      return null;
    }
  }

  async getDeviceCount(): Promise<DeviceCountResponse | null> {
    try {
      const response = await this.instance.get<DeviceCountResponse>(
        "/api/system/device_count"
      );
      return response.data;
    } catch (error) {
      log("huawei.api.get_device_count.failure", error);
      return null;
    }
  }

  async getWizardWifi(): Promise<WizardWifiResponse | null> {
    try {
      const response = await this.instance.get<WizardWifiResponse>(
        "/api/system/wizard_wifi"
      );
      return response.data;
    } catch (error) {
      log("huawei.api.get_wizard_wifi.failure", error);
      return null;
    }
  }

  async getDeviceinfo(): Promise<DeviceinfoResponse | null> {
    try {
      const response = await this.instance.get<DeviceinfoResponse>(
        "/api/system/deviceinfo"
      );
      return response.data;
    } catch (error) {
      log("huawei.api.get_deviceinfo.failure", error);
      return null;
    }
  }

  async getWandetect(): Promise<WandetectResponse | null> {
    try {
      const response = await this.instance.get<WandetectResponse>(
        "/api/ntwk/wandetect"
      );
      return response.data;
    } catch (error) {
      log("huawei.api.get_wandetect.failure", error);
      return null;
    }
  }

  async all() {
    return {
      DiagnoseInternet: await this.getDiagnoseInternet(),
      DeviceCount: await this.getDeviceCount(),
      WizardWifi: await this.getWizardWifi(),
      Deviceinfo: await this.getDeviceinfo(),
      Wandetect: await this.getWandetect(),
    };
  }

  async summary() {
    const internet = await this.getDiagnoseInternet();
    const device = await this.getDeviceinfo();

    return {
      Internet: {
        ConnectionStatus: internet?.ConnectionStatus,
        ErrReason: internet?.ErrReason,
        Uptime: internet?.Uptime,
      },
      Device: {
        DeviceName: device?.DeviceName,
        SerialNumber: device?.SerialNumber,
        UpTime: device?.UpTime,
        SoftwareVersion: device?.SoftwareVersion,
        HardwareVersion: device?.HardwareVersion,
      },
    };
  }
}
