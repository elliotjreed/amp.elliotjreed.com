const { resolve } = require("path");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const purgecss = require("@fullhuman/postcss-purgecss");

module.exports = {
  name: "client",
  context: resolve(__dirname, "./"),
  entry: {
    client: "./index.ts"
  },
  target: "web",
  output: {
    filename: "[name].[contenthash].min.js",
    path: resolve(__dirname, "./templates/css"),
    publicPath: "/"
  },
  module: {
    rules: [
      {
        exclude: /node_modules/,
        include: resolve(__dirname, "./src/index.ts"),
        test: /\.ts$/,
        loader: "ts-loader",
        options: {
          configFile: "tsconfig.json"
        }
      },
      {
        test: /\.css$/,
        use: [
          MiniCssExtractPlugin.loader,
          { loader: "css-loader", options: { sourceMap: true, url: false } },
          {
            loader: "postcss-loader",
            options: {
              sourceMap: true,
              postcssOptions: {
                plugins: [
                  "tailwindcss",
                  "autoprefixer",
                  "cssnano",
                  purgecss({
                    content: ["./templates/html/*.{html,js,ts}"],
                    defaultExtractor: (content) => content.match(/[\w-/:]+(?<!:)/g) || []
                  })
                ]
              }
            }
          }
        ]
      }
    ]
  },
  performance: {
    hints: "warning"
  },
  plugins: [
    new MiniCssExtractPlugin({
      chunkFilename: "[id].[contenthash].css",
      filename: "[name].css"
    }),
  ],
  resolve: {
    extensions: [".ts"],
    symlinks: false,
    fallback: {
      fs: false,
      tls: false,
      net: false,
      path: false,
      zlib: false,
      http: false,
      https: false,
      stream: false,
      crypto: false,
      url: false,
      util: false
    }
  }
};
