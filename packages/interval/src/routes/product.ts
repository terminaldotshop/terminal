import { Action, Page, io, ctx, Layout } from "@interval/sdk";
import { eq } from "@terminal/core/drizzle/index";
import { useTransaction } from "@terminal/core/drizzle/transaction";
import { inventoryTable } from "@terminal/core/inventory/inventory.sql";
import { Product } from "@terminal/core/product/index";
import { productVariantInventoryTable } from "@terminal/core/product/product.sql";

async function selectProduct() {
  let { productID } = ctx.params;
  const products = await Product.list();
  if (!productID) {
    const selected = await io.select.single("select product", {
      options: products.map((p) => ({
        label: p.name,
        value: p.id,
      })),
    });
    productID = selected.value;
  }
  return products.find((p) => p.id === productID)!;
}

async function selectVariant() {
  const product = await selectProduct();
  let { variantID } = ctx.params;
  if (!variantID) {
    const selected = await io.select.single("select variant", {
      options: product.variants.map((p) => ({
        label: p.name,
        value: p.id,
      })),
    });
    variantID = selected.value;
  }
  return {
    product,
    variant: product.variants.find((p) => p.id === variantID)!,
  };
}

export default new Page({
  name: "Product",
  handler: async (input) => {
    const products = await Product.list();
    return new Layout({
      title: "Products",
      menuItems: [
        {
          label: "Create product",
          route: "product/create",
        },
      ],
      children: [
        io.display.table("", {
          data: products.map((p) => ({
            name: p.name,
            description: p.description,
            id: p.id,
            variants: p.variants,
          })),
          columns: [
            {
              label: "name",
              renderCell: (row) => ({
                label: row.name,
                route: "product/detail",
                params: {
                  productID: row.id,
                },
              }),
            },
            "id",
            {
              label: "variants",
              renderCell: (row) => ({
                label: row.variants.length + " variants",
              }),
            },
          ],
          isFilterable: false,
        }),
      ],
    });
  },
  routes: {
    detail: new Page({
      name: "Product Detail",
      unlisted: true,
      async handler() {
        const product = await Product.fromID(ctx.params.productID as string);
        if (!product) throw new Error("Product not found");
        return new Layout({
          title: product.name,
          menuItems: [
            {
              label: "Edit",
              route: "product/edit",
              params: {
                productID: product.id,
              },
            },
          ],
          children: [
            io.display.metadata("", {
              layout: "grid",
              data: [
                {
                  label: "Name",
                  value: product.name,
                },
                {
                  label: "Description",
                  value: product.description,
                },
              ],
            }),
            io.display.link("Add Variant", {
              route: "product/variant/create",
              params: {
                productID: product.id,
              },
            }),
            io.display.table("", {
              data: product.variants.map((item) => ({
                name: item.name,
                price: item.price,
                id: item.id,
              })),
              columns: [
                {
                  label: "name",
                  renderCell: (row) => ({
                    label: row.name,
                    route: "product/variant/edit",
                    params: {
                      productID: product.id,
                      variantID: row.id,
                    },
                  }),
                },
                {
                  label: "price",
                  renderCell: (row) => ({
                    label: "$" + row.price / 100,
                  }),
                },
                "id",
              ],
              isFilterable: false,
            }),
          ],
        });
      },
    }),
    edit: new Action({
      name: "Edit product",
      unlisted: true,
      async handler() {
        const product = await selectProduct();
        const [name, description, order] = await io.group([
          io.input.text("name", {
            defaultValue: product.name,
          }),
          io.input.text("description", {
            defaultValue: product.description,
            multiline: true,
          }),
          io.input.number("order", {
            defaultValue: product.order,
          }),
        ]);
        await Product.edit({
          id: product.id,
          name,
          description,
          order,
        });
        await ctx.redirect({
          route: "product/detail",
          params: { productID: product.id },
        });
      },
    }),
    create: new Action({
      name: "Create product",
      unlisted: true,
      async handler() {
        const [name, description] = await io.group([
          io.input.text("name"),
          io.input.text("description", {
            multiline: true,
          }),
        ]);
        const productID = await Product.create({
          name,
          description,
        });
        ctx.redirect({
          route: "product/detail",
          params: {
            productID: productID,
          },
        });
      },
    }),
    variant: new Page({
      name: "Variant",
      unlisted: true,
      routes: {
        edit: new Action({
          name: "Edit Variant",
          handler: async () => {
            const { product, variant } = await selectVariant();
            const allInventory = await useTransaction((tx) =>
              tx.select().from(inventoryTable),
            );
            const existingInventory = await useTransaction((tx) =>
              tx
                .select()
                .from(productVariantInventoryTable)
                .where(
                  eq(productVariantInventoryTable.productVariantID, variant.id),
                ),
            );
            const [name, price, productIDs] = await io.group([
              io.input.text("name", {
                defaultValue: variant.name,
              }),
              io.input.number("price", {
                defaultValue: variant.price / 100,
                currency: "USD",
              }),
              io.select.multiple("inventory", {
                options: allInventory.map((item) => ({
                  label: item.name,
                  value: item.id,
                })),
                defaultValue: existingInventory.map((item) => ({
                  label: item.inventoryID,
                  value: item.inventoryID,
                })),
              }),
            ]);
            ctx.log("edit variant", { name, price });
            await Product.editVariant({
              id: variant.id,
              name,
              price: price * 100,
              inventoryIDs: productIDs.map((item) => item.value),
            });
            ctx.redirect({
              route: "product/detail",
              params: {
                productID: product.id,
              },
            });
          },
        }),
        create: new Action({
          name: "Create Variant",
          handler: async () => {
            const product = await selectProduct();
            const [name, price] = await io.group([
              io.input.text("name", {
                minLength: 1,
              }),
              io.input.number("price", {
                currency: "USD",
              }),
            ]);
            await Product.addVariant({
              productID: product.id,
              name,
              price: price * 100,
            });
            ctx.redirect({
              route: "product/detail",
              params: {
                productID: product.id,
              },
            });
          },
        }),
      },
    }),
  },
});
