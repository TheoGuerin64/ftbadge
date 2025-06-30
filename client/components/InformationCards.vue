<script setup lang="ts">
import { Clipboard, ExternalLink, Zap } from "lucide-vue-next";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";

interface CardProps {
  title: string;
  description: string;
  icon: any;
  iconColor: string;
}

const cards: CardProps[] = [
  {
    title: "Fast & Reliable",
    description: "Lightning-fast SVG generation with 99.9% uptime.",
    icon: Zap,
    iconColor: "text-blue-400",
  },
  {
    title: "Smart Caching",
    description: "Intelligent caching system ensures rapid delivery.",
    icon: ExternalLink,
    iconColor: "text-green-400",
  },
  {
    title: "Effortless Integration",
    description: "Simple to use, just copy and paste into your README.",
    icon: Clipboard,
    iconColor: "text-green-400",
  },
];
</script>

<template>
  <div class="grid gap-6 md:grid-cols-3">
    <Card v-for="card in cards" class="border-gray-800 bg-gray-900">
      <CardHeader>
        <CardTitle class="flex items-center gap-2 text-white">
          <component :is="card.icon" :size="20" :class="card.iconColor" />
          {{ card.title }}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <p class="text-gray-400">
          {{ card.description }}
        </p>
      </CardContent>
    </Card>
  </div>
</template>
