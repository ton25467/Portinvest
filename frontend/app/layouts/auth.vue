<script setup lang="ts">
import { ActivityIcon } from 'lucide-vue-next'
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-muted/30 px-4 py-12 sm:px-6 lg:px-8 font-sans">
    <div class="max-w-md w-full space-y-8 bg-card p-8 rounded-xl border border-border shadow-lg">
      <div class="flex flex-col items-center">
        <ActivityIcon class="h-12 w-12 text-primary" />
        <h1 class="mt-4 text-3xl font-extrabold tracking-tight text-foreground">PortInves</h1>
        <p class="mt-2 text-sm text-muted-foreground">Investment Portfolio & Infrastructure Monitoring</p>
      </div>
      <slot />
    </div>
  </div>
</template>
