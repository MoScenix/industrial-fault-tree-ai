<template>
  <a-space direction="vertical" style="width: 100%">
    <div class="panel-actions">
      <a-button type="primary" @click="showCreateModal = true">新建版本</a-button>
    </div>
    <a-list :data-source="versions" size="small">
      <template #renderItem="{ item }">
        <a-list-item class="version-row" @click="$emit('switch', item.version)">
          <div class="version-main">
            <div class="version-title">
              {{ item.versionName || item.version }}
              <a-tag v-if="item.version === currentVersionLabel" color="blue">当前</a-tag>
            </div>
            <div class="version-meta">{{ item.version }}</div>
          </div>
          <div class="version-actions">
            <a-button type="link" size="small" @click.stop="openRename(item)">重命名</a-button>
            <a-button
              type="link"
              size="small"
              danger
              :disabled="item.version === currentVersionLabel"
              @click.stop="remove(item)"
            >
              删除
            </a-button>
          </div>
        </a-list-item>
      </template>
    </a-list>

    <a-modal
      v-model:open="showCreateModal"
      title="新建版本"
      ok-text="创建"
      cancel-text="取消"
      @ok="handleCreate"
    >
      <a-input v-model:value="versionDraft" placeholder="例如：v002 或 仿真优化版" />
    </a-modal>

    <a-modal
      v-model:open="showRenameModal"
      title="重命名版本"
      ok-text="保存"
      cancel-text="取消"
      @ok="handleRename"
    >
      <a-input v-model:value="versionDraft" placeholder="输入新的版本名称" />
    </a-modal>
  </a-space>
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
.version-row {
  cursor: pointer;
  border-radius: 12px;
  transition: background 0.2s ease;
}
.version-row:hover {
  background: #f8fafc;
}
.version-main {
  flex: 1;
}
.version-title {
  color: #0f172a;
  font-weight: 600;
}
.version-meta {
  color: #94a3b8;
  font-size: 12px;
}
.version-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}
.panel-actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
  margin-top: 12px;
}
</style>
