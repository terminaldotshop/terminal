/// <reference path="./.sst/platform/config.d.ts" />
import { resolve } from "path";
export default $config({
  app(input) {
    return {
      name: "terminal-shop",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "cloudflare",
      providers: {
        aws: {
          profile: process.env.GITHUB_ACTIONS
            ? undefined
            : input.stage === "production"
              ? "terminal-production"
              : "terminal-dev",
        },
        random: true,
        docker: true,
        tls: true,
      },
    };
  },
  async run() {
    const isPermanentStage =
      $app.stage === "production" || $app.stage === "dev";
    const domain =
      $app.stage === "production"
        ? "terminal.shop"
        : $app.stage + ".dev.terminal.shop";

    const secrets = {
      SwellSecret: new sst.Secret("SwellSecret"),
      AirtableSecret: new sst.Secret("AirtableSecret"),
      StripeSecret: new sst.Secret("StripeSecret"),
      ShippoSecret: new sst.Secret("ShippoSecret"),
    };
    const auth = new sst.cloudflare.Auth("Auth", {
      authenticator: {
        link: [secrets.SwellSecret, secrets.StripeSecret],
        handler: "./packages/workers/src/auth.ts",
        domain: "auth." + domain,
      },
    });
    const api = new sst.cloudflare.Worker("Api", {
      handler: "./packages/workers/src/api.ts",
      link: [
        secrets.SwellSecret,
        secrets.AirtableSecret,
        secrets.StripeSecret,
        secrets.ShippoSecret,
        auth,
      ],
      domain: "api." + domain,
    });
    const www = new sst.cloudflare.StaticSite("Www", {
      domain: "www." + domain,
      path: "./packages/www",
      environment: {
        PUBLIC_API_URL: api.url.apply((u) => u!),
      },
      build: {
        command: "bun run build",
        output: "./dist",
      },
    });
    if (isPermanentStage) {
      const github = new aws.iam.OpenIdConnectProvider("GithubOidc", {
        url: "https://token.actions.githubusercontent.com",
        clientIdLists: ["sts.amazonaws.com"],
        thumbprintLists: [
          "6938fd4d98bab03faadb97b34396831e3780aea1",
          "1c58a3a8518e8759bf075b76b750d4f2df264fcd",
        ],
      });
      const githubRole = new aws.iam.Role("GithubRole", {
        name: [$app.name, $app.stage, "github"].join("-"),
        assumeRolePolicy: {
          Version: "2012-10-17",
          Statement: [
            {
              Effect: "Allow",
              Principal: {
                Federated: github.arn,
              },
              Action: "sts:AssumeRoleWithWebIdentity",
              Condition: {
                StringLike: github.url.apply((url) => ({
                  [`${url}:sub`]: "repo:terminaldotshop/terminal:*",
                })),
              },
            },
          ],
        },
      });
      new aws.iam.RolePolicyAttachment("GithubRolePolicy", {
        policyArn: "arn:aws:iam::aws:policy/AdministratorAccess",
        role: githubRole.name,
      });
    }
    const vpc = new sst.aws.Vpc("Vpc");
    const cluster = new sst.aws.Cluster("Cluster", {
      vpc,
    });
    const ssh = cluster.addService({
      name: "SSH",
      cpu: "2 vCPU",
      memory: "4 GB",
      image: {
        context: "./go",
      },
      public: {
        ports: [
          { listen: "22/tcp", forward: "2222/tcp" },
          { listen: "80/tcp", forward: "8000/tcp" },
        ],
      },
      scaling: {
        min: 2,
        max: 10,
      },
    });
    // new cloudflare.SpectrumApplication("SpectrumSSH", {
    //   dns: {
    //     name: domain,
    //     type: "CNAME",
    //   },
    //   zoneId: sst.cloudflare.DEFAULT_ACCOUNT_ID,
    //   protocol: "ssh",
    //   originDns: {
    //     name: ssh.url,
    //   },
    //   originPort: 22,
    // });

    return {
      api: api.url,
      auth: auth.url,
    };
  },
});
