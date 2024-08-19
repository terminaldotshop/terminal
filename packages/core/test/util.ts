import { it } from "bun:test";
import { ActorContext } from "../src/actor";
import { User } from "../src/user";
import { nanoid } from "nanoid/non-secure";

export function withTestUser(name: string, cb: (id: string) => Promise<any>) {
  return it(name, async () => {
    const user = await User.create({
      fingerprint: "test+" + nanoid(),
    });
    await ActorContext.with(
      { type: "user", properties: { userID: user } },
      async () => {
        await cb(user);
      },
    );
  });
}
