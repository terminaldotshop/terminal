import { Hono, MiddlewareHandler } from "hono";
import { logger } from "hono/logger";
import { session } from "./session";
import { createContext } from "./context";
import { cors } from "hono/cors";
import { vValidator } from "@hono/valibot-validator";
import {
  array,
  email,
  integer,
  length,
  minLength,
  number,
  object,
  optional,
  string,
} from "valibot";
import { Resource } from "sst";
import { stripe } from "./stripe";
import { shippo } from "./shippo";

const SessionContext = createContext<typeof session.$typeValues>();

const from = {
  name: "Terminal Products, Inc.",
  street1: "7969 NW 2nd Street",
  street2: "#1129",
  city: "Miami",
  state: "FL",
  zip: "33126",
  country: "US",
};

// single 12oz bag
const small = {
  length: 7,
  width: 4.5,
  height: 2.75,
  distance_unit: "in",
  weight: 1,
  mass_unit: "lb",
};

// 1-3 12oz bags
const large = {
  length: 11.25,
  width: 8,
  height: 3,
  distance_unit: "in",
  weight: 1,
  mass_unit: "lb",
};

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
          line2: optional(string()),
          city: string(),
          state: optional(string([length(2)])),
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

      // TODO: handle other products, only nil blend for now
      const quantity = body.products.reduce(
        (total, product) => total + product.quantity,
        0,
      );

      if (quantity <= 0) throw new Error("Quantity must be greater than 0");

      let largeBoxesNeeded = Math.floor(quantity / 3);
      let singleBoxesNeeded = quantity % 3;

      if (singleBoxesNeeded === 2) {
        largeBoxesNeeded += 1;
        singleBoxesNeeded = 0; // Used a large box instead of single boxes
      }

      const parcels = [];
      for (let i = 0; i < largeBoxesNeeded; i++) parcels.push({ ...large });
      for (let i = 0; i < singleBoxesNeeded; i++) parcels.push({ ...small });

      const shipment = await shippo("/shipments", {
        method: "POST",
        body: JSON.stringify({
          address_from: from,
          address_to: {
            ...body.shipping,
            line1: undefined,
            line2: undefined,
            street1: body.shipping.line1,
            street2: body.shipping.line2,
          },
          parcels,
          async: false,
        }),
      });

      console.log("shipment", shipment);

      if (shipment.status !== "SUCCESS")
        throw new Error("Failed to get shipping rates.");

      const [rate] = shipment.rates.sort(
        (a, b) => Number.parseFloat(a.amount) - Number.parseFloat(b.amount),
      );
      const shipping = {
        id: rate.object_id,
        name: `${rate.provider} ${rate.servicelevel.name}`,
        cost: Number.parseFloat(rate.amount),
        estimate: rate.duration_terms,
      };

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
        metadata: {
          rate: shipping.id,
        },
        shipping_cost: {
          shipping_rate_data: {
            type: "fixed_amount",
            display_name: shipping.name,
            fixed_amount: {
              currency: "usd",
              amount: shipping.cost * 100,
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
        shipping,
        total: result.amount_due,
      });
    },
  )
  .post(
    "/api/payment",
    vValidator(
      "json",
      object({
        orderID: string([minLength(1)]),
        token: string([minLength(1)]),
      }),
    ),
    async (c) => {
      const body = c.req.valid("json");
      const invoice = await stripe().invoices.retrieve(body.orderID);
      invoice.metadata.rate;
      const paymentMethod = await stripe().paymentMethods.create({
        type: "card",
        card: {
          token: body.token,
        },
      });
      let existing = await stripe()
        .paymentMethods.list()
        .then((result) =>
          result.data.find(
            (pm) => pm.card.fingerprint === paymentMethod.card.fingerprint,
          ),
        );
      if (!existing) {
        existing = await stripe().paymentMethods.attach(paymentMethod.id, {
          customer: useUserID(),
        });
      }
      await stripe().invoices.pay(body.orderID, {
        payment_method: existing.id,
      });

      const label = await shippo("/transactions", {
        method: "POST",
        body: JSON.stringify({
          rate: invoice.metadata.rate,
          async: false,
        }),
      });

      console.log("label", label);

      await stripe().invoices.update(body.orderID, {
        metadata: {
          ...invoice.metadata,
          label: label.label_url,
          trackingNumber: label.tracking_number,
          trackingUrl: label.tracking_url_provider,
        },
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
