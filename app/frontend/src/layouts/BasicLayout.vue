<template>
  <a-layout class="basic-layout">
    <GlobalHeader v-if="!isWorkspace" />
    <a-layout-content class="main-content">
      <router-view />
    </a-layout-content>
    <GlobalFooter v-if="!isWorkspace" />
  </a-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import GlobalHeader from '@/components/GlobalHeader.vue'
import GlobalFooter from '@/components/GlobalFooter.vue'

const route = useRoute()
const isWorkspace = computed(() => route.path.startsWith('/graph/workspace'))
</script>

<style scoped>
.basic-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: none;
}
.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 0;
  background: none;
  margin: 0;
}
</style>
