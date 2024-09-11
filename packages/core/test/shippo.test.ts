import { expect, describe } from "bun:test";
import { withTestUser } from "./util";
import { Shippo } from "../src/shippo";

describe("shippo", async () => {
  const quantity = 2;
  const address = {
    name: "John Smith",
    zip: "33133",
    city: "Miami",
    country: "US",
    street1: "2800 SW 28th Terrace",
    province: "FL",
  };

  withTestUser("estimateShippingRate", async () => {
    const rate = await Shippo.createShipmentRate({
      subtotal: 100,
      ounces: quantity * 12,
      address,
    });

    expect(rate.shippingAmount > 0);
  });
});
