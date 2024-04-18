import type { Config } from 'tailwindcss'
import defaultTheme from 'tailwindcss/defaultTheme'

export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue,svg}'],
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
        highlight: {
          1: 'hsla(0, 0%, 100%, 0.1)',
          2: 'hsla(0, 0%, 100%, 0.15)',
        },
        orange: 'hsla(22, 100%, 50%, 1)',
        gray: {
          1: 'hsla(200, 7%, 9%, 1)',
          5: 'hsla(200, 88%, 93%, 0.11)',
          6: 'hsla(209, 94%, 94%, 0.14)',
          7: 'hsla(203, 6%, 24%, 1)',
          10: 'hsla(210, 100%, 95%, 0.47)',
          11: 'hsla(210, 100%, 97%, 0.62)',
          12: 'hsla(210, 100%, 100%, 0.93)',
        },
        light: {
          8: 'hsla(204, 96%, 10%, 0.24)',
          10: 'hsla(204, 100%, 5%, 0.51)',
          12: 'hsla(202, 24%, 9%, 1)',
        },
        green: {
          5: 'hsla(173, 100%, 50%, 0.14)',
          11: 'hsla(167, 70%, 48%, 1)',
        },
        red: {
          5: 'hsla(5, 99%, 56%, 0.33)',
          11: 'hsl(12, 100%, 75%)',
        },
        blue: {
          5: 'hsla(225, 98%, 62%, 0.42)',
          11: 'hsl(228, 100%, 81%)',
        },
      },
      animation: {
        blink: 'blink 1.45s infinite step-start',
        shake: 'shake 0.52s cubic-bezier(.36,.07,.19,.97) both',
      },
      keyframes: {
        blink: {
          '0%, 25%, 100%': { opacity: '1' },
          '50%, 75%': { opacity: '0' },
        },
        shake: {
          '0%': {
            transform: 'translateX(0)',
          },
          '6.5%': {
            transform: 'translateX(-4px) rotateY(-9deg)',
          },
          '18.5%': {
            transform: 'translateX(3px) rotateY(7deg)',
          },
          '31.5%': {
            transform: 'translateX(-1px) rotateY(-5deg)',
          },
          '43.5%': {
            transform: 'translateX(2px) rotateY(3deg)',
          },
          '50%': {
            transform: 'translateX(0)',
          },
        },
      },
      borderColor: {
        DEFAULT: 'hsla(203, 6%, 24%, 1)',
      },
    },
  },
  plugins: [],
} satisfies Config
