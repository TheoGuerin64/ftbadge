<script setup lang="ts">
import { CodeXml } from "lucide-vue-next";
import { ref, watch } from "vue";

import HtmlTab from "~/components/HtmlTab.vue";
import MarkdownTab from "~/components/MarkdownTab.vue";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";

const props = defineProps<{
  apiUrl: string;
  width?: number;
  height?: number;
}>();

const selectedTab = ref("html");

watch(
  () => [props.width, props.height],
  ([width, height]) => {
    if ((width || height) && selectedTab.value === "markdown") {
      selectedTab.value = "html";
    }
  },
  { immediate: true },
);
</script>

<template>
  <Card class="border-gray-800 bg-gray-900">
    <CardHeader>
      <CardTitle class="flex items-center gap-2 text-white">
        <CodeXml :size="20" />
        Code Snippet
      </CardTitle>
      <CardDescription class="text-gray-400">
        Copy the code in your README or website
      </CardDescription>
    </CardHeader>
    <CardContent>
      <Tabs defaultValue="html" class="w-full" v-model="selectedTab">
        <TabsList class="grid w-full grid-cols-2 bg-gray-800">
          <TabsTrigger
            value="markdown"
            class="data-[state=active]:bg-gray-700"
            :disabled="!!props.width || !!props.height"
          >
            Markdown
          </TabsTrigger>
          <TabsTrigger value="html" class="data-[state=active]:bg-gray-700">
            HTML
          </TabsTrigger>
        </TabsList>

        <TabsContent
          value="markdown"
          class="space-y-3"
          :disabled="!!props.width || !!props.height"
        >
          <MarkdownTab :apiUrl="props.apiUrl" />
        </TabsContent>

        <TabsContent value="html" class="space-y-3">
          <HtmlTab
            :apiUrl="props.apiUrl"
            :width="props.width"
            :height="props.height"
          />
        </TabsContent>
      </Tabs>
    </CardContent>
  </Card>
</template>
