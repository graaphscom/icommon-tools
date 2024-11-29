import type { Config } from 'tailwindcss';

export default {
  content: [
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#1A1A21',
        secondary: '#FFCC00',
      },
      fontFamily: {
        condensed: ['var(--font-ubuntu-condensed)'],
        mono: ['var(--font-ubuntu-mono)'],
        sans: ['var(--font-ubuntu)'],
      },
    },
  },
  plugins: [],
} satisfies Config;
