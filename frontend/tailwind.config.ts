import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        background: "#0f172a",
        foreground: "#f1f5f9",
        accent: "#06b6d4",
        "accent-dark": "#0891b2",
      },
      fontFamily: {
        mono: ['"JetBrains Mono"', "monospace"],
      },
    },
  },
  plugins: [],
};
export default config;