import { Resource } from "sst";
import { eq, sql } from "drizzle-orm";
import { fn } from "../util/fn";
import { useTransaction } from "../drizzle/transaction";
import { orderItemTable, orderTable } from "../order/order.sql";
import { productTable, productVariantTable } from "../product/product.sql";
import { Address } from "../address";
import { VisibleError } from "../error";
import { z } from "zod";
import { userTable } from "../user/user.sql";

const TERMINAL_ADDRESS = {
  name: "Terminal Products Inc",
  street1: "403 Watchung Ave",
  city: "Plainfield",
  state: "NJ",
  country: "US",
  zip: "07060",
  phone: "+17724446678",
};

const CUSTOMS_DECLARATION = {
  contents_type: "MERCHANDISE",
  non_delivery_option: "RETURN",
  certify: true,
  certify_signer: "Dax Raad",
  incoterm: "DDU",
  exporter_identification: {
    // eori_number: "",
    tax_id: {
      number: "365099653",
      type: "EIN",
    },
  },
};
const CUSTOMS_DECLARATION_ITEM = {
  description: "Roasted Coffee Beans",
  quantity: 1,
  mass_unit: "oz",
  value_currency: "USD",
  tariff_number: "0901.21.00",
  origin_country: "US",
};

export module Shippo {
  export const createShipmentRate = fn(
    z.object({
      ounces: z.number(),
      address: Address,
      subtotal: z.number().int(),
    }),
    async (input) => {
      const shipping = input.address;
      const country = shipping.country.toUpperCase();
      const international = country !== "US";
      const shipment = await api("POST", "/shipments", {
        address_from: TERMINAL_ADDRESS,
        address_to: {
          name: shipping.name,
          street1: shipping.street1,
          street2: shipping.street2,
          city: shipping.city,
          state: shipping.province,
          country: shipping.country,
          zip: shipping.zip,
          phone: shipping.phone,
        },
        customs_declaration: international
          ? {
              ...CUSTOMS_DECLARATION,
              eel_pfc: country === "CA" ? "NOEEI_30_36" : "NOEEI_30_37_a",
              items: [
                {
                  ...CUSTOMS_DECLARATION_ITEM,
                  net_weight: input.ounces,
                  value_amount: (input.subtotal / 100).toString(),
                },
              ],
            }
          : undefined,
        parcels: [
          {
            length: 12,
            width: 12,
            height: 12,
            distance_unit: "in",
            weight: input.ounces,
            mass_unit: "oz",
          },
        ],
        async: false,
        extra: { bypass_address_validation: true },
        // carrier_accounts: ["6c1d6acf7cc74ec5a1a4e64a5bd19107"],
      });

      console.error(JSON.stringify(shipment));
      if (shipment.status !== "SUCCESS" || !shipment.rates?.length) {
        throw new VisibleError(
          "input",
          "shipment.rate",
          "Failed to get shipping rates.",
        );
      }

      shipment.rates.sort(
        (a: { amount: string }, b: { amount: string }) =>
          Number.parseFloat(a.amount) - Number.parseFloat(b.amount),
      );
      const rate = shipment.rates[0];
      const shippingAmount = Number.parseFloat(rate.amount) * 100;
      console.log(rate);
      return {
        shippoRateID: rate.object_id,
        shippingAmount,
        shippingService: `${rate.provider} ${rate.servicelevel.name}`,
        shippingDeliveryEstimate: rate.duration_terms,
      };
    },
  );

