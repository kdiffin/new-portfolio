const plugin = require('tailwindcss/plugin')

module.exports = {
  content: [
    "./ui/html/**/*.tmpl"
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          "system-ui",
          "-apple-system",
          "BlinkMacSystemFont",
          "Segoe UI",
          "Roboto",
          "Helvetica Neue",
          "Arial",
          "Noto Sans",
          "Apple Color Emoji",
          "Segoe UI Emoji",
          "Segoe UI Symbol"
        ]
      },
      colors: {
        primary: {
          DEFAULT: '#111111',
          light: '#ffffff'
        },
        secondary: {
          DEFAULT: '#1a1a1a',
          light: '#f9f9f9'
        },
        body: {
          DEFAULT: '#e0e0e0',
          light: '#111111'
        },
        heading: {
          DEFAULT: '#ffffff',
          light: '#000000'
        },
        muted: {
          DEFAULT: '#777777',
          light: '#888888'
        },
        link: {
          DEFAULT: '#cccccc',
          light: '#333333'
        },
        nav: {
          DEFAULT: '#aaaaaa',
          light: '#555555'
        },
        border: {
          DEFAULT: '#333333',
          light: '#dddddd',
          subtle: {
             DEFAULT: '#444444',
             hover: '#888888',
             light: '#cccccc'
          }
        }
      }
    }
  },
  darkMode: "class",
  plugins: [
    plugin(function({ addVariant }) {
      addVariant('light', ['html.light &', 'html.light&'])
    })
  ]
};
