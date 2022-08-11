import { config } from "../config";
import { HuaweiHG659API } from "../devices/HuaweiHG659/api";
import { selfTest } from "../devices/HuaweiHG659/selfTest";

selfTest().then((result) => console.log("self test result", result));

new HuaweiHG659API(config.huawei.hg659.endpoint).all().then(console.log);
