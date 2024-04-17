import { Hono, MiddlewareHandler } from "hono";
import { logger } from "hono/logger";
import { session } from "./session";
import { createContext } from "./context";
import { cors } from "hono/cors";
import { vValidator } from "@hono/valibot-validator";
import { array, email, integer, length, number, object, string } from "valibot";
import { Resource } from "sst";
import { swell } from "./swell";
import { stripe } from "./stripe";

const SessionContext = createContext<typeof session.$typeValues>();

function useUserID() {
  const session = SessionContext.use();
  if (session.type !== "user") throw new Error("User session expected");
  return session.properties.userID;
}

const auth: MiddlewareHandler = async (c, next) => {
  const authHeader = c.req.header("authorization");
  if (authHeader) {
    const match = authHeader.match(/^Bearer (.+)$/);
    if (!match) {
      return c.json(
        { error: "Bearer token not found or improperly formatted" },
        401,
      );
    }
    const bearerToken = match[1];
    const result = await session.verify(bearerToken);
    return SessionContext.with(result, next);
  }
  return SessionContext.with({ type: "public", properties: {} }, next);
};

const app = new Hono()
  .use(logger())
  .use(cors())
  .use(auth)
  .get("/api/user/me", async (c) => {
    return c.json({ userID: useUserID() });
  })
  .get("/api/product", async (c) => {
    const products = await stripe().products.list({
      expand: ["data.default_price"],
    });
    return c.json(
      products.data.map((p) => ({
        id:
          typeof p.default_price === "string"
            ? p.default_price
            : p.default_price.id,
        name: p.name,
        description: p.description,
        price:
          typeof p.default_price === "string"
            ? p.default_price
            : p.default_price.unit_amount,
      })),
    );
  })
  .post(
    "/api/order",
    vValidator(
      "json",
      object({
        email: string(),
        shipping: object({
          name: string(),
          line1: string(),
          line2: string(),
          city: string(),
          state: string(),
          country: string([length(2)]),
          zip: string(),
        }),
        products: array(
          object({
            id: string(),
            quantity: number([integer()]),
          }),
        ),
      }),
    ),
    async (c) => {
      const body = c.req.valid("json");
      console.log(body);
      await stripe().customers.update(useUserID(), {
        email: body.email,
      });
      const invoice = await stripe().invoices.create({
        shipping_details: {
          name: body.shipping.name,
          address: {
            city: body.shipping.city,
            line1: body.shipping.line1,
            line2: body.shipping.line2,
            country: body.shipping.country,
            state: body.shipping.state,
            postal_code: body.shipping.zip,
          },
        },
        shipping_cost: {
          shipping_rate_data: {
            type: "fixed_amount",
            display_name: "Standard Shipping",
            fixed_amount: {
              currency: "usd",
              amount: 1000,
            },
          },
        },
        customer: useUserID(),
      });
      for (const product of body.products) {
        await stripe().invoiceItems.create({
          invoice: invoice.id,
          customer: useUserID(),
          price: product.id,
          quantity: product.quantity,
        });
      }
      const result = await stripe().invoices.retrieve(invoice.id);

      return c.json({
        id: invoice.id,
        tax: result.tax,
        subtotal: result.subtotal,
        shipping: result.shipping_cost.amount_total,
        total: result.amount_due,
      });
    },
  )
  .post(
    "/api/payment",
    vValidator(
      "json",
      object({
        orderID: string([length(1)]),
        token: string([length(1)]),
      }),
    ),
    async (c) => {
      const body = c.req.valid("json");
      const paymentMethod = await stripe().paymentMethods.create({
        type: "card",
        card: {
          token: body.token,
        },
      });
      const attachment = await stripe().paymentMethods.attach(
        paymentMethod.id,
        {
          customer: useUserID(),
        },
      );
      await stripe().invoices.pay(body.orderID, {
        payment_method: attachment.id,
      });
      return c.json(true);
    },
  )
  .post(
    "/api/subscription",
    vValidator(
      "json",
      object({
        email: string([email("Please provide a valid email address")]),
      }),
    ),
    async (c) => {
      const body = c.req.valid("json");
      console.log("subscribing", body.email);
      const result = await fetch(
        `https://api.airtable.com/v0/appKabRJfxfpSDVTo/subscribers`,
        {
          method: "POST",
          body: JSON.stringify({
            fields: body,
          }),
          headers: {
            Authorization: `Bearer ${Resource.AirtableSecret.value}`,
            "Content-Type": "application/json",
          },
        },
      );
      console.log(result, await result.json());
      return c.json({});
    },
  );

export default app;
export type AppType = typeof app;
