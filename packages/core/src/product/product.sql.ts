import {
  text,
  mysqlTable,
  varchar,
  int,
  primaryKey,
} from "drizzle-orm/mysql-core";
import { dollar, id, ulid, timestamps } from "../drizzle/types";
import { inventoryTable } from "../inventory/inventory.sql";
import { z } from "zod";

export const productTable = mysqlTable("product", {
  ...id,
  ...timestamps,
  name: varchar("name", { length: 255 }).notNull(),
  description: text("description").notNull(),
  order: int("order"),
  subscription: varchar("subscription", {
    length: 255,
  }).$type<SubscriptionSetting>(),
});

export const productVariantTable = mysqlTable("product_variant", {
  ...id,
  ...timestamps,
  productID: ulid("product_id")
    .notNull()
    .references(() => productTable.id, {
      onDelete: "cascade",
    }),
  name: varchar("name", { length: 255 }).notNull(),
  price: dollar("price").notNull(),
});

export const productVariantInventoryTable = mysqlTable(
  "product_variant_inventory",
  {
    ...timestamps,
    productVariantID: ulid("product_variant_id")
      .notNull()
      .references(() => productVariantTable.id, {
        onDelete: "cascade",
      }),
    inventoryID: ulid("inventory_id")
      .notNull()
      .references(() => inventoryTable.id, {
        onDelete: "cascade",
      }),
  },
  (table) => ({
    primary: primaryKey({
      columns: [table.productVariantID, table.inventoryID],
    }),
  }),
);

export const SubscriptionSetting = z.union([
  z.literal("allowed"),
  z.literal("required"),
]);
export type SubscriptionSetting = z.infer<typeof SubscriptionSetting>;
