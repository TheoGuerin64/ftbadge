<script setup lang="ts">
import { debouncedRef } from "@vueuse/core";
import { CircleX, Loader2 } from "lucide-vue-next";
import { ref, toRef, watch } from "vue";

enum State {
  LOADING,
  ERROR,
  LOADED,
}

const props = defineProps<{
  apiUrl: string;
}>();
const apiUrl = toRef(() => props.apiUrl);

const imageRef = ref<HTMLImageElement | null>(null);
const state = ref<State>(State.LOADING);
const error = ref<string | null>(null);

const debouncedApiUrl = debouncedRef(apiUrl, 300);

watch(
  debouncedApiUrl,
  async (apiUrl) => {
    state.value = State.LOADING;
    error.value = null;

    if (apiUrl.endsWith("/")) {
      state.value = State.ERROR;
      error.value = "Please provide a valid login";
      return;
    }

    const response = await fetch(apiUrl);
    if (!response.ok) {
      const data = await response.json();
      error.value = data.message || "An error occurred";
      state.value = State.ERROR;
      return;
    }
    const blob = await response.blob();
    const image = URL.createObjectURL(blob);

    imageRef.value!.src = image;
    state.value = State.LOADED;
  },
  { immediate: true },
);
</script>

<template>
  <div class="flex justify-center">
    <a
      href="https://demo.ftbadge.cc"
      target="_blank"
      v-show="state === State.LOADED"
    >
      <img ref="imageRef" alt="generated-badge" />
    </a>
    <Loader2 v-if="state === State.LOADING" :size="20" class="animate-spin" />
    <div v-else-if="state === State.ERROR" class="flex items-center gap-2">
      <CircleX :size="20" class="text-red-500" />
      <span class="text-red-500">{{ error }}</span>
    </div>
  </div>
</template>
