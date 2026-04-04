<template>
  <div class="workspace-page">
    <header class="workspace-header">
      <div class="title-block">
        <a-button type="text" @click="router.push('/graph/manage')">返回项目</a-button>
        <div>
          <p class="eyebrow">Fault Tree Workspace</p>
          <h1>{{ graphInfo?.graphName || '故障树工作台' }}</h1>
          <p class="subtitle">
            当前展示 {{ workingGraph?.isTmp ? '暂存版本' : '目录版本' }} ·
            {{ selectedVersion || graphInfo?.currentVersion || 'v001' }}
          </p>
        </div>
      </div>
      <a-space wrap>
        <a-button type="primary" @click="openSaveModal">保存</a-button>
        <a-upload :show-upload-list="false" accept=".pdf" :before-upload="beforeUploadProjectDoc">
          <a-button>上传项目文档</a-button>
        </a-upload>
        <a-upload :show-upload-list="false" accept=".json" :before-upload="beforeImportGraph">
          <a-button>导入图</a-button>
        </a-upload>
        <a-button @click="exportCurrentGraph">导出当前图</a-button>
        <a-button danger @click="discardTmp">放弃暂存</a-button>
      </a-space>
    </header>

    <div class="workspace-body">
      <aside class="tool-rail">
        <button
          v-for="item in panelItems"
          :key="item.key"
          class="tool-button"
          :class="{ active: activePanel === item.key }"
          @click="togglePanel(item.key)"
        >
          <component :is="item.icon" />
          <span>{{ item.label }}</span>
        </button>
      </aside>

      <section v-if="activePanel" class="side-panel" :style="{ width: `${panelWidth}px` }">
        <div class="panel-header">
          <div>
            <p class="panel-eyebrow">{{ currentPanelMeta.eyebrow }}</p>
            <h3>{{ currentPanelMeta.title }}</h3>
          </div>
          <a-button type="text" @click="activePanel = ''">收起</a-button>
        </div>

        <div class="panel-content">
          <template v-if="activePanel === 'chat'">
            <div class="chat-list">
              <div v-for="item in messages" :key="item.id" class="chat-item">
                <div class="chat-role">{{ item.role === 'assistant' ? 'AI' : '工程师' }}</div>
                <div class="chat-content">{{ item.content }}</div>
              </div>
              <a-empty v-if="!messages.length" description="当前还没有聊天记录" />
            </div>
            <a-textarea
              v-model:value="chatInput"
              :rows="4"
              :maxlength="1000"
              placeholder="告诉 AI 你希望如何修改或解释当前故障树"
            />
            <div class="panel-actions">
              <a-button type="primary" :loading="chatting" @click="sendChatMessage">发起一次对话</a-button>
            </div>
          </template>

          <template v-else-if="activePanel === 'version'">
            <a-space direction="vertical" style="width: 100%">
              <div class="panel-actions">
                <a-button type="primary" @click="showCreateVersionModal = true">新建版本</a-button>
              </div>
              <a-list :data-source="versions" size="small">
                <template #renderItem="{ item }">
                  <a-list-item class="version-row" @click="switchVersion(item.version)">
                    <div class="version-main">
                      <div class="version-title">
                        {{ item.versionName || item.version }}
                        <a-tag v-if="item.isCurrent" color="blue">当前</a-tag>
                      </div>
                      <div class="version-meta">{{ item.version }}</div>
                    </div>
                    <div class="version-actions">
                      <a-button type="link" size="small" @click.stop="openRenameVersion(item)">重命名</a-button>
                      <a-button
                        type="link"
                        size="small"
                        danger
                        :disabled="item.isCurrent"
                        @click.stop="removeVersion(item)"
                      >
                        删除
                      </a-button>
                    </div>
                  </a-list-item>
                </template>
              </a-list>
            </a-space>
          </template>

          <template v-else-if="activePanel === 'suggestion'">
            <div class="suggestion-box">
              <pre>{{ suggestion?.content || '当前版本暂无建议内容。' }}</pre>
            </div>
          </template>
        </div>
        <div class="resize-handle" @mousedown="startResize"></div>
      </section>

      <main class="canvas-panel">
        <div class="canvas-toolbar">
          <div class="canvas-meta">
            <a-tag color="blue">{{ selectedVersion || graphInfo?.currentVersion || 'v001' }}</a-tag>
            <a-tag :color="workingGraph?.isTmp ? 'orange' : 'cyan'">
              {{ workingGraph?.isTmp ? '正在查看暂存版本' : '正在查看目录版本' }}
            </a-tag>
            <span class="meta-text">创建者：{{ graphInfo?.userId || '-' }}</span>
            <span class="meta-text">
              {{ workingGraph?.isTmp ? '打开页面后系统已自动创建该版本的编辑副本' : '当前展示的是目录版本，只读展示' }}
            </span>
          </div>
          <div class="canvas-meta">
            <span class="meta-text">节点 {{ flowNodes.length }}</span>
            <span class="meta-text">边 {{ flowEdges.length }}</span>
            <span v-if="isEditorDirty" class="meta-text">未保存修改</span>
          </div>
        </div>

        <div class="flow-shell">
          <VueFlow
            v-model:nodes="flowNodes"
            v-model:edges="flowEdges"
            fit-view-on-init
            :min-zoom="0.2"
            :max-zoom="1.6"
            :default-edge-options="{ markerEnd: MarkerType.ArrowClosed }"
            @node-drag-stop="markDirty"
            @connect="onConnect"
            @edges-change="markDirty"
            @nodes-change="markDirty"
          >
            <Background pattern-color="#dbeafe" :gap="24" />
            <MiniMap pannable zoomable />
            <Controls />
          </VueFlow>
        </div>

        <div class="editor-footer">
          <div class="editor-footer-head">
            <span class="meta-text">开发调试</span>
            <a-switch v-model:checked="showJsonEditor" checked-children="JSON" un-checked-children="隐藏" />
          </div>
          <a-card :bordered="false" class="json-card">
            <template #title>当前图 JSON 预览</template>
            <a-textarea
              v-if="showJsonEditor"
              v-model:value="workingContent"
              :rows="12"
              @blur="handleJsonEdited"
              placeholder="这里会实时展示当前工作图内容，必要时也可以直接粘贴 JSON 再导入。"
            />
            <a-empty v-else description="默认隐藏 JSON 调试区，避免干扰正式图编辑。" />
          </a-card>
        </div>
      </main>
    </div>

    <a-modal
      v-model:open="showCreateVersionModal"
      title="新建版本"
      ok-text="创建"
      cancel-text="取消"
      @ok="createVersion"
    >
      <a-input v-model:value="versionDraft" placeholder="例如：v002 或 仿真优化版" />
    </a-modal>

    <a-modal
      v-model:open="showRenameVersionModal"
      title="重命名版本"
      ok-text="保存"
      cancel-text="取消"
      @ok="renameVersion"
    >
      <a-input v-model:value="versionDraft" placeholder="输入新的版本名称" />
    </a-modal>

    <a-modal
      v-model:open="showSaveModal"
      title="保存版本"
      ok-text="保存"
      cancel-text="取消"
      @ok="saveCurrent"
    >
      <a-radio-group v-model:value="saveMode" style="margin-bottom: 16px">
        <a-radio value="overwrite">覆盖当前版本</a-radio>
        <a-radio value="new">新建版本保存</a-radio>
      </a-radio-group>
      <a-input
        v-if="saveMode === 'new'"
        v-model:value="saveVersionName"
        placeholder="输入新版本名称，例如：v002 或 优化版"
      />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import type { Connection, Edge, Node } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import { MarkerType, VueFlow } from '@vue-flow/core'
