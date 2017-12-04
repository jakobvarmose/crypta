module.exports = {
  root: true,
  parser: 'babel-eslint',
  parserOptions: {
    sourceType: 'module'
  },
  env: {
    browser: true
  },
  extends: [
    'airbnb-base',
  ],
  // required to lint *.vue files
  plugins: [
    'import',
    'html',
  ],
  globals: {
    'cordova': true,
    'DEV': true,
    'PROD': true,
    '__THEME': true
  },
  // add your custom rules here
  rules: {
    'no-underscore-dangle': 0,
    // allow debugger during development
    'no-debugger': process.env.NODE_ENV === 'production' ? 2 : 0,
    //'import/no-unresolved': 0,
    //'import/extensions': 0,
    'arrow-body-style': 0,
    'no-param-reassign': 0,
  }
}
