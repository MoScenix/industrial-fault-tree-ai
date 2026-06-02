<template>
  <a-modal
    :open="open"
    title="编辑节点"
    ok-text="确定"
    cancel-text="取消"
    class="!rounded-xl"
    :mask-style="{ background: 'rgba(0,0,0,0.3)' }"
    @ok="$emit('save')"
    @cancel="$emit('update:open', false)"
  >
    <template #footer>
      <div class="flex justify-between items-center">
        <a-button danger size="small" class="!rounded-lg" @click="$emit('delete')">删除节点</a-button>
        <div class="flex gap-2">
          <a-button size="small" class="!rounded-lg" @click="$emit('update:open', false)">取消</a-button>
          <a-button type="primary" size="small" class="!rounded-lg" @click="$emit('save')">确定</a-button>
        </div>
      </div>
    </template>
    <a-form layout="vertical">
      <a-form-item label="节点名称">
        <a-input v-model:value="data.label" placeholder="输入节点名称" class="!rounded-lg" />
      </a-form-item>
      <a-form-item label="节点描述">
        <a-textarea v-model:value="data.description" placeholder="输入节点描述" :rows="2" class="!rounded-lg" />
      </a-form-item>
      <a-form-item label="节点类型">
        <a-select v-model:value="data.nodeType" class="!rounded-lg">
          <a-select-option value="top_event">顶事件 (Top Event)</a-select-option>
          <a-select-option value="intermediate_event">中间事件</a-select-option>
          <a-select-option value="basic_event">基础事件</a-select-option>
          <a-select-option value="gate">逻辑门 (Gate)</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item v-if="data.nodeType === 'gate'" label="逻辑门类型">
        <a-select v-model:value="data.gateType" class="!rounded-lg">
          <a-select-option value="AND">AND (与门)</a-select-option>
          <a-select-option value="OR">OR (或门)</a-select-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
interface Props {
  open: boolean
  data: {
    label: string
    description: string
    nodeType: string
    gateType: string
  }
}

defineProps<Props>()
defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'save'): void
  (e: 'delete'): void
}>()
</script>
