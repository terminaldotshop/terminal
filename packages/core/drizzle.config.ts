import { defineConfig } from "drizzle-kit";
import { Resource } from "sst";

export default defineConfig({
  strict: true,
  verbose: true,
  out: "./migrations",
  dialect: "mysql",
  dbCredentials: {
    url: `mysql://${Resource.Database.username}:${Resource.Database.password}@${Resource.Database.host}/${Resource.Database.database}`,
  },
  schema: "./src/**/*.sql.ts",
});
