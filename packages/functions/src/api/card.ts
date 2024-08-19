import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { Result } from "./common";
import { Card } from "@terminal/core/card/index";

export module CardApi {
  export const CardSchema = z.object(Card.Info.shape).openapi("Card");
  export const route = new OpenAPIHono()
    .openapi(
      createRoute({
        security: [
          {
            Bearer: [],
          },
        ],
        method: "get",
        path: "/",
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(CardSchema.array()),
              },
            },
            description: "Returns a list of cards",
          },
        },
      }),
      async (c) => {
        return c.json(
          {
            result: await Card.list(),
          },
          200,
        );
      },
    )
    .openapi(
      createRoute({
        method: "post",
        path: "/",
        request: {
          body: {
            content: {
              "application/json": {
                schema: z.object({ token: z.string() }),
              },
            },
          },
        },
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(z.string()),
              },
            },
            description: "Returns card ID",
          },
        },
      }),
      async (c) => {
        const result = await Card.create(c.req.valid("json"));
        return c.json({ result }, 200);
      },
    );
}
