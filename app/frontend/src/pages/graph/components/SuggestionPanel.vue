<template>
  <div class="flex flex-col h-full overflow-hidden">
    <div class="mb-4 flex-shrink-0 flex justify-between items-center">
      <span class="text-xs text-slate-500">AI 逻辑校验建议</span>
      <a-button type="primary" size="small" :loading="loading" @click="$emit('validate')">
        开始校验
      </a-button>
    </div>
    <div class="suggestion-box flex-1 overflow-auto bg-slate-50 rounded-xl p-4 border border-slate-100">
      <MarkdownRenderer :content="content || '当前版本暂无建议内容。点击下方“开始校验”按钮，让 AI 检查故障树逻辑。'" />
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
.suggestion-box pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  color: #334155;
  line-height: 1.7;
}
</style>
