<script setup lang="ts">
import { ActivityIcon } from 'lucide-vue-next'
</script>

<template>
  <div class="h-screen w-screen overflow-hidden bg-transparent font-sans text-foreground">
    <div class="h-full w-full bg-slate-900/80 backdrop-blur-md border border-slate-700 rounded-xl overflow-hidden shadow-2xl flex flex-col">
      <!-- Title Bar (Draggable region for PWA/desktop) -->
      <div class="h-8 bg-slate-800/50 border-b border-slate-700 flex items-center px-3 drag-region select-none">
        <ActivityIcon class="w-4 h-4 text-primary mr-2" />
        <span class="text-xs font-semibold text-slate-200">PortInves Widget</span>
        <div class="ml-auto flex items-center gap-1 no-drag">
          <!-- Window controls (mock or real if Tauri/Electron) -->
          <div class="w-3 h-3 rounded-full bg-red-500 cursor-pointer hover:bg-red-400"></div>
          <div class="w-3 h-3 rounded-full bg-yellow-500 cursor-pointer hover:bg-yellow-400"></div>
          <div class="w-3 h-3 rounded-full bg-green-500 cursor-pointer hover:bg-green-400"></div>
        </div>
      </div>
      
      <!-- Widget Content -->
      <main class="flex-1 overflow-y-auto custom-scrollbar no-drag">
        <slot />
      </main>
    </div>
  </div>
</template>

<style scoped>
/* App region drag for desktop mode */
.drag-region {
  -webkit-app-region: drag;
}
.no-drag {
  -webkit-app-region: no-drag;
}
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}
</style>
