import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { Result } from "./common";
import { Product } from "@terminal/core/product/index";

export module ProductApi {
  export const ProductVariantSchema = z
    .object(Product.Variant.shape)
    .openapi("ProductVariant");

  export const ProductSchema = z
    .object(Product.Info.shape)
    .extend({
      variants: ProductVariantSchema.array(),
    })
    .openapi("Product");

  export const route = new OpenAPIHono().openapi(
    createRoute({
      security: [
        {
          Bearer: [],
        },
      ],
      method: "get",
      path: "/",
      responses: {
        200: {
          content: {
            "application/json": {
              schema: Result(ProductSchema.array()),
            },
          },
          description: "Returns a list of products",
        },
      },
    }),
    async (c) => {
      c.header("Cache-Control", "s-maxage=60");
      return c.json(
        {
          result: await Product.list(),
        },
        200,
      );
    },
  );
}
