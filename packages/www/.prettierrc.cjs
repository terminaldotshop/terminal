// @ts-check

module.exports = {
  astroAllowShorthand: true,
  htmlWhitespaceSensitivity: 'ignore',
  semi: false,
  singleQuote: true,
  plugins: [require.resolve('prettier-plugin-astro')],
  overrides: [
    {
      files: '*.astro',
      options: {
        parser: 'astro',
      },
    },
  ],
}
