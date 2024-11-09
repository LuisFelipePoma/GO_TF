/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      colors: {
        primary: 'var(--color-primary)',
        secondary: 'var(--color-secondary)',
        tertiary: 'var(--color-tertiary)',
        light: 'var(--color-light)',
        dark: 'var(--color-dark)',
        gray: 'var(--color-gray)'
      },
      fontSize: {
        'body-20': 'var(--font-size-20)',
        'body-16': 'var(--font-size-16)',
        'body-14': 'var(--font-size-14)',
        'body-12': 'var(--font-size-12)'
      }
    }
  },
  plugins: [
    function ({ addUtilities }) {
      addUtilities({
        '.no-scrollbar': {
          '-ms-overflow-style': 'none', // IE and Edge
          'scrollbar-width': 'none' // Firefox
        },
        '.no-scrollbar::-webkit-scrollbar': {
          display: 'none' // Chrome, Safari, and Opera
        }
      })
    }
  ]
}
