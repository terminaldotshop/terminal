import { describe, it, expect } from "bun:test";
import { User } from "../src/user";
import { ActorContext } from "../src/actor";
import { withTestUser } from "./util";

describe("user", () => {
  it("create", async () => {
    const user = await User.create({
      fingerprint: "test",
    });
    expect(await User.fromID(user)).toBeDefined();
    expect(await User.fromFingerprint("test")).toBeDefined();
  });

  withTestUser("add shipping", async (id) => {
    await User.addShipping({
      name: "John Smith",
      zip: "33133",
      city: "Miami",
      country: "US",
      street1: "2800 SW 28th Terrace",
      province: "FL",
    });
    const result = await User.shipping();
    expect(result).toHaveLength(1);
  });

  withTestUser("remove shipping", async () => {
    const address = await User.addShipping({
      name: "John Smith",
      zip: "33133",
      city: "Miami",
      country: "US",
      street1: "2800 SW 28th Terrace",
      province: "FL",
    });
    await User.removeShipping(address);
    const all = await User.shipping();
    expect(all).toHaveLength(0);
  });
});
