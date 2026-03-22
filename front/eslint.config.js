import js from "@eslint/js";
import vue from "eslint-plugin-vue";
import tsParser from "@typescript-eslint/parser";
import vueParser from "vue-eslint-parser";
import globals from "globals";
import prettier from "eslint-config-prettier";

export default [
  {
    ignores: ["node_modules", "dist", "coverage"],
  },

  js.configs.recommended,

  ...vue.configs["flat/recommended"],

  // 🌐 FRONTEND (Vue)
  {
    files: ["**/*.vue"],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: tsParser,
        ecmaVersion: "latest",
        sourceType: "module",
      },
      globals: globals.browser,
    },
  },

  // 🌐 FRONT TS
  {
    files: ["src/**/*.ts"],
    languageOptions: {
      parser: tsParser,
      globals: globals.browser,
    },
  },

  // ⚙️ SERVER (Bun / Node)
  {
    files: ["server.ts"],
    languageOptions: {
      parser: tsParser,
      globals: {
        ...globals.node,
        Bun: "readonly",
      },
    },
  },

  {
    rules: {
      "vue/multi-word-component-names": "off",
    },
  },

  prettier,
];
