const PrefreshPlugin = require("@prefresh/webpack");

module.exports = {
  style: {
    postcss: {
      plugins: [require("tailwindcss"), require("autoprefixer")],
    },
  },

  webpack: {
    configure: {
      module: {
        rules: [
          {
            type: "javascript/auto",
            test: /\.mjs$/,
            include: /node_modules/,
          },
        ],
      },
    },
    alias: {
      react: "preact/compat",
      "react-dom": "preact/compat",
    },
    plugins: {
      add: [new PrefreshPlugin()],
    },
  },
};
