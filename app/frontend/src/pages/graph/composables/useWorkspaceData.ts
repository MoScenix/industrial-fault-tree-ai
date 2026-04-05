import { ref } from 'vue'
import { message } from 'ant-design-vue'
import {
  getGraphVoById,
  listGraphVersion,
  getWorkingGraph,
  getCurrentSuggestion,
} from '@/api/graphController'

export function useWorkspaceData(graphId: number) {
  const graphInfo = ref<API.GraphVO>()
  const workingGraph = ref<API.WorkingGraphVO>()
  const suggestion = ref<API.GraphSuggestionVO>()
  const versions = ref<API.GraphVersionVO[]>([])
  const selectedVersion = ref('')

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

  const loadSuggestion = async (version?: string) => {
    const res = await getCurrentSuggestion({ graphId, version: version || selectedVersion.value })
    if (res.data.code === 0) {
      suggestion.value = res.data.data
    }
  }

  const fetchWorkingGraph = async (version?: string) => {
    const res = await getWorkingGraph({ graphId, version: version || selectedVersion.value })
    if (res.data.code === 0 && res.data.data) {
      workingGraph.value = res.data.data
      return res.data.data
    } else {
      message.error(res.data.message || '获取当前工作图失败')
      return null
    }
  }

  return {
    graphInfo,
    workingGraph,
    suggestion,
    versions,
    selectedVersion,
    loadGraphInfo,
    loadVersions,
    loadSuggestion,
    fetchWorkingGraph
  }
}
