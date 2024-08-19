import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { Result } from "./common";
import { EmailOctopus } from "@terminal/core/email-octopus";

export module EmailApi {
  export const route = new OpenAPIHono().openapi(
    createRoute({
      method: "post",
      path: "/subscription",
      request: {
        body: {
          content: {
            "application/json": {
              schema: z.object({ email: z.string().min(1) }),
            },
          },
        },
      },
      responses: {
        400: {
          content: {
            "application/json": {
              schema: z.object({ error: z.string() }),
            },
          },
          description: "Email is required",
        },
        200: {
          content: {
            "application/json": {
              schema: Result(z.literal("ok")),
            },
          },
          description: "Email subscription was created",
        },
      },
    }),
    async (c) => {
      const body = c.req.valid("json");
      if (!body.email) return c.json({ error: "Email is required" }, 400);
      await EmailOctopus.subscribe({ email: body.email });
      return c.json({ result: "ok" as const }, 200);
    },
  );
}
