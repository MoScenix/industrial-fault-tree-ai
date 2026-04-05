<template>
  <a-modal
    :open="open"
    title="编辑节点"
    ok-text="确定"
    cancel-text="取消"
    @ok="$emit('save')"
    @cancel="$emit('update:open', false)"
  >
    <template #footer>
      <div style="display: flex; justify-content: space-between;">
        <a-button danger @click="$emit('delete')">删除节点</a-button>
        <div>
          <a-button @click="$emit('update:open', false)">取消</a-button>
          <a-button type="primary" @click="$emit('save')">确定</a-button>
        </div>
      </div>
    </template>
    <a-form layout="vertical">
      <a-form-item label="节点名称">
        <a-input v-model:value="data.label" placeholder="输入节点名称" />
      </a-form-item>
      <a-form-item label="节点描述">
        <a-textarea v-model:value="data.description" placeholder="输入节点描述" :rows="3" />
      </a-form-item>
      <a-form-item label="节点类型">
        <a-select v-model:value="data.nodeType">
          <a-select-option value="top_event">顶事件 (Top Event)</a-select-option>
          <a-select-option value="intermediate_event">中间事件 (Intermediate Event)</a-select-option>
          <a-select-option value="basic_event">基础事件 (Basic Event)</a-select-option>
          <a-select-option value="gate">逻辑门 (Gate)</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item v-if="data.nodeType === 'gate'" label="逻辑门类型">
        <a-select v-model:value="data.gateType">
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
