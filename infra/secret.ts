export const secret = {
  AirtableSecret: new sst.Secret("AirtableSecret"),
  StripeSecret: new sst.Secret("StripeSecret", process.env.STRIPE_API_KEY),
  StripePublic: new sst.Secret("StripePublic"),
  ShippoSecret: new sst.Secret("ShippoSecret"),
  EmailOctopusSecret: new sst.Secret("EmailOctopusSecret"),
  IntervalKey: new sst.Secret("IntervalKey"),
  GithubClientID: new sst.Secret("GithubClientID"),
  GithubClientSecret: new sst.Secret("GithubClientSecret"),
};

export const allSecrets = Object.values(secret);
