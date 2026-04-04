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
            <span class="font-medium text-slate-700">{{ selectedVersion || graphInfo?.currentVersion || 'v001' }}</span>
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
            <a-tag color="blue">{{ selectedVersion || graphInfo?.currentVersion || 'v001' }}</a-tag>
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
            <a-button type="primary" size="small" @click="handleContextMenuAction('add')">添加节点</a-button>
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
            @node-drag-stop="markDirty"
            @connect="onConnect"
            @edges-change="markDirty"
            @nodes-change="markDirty"
            @pane-context-menu="onPaneContextMenu"
            @node-context-menu="onNodeContextMenu"
            @edge-context-menu="onEdgeContextMenu"
            @pane-click="closeContextMenu"
          >
            <template #node-custom="props">
              <div class="custom-node-content h-full w-full flex flex-col items-center justify-center p-2 relative box-border">
                <Handle type="target" :position="Position.Top" />
                <div class="custom-node-label text-center font-bold break-all w-full line-clamp-2">{{ props.data.label }}</div>
                <div v-if="props.data.description" class="custom-node-desc text-[11px] text-gray-500 text-center mt-1 leading-tight break-all w-full line-clamp-2">
                  {{ props.data.description }}
                </div>
                <Handle type="source" :position="Position.Bottom" />
              </div>
            </template>
            <Background pattern-color="#dbeafe" :gap="24" />
            <MiniMap pannable zoomable />
            <Controls />
          </VueFlow>

          <!-- Developer Tools absolute positioned at bottom-left -->
          <!-- 隐藏原本的开发调试浮窗，移到右侧侧边栏了 -->
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
          <template v-if="activePanel === 'chat'">
            <div class="chat-list flex flex-col gap-4 overflow-y-auto pb-4 px-2">
              <div v-for="item in messages" :key="item.id" class="flex flex-col group">
                <div 
                  class="flex max-w-[85%]"
                  :class="item.role === 'assistant' ? 'self-start' : 'self-end'"
                >
                  <div class="flex items-start gap-3" :class="item.role === 'assistant' ? 'flex-row' : 'flex-row-reverse'">
                    <!-- 头像区域 -->
                    <div class="w-8 h-8 rounded-full flex-shrink-0 flex items-center justify-center text-xs font-bold text-white shadow-sm"
                         :class="item.role === 'assistant' ? 'bg-gradient-to-br from-indigo-500 to-purple-600' : 'bg-gradient-to-br from-emerald-500 to-teal-500'">
                      {{ item.role === 'assistant' ? 'AI' : '我' }}
                    </div>
                    
                    <!-- 消息内容气泡 -->
                    <div
                      class="px-4 py-3 shadow-sm text-[14px] leading-relaxed relative"
                      :class="[
                        item.role === 'assistant' 
                          ? 'bg-white text-slate-700 border border-slate-100 rounded-2xl rounded-tl-sm' 
                          : 'bg-indigo-50 text-indigo-900 border border-indigo-100 rounded-2xl rounded-tr-sm'
                      ]"
                    >
                      <div class="whitespace-pre-wrap font-sans break-words">{{ item.content }}</div>
                    </div>
                  </div>
                </div>
              </div>
              <a-empty v-if="!messages.length" description="输入需求，AI 将辅助您构建故障树" class="mt-10" />
            </div>
            
            <div class="mt-auto pt-4 border-t border-slate-100 flex-shrink-0 bg-white">
              <div class="relative bg-white rounded-xl border border-slate-200 shadow-sm focus-within:border-indigo-500 focus-within:ring-1 focus-within:ring-indigo-500 transition-all overflow-hidden">
                <a-textarea
                  v-model:value="chatInput"
                  :rows="3"
                  :maxlength="1000"
                  placeholder="例如：帮我添加一个温度过高的基础事件..."
                  class="w-full !border-none !shadow-none !bg-transparent resize-none py-3 px-4 focus:!shadow-none text-sm"
                  @pressEnter="(e: KeyboardEvent) => { if (!e.shiftKey) { e.preventDefault(); sendChatMessage() } }"
                />
                <div class="flex justify-between items-center px-3 py-2 bg-slate-50 border-t border-slate-100">
                  <span class="text-xs text-slate-400">Shift + Enter 换行，Enter 发送</span>
                  <a-button type="primary" shape="circle" :loading="chatting" @click="sendChatMessage" 
                            class="!flex !items-center !justify-center !w-8 !h-8 !min-w-0 !bg-indigo-600 hover:!bg-indigo-500 border-none shadow-md">
                    <svg v-if="!chatting" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-4 h-4">
                      <path d="M3.478 2.404a.75.75 0 0 0-.926.941l2.432 7.905H13.5a.75.75 0 0 1 0 1.5H4.984l-2.432 7.905a.75.75 0 0 0 .926.94 60.519 60.519 0 0 0 18.445-8.986.75.75 0 0 0 0-1.218A60.517 60.517 0 0 0 3.478 2.404Z" />
                    </svg>
                  </a-button>
                </div>
              </div>
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
                        <a-tag v-if="item.version === currentVersionLabel" color="blue">当前</a-tag>
                      </div>
                      <div class="version-meta">{{ item.version }}</div>
                    </div>
                    <div class="version-actions">
                      <a-button type="link" size="small" @click.stop="openRenameVersion(item)">重命名</a-button>
                      <a-button
                        type="link"
                        size="small"
                        danger
                        :disabled="item.version === currentVersionLabel"
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

    <a-modal
      v-model:open="showEditNodeModal"
      title="编辑节点"
      ok-text="确定"
      cancel-text="取消"
      @ok="saveNodeEdit"
      @cancel="showEditNodeModal = false"
    >
      <template #footer>
        <div style="display: flex; justify-content: space-between;">
          <a-button danger @click="deleteNode">删除节点</a-button>
          <div>
            <a-button @click="showEditNodeModal = false">取消</a-button>
            <a-button type="primary" @click="saveNodeEdit">确定</a-button>
          </div>
        </div>
      </template>
      <a-form layout="vertical">
        <a-form-item label="节点名称">
          <a-input v-model:value="editingNodeData.label" placeholder="输入节点名称" />
        </a-form-item>
        <a-form-item label="节点描述">
          <a-textarea v-model:value="editingNodeData.description" placeholder="输入节点描述" :rows="3" />
        </a-form-item>
        <a-form-item label="节点类型">
          <a-select v-model:value="editingNodeData.nodeType">
            <a-select-option value="top_event">顶事件 (Top Event)</a-select-option>
            <a-select-option value="intermediate_event">中间事件 (Intermediate Event)</a-select-option>
            <a-select-option value="basic_event">基础事件 (Basic Event)</a-select-option>
            <a-select-option value="gate">逻辑门 (Gate)</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="editingNodeData.nodeType === 'gate'" label="逻辑门类型">
          <a-select v-model:value="editingNodeData.gateType">
            <a-select-option value="AND">AND (与门)</a-select-option>
            <a-select-option value="OR">OR (或门)</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

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
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import type { Connection, Edge, Node } from '@vue-flow/core'
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

