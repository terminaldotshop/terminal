import type { Config } from "tailwindcss"
import defaultTheme from "tailwindcss/defaultTheme"

export default {
  content: ["./src/**/*.{html,js}"],
  theme: {
    extend: {
      lineHeight: {
        normal: "180%",
      },
      letterSpacing: {
        normal: "-0.32px",
      },
      fontFamily: {
        mono: ["geist", "geist-fallback", ...defaultTheme.fontFamily.sans],
      },
      colors: {
        black: "#000000",
        white: "#FFFFFF",
        gray: "rgb(var(--color-gray))",
        muted: "rgb(var(--color-gray) / 40%)",
        orange: "#FF5C00",
      },
      animation: {
        blink: "blink 1.45s infinite step-start",
      },
      keyframes: {
        blink: {
          "0%, 25%, 100%": { opacity: "1" },
          "50%, 75%": { opacity: "0" },
        },
      },
    },
  },
  plugins: [],
} satisfies Config
