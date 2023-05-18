import axios from "axios";
import { inspect } from "util";
import { config } from "../config";

export const sleep = (timeout: number) =>
  new Promise((resolve) => setTimeout(resolve, timeout));

export const asyncInterval = (fn: () => Promise<void>, minDelay: number) => {
  let next = true;
  const stop = () => {
    next = false;
  };

  const iterate = async () => {
    const now = Date.now();
    await fn();
    const duration = Date.now() - now;
    const timeout = Math.max(minDelay - duration, 0);

    if (timeout > 0) {
      await sleep(timeout);
    }

    if (next) {
      setImmediate(iterate);
    }
  };

  return [iterate, stop] as const;
};

export const unpackError = (error: any) => {
  const { code, name, message } = error ?? {};
  const req = axios.isAxiosError(error) && {
    status: error.response?.status,
    statusText: error.response?.statusText,
    data: error.response?.data,
  };
  return { ...req, code, name, message };
};

export const verbose = (...parameters: Parameters<typeof log>) => {
  if (!config.verbose) {
    return;
  }

  log(...parameters);
};

export const log = (namespace: string, data?: unknown, metadata?: unknown) => {
  const unpacked = data instanceof Error ? unpackError(data) : data;

  if (config.debug) {
    const fields = [`[${namespace}]`];
    unpacked &&
      fields.push(inspect(unpacked, { breakLength: Infinity, compact: true }));
    metadata &&
      fields.push(inspect(metadata, { breakLength: Infinity, compact: true }));

    console.log(...fields);
  } else {
    console.log(JSON.stringify({ namespace, data, metadata }));
  }
};

export const midnight = () => {
  const date = new Date();
  date.setHours(0);
  date.setMinutes(0);
  date.setSeconds(0);
  date.setMilliseconds(0);
  return date.toISOString();
};

export const slug = (seperator: string, ...input: string[]) =>
  input.map((item) => item.replace(/[^a-zA-Z0-9]/g, "")).join(seperator);

export const clamp = (value: number, max: number, min = 0) =>
  Math.min(max, Math.max(min, value));
