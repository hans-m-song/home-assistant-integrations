import axios from "axios";
import { config } from "./config";
import { unpackError } from "./utils";

export const pull = async () => {
  const response = await axios
    .get(config.sourceEndpoint, { timeout: 5000 })
    .catch((error) => {
      console.error(unpackError(error));
      return null;
    });

  return response?.data;
};
