import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { Result } from "./common";
import { User } from "@terminal/core/user/index";
import { useUserID } from "@terminal/core/actor";
import { Address } from "@terminal/core/address";

export module UserApi {
  export const AddressSchema = z.object(Address.shape).openapi("Address");
  export const UserSchema = z.object(User.Info.shape).openapi("User");
  export const UserShippingSchema = z
    .object(User.Shipping.shape)
    .extend({ address: AddressSchema })
    .openapi("UserShipping");

  export const route = new OpenAPIHono()
    .openapi(
      createRoute({
        method: "get",
        path: "/me",
        responses: {
          404: {
            content: {
              "application/json": {
                schema: z.object({ error: z.string() }),
              },
            },
            description: "User not found",
          },
          200: {
            content: {
              "application/json": {
                schema: Result(UserSchema),
              },
            },
            description: "Returns user",
          },
        },
      }),
      async (c) => {
        const result = await User.fromID(useUserID());
        if (!result) {
          return c.json({ error: "User not found" }, 404);
        }
        return c.json({ result }, 200);
      },
    )
    .openapi(
      createRoute({
        method: "put",
        path: "/me",
        request: {
          body: {
            content: {
              "application/json": {
                schema: User.update.schema,
              },
            },
          },
        },
        responses: {
          404: {
            content: {
              "application/json": {
                schema: z.object({ error: z.string() }),
              },
            },
            description: "User not found",
          },
          200: {
            content: {
              "application/json": {
                schema: Result(UserSchema),
              },
            },
            description: "Returns user",
          },
        },
      }),
      async (c) => {
        await User.update(c.req.valid("json"));
        const user = await User.fromID(useUserID());
        if (!user) return c.json({ error: "User not found" }, 404);
        return c.json({ result: user }, 200);
      },
    )
    .openapi(
      createRoute({
        method: "get",
        path: "/shipping",
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(UserShippingSchema.array()),
              },
            },
            description: "Returns shipping addresses",
          },
        },
      }),
      async (c) => {
        const result = await User.shipping();
        return c.json({ result }, 200);
      },
    )
    .openapi(
      createRoute({
        method: "post",
        path: "/shipping",
        request: {
          body: {
            content: {
              "application/json": {
                schema: User.addShipping.schema,
              },
            },
          },
        },
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(UserShippingSchema.shape.id),
              },
            },
            description: "Returns shipping address ID",
          },
        },
      }),
      async (c) => {
        const shippingID = await User.addShipping(c.req.valid("json"));
        return c.json({ result: shippingID }, 200);
      },
    )
    .openapi(
      createRoute({
        method: "delete",
        path: "/shipping/{id}",
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(z.literal("ok")),
              },
            },
            description: "Shipping address was deleted successfully",
          },
        },
      }),
      async (c) => {
        await User.removeShipping(c.req.param("id"));
        return c.json({ result: "ok" as const }, 200);
      },
    );
}
