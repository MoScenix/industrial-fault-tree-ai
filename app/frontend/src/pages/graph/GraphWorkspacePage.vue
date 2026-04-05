<template>
  <div class="workspace-page">
    <header class="workspace-header">
      <div class="title-block">
        <a-button type="text" @click="router.push('/graph/manage')">返回项目</a-button>
        <div>
          <p class="eyebrow">Fault Tree Workspace</p>
          <h1>{{ graphInfo?.graphName || '故障树工作台' }}</h1>
          <p class="subtitle text-sm mt-1 text-slate-500">
            当前展示 {{ workingGraph?.isTmp ? '非正式版本' : '正式版本' }} ·
            <span class="font-medium text-slate-700">{{ currentVersionLabel }}</span>
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
      <main class="canvas-panel">
        <div class="canvas-toolbar">
          <div class="canvas-meta">
            <a-tag color="blue">{{ currentVersionLabel }}</a-tag>
            <a-tag :color="workingGraph?.isTmp ? 'orange' : 'cyan'">
              {{ workingGraph?.isTmp ? '非正式版本' : '正式版本' }}
            </a-tag>
            <span class="meta-text">创建者：{{ graphInfo?.userId || '-' }}</span>
            <span class="meta-text">
              {{ syncStatusText }}
            </span>
          </div>
          <div class="canvas-meta">
            <span class="meta-text">节点 {{ flowNodes.length }}</span>
            <span class="meta-text">边 {{ flowEdges.length }}</span>
          </div> 
        </div>

        <div class="flow-shell relative" @click="closeContextMenu">
          <VueFlow
            v-model:nodes="flowNodes"
            v-model:edges="flowEdges"
            fit-view-on-init
            :min-zoom="0.2"
            :max-zoom="1.6"
            :default-edge-options="{ markerEnd: MarkerType.ArrowClosed }"
            @node-drag-stop="handleMarkDirty"
            @connect="onConnect"
            @edges-change="handleMarkDirty"
            @nodes-change="handleMarkDirty"
            @pane-context-menu="onPaneContextMenu"
            @node-context-menu="onNodeContextMenu"
            @edge-context-menu="onEdgeContextMenu"
            @pane-click="closeContextMenu"
          >
            <template #node-custom="props">
              <div class="custom-node-content h-full w-full flex flex-col items-center justify-center p-2 relative box-border">
                <Handle type="source" :position="Position.Top" />
                <div class="custom-node-label text-center font-bold break-all w-full line-clamp-2">{{ props.data.label }}</div>
                <div v-if="props.data.description" class="custom-node-desc text-[11px] text-gray-500 text-center mt-1 leading-tight break-all w-full line-clamp-2">
                  {{ props.data.description }}
                </div>
                <Handle type="target" :position="Position.Bottom" />
              </div>
            </template>
            <Background pattern-color="#dbeafe" :gap="24" />
            <MiniMap pannable zoomable />
            <Controls />
          </VueFlow>
        </div>
      </main>

      <section v-if="activePanel" class="side-panel" :style="{ width: `${panelWidth}px` }">
        <div class="panel-header">
          <div>
            <p class="panel-eyebrow">{{ currentPanelMeta.eyebrow }}</p>
            <h3>{{ currentPanelMeta.title }}</h3>
          </div>
          <a-button type="text" @click="activePanel = ''">收起</a-button>
        </div>

        <div class="panel-content flex flex-col">
          <ChatPanel 
            v-if="activePanel === 'chat'" 
            ref="chatPanelRef"
            :graph-id="graphId" 
            :selected-version="selectedVersion"
            :current-version="graphInfo?.currentVersion"
            :is-editor-dirty="isEditorDirty"
            :export-graph-model="() => exportGraphModel(graphInfo?.graphName)"
            @success="handleChatSuccess"
          />

          <VersionPanel 
            v-else-if="activePanel === 'version'" 
            :graph-id="graphId"
            :versions="versions"
            :current-version-label="currentVersionLabel"
            @refresh="loadVersions"
            @switch="switchVersion"
          />

          <SuggestionPanel 
            v-else-if="activePanel === 'suggestion'"
            :content="suggestion?.content"
            :loading="validating"
            @validate="handleValidate"
          />

          <template v-else-if="activePanel === 'json'">
            <div class="flex-1 flex flex-col h-full overflow-hidden">
              <a-textarea
                v-model:value="workingContent"
                class="!text-xs font-mono flex-1 h-full"
                @blur="handleJsonEdited"
                placeholder="这里会实时展示当前工作图内容..."
              />
            </div>
          </template>
        </div>
        <div class="resize-handle" @mousedown="startResize" style="left: 0; right: auto; cursor: ew-resize;"></div>
      </section>

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
    </div>

    <a-modal
      v-model:open="showSaveModal"
      title="保存版本"
      ok-text="保存"
      cancel-text="取消"
      @ok="handleSaveCurrent"
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

    <NodeEditModal 
      v-model:open="showEditNodeModal"
      :data="editingNodeData"
      @save="handleSaveNodeEdit"
      @delete="handleDeleteNode"
    />

    <!-- Context Menu -->
    <div
      v-if="contextMenu.show"
      :style="{ top: `${contextMenu.y}px`, left: `${contextMenu.x}px` }"
      class="fixed z-50 bg-white shadow-lg rounded-md border border-gray-200 py-1 w-32"
    >
      <div
        v-if="contextMenu.type === 'pane'"
        class="px-4 py-2 text-sm text-gray-700 hover:bg-blue-50 cursor-pointer"
        @click.stop="handleContextMenuAction('add')"
      >
        新建节点
      </div>
      <div
        v-if="contextMenu.type === 'node'"
        class="px-4 py-2 text-sm text-gray-700 hover:bg-blue-50 cursor-pointer"
        @click.stop="handleContextMenuAction('edit')"
      >
        编辑节点
      </div>
      <div
        v-if="contextMenu.type === 'node'"
        class="px-4 py-2 text-sm text-red-600 hover:bg-red-50 cursor-pointer"
        @click.stop="handleContextMenuAction('delete')"
      >
        删除节点
      </div>
      <div
        v-if="contextMenu.type === 'edge'"
        class="px-4 py-2 text-sm text-red-600 hover:bg-red-50 cursor-pointer"
        @click.stop="handleContextMenuAction('delete_edge')"
      >
        删除连线
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import type { Connection } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import { MarkerType, VueFlow, Handle, Position } from '@vue-flow/core'
import {
  CommentOutlined,
  DeploymentUnitOutlined,
  SafetyCertificateOutlined,
  CodeOutlined,
} from '@ant-design/icons-vue'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'
import ChatPanel from './components/ChatPanel.vue'
import VersionPanel from './components/VersionPanel.vue'
import SuggestionPanel from './components/SuggestionPanel.vue'
import NodeEditModal from './components/NodeEditModal.vue'
import { useGraphEditor } from './composables/useGraphEditor'
import { useWorkspaceData } from './composables/useWorkspaceData'
import {
  discardWorkingGraph,
  saveGraph,
  validateGraph,
} from '@/api/graphController'
import { uploadProjectDocument } from '@/api/documentController'

