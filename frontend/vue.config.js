module.exports = {
    publicPath: './',
    outputDir: '../backend/public',
    lintOnSave: true,
    devServer: {
        port: 3000,
        disableHostCheck: true,
        proxy: 'http://localhost:8080',
    },
};
