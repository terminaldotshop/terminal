import { eq } from "drizzle-orm";
import { createContext } from "./context";
import { useTransaction } from "./drizzle/transaction";
import { UserFlags, userTable } from "./user/user.sql";
import { VisibleError } from "./error";

export interface UserActor {
  type: "user";
  properties: {
    userID: string;
  };
}

export interface PublicActor {
  type: "public";
  properties: {};
}

type Actor = UserActor | PublicActor;
export const ActorContext = createContext<Actor>();

export function useUserID() {
  const actor = ActorContext.use();
  if (actor.type === "user") return actor.properties.userID;
  throw new Error(`Actor is "${actor.type}" not UserActor`);
}

export async function assertFlag(flag: keyof UserFlags) {
  return useTransaction((tx) =>
    tx
      .select({ flags: userTable.flags })
      .from(userTable)
      .where(eq(userTable.id, useUserID()))
      .then((rows) => {
        const flags = rows[0]?.flags;
        if (!flags)
          throw new VisibleError(
            "auth",
            "user.flags",
            "Actor does not have " + flag + " flag",
          );
      }),
  );
}

export function useActor() {
  try {
    return ActorContext.use();
  } catch {
    return { type: "public", properties: {} } as PublicActor;
  }
}

export function assertActor<T extends Actor["type"]>(type: T) {
  const actor = useActor();
  if (actor.type !== type)
    throw new VisibleError("auth", "actor.invalid", `Actor is not "${type}"`);
  return actor as Extract<Actor, { type: T }>;
}