import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'
import '@vue-flow/minimap/dist/style.css'

type PanelKey = '' | 'chat' | 'version' | 'suggestion' | 'json'

const route = useRoute()
const router = useRouter()
const graphId = Number(route.params.id)

const {
  graphInfo,
  workingGraph,
  suggestion,
  versions,
  selectedVersion,
  loadGraphInfo,
  loadVersions,
  loadSuggestion,
  fetchWorkingGraph
} = useWorkspaceData(graphId)

const currentVersionLabel = computed(
  () => selectedVersion.value || graphInfo.value?.currentVersion || 'v001',
)

const {
  flowNodes,
  flowEdges,
  isEditorDirty,
  workingContent,
  syncStatusText,
  nodeColor,
  exportGraphModel,
  parseGraphContent,
  markDirty,
  clearDraftCache,
  scheduleDraftSync,
  restoreDraftCache,
  draftSyncTimer
} = useGraphEditor(graphId, currentVersionLabel)

const activePanel = ref<PanelKey>('')
const panelWidth = ref(380)
const isResizing = ref(false)
const validating = ref(false)
const showSaveModal = ref(false)
const saveMode = ref<'overwrite' | 'new'>('overwrite')
const saveVersionName = ref('')
const chatPanelRef = ref<any>(null)

const showEditNodeModal = ref(false)
const editingNodeId = ref('')
const editingNodeData = ref({ label: '', nodeType: 'intermediate_event', description: '', gateType: '' })

const contextMenu = ref({ show: false, x: 0, y: 0, type: '', nodeId: '' })

const closeContextMenu = () => {
  contextMenu.value.show = false
}

const onPaneContextMenu = (event: MouseEvent) => {
  event.preventDefault()
  contextMenu.value = { show: true, x: event.clientX, y: event.clientY, type: 'pane', nodeId: '' }
}

