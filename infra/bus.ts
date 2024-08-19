import { database } from "./database";
import { email, shortDomainEmail } from "./email";
import { allSecrets } from "./secret";

export const bus = new sst.aws.Bus("Bus");

bus.subscribe({
  handler: "./packages/functions/src/event/event.handler",
  link: [database, email, shortDomainEmail, ...allSecrets],
  timeout: "5 minutes",
  permissions: [
    {
      actions: ["ses:SendEmail"],
      resources: ["*"],
    },
  ],
});
