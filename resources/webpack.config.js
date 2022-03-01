const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const TerserPlugin = require("terser-webpack-plugin");

module.exports = {
    mode: 'development',
    entry: './js/pathview.js',
    devtool: 'inline-source-map',
    optimization: {
        minimize: true,
        minimizer: [new TerserPlugin()],
    },
    output: {
        filename: 'pathview.js',
        path: path.resolve(__dirname, '../static'),
    },
    plugins: [new MiniCssExtractPlugin({
        filename: "pathview.css",
    })],
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: [MiniCssExtractPlugin.loader, "css-loader", "postcss-loader"],
            },
        ],
    },
};

