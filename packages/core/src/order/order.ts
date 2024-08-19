import {
  afterTx,
  createTransaction,
  useTransaction,
} from "../drizzle/transaction";
import { createID } from "../util/id";
import { orderItemTable, orderTable } from "./order.sql";
import { assertFlag, useUserID } from "../actor";
import { userShippingTable, userTable } from "../user/user.sql";
import {
  and,
  eq,
  getTableColumns,
  inArray,
  isNotNull,
  isNull,
  sql,
  sum,
} from "drizzle-orm";
import { cartItemTable, cartTable } from "../cart/cart.sql";
import {
  productVariantInventoryTable,
  productVariantTable,
} from "../product/product.sql";
import { z } from "zod";
import { fn } from "../util/fn";
import { cardTable } from "../card/card.sql";
import { defineEvent } from "../event";
import { bus } from "sst/aws/bus";
import { Resource } from "sst";
import { stripe } from "../stripe";
import { Stripe } from "stripe";
import { Address } from "../address";
import { Shippo } from "../shippo/index";
import { VisibleError } from "../error";
import { inventoryRecordTable } from "../inventory/inventory.sql";

export module Order {
  export const Item = z.object({
    id: z.string(),
    description: z.string().optional(),
    amount: z.number(),
    quantity: z.number().int().gte(0),
    productVariantID: z.string().optional(),
  });

  export const Info = z.object({
    id: z.string(),
    shipping: Address,
    amount: z.object({
      shipping: z.number(),
      subtotal: z.number(),
    }),
    tracking: z.object({
      number: z.string().optional(),
      url: z.string().optional(),
    }),
    items: z.array(Item),
  });
  export type Info = z.infer<typeof Info>;

  export const Event = {
    Created: defineEvent(
      "order.created",
      z.object({
        orderID: Info.shape.id,
      }),
    ),
  };

  export const fromID = fn(Info.shape.id, (input) =>
    useTransaction((tx) =>
      tx
        .select()
        .from(orderTable)
        .innerJoin(orderItemTable, eq(orderTable.id, orderItemTable.orderID))
        .leftJoin(
          productVariantTable,
          eq(orderItemTable.productVariantID, productVariantTable.id),
        )
        .where(eq(orderTable.id, input))
        .then(
          (rows): Info => ({
            id: rows[0]!.order.id,
            shipping: rows[0]!.order.shippingAddress,
            amount: {
              shipping: rows[0]!.order.shippingAmount,
              subtotal: rows.reduce(
                (acc, row) => acc + row.order_item.amount,
                0,
              ),
            },
            tracking: {
              number: rows[0]!.order.trackingNumber || undefined,
              url: rows[0]!.order.trackingURL || undefined,
            },
            items: rows.map((row) => ({
              id: row.order_item.id,
              amount: row.order_item.amount,
              quantity: row.order_item.quantity,
              productVariantID: row.product_variant?.id,
            })),
          }),
        ),
    ),
  );

  export async function convertCart() {
    const userID = useUserID();
    const { items, cart } = await useTransaction(async (tx) => {
      const items = await tx
        .select({
          productVariantID: cartItemTable.productVariantID,
          quantity: cartItemTable.quantity,
          subtotal: sql`(${cartItemTable.quantity} * ${productVariantTable.price})`,
        })
        .from(cartItemTable)
        .innerJoin(
          productVariantTable,
          eq(cartItemTable.productVariantID, productVariantTable.id),
        )
        .where(eq(cartItemTable.userID, userID))
        .then((rows) =>
          rows.map((row) => ({
            ...row,
            subtotal: z.coerce.number().int().parse(row.subtotal),
          })),
        );

      const cart = await tx
        .select({
          shipping: userShippingTable.address,
          card: getTableColumns(cardTable),
          stripeCustomerID: userTable.stripeCustomerID,
          email: userTable.email,
          shippingAmount: cartTable.shippingAmount,
          shippoRateID: cartTable.shippoRateID,
        })
        .from(cartTable)
        .innerJoin(
          userShippingTable,
          eq(cartTable.shippingID, userShippingTable.id),
        )
        .innerJoin(cardTable, eq(cartTable.cardID, cardTable.id))
        .innerJoin(userTable, eq(cartTable.userID, userTable.id))
        .where(eq(cartTable.userID, userID))
        .then((rows) => rows[0]);
      return { items, cart };
    });
    if (!cart) throw new Error("No cart found");
    const orderID = createID("order");
    const subtotal = items.reduce((acc, item) => acc + item.subtotal, 0);
    const shipping = cart.shippingAmount;
    if (shipping === null) throw new Error("Shipping amount not set");
    try {
      const payment = await stripe.paymentIntents.create({
        amount: subtotal + shipping,
        automatic_payment_methods: {
          enabled: true,
          allow_redirects: "never",
        },
        confirm: true,
        currency: "usd",
        shipping: {
          name: cart.shipping.name,
          address: {
            city: cart.shipping.city,
            line1: cart.shipping.street1,
            line2: cart.shipping.street2,
            postal_code: cart.shipping.zip,
            state: cart.shipping.province,
            country: cart.shipping.country,
          },
        },
        customer: cart.stripeCustomerID,
        metadata: {
          orderID,
        },
        payment_method: cart.card.stripePaymentMethodID,
      });
      return createTransaction(async (tx) => {
        await tx.insert(orderTable).values({
          id: orderID,
          userID,
          email: cart.email,
          stripePaymentIntentID: payment.id,
          shippingAddress: cart.shipping,
          shippingAmount: shipping,
          shippoRateID: cart.shippoRateID,
          card: {
            brand: cart.card.brand,
            last4: cart.card.last4,
            expiration: {
              month: cart.card.expirationMonth,
              year: cart.card.expirationYear,
            },
          },
        });
        await tx.insert(orderItemTable).values(
          items.map((item) => ({
            id: createID("cartItem"),
            amount: item.subtotal,
            orderID: orderID,
            productVariantID: item.productVariantID,
            quantity: item.quantity,
          })),
        );
        await tx.delete(cartItemTable).where(eq(cartItemTable.userID, userID));
        await afterTx(() =>
          bus.publish(Resource.Bus, Event.Created, { orderID }),
        );
        return orderID;
      });
    } catch (ex: unknown) {
      if (ex instanceof Stripe.errors.StripeCardError) {
        throw new VisibleError("input", "payment.invalid", ex.message);
      }
      throw ex;
    }
  }

