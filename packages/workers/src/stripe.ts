import { Resource } from "sst";
import { Stripe } from "stripe";

export function stripe() {
  return new Stripe(Resource.StripeSecret.value, {
    httpClient: Stripe.createFetchHttpClient(),
  });
}
