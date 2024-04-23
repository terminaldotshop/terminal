import { AuthHandler } from "sst/auth";
import { session } from "./session";
import { stripe } from "./stripe";

export default AuthHandler({
  providers: {
    ssh: (route, ctx) => {
      route.post("/login", async (c) => {
        const fingerprint = await c.req.json().then((x) => x["fingerprint"]);
        if (!fingerprint) {
          return c.json({ error: "Fingerprint is required" }, 400);
        }

        console.log("searching for user with fingerprint", fingerprint);
        const search = await stripe().customers.search({
          query: `metadata["fingerprint"]:"${fingerprint}"`,
        });
        let user = search.data[0];
        if (!user) {
          console.log("creating user");
          user = await stripe().customers.create({
            metadata: { fingerprint },
          });
        }
        return c.json({
          userID: user.id,
          email: user.email || undefined,
          accessToken: await session.create({
            type: "user",
            properties: { userID: user.id },
          }),
        });
      });
    },
  },
  session,
  callbacks: {
    auth: {
      async allowClient() {
        return true;
      },
      async success(ctx, input) {
        return ctx.session({
          type: "user",
          properties: {
            userID: input.fingerprint,
          },
        });
      },
    },
  },
});
