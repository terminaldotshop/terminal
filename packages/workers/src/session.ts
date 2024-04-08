import { createSessionBuilder } from "sst/auth";

export const session = createSessionBuilder<{
  user: {
    userID: string;
  };
}>();
