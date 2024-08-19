import { bus } from "./bus";
import { cluster } from "./cluster";
import { database } from "./database";
import { allSecrets } from "./secret";

const bucket = new sst.aws.Bucket("IntervalBucket");
cluster.addService("Interval", {
  link: [...allSecrets, database, bucket, bus],
  image: {
    dockerfile: "packages/interval/Dockerfile",
  },
  environment: {
    DRIZZLE_LOG: "true",
  },
  dev: {
    command: "bun dev",
  },
});
