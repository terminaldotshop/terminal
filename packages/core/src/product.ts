import { stripe } from "./stripe";

/** @deprecated */
export module Product {
  export interface Info {
    id: string;
    name: string;
    description: string;
    price: number;
    inventory?: number;
    type?: string;
  }

  export async function list() {
    const products = await stripe.products.list({
      active: true,
      expand: ["data.default_price"],
    });

    return products.data.flatMap((product): Info[] => {
      if (!product.default_price) return [];
      if (typeof product.default_price === "string") return [];
      return [
        {
          id: product.default_price.id,
          price: product.default_price.unit_amount!,
          name: product.name,
          inventory: product.metadata.inventory
            ? parseInt(product.metadata.inventory)
            : undefined,
          description:
            product.metadata.description ?? product.description ?? "",
          type: product.metadata.type,
        },
      ];
    });
  }

  export module Variant {
    export interface Info {
      id: string;
      productID: string;
      name: string;
      description: string;
      price: number;
    }
  }
}
