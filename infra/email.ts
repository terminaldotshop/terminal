import { domain, shortDomain } from "./dns";

export const email = new sst.aws.Email("Email", {
  sender: domain,
  dns: sst.cloudflare.dns(),
});

export const shortDomainEmail = new sst.aws.Email("ShortDomainEmail", {
  sender: shortDomain,
  dns: sst.cloudflare.dns(),
});
