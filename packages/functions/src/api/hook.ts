import { Card } from "@terminal/core/card/index";
import { stripe } from "@terminal/core/stripe";
import { Hono } from "hono";
import { Resource } from "sst";

export module Hook {
  export const route = new Hono().post("/stripe", async (ctx) => {
    const sig = ctx.req.header("stripe-signature");
    console.log({
      sig,
      secret: Resource.StripeWebhook.secret,
      id: Resource.StripeWebhook.id,
    });
    const evt = await stripe.webhooks.constructEventAsync(
      await ctx.req.text(),
      sig!,
      Resource.StripeWebhook.secret,
    );
    console.log(evt);
    switch (evt.type) {
      case "payment_method.updated":
      case "payment_method.attached":
      case "payment_method.detached":
        if (
          evt.data.object.customer &&
          typeof evt.data.object.customer == "string"
        )
          await Card.sync(evt.data.object.customer);
    }
    return ctx.json({});
  });
}
