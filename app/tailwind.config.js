/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'construction-blue': '#2563eb',
        'construction-green': '#16a34a',
        'construction-orange': '#ea580c',
      },
    },
  },
  plugins: [],
}
