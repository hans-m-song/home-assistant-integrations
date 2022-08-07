import axios from "axios";
import { log } from "./utils";

const normaliseDecimal = (input: string) => {
  const [integer, fractional] = input.split(".");
  return `${integer}.${fractional.padStart(2, "0")}`;
};

const parseTimestamp = (input: string): string | null => {
  const match = input.match(
    /(?<hour>\d{2}):(?<minute>\d{2}) (?<day>\d{2})\/(?<month>\d{2})\/(?<year>\d{4})/
  );

  const { groups } = match ?? {};

  if (
    !groups?.hour ||
    !groups?.minute ||
    !groups?.day ||
    !groups?.month ||
    !groups?.year
  ) {
    log("poll.warn", "could not parse timestamp", { input, groups });
    return null;
  }

  const minute = Number(groups.minute) - 1;
  const hour = Number(groups.hour) - 1;
  const day = Number(groups.day) - 1;
  const month = Number(groups.month) - 1;
  const year = Number(groups.year);

  return new Date(year, month, day, hour, minute).toISOString();
};

export type DataPoint = ReturnType<typeof parse>;
const parse = (raw: string) => {
  const fields = raw.trim().split(/\r?\n/g);
  if (fields.length !== 14) {
    log("poll.warn", "data was in an unexpected format", { raw });
  }

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
    powerAc,
    energyToday,
    status,
  ] = fields;

  return {
    registryID,
    registryKey,
    hardwareVersion,
    softwareVersion,
    dateTime: parseTimestamp(dateTime),
    zeverCloudStatus,
    serialNumber,
    powerAc,
    energyToday: normaliseDecimal(energyToday),
    status,
    fields,
  };
};

export const poll = async (endpoint: string): Promise<DataPoint | null> => {
  try {
    const response = await axios.get(endpoint, { timeout: 5000 });
    const { status, data: raw } = response;
    const data = parse(raw);
    log("poll.success", { status, data });
    return data;
  } catch (error) {
    log("poll.error", error);
    return null;
  }
};
