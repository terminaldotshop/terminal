import { Interval } from "@interval/sdk";
import { Resource } from "sst";

import Product from "./routes/product";
import { Cart } from "./routes/cart";
import { User } from "./routes/user";
import { Order } from "./routes/order";
import { InventoryPage } from "./routes/inventory";

const interval = new Interval({
  apiKey: Resource.IntervalKey.value,
  routes: {
    product: Product,
    cart: Cart,
    user: User,
    order: Order,
    inventory: InventoryPage,
  },
});

interval.listen();
