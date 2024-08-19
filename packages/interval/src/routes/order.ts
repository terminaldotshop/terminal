import { Action, Layout, Page, ctx, io } from "@interval/sdk";
import { useTransaction } from "@terminal/core/drizzle/transaction";
import { entries, groupBy, map, pipe } from "remeda";
import {
  and,
  count,
  desc,
  eq,
  inArray,
  isNotNull,
  isNull,
  like,
  or,
  sql,
  sum,
} from "@terminal/core/drizzle/index";
import { orderItemTable, orderTable } from "@terminal/core/order/order.sql";
import { Order as OrderM } from "@terminal/core/order/order";
import { PDFDocument } from "pdf-lib";
import { Resource } from "sst";
import {
  S3Client,
  PutObjectCommand,
  GetObjectCommand,
} from "@aws-sdk/client-s3";
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";
import { Product } from "@terminal/core/product/index";
import {
  productTable,
  productVariantTable,
} from "@terminal/core/product/product.sql";
import { Shippo } from "@terminal/core/shippo/index";

export const Order = new Page({
  name: "Order",
  handler: async () => {
    const needPrinting = await useTransaction((tx) =>
      tx
        .select({ count: count(orderTable.id) })
        .from(orderTable)
        .where(
          and(isNull(orderTable.timePrinted), isNotNull(orderTable.labelURL)),
        )
        .then((rows) => rows[0]!.count),
    );

    const totals = await useTransaction((tx) =>
      tx
        .select({
          product: sql<string>`CONCAT(${productTable.name}, ' - ', ${productVariantTable.name})`,
          total: sql<number>`SUM(${orderItemTable.quantity})`,
        })
        .from(orderItemTable)
        .innerJoin(orderTable, eq(orderItemTable.orderID, orderTable.id))
        .innerJoin(
          productVariantTable,
          eq(orderItemTable.productVariantID, productVariantTable.id),
        )
        .innerJoin(
          productTable,
          eq(productVariantTable.productID, productTable.id),
        )
        .where(
          and(isNull(orderTable.timePrinted), isNotNull(orderTable.labelURL)),
        )
        .groupBy(
          productTable.id,
          productVariantTable.id,
          productTable.name,
          productVariantTable.name,
        ),
    );

    return new Layout({
      title: "Order",
      menuItems:
        needPrinting > 0
          ? [
              {
                label: `Print ${needPrinting} orders`,
                route: "order/print",
              },
            ]
          : [],
      children: [
        io.display.table("Unprinted Items", {
          isFilterable: false,
          data: totals,
        }),
        io.display.table("Orders", {
          getData: async (input) => {
            return useTransaction(async (tx) => ({
              data: await tx
                .select({
                  id: orderTable.id,
                  created: orderTable.timeCreated,
                  printed: orderTable.timePrinted,
                  tracking: orderTable.trackingURL,
                  label: orderTable.labelURL,
                  address: orderTable.shippingAddress,
                  amount: sql<string>`COALESCE(${tx
                    .select({
                      amount: sql<number>`SUM(${orderItemTable.amount})`,
                    })
                    .from(orderItemTable)
                    .where(eq(orderItemTable.orderID, orderTable.id))}, 0)`,
                })
                .from(orderTable)
                .where(
                  input.queryTerm
                    ? or(
                        sql`lower(${orderTable.shippingAddress}->>'$.name') LIKE ${"%" + input.queryTerm.toLowerCase().replaceAll(" ", "%") + "%"}`,
                        like(orderTable.email, "%" + input.queryTerm + "%"),
                      )
                    : undefined,
                )
                .orderBy(desc(orderTable.id))
                .offset(input.offset)
                .limit(input.pageSize),
            }));
          },
          rowMenuItems: (row) =>
            [
              row.label && {
                label: "Label",
                url: row.label!,
              },
              row.tracking && {
                label: "Tracking",
                url: row.tracking!,
              },
            ].filter(Boolean) as any,
          columns: [
            "id",
            {
              label: "amount",
              renderCell: (row) => ({
                label: `$${parseInt(row.amount) / 100}`,
              }),
            },
            {
              label: "name",
              renderCell: (row) => ({
                label: row.address!.name,
              }),
            },
            "created",
            "printed",
          ],
          isSortable: false,
        }),
      ],
    });
  },
  routes: {
    print: new Action({
      name: "Print Orders",
      unlisted: true,
      async handler() {
        await ctx.loading.start({
          label: "Generating labels",
          description: "This may take a few minutes",
        });

        const orders = await useTransaction((tx) =>
          tx
            .select({
              id: orderTable.id,
              label: orderTable.labelURL,
              count: sum(orderItemTable.quantity),
            })
            .from(orderTable)
            .innerJoin(
              orderItemTable,
              eq(orderItemTable.orderID, orderTable.id),
            )
            .where(
              and(
                isNull(orderTable.timePrinted),
                isNotNull(orderTable.labelURL),
              ),
            )
            .groupBy(orderTable.id)
            .orderBy(sum(orderItemTable.id)),
        );
        const grouped = pipe(
          orders,
          groupBy((x) => x.count || "0"),
          entries(),
          map(async ([count, group]) => {
            const mergedPdf = await PDFDocument.create();
            for (const order of group) {
              if (!order.label) continue;
              const response = await fetch(order.label!);
              if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
              }
              const arrayBuffer = await response.arrayBuffer();
              const pdf = await PDFDocument.load(arrayBuffer);
              const copiedPages = await mergedPdf.copyPages(
                pdf,
                pdf.getPageIndices(),
              );
              copiedPages.forEach((page) => mergedPdf.addPage(page));
            }
            const bytes = await mergedPdf.save();
            const s3 = new S3Client({
              region: "us-east-2",
            });
            const key = `labels/${new Date().toISOString()}-${count}.pdf`;
            await s3.send(
              new PutObjectCommand({
                Bucket: Resource.IntervalBucket.name,
                Key: key,
                Body: bytes,
                ContentType: "application/pdf",
              }),
            );
            const command = new GetObjectCommand({
              Bucket: Resource.IntervalBucket.name,
              Key: key,
            });
            const presigned = await getSignedUrl(s3, command, {
              expiresIn: 3600,
            });
            return {
              count,
              label: presigned,
            };
          }),
        );
        const labels = await Promise.all(grouped);
        await io.display.metadata("Order", {
          layout: "list",
          data: [
            {
              label: "Orders",
              value: orders.length,
            },
            ...labels.map((item: any) => ({
              label: `${item.count} count`,
              value: item.label,
            })),
          ],
        });
        const result = await io.confirm("Confirm these orders as printed?");
        if (result) {
          useTransaction((tx) =>
            tx
              .update(orderTable)
              .set({
                timePrinted: sql`CURRENT_TIMESTAMP(3)`,
              })
              .where(
                inArray(
                  orderTable.id,
                  orders.map((x) => x.id),
                ),
              ),
          );
        }
        await ctx.redirect({
          route: "order",
        });
      },
    }),
    create: new Action({
      name: "Create",
      handler: async () => {
        const products = await Product.list();
        const results = await io.group(
          products.flatMap((product) => {
            return product.variants.map((variant) =>
              io.input.number(`${product.name} - ${variant.name}`, {
                defaultValue: 0,
              }),
            );
          }),
        );
        const items = {} as Record<string, number>;
        let total = 0;
        for (const [number, product] of products.entries()) {
          const amount = results[number] as number;
          if (amount === 0) continue;
          total += amount;
          items[product.variants[0]!.id] = amount;
        }
        if (total === 0) return;
        console.log(items);
        const [
          email,
          name,
          street1,
          street2,
          city,
          province,
          zip,
          country,
          phone,
        ] = await io.group([
          io.input.text("Email"),
          io.input.text("Name"),
          io.input.text("Street 1"),
          io.input.text("Street 2").optional(),
          io.input.text("City"),
          io.input.text("State / Province"),
          io.input.text("Zip"),
          io.input.text("Country"),
          io.input.text("Phone").optional(),
        ]);
        await OrderM.createInternal({
          email,
          items,
          address: {
            name,
            street1,
            street2,
            city,
            province,
            zip,
            country,
            phone,
          },
        });
        await ctx.redirect({
          route: "order",
        });
      },
    }),
    shipping: new Action({
      name: "Shipping",
      async handler() {
        const [
          email,
          name,
          street1,
          street2,
          city,
          province,
          zip,
          country,
          phone,
          weight,
        ] = await io.group([
          io.input.text("Email"),
          io.input.text("Name"),
          io.input.text("Street 1"),
          io.input.text("Street 2").optional(),
          io.input.text("City"),
          io.input.text("State / Province"),
          io.input.text("Zip"),
          io.input.text("Country"),
          io.input.text("Phone").optional(),
          io.input.number("Weight (oz)"),
        ]);

        const shipping = await Shippo.createShipmentRate({
          subtotal: 0,
          address: {
            name,
            street1,
            street2,
            city,
            province,
            zip,
            country,
            phone,
          },
          ounces: weight,
        });

        await io.display.object("shipping", { data: shipping });
      },
    }),
  },
});
