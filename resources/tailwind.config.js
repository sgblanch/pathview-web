module.exports = {
  content: ["./js/*.js", "../template/*.html", "../template/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/forms')({
      strategy: 'class'
    }),
  ],
}
