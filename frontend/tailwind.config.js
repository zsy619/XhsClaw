/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#fff1f2',
          100: '#ffe4e6',
          200: '#fecdd3',
          300: '#fda4af',
          400: '#fb7185',
          500: '#ff2442',
          600: '#e11d48',
          700: '#be123c',
          800: '#9f1239',
          900: '#881337',
        },
        xiaohongshu: {
          red: '#ff2442',
          redLight: '#ff4d64',
          redDark: '#e11d48',
          pink: '#ff4d64',
          dark: '#1a1a1a',
          gray: '#666666',
          grayLight: '#999999',
          grayDark: '#333333',
          light: '#f5f5f5',
          lighter: '#fafafa',
          bg: '#f7f8fa',
        }
      },
      fontFamily: {
        sans: ['-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', 'Arial', 'sans-serif'],
      },
      boxShadow: {
        'xiaohongshu': '0 2px 12px rgba(0, 0, 0, 0.08)',
        'xiaohongshu-lg': '0 4px 20px rgba(0, 0, 0, 0.12)',
        'xiaohongshu-xl': '0 8px 30px rgba(0, 0, 0, 0.15)',
      },
      borderRadius: {
        'xl': '12px',
        '2xl': '16px',
        '3xl': '20px',
      },
      screens: {
        'xs': '360px',
        'sm': '640px',
        'md': '768px',
        'lg': '1024px',
        'xl': '1280px',
        '2xl': '1536px',
      },
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
        '128': '32rem',
      }
    },
  },
  plugins: [],
}
