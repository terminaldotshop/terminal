import { Layout, Page, io } from "@interval/sdk";
import { useTransaction } from "@terminal/core/drizzle/transaction";
import { userTable } from "@terminal/core/user/user.sql";

export const User = new Page({
  name: "User",
  handler: async (c) => {
    return new Layout({
      title: "User",
      children: [
        io.display.table("", {
          getData: async (input) => {
            return useTransaction(async (tx) => ({
              data: await tx
                .select()
                .from(userTable)
                .offset(input.offset)
                .limit(input.pageSize),
            }));
          },
          isSortable: false,
          isFilterable: false,
        }),
      ],
    });
  },
});
