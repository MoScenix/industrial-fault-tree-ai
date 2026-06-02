<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="flex items-center justify-between px-4 py-2 border-b border-gray-100">
      <span class="text-xs font-medium text-gray-500 uppercase tracking-wider">AI 校验建议</span>
      <button
        class="text-xs font-medium text-indigo-600 hover:text-indigo-500 transition-colors disabled:opacity-40"
        :disabled="loading"
        @click="$emit('validate')"
      >
        {{ loading ? '校验中...' : '开始校验' }}
      </button>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4">
      <div
        class="prose prose-sm max-w-none text-sm text-gray-600 leading-relaxed"
        :class="{ 'opacity-50': loading }"
      >
        <MarkdownRenderer :content="content || '点击上方「开始校验」让 AI 检查故障树逻辑。'" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'

interface Props {
  content?: string
  loading: boolean
}

defineProps<Props>()
defineEmits<{
  (e: 'validate'): void
}>()
</script>

<style scoped>
.prose :deep(p) {
  margin: 0 0 0.75em;
}
.prose :deep(p:last-child) {
  margin: 0;
}
</style>
