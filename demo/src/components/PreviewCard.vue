<script setup lang="ts">
import { Play } from "lucide-vue-next";

import BadgeLoader from "~/components/BadgeLoader.vue";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";

const props = defineProps<{
  apiUrl: string;
}>();
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
      <BadgeLoader
        class="rounded-lg border border-gray-700 bg-gray-800 p-4 text-white"
        :api-url="props.apiUrl"
      />
    </CardContent>
  </Card>
</template>