import {
  CommentOutlined,
  DeploymentUnitOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons-vue'
import {
  chatToModifyGraph,
  createGraphVersion,
  deleteGraphVersion,
  discardWorkingGraph,
  getCurrentSuggestion,
  getGraphVoById,
  getWorkingGraph,
  listGraphMessage,
  listGraphVersion,
  renameGraphVersion,
  saveGraph,
  startEdit,
} from '@/api/graphController'
import { uploadProjectDocument } from '@/api/documentController'
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'
import '@vue-flow/minimap/dist/style.css'

type PanelKey = '' | 'chat' | 'version' | 'suggestion'

const route = useRoute()
const router = useRouter()
const graphId = Number(route.params.id)

const graphInfo = ref<API.GraphVO>()
const workingGraph = ref<API.WorkingGraphVO>()
const suggestion = ref<API.GraphSuggestionVO>()
const versions = ref<API.GraphVersionVO[]>([])
const messages = ref<API.GraphMessageVO[]>([])
const selectedVersion = ref('')
const workingContent = ref('')
const flowNodes = ref<Node[]>([])
const flowEdges = ref<Edge[]>([])
const isEditorDirty = ref(false)
const activePanel = ref<PanelKey>('')
const panelWidth = ref(380)
const isResizing = ref(false)
const originalGraphMeta = ref<Record<string, any>>({})
const chatInput = ref('')
const chatting = ref(false)
const showCreateVersionModal = ref(false)
const showRenameVersionModal = ref(false)
const showSaveModal = ref(false)
const versionDraft = ref('')
const editingVersion = ref<API.GraphVersionVO>()
const showJsonEditor = ref(false)
const saveMode = ref<'overwrite' | 'new'>('overwrite')
const saveVersionName = ref('')
let draftSyncTimer: ReturnType<typeof setTimeout> | undefined

const panelItems = [
  { key: 'chat', label: '对话', icon: CommentOutlined },
  { key: 'version', label: '版本', icon: DeploymentUnitOutlined },
  { key: 'suggestion', label: '校验', icon: SafetyCertificateOutlined },
] as const

const panelMeta = {
  chat: { eyebrow: 'Dialogue', title: 'AI 对话记录' },
  version: { eyebrow: 'Version', title: '版本列表' },
  suggestion: { eyebrow: 'Validate', title: '当前建议' },
}

const currentVersionLabel = computed(
  () => selectedVersion.value || graphInfo.value?.currentVersion || 'v001',
)
const draftCacheKey = computed(
  () => `graph-workspace:${graphId}:${currentVersionLabel.value || 'v001'}`,
)
const currentPanelMeta = computed(() => {
  if (activePanel.value && panelMeta[activePanel.value]) {
    return panelMeta[activePanel.value]
  }
  return { eyebrow: '', title: '' }
})

const buildDefaultGraph = () => {
  const fallback = {
    schema_version: 'fault-tree/v1',
    tree: { name: '新建故障树', top_node_id: 'node-top' },
    nodes: [
      {
        node_id: 'node-top',
        node_type: 'top_event',
        label: '系统故障',
        description: '顶事件',
        gate_type: '',
        points_to: ['node-gate-1'],
        pointed_by: [],
      },
      {
        node_id: 'node-gate-1',
        node_type: 'gate',
        label: 'OR',
        description: '逻辑门',
        gate_type: 'OR',
        points_to: ['node-left', 'node-right'],
        pointed_by: ['node-top'],
      },
      {
        node_id: 'node-left',
        node_type: 'basic_event',
        label: '传感器失效',
        description: '基础事件',
        gate_type: '',
        points_to: [],
        pointed_by: ['node-gate-1'],
      },
      {
        node_id: 'node-right',
        node_type: 'basic_event',
        label: '电源波动',
        description: '基础事件',
        gate_type: '',
        points_to: [],
        pointed_by: ['node-gate-1'],
      },
    ],
    meta: { version: currentVersionLabel.value },
  }
  return fallback
}

const nodeColor = (type?: string) => {
  if (type === 'top_event') return '#dbeafe'
  if (type === 'gate') return '#ccfbf1'
  if (type === 'basic_event') return '#ffffff'
  return '#eff6ff'
}

const parseGraphContent = (content?: string) => {
  let parsed: any
  try {
    parsed = content ? JSON.parse(content) : buildDefaultGraph()
  } catch {
    parsed = buildDefaultGraph()
  }

  originalGraphMeta.value = {
    schema_version: parsed.schema_version || 'fault-tree/v1',
    tree: parsed.tree || { name: graphInfo.value?.graphName || '故障树', top_node_id: '' },
    meta: parsed.meta || {},
  }

  const nodes = Array.isArray(parsed.nodes) ? parsed.nodes : []
  flowNodes.value = nodes.map((item: any, index: number) => ({
    id: item.node_id || `node-${index + 1}`,
    position: item.position || {
      x: 100 + (index % 3) * 260,
      y: 80 + Math.floor(index / 3) * 160,
    },
    data: {
      label: item.label || item.node_id || `节点 ${index + 1}`,
      nodeType: item.node_type || 'intermediate_event',
      description: item.description || '',
      gateType: item.gate_type || '',
    },
    style: {
      background: nodeColor(item.node_type),
      border: item.node_type === 'gate' ? '1px solid #2dd4bf' : '1px solid #93c5fd',
      borderRadius: '16px',
      padding: '10px 14px',
      color: '#0f172a',
      width: item.node_type === 'gate' ? '120px' : '180px',
      fontWeight: 600,
      boxShadow: '0 10px 24px rgba(15, 23, 42, 0.08)',
    },
  }))

  const edges: Edge[] = []
  nodes.forEach((item: any) => {
    const source = item.node_id
    const targets = Array.isArray(item.points_to) ? item.points_to : []
    targets.forEach((targetId: string, edgeIndex: number) => {
      edges.push({
        id: `${source}-${targetId}-${edgeIndex}`,
        source,
        target: targetId,
        animated: item.node_type === 'gate',
      })
    })
  })
  flowEdges.value = edges
  workingContent.value = JSON.stringify(parsed, null, 2)
  isEditorDirty.value = false
}

const exportGraphModel = () => {
  const edgeMap = new Map<string, string[]>()
  const reverseEdgeMap = new Map<string, string[]>()

  flowEdges.value.forEach((edge) => {
    const pointsTo = edgeMap.get(edge.source) || []
    pointsTo.push(edge.target)
    edgeMap.set(edge.source, pointsTo)

    const pointedBy = reverseEdgeMap.get(edge.target) || []
    pointedBy.push(edge.source)
    reverseEdgeMap.set(edge.target, pointedBy)
  })

  const nodes = flowNodes.value.map((node) => ({
    node_id: node.id,
    node_type: String(node.data?.nodeType || 'intermediate_event'),
    label: String(node.data?.label || node.id),
    description: String(node.data?.description || ''),
    gate_type: String(node.data?.gateType || ''),
    points_to: edgeMap.get(node.id) || [],
    pointed_by: reverseEdgeMap.get(node.id) || [],
    position: node.position,
  }))

  return {
    schema_version: originalGraphMeta.value.schema_version || 'fault-tree/v1',
    tree: {
      ...(originalGraphMeta.value.tree || {}),
      name: graphInfo.value?.graphName || originalGraphMeta.value.tree?.name || '故障树',
      top_node_id:
        originalGraphMeta.value.tree?.top_node_id ||
        nodes.find((item) => item.node_type === 'top_event')?.node_id ||
        nodes[0]?.node_id ||
        '',
    },
    nodes,
    meta: {
      ...(originalGraphMeta.value.meta || {}),
      version: currentVersionLabel.value,
      generated_at: new Date().toISOString(),
    },
  }
}

const markDirty = () => {
  isEditorDirty.value = true
  workingContent.value = JSON.stringify(exportGraphModel(), null, 2)
}

const persistDraftCache = () => {
  if (!workingContent.value || !isEditorDirty.value) return
  localStorage.setItem(
    draftCacheKey.value,
    JSON.stringify({
      content: workingContent.value,
      updatedAt: new Date().toISOString(),
    }),
  )
}

const clearDraftCache = () => {
  localStorage.removeItem(draftCacheKey.value)
}

const scheduleDraftSync = () => {
  if (draftSyncTimer) clearTimeout(draftSyncTimer)
  draftSyncTimer = setTimeout(() => {
    persistDraftCache()
  }, 3000)
}

const restoreDraftCache = () => {
  const raw = localStorage.getItem(draftCacheKey.value)
  if (!raw) return false
  try {
    const parsed = JSON.parse(raw)
    if (parsed?.content) {
      parseGraphContent(parsed.content)
      isEditorDirty.value = true
      message.info('已恢复当前版本的本地缓存编辑内容')
      return true
    }
  } catch {
    localStorage.removeItem(draftCacheKey.value)
  }
  return false
}

const loadGraphInfo = async () => {
  const res = await getGraphVoById({ id: graphId })
  if (res.data.code === 0) {
    graphInfo.value = res.data.data
    selectedVersion.value = selectedVersion.value || res.data.data?.currentVersion || 'v001'
  } else {
    message.error(res.data.message || '获取项目详情失败')
  }
}

const loadVersions = async () => {
  const res = await listGraphVersion({ graphId })
  if (res.data.code === 0 && res.data.data) {
    versions.value = res.data.data.records || []
  }
}

const loadWorkingGraph = async () => {
  workingContent.value = ''
  const res = await getWorkingGraph({ graphId, version: selectedVersion.value })
  if (res.data.code === 0 && res.data.data) {
    workingGraph.value = res.data.data
    parseGraphContent(res.data.data.content || '')
    restoreDraftCache()
  } else {
    message.error(res.data.message || '获取当前工作图失败')
  }
}

const loadSuggestion = async () => {
  const res = await getCurrentSuggestion({ graphId, version: selectedVersion.value })
  if (res.data.code === 0) {
    suggestion.value = res.data.data
  }
}

const loadMessages = async () => {
  const res = await listGraphMessage({ graphId, pageSize: 10 })
  if (res.data.code === 0 && res.data.data) {
    messages.value = res.data.data.records || []
  }
}

const refreshAll = async () => {
  await Promise.all([loadGraphInfo(), loadVersions()])
  await ensureEditReady()
  await Promise.all([loadWorkingGraph(), loadSuggestion(), loadMessages()])
}

const switchVersion = async (version?: string) => {
  selectedVersion.value = version || graphInfo.value?.currentVersion || 'v001'
  await ensureEditReady()
  await Promise.all([loadWorkingGraph(), loadSuggestion()])
}

const ensureEditReady = async () => {
  const res = await startEdit({
    graphId,
    version: selectedVersion.value || graphInfo.value?.currentVersion || 'v001',
  })
  if (res.data.code !== 0) {
    throw new Error(res.data.message || '自动准备编辑副本失败')
  }
}

const openSaveModal = () => {
  saveMode.value = 'overwrite'
  saveVersionName.value = ''
  showSaveModal.value = true
}

const discardTmp = async () => {
  const res = await discardWorkingGraph({ graphId, version: selectedVersion.value })
  if (res.data.code === 0) {
    message.success('已放弃当前暂存内容')
    clearDraftCache()
    await refreshAll()
  } else {
    message.error(res.data.message || '放弃暂存失败')
  }
}

const saveCurrent = async () => {
  const fromVersion = currentVersionLabel.value
  const toVersion = saveMode.value === 'new' ? saveVersionName.value.trim() : fromVersion
  if (!toVersion) {
    message.warning('请输入目标版本名称')
    return
  }
  const content = JSON.stringify(exportGraphModel(), null, 2)
  workingContent.value = content
  const useTmp = !isEditorDirty.value && Boolean(workingGraph.value?.isTmp)
  const res = await saveGraph({
    graphId,
    fromVersion,
    toVersion,
    remark: '前端工作台保存',
    useTmp,
    content: useTmp ? '' : content,
  })
  if (res.data.code === 0) {
    message.success(res.data.data?.message || '保存成功')
    showSaveModal.value = false
    clearDraftCache()
    if (saveMode.value === 'new') {
      selectedVersion.value = toVersion
    }
    await refreshAll()
  } else {
    message.error(res.data.message || '保存失败')
  }
}

const sendChatMessage = async () => {
  if (!chatInput.value.trim()) {
    message.warning('请先输入对话内容')
    return
  }
  chatting.value = true
  try {
    const content = chatInput.value.trim()
    const res = await chatToModifyGraph({ graphId, message: content })
    if (res.data?.d) {
      message.success('AI 回复成功')
      chatInput.value = ''
      await loadMessages()
      activePanel.value = 'chat'
    } else {
      message.error(res.data?.message || 'AI 对话失败')
    }
  } finally {
    chatting.value = false
  }
}

const createVersion = async () => {
  if (!versionDraft.value.trim()) {
    message.warning('请输入版本名称')
    return
  }
  const res = await createGraphVersion({
    graphId,
    versionName: versionDraft.value.trim(),
  })
  if (res.data.code === 0) {
    message.success('版本创建成功')
    showCreateVersionModal.value = false
    versionDraft.value = ''
    await loadVersions()
  } else {
    message.error(res.data.message || '版本创建失败')
  }
}

const openRenameVersion = (item: API.GraphVersionVO) => {
  editingVersion.value = item
  versionDraft.value = item.versionName || item.version || ''
  showRenameVersionModal.value = true
}

const renameVersion = async () => {
  if (!editingVersion.value?.version) return
  const res = await renameGraphVersion({
    graphId,
    version: editingVersion.value.version,
    versionName: versionDraft.value.trim(),
  })
  if (res.data.code === 0) {
    message.success('版本重命名成功')
    showRenameVersionModal.value = false
    versionDraft.value = ''
    editingVersion.value = undefined
    await loadVersions()
  } else {
    message.error(res.data.message || '版本重命名失败')
  }
}

const removeVersion = async (item: API.GraphVersionVO) => {
  if (!item.version) return
  const res = await deleteGraphVersion({
    graphId,
    version: item.version,
  })
  if (res.data.code === 0) {
    message.success('版本删除成功')
    await loadVersions()
  } else {
    message.error(res.data.message || '版本删除失败')
  }
}

const beforeUploadProjectDoc = async (file: File) => {
  const res = await uploadProjectDocument(graphId, file)
  if (res.data.code === 0) {
    message.success('项目文档上传成功')
  } else {
    message.error(res.data.message || '项目文档上传失败')
  }
  return false
}

const beforeImportGraph = async (file: File) => {
  const text = await file.text()
  parseGraphContent(text)
  isEditorDirty.value = true
  scheduleDraftSync()
  message.success('图文件已导入到当前编辑器')
  return false
}

const exportCurrentGraph = () => {
  const content = JSON.stringify(exportGraphModel(), null, 2)
  const blob = new Blob([content], { type: 'application/json;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `${graphInfo.value?.graphName || 'fault-tree'}-${currentVersionLabel.value}.json`
  link.click()
  URL.revokeObjectURL(url)
}

const handleJsonEdited = () => {
  try {
    JSON.parse(workingContent.value)
    parseGraphContent(workingContent.value)
    isEditorDirty.value = true
    scheduleDraftSync()
  } catch (error) {
    console.error(error)
    message.warning('JSON 格式暂时不合法，先检查内容再继续')
  }
}

const onConnect = (connection: Connection) => {
  if (!connection.source || !connection.target) {
    return
  }
  flowEdges.value = [
    ...flowEdges.value,
    {
      id: `${connection.source}-${connection.target}-${Date.now()}`,
      source: connection.source,
      target: connection.target,
      markerEnd: MarkerType.ArrowClosed,
    },
  ]
  markDirty()
}

const togglePanel = (key: PanelKey) => {
  activePanel.value = activePanel.value === key ? '' : key
}

const startResize = () => {
  isResizing.value = true
}

const handleMouseMove = (event: MouseEvent) => {
  if (!isResizing.value) return
  panelWidth.value = Math.max(320, Math.min(560, event.clientX - 72))
}

const stopResize = () => {
  isResizing.value = false
}

onMounted(() => {
  refreshAll()
  window.addEventListener('mousemove', handleMouseMove)
  window.addEventListener('mouseup', stopResize)
})

onBeforeUnmount(() => {
  if (draftSyncTimer) {
    clearTimeout(draftSyncTimer)
  }
  window.removeEventListener('mousemove', handleMouseMove)
  window.removeEventListener('mouseup', stopResize)
})

watch(
  [flowNodes, flowEdges],
  () => {
    if (!workingGraph.value || !isEditorDirty.value) return
    scheduleDraftSync()
  },
  { deep: true },
)
</script>

<style scoped>
.workspace-page {
  min-height: calc(100vh - 120px);
  padding: 24px;
  background:
    radial-gradient(circle at top right, rgba(8, 145, 178, 0.14), transparent 24%),
    linear-gradient(180deg, #f8fafc 0%, #eef2ff 100%);
}

.workspace-header {
  display: flex;
  justify-content: space-between;
  align-items: end;
  gap: 24px;
  margin-bottom: 20px;
}

.title-block {
  display: flex;
  align-items: start;
  gap: 12px;
}

.eyebrow {
  margin: 0 0 8px;
  color: #0891b2;
  text-transform: uppercase;
  letter-spacing: 0.14em;
  font-weight: 700;
  font-size: 12px;
}

.workspace-header h1 {
  margin: 0 0 8px;
  font-size: 30px;
  color: #0f172a;
}

.subtitle {
  margin: 0;
  color: #64748b;
}

.workspace-body {
  display: grid;
  grid-template-columns: 64px auto minmax(0, 1fr);
  gap: 16px;
  min-height: calc(100vh - 220px);
}

.tool-rail {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px 0;
}

.tool-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 74px;
  border: 1px solid #dbeafe;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.88);
  color: #334155;
  cursor: pointer;
  transition: all 0.22s ease;
}

.tool-button:hover,
.tool-button.active {
  border-color: #60a5fa;
  background: #dbeafe;
  color: #1d4ed8;
}

.tool-button span {
  font-size: 12px;
}

.side-panel {
  position: relative;
  display: flex;
  flex-direction: column;
  min-width: 280px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid #dbeafe;
  border-radius: 24px;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.06);
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  padding: 20px 20px 16px;
  border-bottom: 1px solid #eff6ff;
}

.panel-eyebrow {
  margin: 0 0 6px;
  color: #2563eb;
  font-size: 11px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.panel-header h3 {
  margin: 0;
  color: #0f172a;
}

.panel-content {
  flex: 1;
  padding: 16px 20px 20px;
  overflow: auto;
}

.panel-actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
  margin-top: 12px;
}

.resize-handle {
  position: absolute;
  top: 0;
  right: 0;
  width: 8px;
  height: 100%;
  cursor: ew-resize;
}

.canvas-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.canvas-toolbar {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 18px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.86);
  border: 1px solid #dbeafe;
}

.canvas-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

.meta-text {
  color: #64748b;
  font-size: 13px;
}

.flow-shell {
  height: 620px;
  border-radius: 28px;
  overflow: hidden;
  border: 1px solid #dbeafe;
  background: linear-gradient(180deg, #f8fbff 0%, #eff6ff 100%);
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.08);
}

.editor-footer {
  display: block;
}

.editor-footer-head {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
}

.json-card {
  border-radius: 22px;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.06);
}

.chat-list {
  margin-bottom: 16px;
}

.chat-item + .chat-item {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #eef2f7;
}

.chat-role {
  margin-bottom: 6px;
  color: #2563eb;
  font-size: 12px;
  font-weight: 700;
}

.chat-content {
  color: #334155;
  line-height: 1.7;
  white-space: pre-wrap;
}

.version-row {
  cursor: pointer;
  border-radius: 12px;
  transition: background 0.2s ease;
}

.version-main {
  flex: 1;
}

.version-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.version-row:hover {
  background: #f8fafc;
}

.version-title {
  color: #0f172a;
  font-weight: 600;
}

.version-meta {
  color: #94a3b8;
  font-size: 12px;
}

.suggestion-box pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  color: #334155;
  line-height: 1.7;
}

:deep(.vue-flow__node) {
  font-size: 13px;
}

:deep(.vue-flow__controls) {
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.12);
}
</style>
