import { z } from "zod";
import { createTransaction, useTransaction } from "../drizzle/transaction";
import { fn } from "../util/fn";
import { cartItemTable, cartTable } from "./cart.sql";
import { createID } from "../util/id";
import { productVariantTable } from "../product/product.sql";
import { and, eq, getTableColumns, ne, sql, sum } from "drizzle-orm";
import { useUserID } from "../actor";
import { userShippingTable } from "../user/user.sql";
import { cardTable } from "../card/card.sql";
import { Address } from "../address";
import { Shippo } from "../shippo/";

export module Cart {
  export const Item = z.object({
    id: z.string(),
    productVariantID: z.string(),
    quantity: z.number().int().gte(0),
    subtotal: z.number().int(),
  });
  export type Item = z.infer<typeof Item>;

  export const Info = z.object({
    items: z.array(Item),
    subtotal: z.number().int().gte(0),
    shippingID: z.string().optional(),
    cardID: z.string().optional(),
    amount: z.object({
      subtotal: z.number().int(),
      shipping: z.number().int().optional(),
    }),
    shipping: z
      .object({
        service: z.string().optional(),
        timeframe: z.string().optional(),
      })
      .optional(),
  });
  type Info = z.infer<typeof Info>;

  export async function get() {
    return createTransaction(async (tx): Promise<Info> => {
      const cart = await tx
        .select({
          cardID: cardTable.id,
          shippingID: userShippingTable.id,
          shippingAmount: cartTable.shippingAmount,
          shippingService: cartTable.shippingService,
          shippingDeliveryEstimate: cartTable.shippingDeliveryEstimate,
        })
        .from(cartTable)
        .leftJoin(cardTable, eq(cartTable.cardID, cardTable.id))
        .leftJoin(
          userShippingTable,
          eq(cartTable.shippingID, userShippingTable.id),
        )
        .where(eq(cartTable.userID, useUserID()))
        .then((rows) => rows[0]);
      if (!cart)
        return {
          items: [],
          amount: {
            shipping: 0,
            subtotal: 0,
          },
          subtotal: 0,
        };
      const items = await list();
      const subtotal = items.reduce((acc, item) => item.subtotal + acc, 0);
      return {
        items,
        subtotal,
        amount: {
          subtotal,
          shipping: cart.shippingAmount ?? undefined,
        },
        cardID: cart.cardID || undefined,
        shippingID: cart.shippingID || undefined,
        shipping: {
          service: cart.shippingService || undefined,
          timeframe: cart.shippingDeliveryEstimate || undefined,
        },
      };
    });
  }

  const FREE_SHIPPING_THRESHOLD = 40 * 100;
  export async function calculateShipping(
    subtotal: number,
    ounces: number,
    address: Address,
  ) {
    const rate = await Shippo.createShipmentRate({ ounces, address, subtotal });
    if (address.country === "US") {
      return {
        ...rate,
        shippingAmount: subtotal >= FREE_SHIPPING_THRESHOLD ? 0 : 800,
      };
    }
    return rate;
  }

  export const list = () =>
    useTransaction(async (tx) => {
      return tx
        .select({
          cartItem: getTableColumns(cartItemTable),
          productVariant: getTableColumns(productVariantTable),
          subtotal: sql<string>`(${cartItemTable.quantity} * ${productVariantTable.price})`,
        })
        .from(cartItemTable)
        .innerJoin(
          productVariantTable,
          eq(cartItemTable.productVariantID, productVariantTable.id),
        )
        .where(eq(cartItemTable.userID, useUserID()))
        .then((rows): Item[] =>
          rows.map((row) => ({
            id: row.cartItem.id,
            productVariantID: row.productVariant.id,
            quantity: row.cartItem.quantity,
            subtotal: parseInt(row.subtotal, 10),
          })),
        );
    });

  export const setShipping = fn(z.string(), async (userShippingID) => {
    const shippingInfo = await useTransaction(async (tx) => {
      const response = await tx
        .select({
          count: sum(cartItemTable.quantity).mapWith(Number),
          subtotal:
            sql`sum(${productVariantTable.price} * ${cartItemTable.quantity})`.mapWith(
              Number,
            ),
          address: userShippingTable.address,
        })
        .from(cartItemTable)
        .innerJoin(
          productVariantTable,
          eq(productVariantTable.id, cartItemTable.productVariantID),
        )
        .innerJoin(userShippingTable, eq(userShippingTable.id, userShippingID))
        .where(eq(cartItemTable.userID, useUserID()))
        .then((rows) => rows[0]!);

      const weight = response.count * 12;
      const address = response.address;
      return await calculateShipping(response.subtotal, weight, address);
    });

    await useTransaction(async (tx) => {
      const shippingID = await tx
        .select({
          shippingID: userShippingTable.id,
        })
        .from(userShippingTable)
        .where(eq(userShippingTable.id, userShippingID))
        .then((rows) => rows[0]!.shippingID);
      await tx
        .insert(cartTable)
        .values({
          userID: useUserID(),
          shippingID,
          ...shippingInfo,
          id: createID("cart"),
        })
        .onDuplicateKeyUpdate({
          set: {
            shippingID,
            ...shippingInfo,
          },
        });
    });
  });

  export const setCard = fn(z.string(), (input) =>
    useTransaction(async (tx) => {
      const cardID = await tx
        .select({
          cardID: cardTable.id,
        })
        .from(cardTable)
        .where(and(eq(cardTable.id, input), eq(cardTable.userID, useUserID())))
        .then((rows) => rows[0]?.cardID);
      if (!cardID) {
        throw new Error("card not found");
      }
      await tx
        .insert(cartTable)
        .values({
          userID: useUserID(),
          shippingID: cardID,
          id: createID("cart"),
        })
        .onDuplicateKeyUpdate({
          set: {
            cardID,
          },
        });
    }),
  );

  export const setItem = fn(
    z.object({
      id: z.string().optional(),
      productVariantID: Item.shape.productVariantID,
      quantity: Item.shape.quantity,
    }),
    async (input) => {
      return useTransaction(async (tx) => {
        const id = input.id || createID("cartItem");
        const variant = await tx
          .select({ id: productVariantTable.id })
          .from(productVariantTable)
          .where(eq(productVariantTable.id, input.productVariantID))
          .then((rows) => rows[0]);
        if (!variant) {
          throw new Error("variant not found");
        }
        if (input.quantity <= 0) {
          await tx
            .delete(cartItemTable)
            .where(
              and(
                eq(cartItemTable.productVariantID, variant.id),
                eq(cartItemTable.userID, useUserID()),
              ),
            );
          return;
        }
        await tx
          .insert(cartItemTable)
          .values({
            id,
            quantity: input.quantity,
            productVariantID: variant.id,
            userID: useUserID(),
          })
          .onDuplicateKeyUpdate({
            set: { quantity: input.quantity },
          });
        await tx
          .insert(cartTable)
          .ignore()
          .values({
            userID: useUserID(),
            id: createID("cart"),
          });
      });
    },
  );

  export async function clear() {
    return useTransaction(async (tx) =>
      tx.delete(cartItemTable).where(eq(cartItemTable.userID, useUserID())),
    );
  }

  export async function subtotal() {}
}
