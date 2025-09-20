<script setup lang="ts">
import { computed, ref } from "vue";

import ConfigurationCard from "~/components/ConfigurationCard.vue";
import HomeFooter from "~/components/HomeFooter.vue";
import HomeHeader from "~/components/HomeHeader.vue";
import PreviewCard from "~/components/PreviewCard.vue";
import SnippetCard from "~/components/SnippetCard.vue";
import { env } from "~/lib/env";

const login = ref<string>("test");
const width = ref<number | undefined>();
const height = ref<number | undefined>();

const apiUrl = computed<string>(() => {
  return `${env.VITE_API_URL}/profile/${login.value.trim()}`;
});
</script>

<template>
  <div class="bg-gray-950">
    <div class="min-h-screen">
      <HomeHeader />

      <main class="mx-auto max-w-6xl px-4 py-8">
        <div class="mb-12 text-center">
          <h2
            class="mb-4 bg-gradient-to-r from-white to-gray-300 bg-clip-text text-4xl font-bold text-transparent"
          >
            Generate Your Custom Badge
          </h2>
          <p class="mx-auto max-w-2xl text-xl text-gray-400">
            Create a personalized badge for your 42 profile with dynamic data.
          </p>
        </div>

        <div class="grid gap-6 lg:grid-cols-2">
          <div>
            <ConfigurationCard
              v-model:login="login"
              v-model:width="width"
              v-model:height="height"
              :apiUrl="apiUrl"
            />
          </div>

          <div class="space-y-6">
            <PreviewCard :api-url="apiUrl" :width="width" :height="height" />
            <SnippetCard :api-url="apiUrl" :width="width" :height="height" />
          </div>
        </div>
      </main>
    </div>

    <HomeFooter />
  </div>
</template>
