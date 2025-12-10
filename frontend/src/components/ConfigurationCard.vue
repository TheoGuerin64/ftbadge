<script setup lang="ts">
import { Settings } from "lucide-vue-next";
import { computed } from "vue";

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
import { isAlphanum } from "~/lib/utils";

const props = defineProps<{
  apiUrl: string;
}>();

const login = defineModel<string>("login");
const width = defineModel<number | undefined>("width");
const height = defineModel<number | undefined>("height");

const loginEmpty = computed(() => !login.value || !login.value.trim());
const loginAlphanum = computed(() => !login.value || isAlphanum(login.value));
</script>

<template>
  <Card class="border-gray-800 bg-gray-900">
    <CardHeader>
      <CardTitle class="flex items-center gap-2 text-white">
        <Settings :size="20" />
        Configuration
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
          maxlength="32"
          required
          placeholder="Enter your login"
          class="border-gray-700 bg-gray-800 text-white placeholder-gray-400"
          :aria-invalid="loginEmpty || !loginAlphanum"
        />
        <p v-if="loginEmpty" class="text-sm text-red-400">Login is required</p>
        <p v-if="!loginAlphanum" class="text-sm text-red-400">
          Login must be alphanumeric (letters and numbers only)
        </p>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div class="space-y-2">
          <Label htmlFor="width" class="text-gray-200">
            Width (px)
            <span class="text-xs text-gray-400">(optional)</span>
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
            <span class="text-xs text-gray-400">(optional)</span>
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
        <h4 class="text-sm text-gray-200">API URL</h4>
        <div class="mt-2 rounded-md border border-gray-700 bg-gray-800 p-3">
          <code class="text-sm break-all text-green-400">{{
            props.apiUrl
          }}</code>
        </div>
        <CopyButton label="URL" :copyText="props.apiUrl" class="mt-3" />
      </div>
    </CardContent>
  </Card>
</template>
