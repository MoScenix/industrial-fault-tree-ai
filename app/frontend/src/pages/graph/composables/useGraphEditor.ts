import { ref, computed } from 'vue'
import type { Node, Edge } from '@vue-flow/core'
import { MarkerType } from '@vue-flow/core'
import { message } from 'ant-design-vue'
import { startEdit } from '@/api/graphController'

export function useGraphEditor(graphId: number, currentVersionLabel: any) {
  const flowNodes = ref<Node[]>([])
  const flowEdges = ref<Edge[]>([])
  const isEditorDirty = ref(false)
  const workingContent = ref('')
  const syncStatusText = ref('缓存已同步')
  const originalGraphMeta = ref<Record<string, any>>({})
  
  let draftSyncTimer: ReturnType<typeof setTimeout> | undefined

  const draftCacheKey = computed(
    () => `graph-workspace:${graphId}:${currentVersionLabel.value || 'v001'}`,
  )

  const nodeColor = (type?: string) => {
    if (type === 'top_event') return '#dbeafe'
    if (type === 'gate') return '#ccfbf1'
    if (type === 'basic_event') return '#ffffff'
    return '#eff6ff'
  }

  const buildDefaultGraph = () => {
    return {
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
  }

  const exportGraphModel = (graphName?: string) => {
    const pointsToMap = new Map<string, string[]>()
    const pointedByMap = new Map<string, string[]>()

    flowEdges.value.forEach((edge) => {
      const childId = edge.source
      const parentId = edge.target
      
      const pts = pointsToMap.get(parentId) || []
      pts.push(childId)
      pointsToMap.set(parentId, pts)

      const pby = pointedByMap.get(childId) || []
      pby.push(parentId)
      pointedByMap.set(childId, pby)
    })

    const nodes = flowNodes.value.map((node) => ({
      node_id: node.id,
      node_type: String(node.data?.nodeType || 'intermediate_event'),
      label: String(node.data?.label || node.id),
      description: String(node.data?.description || ''),
      gate_type: String(node.data?.gateType || ''),
      points_to: pointsToMap.get(node.id) || [],
      pointed_by: pointedByMap.get(node.id) || [],
      position: node.position,
    }))

    return {
      schema_version: originalGraphMeta.value.schema_version || 'fault-tree/v1',
      tree: {
        ...(originalGraphMeta.value.tree || {}),
        name: graphName || originalGraphMeta.value.tree?.name || '故障树',
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

  const parseGraphContent = (content?: string, graphName?: string) => {
    let parsed: any
    try {
      parsed = content ? JSON.parse(content) : buildDefaultGraph()
    } catch {
      parsed = buildDefaultGraph()
    }

    originalGraphMeta.value = {
      schema_version: parsed.schema_version || 'fault-tree/v1',
      tree: parsed.tree || { name: graphName || '故障树', top_node_id: '' },
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
      const parentId = item.node_id
      const childrenIds = Array.isArray(item.points_to) ? item.points_to : []
      childrenIds.forEach((childId: string, edgeIndex: number) => {
        edges.push({
          id: `${childId}-${parentId}-${edgeIndex}`,
          source: childId,
          target: parentId,
          animated: item.node_type === 'gate',
        })
      })
    })
    flowEdges.value = edges
    workingContent.value = JSON.stringify(parsed, null, 2)
    isEditorDirty.value = false
  }

  const ensureEditReady = async (selectedVersion: string, currentVersion: string) => {
    const res = await startEdit({
      graphId,
      version: selectedVersion || currentVersion || 'v001',
    })
    if (res.data.code !== 0) {
      throw new Error(res.data.message || '自动准备编辑副本失败')
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

  const markDirty = async (workingGraph: any, selectedVersion: string, currentVersion: string, graphName?: string) => {
    const alreadyDirty = isEditorDirty.value
    isEditorDirty.value = true
    workingContent.value = JSON.stringify(exportGraphModel(graphName), null, 2)

    if (!alreadyDirty && workingGraph && !workingGraph.isTmp) {
      try {
        await ensureEditReady(selectedVersion, currentVersion)
        if (workingGraph) {
          workingGraph.isTmp = true
        }
      } catch (e) {
        console.error('自动准备编辑副本失败:', e)
      }
    }
  }

  const restoreDraftCache = (graphName?: string) => {
    const raw = localStorage.getItem(draftCacheKey.value)
    if (!raw) return false
    try {
      const parsed = JSON.parse(raw)
      if (parsed?.content) {
        parseGraphContent(parsed.content, graphName)
        isEditorDirty.value = true
        message.info('已恢复当前版本的本地缓存编辑内容')
        return true
      }
    } catch {
      localStorage.removeItem(draftCacheKey.value)
    }
    return false
  }

  return {
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
  }
}
