<script setup lang="ts">
import { RouterView } from "vue-router";
import { Analytics } from "@vercel/analytics/vue";
import { SpeedInsights } from "@vercel/speed-insights/vue";
</script>

<template>
  <Analytics />
  <SpeedInsights />
  <RouterView />
</template>
