import { Hono } from "hono";
import { logger } from "hono/logger";
import { session } from "./session";
import { createContext } from "./context";

const SessionContext = createContext<typeof session.$typeValues>();

function useUserID() {
  const session = SessionContext.use();
  if (session.type !== "user") throw new Error("User session expected");
  return session.properties.userID;
}

const app = new Hono();
app.use(logger());
app.use(async (c, next) => {
  const authHeader = c.req.header("authorization");
  if (!authHeader) {
    return c.json({ error: "Authorization header is missing" }, 401);
  }

  const match = authHeader.match(/^Bearer (.+)$/);
  if (!match) {
    return c.json(
      { error: "Bearer token not found or improperly formatted" },
      401,
    );
  }

  const bearerToken = match[1];
  const result = await session.verify(bearerToken);
  return SessionContext.with(result, next);
});

app.get("/user/me", async (c) => {
  const userID = useUserID();
});

export default app;
