import { Resource } from "sst";
import crypto from "crypto";
import { useTransaction } from "./drizzle/transaction";
import { orderItemTable, orderTable } from "./order/order.sql";
import { eq } from "drizzle-orm";
import { productTable, productVariantTable } from "./product/product.sql";
import { userTable } from "./user/user.sql";

const CUSTOMERS_LIST = "3a96a118-0250-11ef-a951-f72c9a1b1995";
const SUBSCRIBERS_LIST = "e4d1931a-019a-11ef-8784-674ae09367af";

export module EmailOctopus {
  export async function find(props: {
    email: string;
    list?: "customers" | "subscribers";
  }) {
    const listId =
      props.list === "customers" ? CUSTOMERS_LIST : SUBSCRIBERS_LIST;
    const hashedEmail = crypto
      .createHash("md5")
      .update(props.email.toLowerCase())
      .digest("hex");

    return fetch(
      `https://emailoctopus.com/api/1.6/lists/${listId}/contacts/${hashedEmail}?api_key=${Resource.EmailOctopusSecret.value}`,
    ).then((res) => res.json() as unknown as { id: string | null });
  }

  export async function subscribe(props: {
    email: string;
    fields?: unknown;
    tags?: string[];
    list?: "customers" | "subscribers";
  }) {
    const listId =
      props.list === "customers" ? CUSTOMERS_LIST : SUBSCRIBERS_LIST;
    const response = await fetch(
      `https://emailoctopus.com/api/1.6/lists/${listId}/contacts`,
      {
        method: "POST",
        body: JSON.stringify({
          api_key: Resource.EmailOctopusSecret.value,
          email_address: props.email,
          fields: props.fields,
          tags: props.tags,
        }),
        headers: { "Content-Type": "application/json" },
      },
    ).then(
      (res) =>
        res.json() as unknown as {
          id: string | null;
          error?: { code: string };
        },
    );

    console.log("EmailOctopus response", response);

    if (response.error?.code === "MEMBER_EXISTS_WITH_EMAIL_ADDRESS") {
      console.log(
        "EmailOctopus contact not found, looking up by email: " + props.email,
      );
      const contact = await find({ email: props.email, list: props.list });
      console.log("EmailOctopus contact found?", contact);
      return contact;
    }

    return response;
  }

  export async function update(props: {
    memberId: string;
    email?: string;
    fields?: unknown;
    tags?: Record<string, boolean>;
    list: "customers" | "subscribers";
  }) {
    const listId =
      props.list === "customers" ? CUSTOMERS_LIST : SUBSCRIBERS_LIST;
    return fetch(
      `https://emailoctopus.com/api/1.6/lists/${listId}/contacts/${props.memberId}`,
      {
        method: "PUT",
        body: JSON.stringify({
          api_key: Resource.EmailOctopusSecret.value,
          email_address: props.email,
          fields: props.fields,
          tags: props.tags,
        }),
        headers: { "Content-Type": "application/json" },
      },
    ).then((res) => res.json());
  }

  export async function addToCustomersList(orderID: string) {
    const items = await useTransaction((tx) =>
      tx
        .select({
          email: userTable.email,
          emailOctopusID: userTable.emailOctopusID,
          productName: productTable.name,
          quantity: orderItemTable.quantity,
        })
        .from(orderTable)
        .innerJoin(orderItemTable, eq(orderItemTable.orderID, orderTable.id))
        .innerJoin(userTable, eq(userTable.id, orderTable.userID))
        .innerJoin(
          productVariantTable,
          eq(productVariantTable.id, orderItemTable.productVariantID),
        )
        .innerJoin(
          productTable,
          eq(productTable.id, productVariantTable.productID),
        )
        .where(eq(orderTable.id, orderID))
        .execute(),
    );
    if (!items?.length) {
      console.error("No order found: " + orderID);
      return;
    }

    const order = items[0];
    if (order?.emailOctopusID) {
      await update({
        memberId: order.emailOctopusID,
        list: "customers",
        tags: { purchased: true },
      });
    }
  }

  export async function addToMarketingList(userID: string) {
    const user = await useTransaction((tx) =>
      tx
        .select({
          email: userTable.email,
          emailOctopusID: userTable.emailOctopusID,
        })
        .from(userTable)
        .where(eq(userTable.id, userID))
        .execute()
        .then((r) => r.at(0)),
    );
    if (!user) {
      console.error("No user found: " + userID);
      return;
    }

    if (!user.emailOctopusID && user.email) {
      const contact = await EmailOctopus.subscribe({
        email: user.email,
        list: "subscribers",
      });
      await useTransaction((tx) =>
        tx
          .update(userTable)
          .set({ emailOctopusID: contact.id })
          .where(eq(userTable.id, userID)),
      );
    }
  }
}
