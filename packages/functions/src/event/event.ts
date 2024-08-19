import { bus } from "sst/aws/bus";
import { Order } from "@terminal/core/order/order";
import { Shippo } from "@terminal/core/shippo/index";
import { User } from "@terminal/core/user/index";
import { Stripe } from "@terminal/core/stripe";
import { Template } from "@terminal/core/email/template";
import { EmailOctopus } from "@terminal/core/email-octopus";

export const handler = bus.subscriber(
  [Order.Event.Created, User.Events.Updated],
  async (event) => {
    console.log(event.type, event.properties, event.metadata);
    switch (event.type) {
      case "order.created": {
        await Shippo.createShipment(event.properties.orderID);
        await Template.sendOrderConfirmation(event.properties.orderID);
        await EmailOctopus.addToCustomersList(event.properties.orderID);
        break;
      }
      case "user.updated": {
        await Stripe.syncUser(event.properties.userID);
        await EmailOctopus.addToMarketingList(event.properties.userID);
        break;
      }
    }
  },
);
