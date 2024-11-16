/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: '#ba0c0c',
        secondary: '#980a0a',
        tertiary: '#564d4d',
        light: '#ffffff',
        dark: '#0B0000',
        gray: '#acb2b7'
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
