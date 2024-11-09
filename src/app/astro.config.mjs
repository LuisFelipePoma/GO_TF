// @ts-check
import { defineConfig, envField } from 'astro/config'

import react from '@astrojs/react'
import tailwind from '@astrojs/tailwind'
import node from '@astrojs/node'

// https://astro.build/config
export default defineConfig({
  output: 'server',
  adapter: node({ mode: 'standalone' }),
  integrations: [react(), tailwind()],
  experimental: {
    env: {
      schema: {
        API_TMDB: envField.string({
          context: 'server',
          access: 'secret'
        })
      }
    }
  }
})
