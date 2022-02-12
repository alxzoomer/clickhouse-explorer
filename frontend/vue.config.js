// vue.config.js

/**
 * @type {import('@vue/cli-service').ProjectOptions}
 */
module.exports = {
  devServer: {
    port: 3001,
    proxy: 'http://localhost:8000',
  },

  pluginOptions: {
    quasar: {
      importStrategy: 'kebab',
      rtlSupport: false,
    },
  },

  transpileDependencies: [
    'quasar',
  ],

  pages: {
    index: {
      entry: 'src/main.ts',
      title: 'ClickHouse Query Executor',
    },
  },
};
