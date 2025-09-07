<script setup lang="ts">
import { computed } from "vue";

import CopyButton from "~/components/CopyButton.vue";
import { currentUrl } from "~/lib/utils";

const props = defineProps<{
  apiUrl: string;
}>();

const code = computed<string>(() => {
  return `[![ftbadge](${props.apiUrl})](${currentUrl()})`;
});
</script>

<template>
  <div class="rounded-md border border-gray-700 bg-gray-800 p-4">
    <code class="text-sm break-all text-orange-400">{{ code }}</code>
  </div>
  <CopyButton label="Markdown" :copyText="code" class="mt-3" />
</template>