type PanelKey = '' | 'chat' | 'version' | 'suggestion' | 'json'

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

const isLocalDraft = computed(() => isEditorDirty.value) // 虽然不再在界面直接显示该文字，但保留此状态供逻辑判断
const syncStatusText = ref('缓存已同步')

const showEditNodeModal = ref(false)
const editingNodeId = ref('')
const editingNodeData = ref({ label: '', nodeType: 'intermediate_event', description: '', gateType: '' })

const contextMenu = ref({ show: false, x: 0, y: 0, type: '', nodeId: '' })

const closeContextMenu = () => {
  contextMenu.value.show = false
}

const onPaneContextMenu = (event: MouseEvent) => {
  event.preventDefault()
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    type: 'pane',
    nodeId: '',
  }
}

const onNodeContextMenu = ({ event, node }: any) => {
  event.preventDefault()
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    type: 'node',
    nodeId: node.id,
  }
}

const onEdgeContextMenu = ({ event, edge }: any) => {
  event.preventDefault()
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    type: 'edge',
    nodeId: edge.id, // 用 nodeId 字段复用存 edgeId
  }
}

const handleContextMenuAction = (action: string) => {
  if (action === 'add') {
    const id = `node-${Date.now()}`
    
    // 如果没有使用特定的坐标转换，至少在视觉中心或者点击处
    // 此处简化为获取到点击的 contextMenu 坐标
    // 假设 contextMenu 的坐标就是点击坐标（其实相对于视口，可能需要减去 flow 容器的偏移，这里简单处理）
    flowNodes.value.push({
      id,
      position: { x: 200, y: 200 }, // 可以进一步优化根据点击坐标换算
      type: 'custom',
      data: {
        label: '新节点',
        nodeType: 'intermediate_event',
        description: '',
        gateType: '',
      },
      style: {
        background: nodeColor('intermediate_event'),
        width: '180px',
      },
    })
    markDirty()
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
    deleteNode()
  } else if (action === 'delete_edge') {
    flowEdges.value = flowEdges.value.filter((e) => e.id !== contextMenu.value.nodeId)
    markDirty()
  }
  closeContextMenu()
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
    type: 'custom',
    style: {
      background: nodeColor(item.node_type),
      width: item.node_type === 'gate' ? '120px' : '180px',
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

const markDirty = async () => {
  const alreadyDirty = isEditorDirty.value
  isEditorDirty.value = true
  workingContent.value = JSON.stringify(exportGraphModel(), null, 2)

  // 只有当前是正式版本，且是本次编辑的第一次操作时，才异步通知后端创建 tmp
  if (!alreadyDirty && workingGraph.value && !workingGraph.value.isTmp) {
    try {
      await ensureEditReady()
      // 同步更新本地状态，确保界面立即反映为“非正式版本”
      if (workingGraph.value) {
        workingGraph.value.isTmp = true
      }
    } catch (e) {
      console.error('自动准备编辑副本失败:', e)
    }
  }
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
  syncStatusText.value = '正在同步...'
  if (draftSyncTimer) clearTimeout(draftSyncTimer)
  draftSyncTimer = setTimeout(() => {
    persistDraftCache()
    syncStatusText.value = '缓存已同步'
  }, 5000)
}

const saveNodeEdit = () => {
  const node = flowNodes.value.find((n) => n.id === editingNodeId.value)
  if (node) {
    node.data = { ...node.data, ...editingNodeData.value }
    node.style = {
      ...node.style,
      background: nodeColor(editingNodeData.value.nodeType),
      width: editingNodeData.value.nodeType === 'gate' ? '120px' : '180px',
    }
    markDirty()
  }
  showEditNodeModal.value = false
}

const deleteNode = () => {
  flowNodes.value = flowNodes.value.filter((n) => n.id !== editingNodeId.value)
  flowEdges.value = flowEdges.value.filter(
    (e) => e.source !== editingNodeId.value && e.target !== editingNodeId.value
  )
  markDirty()
  showEditNodeModal.value = false
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
  await Promise.all([loadWorkingGraph(), loadSuggestion(), loadMessages()])
}

const switchVersion = async (version?: string) => {
  selectedVersion.value = version || graphInfo.value?.currentVersion || 'v001'
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
      
      // AI修改完成后刷新图数据
      clearDraftCache()
      await loadWorkingGraph()
      isEditorDirty.value = false
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
  markDirty()
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
    markDirty()
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
  height: calc(100vh - 64px); /* assuming there is a top nav */
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
  font-size: 13px;
  line-height: 1.6;
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
:deep(.vue-flow__node-custom) {
  border-radius: 12px;
  border: 1px solid #94a3b8;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease;
  overflow: visible !important; /* handles should be visible */
  background: white; /* fallback */
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

/* 连接点样式优化 */
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

/* 连线悬停效果 */
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
