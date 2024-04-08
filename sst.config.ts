/// <reference path="./.sst/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "terminal-shop",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "cloudflare",
    };
  },
  async run() {
    const domain =
      $app.stage === "production"
        ? "terminal.shop"
        : $app.stage + ".dev.terminal.shop";

    const secrets = {
      SwellSecret: new sst.Secret("SwellSecret"),
    };

    const auth = new sst.cloudflare.Auth("Auth", {
      authenticator: {
        link: [secrets.SwellSecret],
        handler: "./packages/workers/src/auth.ts",
        domain: "auth." + domain,
      },
    });

    const api = new sst.cloudflare.Worker("Api", {
      handler: "./packages/workers/src/api.ts",
      link: [secrets.SwellSecret, auth],
      domain: "api." + domain,
    });

    const www = new sst.cloudflare.StaticSite("Www", {
      domain: "www." + domain,
      path: "./packages/www",
      build: {
        command: "bun run build",
        output: "./dist",
      },
    });

    return {
      api: api.url,
      auth: auth.url,
      www: www.url,
    };
  },
});
