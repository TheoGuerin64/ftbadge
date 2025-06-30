<script setup lang="ts">
import { computed, ref } from "vue";
import BadgeConfiguration from "~/components/BadgeConfiguration.vue";
import Footer from "~/components/Footer.vue";
import Header from "~/components/Header.vue";
import Implementation from "~/components/Implementation.vue";
import InformationCards from "~/components/InformationCards.vue";
import LivePreview from "~/components/LivePreview.vue";

const login = ref<string>("tguerin");
const width = ref<number | undefined>();
const height = ref<number | undefined>();

const apiURL = computed<string>(() => {
  return window.location.href + login.value;
});
</script>

<template>
  <div class="bg-gray-950">
    <Header />

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
          <BadgeConfiguration
            v-model:login="login"
            v-model:width="width"
            v-model:height="height"
            :apiURL="apiURL"
          />
        </div>

        <div class="space-y-6">
          <LivePreview :api-url="apiURL" :width="width" :height="height" />
          <Implementation :api-url="apiURL" :width="width" :height="height" />
        </div>
      </div>

      <div class="mt-16">
        <h3 class="mb-8 text-center text-2xl font-bold text-white">
          Why Choose ftbadge?
        </h3>
        <InformationCards />
      </div>
    </main>

    <Footer />
  </div>
</template>
