import { z } from "@hono/zod-openapi";

export function Result<T extends z.ZodTypeAny>(schema: T) {
  return z.object({
    result: schema,
  });
}
