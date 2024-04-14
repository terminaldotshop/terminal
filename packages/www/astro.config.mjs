import { defineConfig } from 'astro/config'
import tailwind from '@astrojs/tailwind'
import solid from '@astrojs/solid-js'

export default defineConfig({
  integrations: [tailwind({ applyBaseStyles: false }), solid()],
  server: { host: true },
})
