import { z } from "zod";
import { fn } from "../util/fn";
import { inventoryTable, inventoryRecordTable } from "./inventory.sql";
import { useTransaction } from "../drizzle/transaction";
import { createID } from "../util/id";

export module Inventory {
  export const create = fn(
    z.object({
      name: z.string(),
      description: z.string().optional(),
    }),
    async (input) =>
      useTransaction(async (tx) =>
        tx.insert(inventoryTable).values({
          id: createID("inventory"),
          name: input.name,
          description: input.description || null,
        }),
      ),
  );

  export const record = fn(
    z.object({
      inventoryID: z.string(),
      quantity: z.number(),
      notes: z.string().optional(),
    }),
    async (input) =>
      useTransaction(async (tx) =>
        tx.insert(inventoryRecordTable).values({
          id: createID("inventoryRecord"),
          quantity: input.quantity,
          inventoryID: input.inventoryID,
          notes: input.notes || null,
        }),
      ),
  );
}
