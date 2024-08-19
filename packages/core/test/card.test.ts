import { describe, expect } from "bun:test";
import { withTestUser } from "./util";
import { Card } from "../src/card";

describe("card", () => {
  withTestUser("create", async () => {
    await Card.create({
      token: "tok_visa",
    });
    const list = await Card.list();
    expect(list.length).toBe(1);
    expect(list[0]!.brand).toBe("visa");
    expect(list[0]!.last4).toBe("4242");
  });
});
