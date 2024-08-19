import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { Result } from "./common";
import { Cart } from "@terminal/core/cart/index";
import { UserApi } from "./user";
import { CardApi } from "./card";

export module CartApi {
  export const CartItemSchema = z.object(Cart.Item.shape).openapi("CartItem");
  export const CartSchema = z
    .object(Cart.Info.shape)
    .extend({ items: CartItemSchema.array() })
    .openapi("Cart");

  export const route = new OpenAPIHono()
    .openapi(
      createRoute({
        security: [{ Bearer: [] }],
        method: "get",
        path: "/",
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(CartSchema),
              },
            },
            description: "Returns a list of items in the cart",
          },
        },
      }),
      async (c) => {
        return c.json(
          {
            result: await Cart.get(),
          },
          200,
        );
      },
    )
    .openapi(
      createRoute({
        security: [{ Bearer: [] }],
        method: "put",
        path: "/item",
        request: {
          body: {
            content: {
              "application/json": {
                schema: z.object({
                  productVariantID: CartItemSchema.shape.productVariantID,
                  quantity: CartItemSchema.shape.quantity,
                }),
              },
            },
          },
        },
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(CartSchema),
              },
            },
            description: "Returns the cart",
          },
        },
      }),
      async (c) => {
        const body = c.req.valid("json");
        await Cart.setItem(body);
        return c.json({ result: await Cart.get() }, 200);
      },
    )
    .openapi(
      createRoute({
        security: [{ Bearer: [] }],
        method: "put",
        path: "/shipping",
        request: {
          body: {
            content: {
              "application/json": {
                schema: z.object({
                  shippingID: UserApi.UserShippingSchema.shape.id,
                }),
              },
            },
          },
        },
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(z.literal("ok")),
              },
            },
            description: "Shipping address was set successfully",
          },
        },
      }),
      async (c) => {
        const body = c.req.valid("json");
        await Cart.setShipping(body.shippingID);
        return c.json({ result: "ok" as const }, 200);
      },
    )
    .openapi(
      createRoute({
        security: [{ Bearer: [] }],
        method: "put",
        path: "/card",
        request: {
          body: {
            content: {
              "application/json": {
                schema: z.object({
                  cardID: CardApi.CardSchema.shape.id,
                }),
              },
            },
          },
        },
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(z.literal("ok")),
              },
            },
            description: "Card was set successfully",
          },
        },
      }),
      async (c) => {
        const body = c.req.valid("json");
        console.log({ body });
        await Cart.setCard(body.cardID);
        return c.json({ result: "ok" as const }, 200);
      },
    );
}
