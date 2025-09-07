<script setup lang="ts">
import { Check, Copy } from "lucide-vue-next";
import { ref } from "vue";

import { Button } from "~/components/ui/button";

const props = defineProps<{
  label: string;
  copyText: string;
}>();

const copied = ref(false);

async function handleCopy() {
  await navigator.clipboard.writeText(props.copyText);
  copied.value = true;
  setTimeout(() => {
    copied.value = false;
  }, 1500);
}
</script>

<template>
  <Button
    variant="ghost"
    size="sm"
    class="text-gray-400 hover:bg-gray-700 hover:text-white"
    :onclick="handleCopy"
  >
    <component :is="copied ? Check : Copy" :size="16" class="mr-2" />
    {{ copied ? "Copied!" : `Copy ${label}` }}
  </Button>
</template>
