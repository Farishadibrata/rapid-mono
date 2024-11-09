import fluid, { extract } from 'fluid-tailwind'

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    files: ["./view/**/*.templ"],
    extract
  },
  theme: {
    extend: {
      colors: {
        mariner: {
          '50': '#eff9ff',
          '100': '#dbf1fe',
          '200': '#bfe7fe',
          '300': '#92dafe',
          '400': '#5fc3fb',
          '500': '#3aa6f7',
          '600': '#2489ec',
          '700': '#1c71d8',
          '800': '#1d5cb0',
          '900': '#1d4f8b',
          '950': '#173154',
        },
      },
    },
  },
  plugins: [fluid],
}

