<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Minimal Header -->
    <header class="h-14 flex items-center justify-between px-4 border-b border-gray-200 bg-white flex-shrink-0">
      <div class="flex items-center gap-3">
        <button
          class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-gray-100 text-gray-500 transition-colors"
          @click="router.push('/graph/manage')"
        >
          ←
        </button>
        <div class="flex items-center gap-2">
          <h1 class="text-[15px] font-semibold text-gray-900 m-0">{{ graphInfo?.graphName || '故障树工作台' }}</h1>
          <span class="px-2 py-0.5 text-[11px] font-medium bg-indigo-50 text-indigo-600 rounded-md flex-shrink-0">{{ currentVersionLabel }}</span>
          <span v-if="workingGraph?.isTmp" class="text-[11px] text-amber-600 bg-amber-50 px-2 py-0.5 rounded-md flex-shrink-0">暂存</span>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-xs text-gray-400 hidden sm:inline mr-1">{{ flowNodes.length }} 节点 · {{ flowEdges.length }} 边</span>
        <a-button type="primary" size="small" class="!rounded-lg !h-8 !px-4 !text-[13px]" @click="openSaveModal">保存</a-button>
        <a-dropdown placement="bottomRight">
          <a-button size="small" class="!rounded-lg !h-8 !border-gray-200 !px-2 !inline-flex !items-center !justify-center">
            <EllipsisOutlined />
          </a-button>
          <template #overlay>
            <a-menu>
              <a-menu-item @click="handleAutoArrange">自动整理</a-menu-item>
              <a-menu-item>
                <a-upload :show-upload-list="false" accept=".pdf" :before-upload="beforeUploadProjectDoc" class="w-full">
                  上传项目文档
                </a-upload>
              </a-menu-item>
              <a-menu-item>
                <a-upload :show-upload-list="false" accept=".json" :before-upload="beforeImportGraph" class="w-full">
                  导入图
                </a-upload>
              </a-menu-item>
              <a-menu-item @click="exportCurrentGraph">导出当前图</a-menu-item>
              <a-menu-divider />
              <a-menu-item danger @click="discardTmp">放弃暂存</a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
    </header>

    <!-- Body -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Left Sidebar -->
      <div class="flex flex-shrink-0">
        <!-- Icon Rail (always visible) -->
        <div class="w-11 flex flex-col items-center py-2 gap-1 bg-white border-r border-gray-200">
          <button
            v-for="item in panelItems"
            :key="item.key"
            class="w-9 h-9 flex items-center justify-center rounded-lg text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 transition-colors"
            :class="{ '!text-indigo-600 !bg-indigo-50': activePanel === item.key }"
            :title="item.label"
            @click="togglePanel(item.key)"
          >
            <component :is="item.icon" :style="{ fontSize: '18px' }" />
          </button>
        </div>

        <!-- Expanded Panel -->
        <Transition name="panel-slide">
          <div
            v-if="activePanel"
            class="flex bg-white border-r border-gray-200"
            :style="{ width: sidebarWidth + 'px' }"
          >
            <div class="flex flex-col flex-1 overflow-hidden min-w-0">
              <!-- Panel Header -->
              <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 flex-shrink-0">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wider">{{ currentPanelMeta.title }}</span>
                <button class="w-6 h-6 flex items-center justify-center rounded hover:bg-gray-100 text-gray-400 transition-colors" @click="activePanel = ''">
                  <CloseOutlined :style="{ fontSize: '12px' }" />
                </button>
              </div>
              <!-- Panel Content -->
              <div class="flex-1 overflow-hidden">
                <KeepAlive>
                  <ChatPanel
                    v-if="activePanel === 'chat'"
                    ref="chatPanelRef"
                    :graph-id="graphId"
                    :selected-version="selectedVersion"
                    :current-version="graphInfo?.currentVersion"
                    :is-editor-dirty="isEditorDirty"
                    @success="handleChatSuccess"
                  />
                </KeepAlive>
                <VersionPanel
                  v-if="activePanel === 'version'"
                  :graph-id="graphId"
                  :versions="versions"
                  :current-version-label="currentVersionLabel"
                  @refresh="loadVersions"
                  @switch="switchVersion"
                />
                <SuggestionPanel
                  v-if="activePanel === 'suggestion'"
                  :content="suggestion?.content"
                  :loading="validating"
                  @validate="handleValidate"
                />
                <div v-if="activePanel === 'json'" class="h-full flex flex-col p-3">
                  <a-textarea
                    v-model:value="workingContent"
                    class="flex-1 !text-xs !font-mono !rounded-lg !border-gray-200 resize-none"
                    @blur="handleJsonEdited"
                    placeholder="当前工作图 JSON 内容..."
                  />
                </div>
              </div>
            </div>
            <!-- Resize Handle -->
            <div
              class="w-1.5 cursor-ew-resize flex-shrink-0 relative z-10 hover:bg-indigo-400 transition-colors bg-transparent"
              @mousedown="startResize"
            ></div>
          </div>
        </Transition>
      </div>

      <!-- Canvas -->
      <main class="flex-1 min-w-0 p-3">
        <div class="h-full rounded-xl border border-gray-200 overflow-hidden bg-white shadow-sm" @click="closeContextMenu">
          <VueFlow
            v-model:nodes="flowNodes"
            v-model:edges="flowEdges"
            fit-view-on-init
            :min-zoom="0.2"
            :max-zoom="1.6"
            :default-edge-options="{ markerEnd: MarkerType.ArrowClosed, type: 'smoothstep' }"
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
              <div
                class="relative box-border"
                :class="{ 'w-[92px] h-[72px] flex items-center justify-center': props.data.nodeType === 'gate', 'min-w-[160px]': props.data.nodeType !== 'gate' }"
              >
                <Handle type="target" :position="Position.Left" />
                <template v-if="props.data.nodeType === 'gate'">
                  <div class="w-[72px] h-[56px]" :title="props.data.gateType || 'GATE'">
                    <svg
                      v-if="props.data.gateType === 'AND'"
                      viewBox="0 0 72 56"
                      class="w-[72px] h-[56px]"
                      aria-hidden="true"
                    >
                      <path d="M8 4 H32 A20 20 0 0 1 32 52 H8 Z" fill="#fff" stroke="#0f766e" stroke-width="2.5" />
                    </svg>
                    <svg
                      v-else
                      viewBox="0 0 72 56"
                      class="w-[72px] h-[56px]"
                      aria-hidden="true"
                    >
                      <path d="M10 4 Q38 4 64 28 Q38 52 10 52 Q20 40 20 28 Q20 16 10 4 Z" fill="#fff" stroke="#0f766e" stroke-width="2.5" />
                    </svg>
                  </div>
                </template>
                <template v-else>
                  <div class="min-h-[68px] flex flex-col items-center justify-center px-3 py-2">
                    <div class="text-sm font-semibold text-gray-800 text-center break-all line-clamp-2">{{ props.data.label }}</div>
                    <div v-if="props.data.description" class="text-[11px] text-gray-500 text-center mt-1 leading-tight line-clamp-2">
                      {{ props.data.description }}
                    </div>
                  </div>
                </template>
                <Handle type="source" :position="Position.Right" />
              </div>
            </template>
            <Background pattern-color="#e5e7eb" :gap="24" />
            <MiniMap pannable zoomable />
            <Controls />
          </VueFlow>
        </div>
      </main>
    </div>

    <!-- Save Modal -->
    <a-modal
      v-model:open="showSaveModal"
      title="保存版本"
      ok-text="保存"
      cancel-text="取消"
      @ok="handleSaveCurrent"
    >
      <a-radio-group v-model:value="saveMode" class="mb-4">
        <a-radio value="overwrite">覆盖当前版本</a-radio>
        <a-radio value="new">新建版本保存</a-radio>
      </a-radio-group>
      <a-input
        v-if="saveMode === 'new'"
        v-model:value="saveVersionName"
        placeholder="输入新版本名称，例如：v002 或 优化版"
      />
    </a-modal>

    <!-- Node Edit Modal -->
    <NodeEditModal
      v-model:open="showEditNodeModal"
      :data="editingNodeData"
      @save="handleSaveNodeEdit"
      @delete="handleDeleteNode"
    />

    <!-- Context Menu -->
    <div
      v-if="contextMenu.show"
      :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
      class="fixed z-50 bg-white shadow-lg rounded-lg border border-gray-200 py-1 w-32"
    >
      <div
        v-if="contextMenu.type === 'pane'"
        class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer transition-colors"
        @click.stop="handleContextMenuAction('add')"
      >
        新建节点
      </div>
      <div
        v-if="contextMenu.type === 'node'"
        class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer transition-colors"
        @click.stop="handleContextMenuAction('edit')"
      >
        编辑节点
      </div>
      <div
        v-if="contextMenu.type === 'node'"
        class="px-4 py-2 text-sm text-red-600 hover:bg-red-50 cursor-pointer transition-colors"
        @click.stop="handleContextMenuAction('delete')"
      >
        删除节点
      </div>
      <div
        v-if="contextMenu.type === 'edge'"
        class="px-4 py-2 text-sm text-red-600 hover:bg-red-50 cursor-pointer transition-colors"
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
  EllipsisOutlined,
  CloseOutlined,
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
  autoArrangeNodes,
  markDirty,
  clearDraftCache,
  scheduleDraftSync,
  restoreDraftCache,
  draftSyncTimer
} = useGraphEditor(graphId, currentVersionLabel)