const onNodeContextMenu = ({ event, node }: any) => {
  event.preventDefault()
  contextMenu.value = { show: true, x: event.clientX, y: event.clientY, type: 'node', nodeId: node.id }
}

const onEdgeContextMenu = ({ event, edge }: any) => {
  event.preventDefault()
  contextMenu.value = { show: true, x: event.clientX, y: event.clientY, type: 'edge', nodeId: edge.id }
}

const handleMarkDirty = () => {
  markDirty(workingGraph.value, selectedVersion.value, graphInfo.value?.currentVersion || 'v001', graphInfo.value?.graphName)
}

const handleValidate = async () => {
  validating.value = true
  try {
    const res = await validateGraph({
      graphId,
      version: currentVersionLabel.value,
    })
    if (res.data.code === 0) {
      message.success('校验请求已提交，正在生成建议...')
      await loadSuggestion()
    } else {
      message.error(res.data.message || '校验失败')
    }
  } finally {
    validating.value = false
  }
}

const addNode = () => {
  const id = `node-${Date.now()}`
  flowNodes.value.push({
    id,
    position: { x: 250, y: 150 },
    type: 'custom',
    data: { label: '新节点', nodeType: 'intermediate_event', description: '', gateType: '' },
    style: { background: nodeColor('intermediate_event'), width: '180px' },
  })
  handleMarkDirty()
}

const handleContextMenuAction = (action: string) => {
  if (action === 'add') {
    addNode()
  } else if (action === 'edit') {
    const node = flowNodes.value.find((n) => n.id === contextMenu.value.nodeId)
    if (node) {
      editingNodeId.value = node.id
      editingNodeData.value = {
        label: node.data?.label || '',
        description: node.data?.description || '',
        nodeType: node.data?.nodeType || 'intermediate_event',
        gateType: node.data?.gateType || '',
      }
      showEditNodeModal.value = true
    }
  } else if (action === 'delete') {
    editingNodeId.value = contextMenu.value.nodeId
    handleDeleteNode()
  } else if (action === 'delete_edge') {
    flowEdges.value = flowEdges.value.filter((e) => e.id !== contextMenu.value.nodeId)
    handleMarkDirty()
  }
  closeContextMenu()
}

const handleSaveNodeEdit = () => {
  const node = flowNodes.value.find((n) => n.id === editingNodeId.value)
  if (node) {
    node.data = { ...node.data, ...editingNodeData.value }
    node.style = {
      ...node.style,
      background: nodeColor(editingNodeData.value.nodeType),
      width: editingNodeData.value.nodeType === 'gate' ? '120px' : '180px',
    }
    handleMarkDirty()
  }
  showEditNodeModal.value = false
}

const handleDeleteNode = () => {
  flowNodes.value = flowNodes.value.filter((n) => n.id !== editingNodeId.value)
  flowEdges.value = flowEdges.value.filter(
    (e) => e.source !== editingNodeId.value && e.target !== editingNodeId.value
  )
  handleMarkDirty()
  showEditNodeModal.value = false
}

const loadWorkingGraphData = async (version?: string) => {
  const data = await fetchWorkingGraph(version)
  if (data) {
    parseGraphContent(data.content || '', graphInfo.value?.graphName)
    restoreDraftCache(graphInfo.value?.graphName)
  }
}

const refreshAll = async () => {
  await Promise.all([loadGraphInfo(), loadVersions()])
  await Promise.all([loadWorkingGraphData(), loadSuggestion()])
}

