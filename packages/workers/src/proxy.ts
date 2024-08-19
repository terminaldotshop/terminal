export default {
  async fetch(req: Request, env: any, ctx: ExecutionContext) {
    const shouldCache = !env.NO_CACHE;
    if (shouldCache) {
      const match = await caches.default.match(req);
      if (match) return match;
    }
    const url = new URL(req.url);
    const target = new URL(env.ORIGIN_URL);
    const response = await fetch(
      target.origin + url.pathname + url.search,
      req,
    );
    if (shouldCache && response.ok) {
      ctx.waitUntil(caches.default.put(req, response.clone()));
    }
    return response;
  },
};
