import { SESv2Client, SendEmailCommand } from "@aws-sdk/client-sesv2";
import { Resource } from "sst";

export module Email {
  const ses = new SESv2Client({});

  export async function send(
    from: string,
    to: string,
    subject: string,
    body: string,
  ) {
    from = from + "@" + Resource.ShortDomainEmail.sender;
    console.log("sending email", subject, from, to);
    await ses.send(
      new SendEmailCommand({
        Destination: {
          ToAddresses: [to],
        },
        Content: {
          Simple: {
            Body: {
              Text: {
                Data: body,
              },
            },
            Subject: {
              Data: subject,
            },
          },
        },
        FromEmailAddress: `Terminal <${from}>`,
      }),
    );
  }
}
