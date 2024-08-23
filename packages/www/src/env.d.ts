/// <reference path="../.astro/types.d.ts" />
/// <reference types="astro/client" />

interface ImportMetaEnv {
  readonly PUBLIC_API_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
