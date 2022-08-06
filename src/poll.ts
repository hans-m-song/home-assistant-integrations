import axios from "axios";
import { log } from "./utils";

const normaliseDecimal = (input: string) => {
  const [integer, fractional] = input.split(".");
  return `${integer}.${fractional.padStart(2, "0")}`;
};

export type DataPoint = ReturnType<typeof parse>;
const parse = (raw: string) => {
  const fields = raw.trim().split(/\r?\n/g);
  const [
    ,
    ,
    registryID,
    registryKey,
    hardwareVersion,
    softwareVersion,
    dateTime,
    zeverCloudStatus,
    ,
    serialNumber,
    pacW,
    eTodayKWh,
    status,
  ] = fields;

  return {
    registryID,
    registryKey,
    hardwareVersion,
    softwareVersion,
    dateTime,
    zeverCloudStatus,
    serialNumber,
    pacW,
    eTodayKWh: normaliseDecimal(eTodayKWh),
    status,
    fields,
  };
};

export const poll = async (endpoint: string): Promise<DataPoint | null> => {
  try {
    const response = await axios.get(endpoint, { timeout: 5000 });
    const { status, statusText, data: raw } = response;
    const data = parse(raw);
    log("pull.success", { status, statusText, data });
    return data;
  } catch (error) {
    log("pull.error", error);
    return null;
  }
};
