if ($app.stage === "production") {
  const zone = cloudflare.getZoneOutput({
    name: "trm.sh",
  });
  new cloudflare.PageRule("ShortRedirect", {
    zoneId: zone.id,
    target: "trm.sh/*",
    actions: {
      forwardingUrl: {
        url: "https://www.terminal.shop/$1",
        statusCode: 302,
      },
    },
    priority: 1,
  });
}
