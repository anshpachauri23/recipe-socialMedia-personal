/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        // Earthy tone palette
        earth: {
          50: '#f7f5f0',
          100: '#ede8dd',
          200: '#d9d0b8',
          300: '#c4b693',
          400: '#b19f6e',
          500: '#9d8a4a',
          600: '#7d6e3a',
          700: '#5d522a',
          800: '#3d361c',
          900: '#1d1a0e',
        },
        sage: {
          50: '#f6f7f4',
          100: '#e8ebe3',
          200: '#d1d7c7',
          300: '#b4c2a5',
          400: '#97ad83',
          500: '#7a9861',
          600: '#5f7a4a',
          700: '#475c37',
          800: '#2f3e25',
          900: '#172012',
        },
        terracotta: {
          50: '#fdf4f1',
          100: '#fbe6e0',
          200: '#f6cdc1',
          300: '#f0b4a2',
          400: '#ea9b83',
          500: '#e48264',
          600: '#b66850',
          700: '#884e3c',
          800: '#5a3428',
          900: '#2c1a14',
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