const switchVersion = async (version?: string) => {
  selectedVersion.value = version || graphInfo.value?.currentVersion || 'v001'
  await Promise.all([loadWorkingGraphData(), loadSuggestion()])
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

const handleSaveCurrent = async () => {
  const fromVersion = currentVersionLabel.value
  const toVersion = saveMode.value === 'new' ? saveVersionName.value.trim() : fromVersion
  if (!toVersion) {
    message.warning('请输入目标版本名称')
    return
  }
  const content = JSON.stringify(exportGraphModel(graphInfo.value?.graphName), null, 2)
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

const handleChatSuccess = async () => {
  clearDraftCache()
  await loadWorkingGraphData()
  isEditorDirty.value = false
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
  parseGraphContent(text, graphInfo.value?.graphName)
  handleMarkDirty()
  scheduleDraftSync()
  message.success('图文件已导入到当前编辑器')
  return false
}

const exportCurrentGraph = () => {
  const content = JSON.stringify(exportGraphModel(graphInfo.value?.graphName), null, 2)
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
    parseGraphContent(workingContent.value, graphInfo.value?.graphName)
    handleMarkDirty()
    scheduleDraftSync()
  } catch (error) {
    console.error(error)
    message.warning('JSON 格式暂时不合法，先检查内容再继续')
  }
}

const onConnect = (connection: Connection) => {
  if (!connection.source || !connection.target) return
  flowEdges.value = [
    ...flowEdges.value,
    {
      id: `${connection.source}-${connection.target}-${Date.now()}`,
      source: connection.source,
      target: connection.target,
      markerEnd: MarkerType.ArrowClosed,
    },
  ]
  handleMarkDirty()
}

const panelItems = [
  { key: 'chat', label: '对话', icon: CommentOutlined },
  { key: 'version', label: '版本', icon: DeploymentUnitOutlined },
  { key: 'suggestion', label: '校验', icon: SafetyCertificateOutlined },
  { key: 'json', label: '源码', icon: CodeOutlined },
] as const

const panelMeta = {
  chat: { eyebrow: 'Dialogue', title: 'AI 对话记录' },
  version: { eyebrow: 'Version', title: '版本列表' },
  suggestion: { eyebrow: 'Validate', title: '当前建议' },
  json: { eyebrow: 'Source', title: '当前图 JSON 预览' },
}

const currentPanelMeta = computed(() => {
  if (activePanel.value && panelMeta[activePanel.value]) {
    return panelMeta[activePanel.value]
  }
  return { eyebrow: '', title: '' }
})

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

watch(activePanel, (newVal) => {
  if (newVal === 'chat') {
    nextTick(() => {
      chatPanelRef.value?.scrollToBottom()
    })
  }
})

onMounted(() => {
  refreshAll()
  window.addEventListener('mousemove', handleMouseMove)
  window.addEventListener('mouseup', stopResize)
})

onBeforeUnmount(() => {
  if (draftSyncTimer) clearTimeout(draftSyncTimer)
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
  height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
  padding: 16px 24px;
  background:
    radial-gradient(circle at top right, rgba(8, 145, 178, 0.14), transparent 24%),
    linear-gradient(180deg, #f8fafc 0%, #eef2ff 100%);
  overflow: hidden;
}

.workspace-header {
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: end;
  gap: 24px;
  margin-bottom: 16px;
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
  font-size: 26px;
  color: #0f172a;
}

.subtitle {
  margin: 0;
  color: #64748b;
}

.workspace-body {
  flex: 1;
  display: flex;
  flex-direction: row;
  gap: 16px;
  min-height: 0;
}

.tool-rail {
  width: 64px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 0;
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
  flex-shrink: 0;
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

.resize-handle {
  position: absolute;
  top: 0;
  right: 0;
  width: 8px;
  height: 100%;
  cursor: ew-resize;
}

.canvas-panel {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.canvas-toolbar {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 18px;
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
  flex: 1;
  border-radius: 24px;
  overflow: hidden;
  border: 1px solid #dbeafe;
  background: linear-gradient(180deg, #f8fbff 0%, #eff6ff 100%);
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.08);
}

:deep(.vue-flow__node) {
  font-size: 13px;
}

:deep(.vue-flow__controls) {
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.12);
}
:deep(.vue-flow__node-custom) {
  border-radius: 12px;
  border: 1px solid #94a3b8;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease;
  overflow: visible !important;
  background: white;
}

:deep(.vue-flow__node-custom.selected) {
  border-color: #3b82f6;
  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.2);
}

:deep(.vue-flow__node-custom:hover) {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  transform: translateY(-2px);
}

.custom-node-label {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 2px;
}

.custom-node-desc {
  font-size: 11px;
  color: #64748b;
  line-height: 1.3;
}

.vue-flow__handle {
  width: 10px;
  height: 10px;
  background-color: #94a3b8;
  border: 2px solid white;
  transition: all 0.2s ease;
  z-index: 10;
}

.vue-flow__handle:hover {
  background-color: #3b82f6;
  transform: scale(1.5);
}

.vue-flow__handle-top {
  top: -6px;
}

.vue-flow__handle-bottom {
  bottom: -6px;
}

.vue-flow__edge-path {
  stroke-width: 2;
  stroke: #94a3b8;
  transition: stroke 0.2s ease, stroke-width 0.2s ease;
}

.vue-flow__edge:hover .vue-flow__edge-path,
.vue-flow__edge.selected .vue-flow__edge-path {
  stroke: #3b82f6;
  stroke-width: 3;
}
</style>
