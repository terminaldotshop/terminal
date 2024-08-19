import { expect, describe, it } from "bun:test";
import { withTestUser } from "./util";
import { Cart } from "../src/cart/index";
import { Product } from "../src/product/index";
import { createID } from "../src/util/id";
import { User } from "../src/user";
import { Card } from "../src/card";
import { Order } from "../src/order/order";

describe("cart", async () => {
  const price = 1000;
  const quantity = Math.round(Math.random() * 10);
  const productID = await Product.create({
    id: createID("product"),
    name: "test-product",
    description: "",
  });
  const variantID = await Product.addVariant({
    productID: productID,
    name: "test-variant",
    price,
  });

  withTestUser("setItem", async () => {
    await Cart.setItem({
      productVariantID: variantID,
      quantity,
    });
    const cart = await Cart.get();
    expect(cart.items[0]?.quantity).toEqual(quantity);
  });

  withTestUser("flow", async () => {
    await Cart.setItem({
      productVariantID: variantID,
      quantity,
    });

    const cart = await Cart.get();
    expect(cart.items).toHaveLength(1);
    expect(cart.items[0]!.productVariantID).toEqual(variantID);
    expect(cart.items[0]!.quantity).toEqual(quantity);
    expect(cart.subtotal).toEqual(quantity * price);

    const cardID = await Card.create({
      token: "tok_visa",
    });

    const shippingID = await User.addShipping({
      name: "John Smith",
      zip: "33133",
      city: "Miami",
      country: "US",
      street1: "2800 SW 28th Terrace",
      province: "FL",
    });

    await Cart.setShipping(shippingID);
    await Cart.setCard(cardID);
    const orderID = await Order.convertCart();
    console.log(orderID);
    const order = await Order.fromID(orderID);

    expect(order.items.length).toEqual(1);
    expect(order.items[0]!.productVariantID).toEqual(variantID);
    expect(order.items[0]!.quantity).toEqual(quantity);
    expect(order.amount.subtotal).toEqual(quantity * price);
  });
});
