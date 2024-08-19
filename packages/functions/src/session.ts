import { auth } from "sst/auth";

export const session = auth.sessions<{
  user: {
    userID: string;
  };
}>();
