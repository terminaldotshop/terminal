import { Order } from "@terminal/core/order/order";
import { Hono } from "hono";

export module Print {
  export const route = new Hono().get("/", async (ctx) => {
    const result = await Order.getNextLabel();
    return ctx.json(result);
  });
}
