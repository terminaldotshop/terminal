import { Resource } from "sst";

export async function shippo(url: string, options: RequestInit) {
  const headers = new Headers(options.headers || {});
  headers.set("authorization", "ShippoToken " + Resource.ShippoSecret.value);
  headers.set("Content-Type", "application/json");
  return fetch("https://api.goshippo.com" + url, {
    ...options,
    headers,
  }).then((res) => res.json() as any);
}
