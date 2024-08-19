import { Order } from "@terminal/core/order/order";

export const handler = async () => {
  await Order.trackInventory();
};
