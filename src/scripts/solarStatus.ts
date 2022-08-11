import { config } from "../config";
import { pull } from "../devices/ZeversolarTLC5000/pull";

pull(config.zeversolar.tlc5000.endpoint).then(console.log);
