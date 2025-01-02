const path = require('path')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const TerserWebpackPlugin = require('terser-webpack-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
  mode: 'production',
  entry: './assets/app.js', // Входной файл
  output: {
    filename: 'build.js',
    path: path.resolve(__dirname, 'build/'), // Выходная папка
  },
  optimization: {
    minimize: true,
    minimizer: [new TerserWebpackPlugin()],
  },
  plugins: [
    new CleanWebpackPlugin(), // Удаляем старую папку с билдом перед каждой новой сборкой
    new MiniCssExtractPlugin({
      filename: '[name].css',
    }),
  ],
  module: {
    rules: [
      {
        test: /\.css$/, // Правило для CSS-файлов
        use: ['style-loader', 'css-loader'], // Загрузчики для CSS
      },
      {
        test: /\.(woff|woff2|eot|ttf|otf)$/,
        loader: 'file-loader',
        options: {
          name: '[name].[ext]',
          outputPath: 'fonts/', // Папка для шрифтов
        },
      },
    ],
  },
};