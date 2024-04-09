import type { Config } from 'tailwindcss'
import defaultTheme from 'tailwindcss/defaultTheme'

export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      lineHeight: {
        normal: '180%',
      },
      letterSpacing: {
        normal: '-0.32px',
      },
      fontFamily: {
        mono: ['geist', 'geist-fallback', ...defaultTheme.fontFamily.sans],
      },
      colors: {
        black: '#000000',
        white: '#FFFFFF',
        orange: 'hsla(22, 100%, 50%, 1)',
        slate: {
          1: 'hsla(200, 7%, 9%, 1)',
          5: 'hsla(200, 88%, 93%, 0.11)',
          7: 'hsla(203, 6%, 24%, 1)',
          10: 'hsla(210, 100%, 95%, 0.47)',
          11: 'hsla(210, 100%, 97%, 0.62)',
        },
      },
      animation: {
        blink: 'blink 1.45s infinite step-start',
      },
      keyframes: {
        blink: {
          '0%, 25%, 100%': { opacity: '1' },
          '50%, 75%': { opacity: '0' },
        },
      },
      borderColor: {
        DEFAULT: 'hsla(203, 6%, 24%, 1)',
      },
    },
  },
  plugins: [],
} satisfies Config
