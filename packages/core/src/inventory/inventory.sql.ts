import { int, mysqlTable, text, unique, varchar } from "drizzle-orm/mysql-core";
import { id, ulid, timestamps } from "../drizzle/types";

export const inventoryTable = mysqlTable("inventory", {
  ...id,
  ...timestamps,
  name: varchar("name", { length: 255 }).notNull().unique(),
  description: text("description"),
});

export const inventoryRecordTable = mysqlTable("inventory_record", {
  ...id,
  ...timestamps,
  inventoryID: ulid("inventory_id")
    .references(() => inventoryTable.id, {
      onDelete: "cascade",
    })
    .notNull(),
  quantity: int("quantity").notNull(),
  notes: text("notes"),
});
