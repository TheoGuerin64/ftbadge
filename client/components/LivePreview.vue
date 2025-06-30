<script setup lang="ts">
import { debouncedRef, toRef } from "@vueuse/core";
import { CircleX, Loader2, Play } from "lucide-vue-next";
import { ref, watch } from "vue";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";

enum ImageState {
  LOADING,
  ERROR,
  LOADED,
}

const props = defineProps<{
  apiUrl: string;
}>();
const apiURL = toRef(props, "apiUrl");

const imageState = ref<ImageState>(ImageState.LOADING);
watch(apiURL, () => {
  imageState.value = ImageState.LOADING;
});

const debouncedApiURL = debouncedRef(apiURL, 300);
</script>

<template>
  <Card class="border-gray-800 bg-gray-900">
    <CardHeader>
      <CardTitle class="flex items-center gap-2 text-white">
        <Play :size="20" />
        Live Preview
      </CardTitle>
      <CardDescription class="text-gray-400">
        See your badge in real-time (except for width/height changes)
      </CardDescription>
    </CardHeader>
    <CardContent>
      <div
        class="flex justify-center rounded-lg border border-gray-700 bg-gray-800 p-4 text-white"
      >
        <a
          href="https://github.com/TheoGuerin64/ftbadge"
          target="_blank"
          v-show="imageState == ImageState.LOADED"
        >
          <img
            :src="debouncedApiURL"
            @load="() => (imageState = ImageState.LOADED)"
            @error="() => (imageState = ImageState.ERROR)"
          />
        </a>
        <Loader2
          v-if="imageState == ImageState.LOADING"
          :size="20"
          class="animate-spin"
        />
        <CircleX
          v-if="imageState == ImageState.ERROR"
          :size="20"
          class="text-red-500"
        />
      </div>
    </CardContent>
  </Card>
</template>
