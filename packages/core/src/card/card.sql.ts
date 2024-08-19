import {
  char,
  int,
  mysqlTable,
  text,
  unique,
  varchar,
} from "drizzle-orm/mysql-core";
import { id, timestamps, ulid } from "../drizzle/types";
import { userTable } from "../user/user.sql";

export const cardTable = mysqlTable(
  "card",
  {
    ...id,
    ...timestamps,
    userID: ulid("user_id")
      .references(() => userTable.id, {
        onDelete: "cascade",
      })
      .notNull(),
    stripePaymentMethodID: varchar("stripe_payment_method_id", {
      length: 255,
    }).notNull(),
    brand: text("brand").notNull(),
    expirationMonth: int("expiration_month").notNull(),
    expirationYear: int("expiration_year").notNull(),
    last4: char("last4", { length: 4 }).notNull(),
  },
  (table) => ({
    unique: unique("unique").on(table.userID, table.stripePaymentMethodID),
  }),
);
