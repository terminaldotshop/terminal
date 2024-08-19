import { z } from "zod";

export const Address = z.object({
  name: z.string(),
  street1: z.string(),
  street2: z.string().optional(),
  city: z.string(),
  province: z.string().optional(),
  country: z
    .string()
    .length(2, "Country must be a 2 character country code (ISO 3166-1)"),
  zip: z.string(),
  phone: z.string().optional(),
});

export type Address = z.infer<typeof Address>;
