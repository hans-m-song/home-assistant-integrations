import "dotenv/config";
export const config = Object.freeze({
  sourceEndpoint: process.env.SOURCE_ENDPOINT ?? "localhost",
  destinationEndpoint: process.env.DESTINATION_ENDPOINT ?? "",
  pullRate: Number(process.env.PULL_RATE) || 5000,
  pushRate: Number(process.env.PUSH_RATE) || 5000,
});

console.log(config);
