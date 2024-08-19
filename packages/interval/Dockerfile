FROM oven/bun:1.1.24

WORKDIR /app
RUN mkdir packages
COPY package.json ./
COPY packages/*/package.json ./packages/
COPY packages ./packages
RUN bun install
WORKDIR /app/packages/interval
CMD bun run ./src/interval.ts
