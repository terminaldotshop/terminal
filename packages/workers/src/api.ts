import { Hono, MiddlewareHandler } from "hono";
import { logger } from "hono/logger";
import { session } from "./session";
import { createContext } from "./context";
import { cors } from "hono/cors";
import { vValidator } from "@hono/valibot-validator";
import { email, object, string } from "valibot";
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
