<script setup lang="ts">
import { GitHubIcon } from "vue3-simple-icons";

import { Button } from "~/components/ui/button";
</script>

<template>
  <header
    class="sticky top-0 z-50 border-b border-gray-800 bg-gray-900/50 backdrop-blur-sm"
  >
    <div class="mx-auto flex max-w-6xl items-center justify-between p-4">
      <div class="flex items-center gap-3">
        <img
          src="/favicon.svg"
          alt="ftbadge icon"
          class="h-10 w-10 rounded-lg"
        />
        <div>
          <h1
            class="bg-linear-to-r from-blue-400 to-purple-400 bg-clip-text text-2xl font-bold text-transparent"
          >
            ftbadge
          </h1>
          <p class="text-sm text-gray-400">Dynamic 42 Profile Badge</p>
        </div>
      </div>
      <Button
        variant="ghost"
        size="icon"
        class="text-gray-400 hover:bg-gray-700 hover:text-white"
        as-child
      >
        <a href="https://github.com/TheoGuerin64/ftbadge" target="_blank">
          <GitHubIcon :size="20" />
        </a>
      </Button>
    </div>
  </header>
</template>
