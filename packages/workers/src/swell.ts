import { Resource } from "sst";

export async function swell(url: string, options: RequestInit) {
  const headers = new Headers(options.headers || {});
  headers.set(
    "authorization",
    "Basic " + btoa("terminal:" + Resource.SwellSecret.value),
  );
  return fetch("https://api.swell.store" + url, {
    ...options,
    headers,
  }).then((res) => res.json() as any);
}
