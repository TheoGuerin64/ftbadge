<script setup lang="ts">
import { computed } from "vue";
import CopyButton from "~/components/CopyButton.vue";

const props = defineProps<{
  apiUrl: string;
  width?: number;
  height?: number;
}>();

const code = computed<string>(() => {
  let attributes = `src="${props.apiUrl}"`;
  if (props.width) {
    attributes += ` width="${props.width}"`;
  }
  if (props.height) {
    attributes += ` height="${props.height}"`;
  }
  return `<a href="https://ftbadge.vercel.app"><img ${attributes}></a>`;
});
</script>

<template>
  <div class="rounded-md border border-gray-700 bg-gray-800 p-4">
    <code class="text-sm break-all text-orange-400">{{ code }}</code>
  </div>
  <CopyButton label="HTML" :copyText="code" class="mt-3" />
</template>
