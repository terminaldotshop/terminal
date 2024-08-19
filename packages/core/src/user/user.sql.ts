import {
  json,
  mysqlTable,
  primaryKey,
  text,
  varchar,
} from "drizzle-orm/mysql-core";
import { address, id, ulid, timestamps } from "../drizzle/types";
import { z } from "zod";

export const UserFlags = z.object({
  printer: z.boolean().optional(),
});
export type UserFlags = z.infer<typeof UserFlags>;

export const userTable = mysqlTable("user", {
  ...id,
  ...timestamps,
  name: varchar("name", { length: 255 }),
  email: varchar("email", { length: 255 }),
  fingerprint: varchar("fingerprint", { length: 255 }).unique(),
  stripeCustomerID: varchar("stripe_customer_id", { length: 255 })
    .unique()
    .notNull(),
  emailOctopusID: text("email_octopus_id"),
  flags: json("flags").$type<UserFlags>().default({}),
});

export const userFingerprintTable = mysqlTable(
  "user_fingerprint",
  {
    userID: ulid("user_id")
      .references(() => userTable.id, {
        onDelete: "cascade",
      })
      .notNull(),
    fingerprint: varchar("fingerprint", { length: 255 }).notNull(),
    ...timestamps,
  },
  (table) => ({
    primary: primaryKey({
      name: "primary",
      columns: [table.userID, table.fingerprint],
    }),
  }),
);

export const userShippingTable = mysqlTable("user_shipping", {
  ...id,
  ...timestamps,
  userID: ulid("user_id")
    .references(() => userTable.id, {
      onDelete: "cascade",
    })
    .notNull(),
  address: address("address").notNull(),
});
