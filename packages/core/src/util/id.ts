import { ulid } from "ulid";

const prefixes = {
  user: "usr",
  inventory: "inv",
  inventoryRecord: "irc",
  userShipping: "shp",
  card: "crd",
  product: "prd",
  productVariant: "var",
  cartItem: "itm",
  cart: "crt",
  order: "ord",
  subscription: "sub",
} as const;

export function createID(prefix: keyof typeof prefixes): string {
  return [prefixes[prefix], ulid()].join("_");
}

const _typeCheckUniquePrefixes: ErrorIfDuplicates<typeof prefixes> = prefixes;

type Duplicate<T extends Record<string, string>> = {
  [K in keyof T]: {
    [L in keyof T]: T[K] extends T[L] ? (K extends L ? never : K) : never;
  }[keyof T];
}[keyof T];

type HasDuplicates<T extends Record<string, string>> =
  Duplicate<T> extends never ? false : true;

type ErrorIfDuplicates<T extends Record<string, string>> =
  HasDuplicates<T> extends true
    ? { error: `Duplicate values found: ${Duplicate<T> & string}` }
    : T;
