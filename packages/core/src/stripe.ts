import { Resource } from "sst";
import { Stripe as StripeClient } from "stripe";
import { z } from "zod";
import { fn } from "./util/fn";
import { useTransaction } from "./drizzle/transaction";
import { userTable } from "./user/user.sql";
import { eq } from "drizzle-orm";

export const stripe = new StripeClient(Resource.StripeSecret.value, {
  httpClient: StripeClient.createFetchHttpClient(),
});

export module Stripe {
  export const client = stripe;

  export const syncUser = fn(z.string(), (id) =>
    useTransaction(async (tx) => {
      const user = await tx
        .select({
          stripeCustomerID: userTable.stripeCustomerID,
          email: userTable.email,
          name: userTable.name,
        })
        .from(userTable)
        .where(eq(userTable.id, id))
        .then((rows) => rows.at(0));
      if (!user) throw new Error("User not found");
      await stripe.customers.update(user.stripeCustomerID, {
        email: user.email || undefined,
        name: user.name || undefined,
        metadata: {
          userID: id,
          fingerprint: "",
        },
      });
    }),
  );
}
