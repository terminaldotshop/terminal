import { z } from "zod";
import { SubscriptionFrequency, subscriptionTable } from "./subscription.sql";
import { useTransaction } from "../drizzle/transaction";
import { eq, sql } from "drizzle-orm";
import { useUserID } from "../actor";
import { createID } from "../util/id";
import { fn } from "../util/fn";
import { productTable, productVariantTable } from "../product/product.sql";

export module Subscription {
  export const Info = z.object({
    id: z.string(),
    productVariantID: z.string(),
    quantity: z.number().int(),
    shippingID: z.string(),
    cardID: z.string(),
    frequency: SubscriptionFrequency,
  });

  export type Info = z.infer<typeof Info>;

  export const list = () =>
    useTransaction(async (tx) =>
      tx
        .select()
        .from(subscriptionTable)
        .where(eq(subscriptionTable.userID, useUserID()))
        .then((rows) =>
          rows.map(
            (r): Info => ({
              id: r.id,
              cardID: r.cardID,
              quantity: r.quantity,
              frequency: r.frequency,
              shippingID: r.shippingID,
              productVariantID: r.productVariantID,
            }),
          ),
        ),
    );

  export const create = fn(
    z.object({
      productVariantID: Info.shape.productVariantID,
      quantity: Info.shape.quantity,
      shippingID: Info.shape.shippingID,
      cardID: Info.shape.cardID,
      frequency: Info.shape.frequency,
    }),
    async (input) =>
      useTransaction(async (tx) => {
        const id = createID("subscription");
        const product = await tx
          .select({
            subscription: productTable.subscription,
          })
          .from(productVariantTable)
          .innerJoin(
            productTable,
            eq(productVariantTable.productID, productTable.id),
          )
          .where(eq(productVariantTable.id, input.productVariantID))
          .then((rows) => rows[0]);
        if (!product?.subscription) {
          throw new Error("Product variant does not allow subscriptions");
        }
        if (
          product.subscription === "required" &&
          input.frequency !== "fixed"
        ) {
          throw new Error(
            "Subscription frequency must be 'fixed' for this product",
          );
        }
        await tx
          .insert(subscriptionTable)
          .values({
            id,
            userID: useUserID(),
            productVariantID: input.productVariantID,
            quantity: input.quantity,
            shippingID: input.shippingID,
            cardID: input.cardID,
            frequency: input.frequency,
          })
          .onDuplicateKeyUpdate({
            set: {
              quantity: sql`VALUES(quantity)`,
              shippingID: sql`VALUES(shippingID)`,
              cardID: sql`VALUES(cardID)`,
              frequency: sql`VALUES(frequency)`,
              timeDeleted: null,
            },
          });
      }),
  );
}
