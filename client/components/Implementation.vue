<script setup lang="ts">
import { CodeXml } from "lucide-vue-next";
import { computed } from "vue";
import CopyButton from "~/components/CopyButton.vue";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";

const props = defineProps<{
  apiUrl: string;
  width?: number;
  height?: number;
}>();

const htmlCode = computed<string>(() => {
  let attributes = `src="${props.apiUrl}"`;
  if (props.width) {
    attributes += ` width="${props.width}"`;
  }
  if (props.height) {
    attributes += ` height="${props.height}"`;
  }
  return `<a href="https://ftbadge.cc"><img ${attributes}></a>`;
});
</script>

<template>
  <Card class="border-gray-800 bg-gray-900">
    <CardHeader>
      <CardTitle class="flex items-center gap-2 text-white">
        <CodeXml :size="20" />
        Implementation
      </CardTitle>
      <CardDescription class="text-gray-400">
        Copy the code in your README or website
      </CardDescription>
    </CardHeader>
    <CardContent>
      <div class="rounded-md border border-gray-700 bg-gray-800 p-4">
        <code class="text-sm break-all text-orange-400">{{ htmlCode }}</code>
      </div>
      <CopyButton label="HTML" :copyText="htmlCode" class="mt-3" />
    </CardContent>
  </Card>
</template>
