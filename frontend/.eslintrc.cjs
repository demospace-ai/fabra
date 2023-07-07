/* eslint-env node */
module.exports = {
  extends: [],
  ignorePatterns: ["postcss.config.cjs", "tailwind.config.cjs", "vite-*.ts"],
  rules: {
    "@typescript-eslint/switch-exhaustiveness-check": "error",
    "@typescript-eslint/quotes": ["error", "double"],
    "no-restricted-imports": ["error", {
      "patterns": [".*"]
    }],
  },
  parser: "@typescript-eslint/parser",
  parserOptions: {
    project: ["./tsconfig.json"],
  },
  plugins: ["@typescript-eslint"],
  root: true,
};
