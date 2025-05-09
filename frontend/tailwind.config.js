// tailwind.config.js
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",  // Указываем путь для всех файлов, в которых будем использовать Tailwind (React-компоненты)
    "./public/index.html",          // Если используете стандартную структуру React
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
