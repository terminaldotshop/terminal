import { mysqlTable, int, unique, text } from "drizzle-orm/mysql-core";
import { dollar, id, timestamps, ulid } from "../drizzle/types";
import { userShippingTable, userTable } from "../user/user.sql";
import { productVariantTable } from "../product/product.sql";
import { cardTable } from "../card/card.sql";

export const cartItemTable = mysqlTable(
  "cart_item",
  {
    ...id,
    ...timestamps,
    userID: ulid("user_id")
      .references(() => userTable.id, {
        onDelete: "cascade",
      })
      .notNull(),
    productVariantID: ulid("product_variant_id")
      .references(() => productVariantTable.id, {
        onDelete: "cascade",
      })
      .notNull(),
    quantity: int("quantity").notNull(),
  },
  (table) => ({
    unique: unique("unique").on(table.userID, table.productVariantID),
  }),
);

export const cartTable = mysqlTable("cart", {
  ...id,
  ...timestamps,
  userID: ulid("user_id")
    .references(() => userTable.id, {
      onDelete: "cascade",
    })
    .notNull()
    .unique(),
  shippingID: ulid("shipping_id").references(() => userShippingTable.id, {
    onDelete: "set null",
  }),
  cardID: ulid("card_id").references(() => cardTable.id, {
    onDelete: "set null",
  }),
  shippingAmount: dollar("shipping_amount"),
  shippingService: text("shipping_service"),
  shippoRateID: text("shippo_rate_id"),
  shippingDeliveryEstimate: text("shipping_delivery_estimate"),
});
