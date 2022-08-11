import axios from "axios";
import { config } from "../../config";
import { log } from "../../lib/utils";

export const selfTest = async (): Promise<{
  success: boolean;
  message?: string;
}> => {
  try {
    await axios.get("https://tpg.com.au", {
      timeout: Math.max(4000, config.huawei.hg659.pollRate - 1000),
    });
    return { success: true };
  } catch (error) {
    if ((error as any).code !== "ECONNABORTED") {
      log("huawei.test_connection.failure", error);
      return { success: true };
    }

    return { success: false, message: (error as any).message };
  }
};
