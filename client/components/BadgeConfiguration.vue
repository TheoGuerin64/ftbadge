<script setup lang="ts">
import { Settings } from "lucide-vue-next";
import CopyButton from "~/components/CopyButton.vue";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";

const props = defineProps<{
  apiURL: string;
}>();

const login = defineModel<string>("login");
const width = defineModel<number | undefined>("width");
const height = defineModel<number | undefined>("height");
</script>

<template>
  <Card class="border-gray-800 bg-gray-900">
    <CardHeader>
      <CardTitle class="flex items-center gap-2 text-white">
        <Settings :size="20" />
        Badge Configuration
      </CardTitle>
      <CardDescription class="text-gray-400">
        Configure your badge login and dimensions
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-6">
      <div class="space-y-2">
        <Label htmlFor="login" class="text-gray-200">Login</Label>
        <Input
          id="login"
          v-model="login"
          type="text"
          required
          placeholder="Enter your login"
          class="border-gray-700 bg-gray-800 text-white placeholder-gray-400"
        />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div class="space-y-2">
          <Label htmlFor="width" class="text-gray-200">
            Width (px)
            <span class="text-xs text-gray-500">(optional)</span>
          </Label>
          <Input
            id="width"
            type="number"
            min="0"
            max="1000"
            v-model.number="width"
            placeholder="Auto"
            class="border-gray-700 bg-gray-800 text-white placeholder-gray-500"
          />
        </div>

        <div class="space-y-2">
          <Label htmlFor="height" class="text-gray-200">
            Height (px)
            <span class="text-xs text-gray-500">(optional)</span>
          </Label>
          <Input
            id="height"
            type="number"
            min="0"
            max="1000"
            v-model.number="height"
            placeholder="Auto"
            class="border-gray-700 bg-gray-800 text-white placeholder-gray-500"
          />
        </div>
      </div>

      <div class="border-t border-gray-800 pt-4">
        <Label class="text-sm text-gray-200">API URL</Label>
        <div class="mt-2 rounded-md border border-gray-700 bg-gray-800 p-3">
          <code class="text-sm break-all text-green-400">{{
            props.apiURL
          }}</code>
        </div>
        <CopyButton label="URL" :copyText="props.apiURL" class="mt-3" />
      </div>
    </CardContent>
  </Card>
</template>