  export const createShipment = fn(z.string(), async (orderID) => {
    const items = await useTransaction((tx) =>
      tx
        .select({
          address: orderTable.shippingAddress,
          name: sql`CONCAT(${productTable.name}, " - ", ${productVariantTable.name})`,
          quantity: orderItemTable.quantity,
          amount: orderItemTable.amount,
          shippoRateID: orderTable.shippoRateID,
          shippoLabelID: orderTable.shippoLabelID,
          existingRateID: orderTable.shippoRateID,
          email: userTable.email,
          price: productVariantTable.price,
        })
        .from(orderTable)
        .innerJoin(orderItemTable, eq(orderItemTable.orderID, orderTable.id))
        .innerJoin(userTable, eq(userTable.id, orderTable.userID))
        .innerJoin(
          productVariantTable,
          eq(productVariantTable.id, orderItemTable.productVariantID),
        )
        .innerJoin(
          productTable,
          eq(productTable.id, productVariantTable.productID),
        )
        .where(eq(orderTable.id, orderID))
        .execute(),
    );
    if (!items.length) throw new Error("Order not found");
    if (items.some((item) => item.shippoLabelID)) {
      console.log("shipment already created", orderID);
      return;
    }
    const terminalOrder = items[0]!;
    if (!terminalOrder.shippoRateID) {
      const rate = await createShipmentRate({
        address: terminalOrder.address,
        ounces: items
          .map((item) => item.quantity * 12)
          .reduce((a, b) => a + b, 0),
        subtotal: items.map((item) => item.amount).reduce((a, b) => a + b, 0),
      });
      terminalOrder.shippoRateID = rate.shippoRateID;
    }
    const shipping = terminalOrder.address;

    const weight = items.reduce((sum, item) => sum + item.quantity * 12, 0);
    const order = await api("POST", "/v1/orders", {
      to_address: {
        name: shipping.name,
        street1: shipping.street1,
        street2: shipping.street2,
        city: shipping.city,
        state: shipping.province,
        country: shipping.country,
        zip: shipping.zip,
        // International orders need a phone, and email helps
        phone: shipping.phone,
        email: terminalOrder.email,
      },
      line_items: items.map((item) => ({
        quantity: item.quantity,
        title: item.name,
        weight: 12,
        weight_unit: "oz",
      })),
      placed_at: new Date().toISOString(),
      weight,
      weight_unit: "oz",
    });
    console.log({ order });
    const rate = await api("GET", "/rates/" + terminalOrder.shippoRateID);
    console.log({ rate });

    const shipment = await api("POST", "/v1/transactions", {
      rate: terminalOrder.shippoRateID,
      metadata: JSON.stringify({
        orderID,
      }),
      order: order.object_id,
    });
    console.log({ shipment });
    if (shipment.object_state !== "VALID") {
      throw new Error("Shipment invalid");
    }
    await useTransaction((tx) =>
      tx
        .update(orderTable)
        .set({
          shippoOrderID: order.object_id,
          shippoLabelID: shipment.object_id,
        })
        .where(eq(orderTable.id, orderID))
        .execute(),
    );

    while (true) {
      const result = await api("GET", "/transactions/" + shipment.object_id);
      console.log({ polling: JSON.stringify(result, null, 2) });
      if (result.status === "ERROR") throw new Error("Shipment invalid");
      if (result.status !== "SUCCESS") {
        await new Promise((r) => setTimeout(r, 1000));
        continue;
      }
      await useTransaction((tx) =>
        tx
          .update(orderTable)
          .set({
            trackingURL: result.tracking_url_provider,
            labelURL: result.label_url,
            trackingNumber: result.tracking_number,
          })
          .where(eq(orderTable.id, orderID))
          .execute(),
      );
      break;
    }
  });

  export const assertValidAddress = fn(
    Address.omit({ country: true }).extend({ country: z.string() }),
    async (input) => {
      const result = await api("POST", "/v1/addresses", {
        name: input.name,
        street1: input.street1,
        street2: input.street2,
        city: input.city,
        state: input.province,
        country: input.country,
        zip: input.zip,
        validate: true,
      });
      if (result.test || result.validation_results?.is_valid) {
        return {
          ...input,
          street1: result.street1,
          street2: result.street2,
          city: result.city,
          province: result.province ?? result.state,
          country: result.country,
          zip: result.zip,
          phone: result.phone,
        };
      }
      console.log({ result });
      throw new AddressInvalidError();
    },
  );

  export class AddressInvalidError extends VisibleError {
    constructor() {
      super("input", "address.invalid", "Address is invalid");
    }
  }

  const ROOT = "https://api.goshippo.com";
  async function api(method: string, path: string, body?: any) {
    return fetch(ROOT + path, {
      method,
      headers: {
        authorization: `ShippoToken ${Resource.ShippoSecret.value}`,
        "content-type": "application/json",
      },
      body: body ? JSON.stringify(body) : undefined,
    }).then((r) => {
      return r.json();
    }) as any;
  }
}
