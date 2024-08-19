import { drizzle } from "drizzle-orm/planetscale-serverless";
import { Resource } from "sst";
import { Client } from "@planetscale/database";
export * from "drizzle-orm";

const client = new Client({
  host: Resource.Database.host,
  username: Resource.Database.username,
  password: Resource.Database.password,
});

export const db = drizzle(client, {
  logger:
    process.env.DRIZZLE_LOG === "true"
      ? {
          logQuery(query, params) {
            console.log("query", query);
            console.log("params", params);
          },
        }
      : undefined,
});
