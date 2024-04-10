import { AuthHandler } from "sst/auth";
import { session } from "./session";
import { swell } from "./swell";

export default AuthHandler({
  providers: {
    ssh: (route, ctx) => {
      route.post("/login", async (c) => {
        const fingerprint = await c.req
          .formData()
          .then((data) => data.get("fingerprint"));

        if (!fingerprint) {
          return c.json({ error: "Fingerprint is required" }, 400);
        }

        console.log("searching for user with fingerprint", fingerprint);
        const search = await swell(
          "/accounts?where[fingerprint]=" + fingerprint,
          {},
        );
        let userID = search.results[0]?.id;
        if (!userID) {
          console.log("creating user");
          userID = await swell("/accounts", {
            method: "POST",
            headers: {
              "content-type": "application/x-www-form-urlencoded",
            },
            body: new URLSearchParams({
              email: "ssh+" + fingerprint + "@terminal.shop",
              fingerprint,
            }).toString(),
          }).then((res) => res.id);
        }
        return c.json({
          accessToken: await session.create({
            type: "user",
            properties: { userID },
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
