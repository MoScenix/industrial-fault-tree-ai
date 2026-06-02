<template>
  <div class="flex flex-col h-full">
    <!-- Header with actions -->
    <div class="flex items-center justify-between px-4 py-2 border-b border-gray-100">
      <span class="text-xs font-medium text-gray-500 uppercase tracking-wider">版本列表</span>
      <button
        class="text-xs font-medium text-indigo-600 hover:text-indigo-500 transition-colors"
        @click="showCreateModal = true"
      >
        + 新建
      </button>
    </div>

    <!-- Version list -->
    <div class="flex-1 overflow-y-auto">
      <div v-if="!versions.length" class="text-center text-gray-400 text-sm mt-8">暂无版本</div>
      <div
        v-for="item in versions"
        :key="item.version"
        class="group flex items-center justify-between px-4 py-3 cursor-pointer border-b border-gray-50 hover:bg-gray-50 transition-colors"
        :class="{ 'bg-indigo-50/40': item.version === currentVersionLabel }"
        @click="$emit('switch', item.version)"
      >
        <div class="flex items-center gap-2 min-w-0">
          <div
            class="w-1.5 h-1.5 rounded-full flex-shrink-0"
            :class="item.version === currentVersionLabel ? 'bg-indigo-600' : 'bg-gray-300'"
          ></div>
          <div class="min-w-0">
            <div class="text-sm font-medium text-gray-900 truncate">{{ item.versionName || item.version }}</div>
          </div>
          <span v-if="item.version === currentVersionLabel" class="text-[11px] text-indigo-600 bg-indigo-50 px-1.5 py-0.5 rounded flex-shrink-0">当前</span>
        </div>
        <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
          <button
            class="px-2 py-1 text-[11px] text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded transition-colors"
            @click.stop="openRename(item)"
          >
            重命名
          </button>
          <button
            v-if="item.version !== currentVersionLabel"
            class="px-2 py-1 text-[11px] text-red-500 hover:text-red-600 hover:bg-red-50 rounded transition-colors"
            @click.stop="remove(item)"
          >
            删除
          </button>
        </div>
      </div>
    </div>

    <!-- Create Modal -->
    <a-modal
      v-model:open="showCreateModal"
      title="新建版本"
      ok-text="创建"
      cancel-text="取消"
      class="!rounded-xl"
      @ok="handleCreate"
    >
      <a-input v-model:value="versionDraft" placeholder="例如：v002 或 仿真优化版" class="!rounded-lg" />
    </a-modal>

    <!-- Rename Modal -->
    <a-modal
      v-model:open="showRenameModal"
      title="重命名版本"
      ok-text="保存"
      cancel-text="取消"
      class="!rounded-xl"
      @ok="handleRename"
    >
      <a-input v-model:value="versionDraft" placeholder="输入新的版本名称" class="!rounded-lg" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import {
  createGraphVersion,
  deleteGraphVersion,
  renameGraphVersion,
} from '@/api/graphController'

interface Props {
  graphId: number
  versions: API.GraphVersionVO[]
  currentVersionLabel: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'refresh'): void
  (e: 'switch', version?: string): void
}>()

const showCreateModal = ref(false)
const showRenameModal = ref(false)
const versionDraft = ref('')
const editingVersion = ref<API.GraphVersionVO>()

const handleCreate = async () => {
  if (!versionDraft.value.trim()) {
    message.warning('请输入版本名称')
    return
  }
  const res = await createGraphVersion({
    graphId: props.graphId,
    versionName: versionDraft.value.trim(),
  })
  if (res.data.code === 0) {
    message.success('版本创建成功')
    showCreateModal.value = false
    versionDraft.value = ''
    emit('refresh')
  } else {
    message.error(res.data.message || '版本创建失败')
  }
}

const openRename = (item: API.GraphVersionVO) => {
  editingVersion.value = item
  versionDraft.value = item.versionName || item.version || ''
  showRenameModal.value = true
}

const handleRename = async () => {
  if (!editingVersion.value?.version) return
  const res = await renameGraphVersion({
    graphId: props.graphId,
    version: editingVersion.value.version,
    versionName: versionDraft.value.trim(),
  })
  if (res.data.code === 0) {
    message.success('版本重命名成功')
    showRenameModal.value = false
    versionDraft.value = ''
    editingVersion.value = undefined
    emit('refresh')
  } else {
    message.error(res.data.message || '版本重命名失败')
  }
}

const remove = async (item: API.GraphVersionVO) => {
  if (!item.version) return
  const res = await deleteGraphVersion({
    graphId: props.graphId,
    version: item.version,
  })
  if (res.data.code === 0) {
    message.success('版本删除成功')
    emit('refresh')
  } else {
    message.error(res.data.message || '版本删除失败')
  }
}
</script>

<style scoped>
/* All styles replaced by Tailwind */
</style>