  export const createInternal = fn(
    z.object({
      email: z.string().email(),
      items: z.record(z.number().int()),
      address: Address,
    }),
    async (input) => {
      await Shippo.assertValidAddress(input.address);
      const shippingInfo = await Shippo.createShipmentRate({
        ounces: 0,
        address: input.address,
        subtotal: 0,
      });
      await useTransaction(async (tx) => {
        const orderID = createID("order");
        await tx.insert(orderTable).values({
          id: orderID,
          email: input.email,
          shippingAmount: 0,
          shippingAddress: input.address,
          shippoRateID: shippingInfo.shippoRateID,
        });
        for (const [productVariantID, quantity] of Object.entries(
          input.items,
        )) {
          if (quantity < 1) throw new Error("Invalid quantity");
          await tx.insert(orderItemTable).values({
            id: createID("cartItem"),
            amount: 0,
            productVariantID,
            quantity,
            orderID,
          });
        }
        await afterTx(() =>
          bus.publish(Resource.Bus, Event.Created, { orderID }),
        );
      });
    },
  );

  export const setPrinted = fn(Info.shape.id, async (input) => {
    assertFlag("printer");
    await useTransaction(async (tx) =>
      tx
        .update(orderTable)
        .set({
          timePrinted: sql`CURRENT_TIMESTAMP(3)`,
        })
        .where(eq(orderTable.id, input)),
    );
  });

  export async function getNextLabel() {
    await assertFlag("printer");
    const result = await useTransaction((tx) =>
      tx
        .select({
          id: orderTable.id,
          label: orderTable.labelURL,
        })
        .from(orderTable)
        .where(
          and(isNull(orderTable.timePrinted), isNotNull(orderTable.labelURL)),
        )
        .orderBy(orderTable.id)
        .limit(1)
        .then((rows) => rows[0]),
    );
    if (!result) return;
    return result;
  }

  export async function trackInventory() {
    await createTransaction(async (tx) => {
      const items = await tx
        .select({
          quantity: sum(orderItemTable.quantity).mapWith(parseInt),
          inventoryID: productVariantInventoryTable.inventoryID,
        })
        .from(orderItemTable)
        .innerJoin(orderTable, eq(orderItemTable.orderID, orderTable.id))
        .innerJoin(
          productVariantInventoryTable,
          eq(
            orderItemTable.productVariantID,
            productVariantInventoryTable.productVariantID,
          ),
        )
        .where(
          and(
            isNull(orderItemTable.timeInventoryTracked),
            isNotNull(orderTable.timePrinted),
          ),
        )
        .groupBy(productVariantInventoryTable.inventoryID);
      if (items.length === 0) {
        console.log("No inventory to track");
        return;
      }
      await tx.insert(inventoryRecordTable).values(
        items.map((item) => ({
          quantity: item.quantity * -1,
          inventoryID: item.inventoryID,
          id: createID("inventoryRecord"),
          notes: "automated",
        })),
      );

      const updated = await tx
        .select({
          id: orderItemTable.id,
        })
        .from(orderTable)
        .innerJoin(orderItemTable, eq(orderItemTable.orderID, orderTable.id))
        .where(
          and(
            isNull(orderItemTable.timeInventoryTracked),
            isNotNull(orderTable.timePrinted),
          ),
        );
      const result = await tx
        .update(orderItemTable)
        .set({
          timeInventoryTracked: sql`CURRENT_TIMESTAMP(3)`,
        })
        .where(
          inArray(
            orderItemTable.id,
            updated.map((row) => row.id),
          ),
        );
      console.log("Tracked inventory", result.rowsAffected);
    }, "repeatable read");
  }
}
