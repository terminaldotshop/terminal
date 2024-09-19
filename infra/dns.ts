export const domain =
  {
    production: "terminal.shop",
    dev: "dev.terminal.shop",
  }[$app.stage] || $app.stage + ".dev.terminal.shop";
export const zone = cloudflare.getZoneOutput({
  name: "terminal.shop",
});

export const shortDomain = domain.replace(/terminal\.shop$/, "trm.sh");

export const shortZone = cloudflare.getZoneOutput({
  name: "trm.sh",
});
