import { defineConfig } from 'astro/config'
import tailwind from '@astrojs/tailwind'
import solid from '@astrojs/solid-js'
import aws from 'astro-sst'

export default defineConfig({
  integrations: [tailwind({ applyBaseStyles: false }), solid()],
  server: { host: true },
  adapter: aws(),
  output: 'server',
  redirects: {
    '/report': 'https://www.youtube.com/watch?v=dQw4w9WgXcQ',
  },
})
