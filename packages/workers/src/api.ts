import { Hono, MiddlewareHandler } from "hono";
import { logger } from "hono/logger";
import { session } from "./session";
import { createContext } from "./context";
import { cors } from "hono/cors";
import { vValidator } from "@hono/valibot-validator";
import { array, email, integer, length, number, object, string } from "valibot";
import { Resource } from "sst";
import { swell } from "./swell";

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
    const products = await swell("/products", {});
    return c.json(products);
  })
  .post(
    "/api/order",
    // vValidator(
    //   "json",
    //   object({
    //     shipping: object({
    //       address1: string(),
    //       address2: string(),
    //       city: string(),
    //       country: string([length(2)]),
    //       name: string(),
    //     }),
    //     products: array(
    //       object({
    //         id: string(),
    //         quantity: number([integer()]),
    //       }),
    //     ),
    //   }),
    // ),
    async (c) => {
      const accountID = "6612082b13c85300127985de";
      console.log("wtf");

      console.log(
        "ok",
        await swell("/orders", {
          method: "POST",
          headers: {
            "content-type": "application/x-www-form-urlencoded",
          },
          body: new URLSearchParams({
            account_id: accountID,
            "items[0][product_id]": "6615f3dc14a8960012e304d9",
            "items[0][quantity]": "1",
            "shipping[price]": "10",
          }).toString(),
        }),
      );
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
