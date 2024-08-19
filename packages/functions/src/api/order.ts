import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { Result } from "./common";
import { Order } from "@terminal/core/order/order";

export module OrderApi {
  export const OrderSchema = z.object(Order.Info.shape).openapi("Order");

  export const route = new OpenAPIHono()
    .openapi(
      createRoute({
        method: "get",
        path: "/{id}",
        responses: {
          404: {
            content: {
              "application/json": {
                schema: z.object({ error: z.string() }),
              },
            },
            description: "Order not found",
          },
          200: {
            content: {
              "application/json": {
                schema: Result(OrderSchema),
              },
            },
            description: "Returns order",
          },
        },
      }),
      async (c) => {
        const order = await Order.fromID(c.req.param("id"));
        if (!order) return c.json({ error: "Order not found" }, 404);
        return c.json({ result: order }, 200);
      },
    )
    .openapi(
      createRoute({
        security: [{ Bearer: [] }],
        method: "post",
        path: "/",
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(OrderSchema),
              },
            },
            description: "Returns the order",
          },
        },
      }),
      async (c) => {
        const orderID = await Order.convertCart();
        return c.json({ result: await Order.fromID(orderID) }, 200);
      },
    );
}
