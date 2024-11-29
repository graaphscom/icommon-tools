import vercelPrettier from '@vercel/style-guide/prettier';

/**
 * @type {import("prettier").Config}
 */
const config = {
  ...vercelPrettier,
  plugins: [...vercelPrettier.plugins, 'prettier-plugin-tailwindcss'],
};

export default config;
