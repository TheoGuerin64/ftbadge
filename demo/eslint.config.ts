import skipFormatting from "@vue/eslint-config-prettier/skip-formatting";
import {
  defineConfigWithVueTs,
  vueTsConfigs,
} from "@vue/eslint-config-typescript";
import pluginVue from "eslint-plugin-vue";
import { globalIgnores } from "eslint/config";

export default defineConfigWithVueTs(
  {
    name: "app/files-to-lint",
    files: ["**/*.{ts,vue}"],
  },

  globalIgnores(["**/dist/**"]),

  pluginVue.configs["flat/essential"],
  vueTsConfigs.recommended,
  skipFormatting,
);