const activePanel = ref<PanelKey>('')
const sidebarWidth = ref(380)
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

const handleAutoArrange = async () => {
  autoArrangeNodes(graphInfo.value?.graphName)
  handleMarkDirty()
  scheduleDraftSync()
  message.success('已按横向树结构自动整理当前图')
}

const addNode = () => {
  const id = `node-${Date.now()}`
  flowNodes.value.push({
    id,
    position: { x: 250, y: 150 },
    type: 'custom',
    data: { label: '新节点', nodeType: 'intermediate_event', description: '', gateType: '' },
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
    style: { background: nodeColor('intermediate_event'), width: '180px', minHeight: '72px' },
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
      ...(editingNodeData.value.nodeType === 'gate'
        ? {
            background: 'transparent',
            width: '92px',
            height: '72px',
            border: 'none',
            boxShadow: 'none',
          }
        : {
            background: nodeColor(editingNodeData.value.nodeType),
            width: '180px',
            minHeight: '72px',
            height: undefined,
          }),
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

const panelMeta: Record<string, { eyebrow: string; title: string }> = {
  chat: { eyebrow: '', title: '' },
  version: { eyebrow: '', title: '版本管理' },
  suggestion: { eyebrow: '', title: '校验建议' },
  json: { eyebrow: '', title: '图源码' },
}

const currentPanelMeta = computed(() => {
  if (activePanel.value && panelMeta[activePanel.value]) {
    return panelMeta[activePanel.value]
  }
  return { eyebrow: '', title: '' }
})

const togglePanel = (key: PanelKey) => {
  if (activePanel.value === key) {
    activePanel.value = ''
  } else {
    activePanel.value = key
  }
}

const startResize = () => {
  isResizing.value = true
}

const handleMouseMove = (event: MouseEvent) => {
  if (!isResizing.value) return
  const maxWidth = window.innerWidth * 0.5
  sidebarWidth.value = Math.max(280, Math.min(maxWidth, event.clientX))
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
/* Panel slide transition */
.panel-slide-enter-active {
  transition: all 0.2s ease-out;
}
.panel-slide-leave-active {
  transition: all 0.15s ease-in;
}
.panel-slide-enter-from {
  opacity: 0;
  transform: translateX(-12px);
}
.panel-slide-leave-to {
  opacity: 0;
  transform: translateX(-12px);
}

/* Minimal overrides for VueFlow */
:deep(.vue-flow__node) {
  font-size: 13px;
}

:deep(.vue-flow__controls) {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  overflow: hidden;
}

:deep(.vue-flow__node-custom) {
  border-radius: 12px;
  border: 1.5px solid #d1d5db;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease;
  overflow: visible !important;
  background: white;
}

:deep(.vue-flow__node-custom.selected) {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.2);
}

:deep(.vue-flow__node-custom:hover) {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1);
}

:deep(.vue-flow__handle) {
  width: 8px;
  height: 8px;
  background-color: #9ca3af;
  border: 2px solid white;
  transition: all 0.2s ease;
  z-index: 10;
}

:deep(.vue-flow__handle:hover) {
  background-color: #6366f1;
  transform: scale(1.4);
}

:deep(.vue-flow__handle-left) {
  left: -5px;
}

:deep(.vue-flow__handle-right) {
  right: -5px;
}

:deep(.vue-flow__edge-path) {
  stroke-width: 2;
  stroke: #9ca3af;
  transition: stroke 0.2s ease, stroke-width 0.2s ease;
}

:deep(.vue-flow__edge:hover .vue-flow__edge-path),
:deep(.vue-flow__edge.selected .vue-flow__edge-path) {
  stroke: #6366f1;
  stroke-width: 2.5;
}
</style>
