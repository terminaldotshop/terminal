import { Layout, Page, io } from "@interval/sdk";
import { useTransaction } from "@terminal/core/drizzle/transaction";
import { cartItemTable, cartTable } from "@terminal/core/cart/cart.sql";
import { userShippingTable, userTable } from "@terminal/core/user/user.sql";
import { cardTable } from "@terminal/core/card/card.sql";
import { eq, sql, sum } from "@terminal/core/drizzle/index";
import { productVariantTable } from "@terminal/core/product/product.sql";

export const Cart = new Page({
  name: "Cart",
  handler: async (c) => {
    return new Layout({
      title: "Cart",
      children: [
        io.display.table("", {
          getData: async (input) => {
            return useTransaction(async (tx) => ({
              data: await tx
                .select({
                  cartID: cartTable.id,
                  userID: userTable.id,
                  email: userTable.email,
                  cardID: cardTable.id,
                  shippingID: userShippingTable.id,
                  items: sum(cartItemTable.quantity),
                  cost: sql`SUM(${cartItemTable.quantity} * ${productVariantTable.price})`,
                })
                .from(cartTable)
                .leftJoin(userTable, eq(cartTable.userID, userTable.id))
                .leftJoin(cardTable, eq(cartTable.cardID, cardTable.id))
                .leftJoin(
                  cartItemTable,
                  eq(cartTable.userID, cartItemTable.userID),
                )
                .leftJoin(
                  productVariantTable,
                  eq(cartItemTable.productVariantID, productVariantTable.id),
                )
                .leftJoin(
                  userShippingTable,
                  eq(cartTable.shippingID, userShippingTable.id),
                )
                .groupBy(
                  cartTable.id,
                  userTable.id,
                  userTable.email,
                  cardTable.id,
                  userShippingTable.id,
                )
                .offset(input.offset)
                .limit(input.pageSize),
            }));
          },
        }),
      ],
    });
  },
});
