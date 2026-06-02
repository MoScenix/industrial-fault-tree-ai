<template>
  <div class="h-full flex flex-col bg-gray-50 px-6 py-6 md:px-8 md:py-8">
    <!-- Top: simple header + button -->
    <div class="flex items-center justify-between gap-6 mb-4 flex-shrink-0">
      <h1 class="text-base font-bold text-gray-900 m-0">{{ isAdminView ? '项目总览' : '我的项目' }}</h1>
      <a-button type="primary" class="!rounded-lg flex-shrink-0" @click="openCreateModal">新建项目</a-button>
    </div>

    <!-- Search bar (fixed, does not scroll) -->
    <div class="flex-shrink-0 mb-6 bg-white rounded-xl border border-gray-200 p-4 shadow-sm">
      <a-form layout="inline" :model="searchParams" @finish="fetchGraphs">
        <a-form-item label="项目名">
          <a-input v-model:value="searchParams.graphName" placeholder="搜索项目名称" allow-clear class="!rounded-lg" />
        </a-form-item>
        <a-form-item v-if="isAdminView" label="用户 ID">
          <a-input-number v-model:value="searchParams.userId" placeholder="创建者 ID" class="!rounded-lg" style="width: 160px" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit" class="!rounded-lg">搜索</a-button>
            <a-button @click="resetSearch" class="!rounded-lg">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- Scrollable list area -->
    <div class="flex-1 overflow-y-auto">
      <div class="space-y-3" v-if="graphs.length">
        <div
          v-for="graph in graphs"
          :key="graph.id"
          class="group bg-white border border-gray-100 rounded-2xl px-5 py-4 flex items-center justify-between hover:shadow-sm hover:-translate-y-0.5 transition-all duration-300 cursor-pointer"
          @click="goWorkspace(graph)"
        >
          <div class="flex items-center gap-4 min-w-0">
            <div class="w-10 h-10 rounded-xl bg-indigo-50 flex items-center justify-center flex-shrink-0">
              <svg viewBox="0 0 24 24" class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" />
              </svg>
            </div>
            <div class="min-w-0">
              <div class="text-sm font-medium text-gray-900 truncate">{{ graph.graphName }}</div>
              <div class="text-xs text-gray-400 mt-0.5 flex items-center gap-2">
                <span>{{ graph.currentVersion || 'v001' }}</span>
                <span v-if="graph.hasTmp" class="text-amber-500">· 暂存</span>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-4 flex-shrink-0">
            <span class="text-xs text-gray-400 hidden sm:inline">{{ graph.createTime?.split('T')[0] }}</span>
            <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
              <span class="px-2.5 py-1 text-xs text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded-lg transition-colors cursor-pointer" @click.stop="showDetail(graph)">详情</span>
              <span class="px-2.5 py-1 text-xs text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded-lg transition-colors cursor-pointer" @click.stop="openEditModal(graph)">编辑</span>
              <span class="px-2.5 py-1 text-xs text-red-500 hover:bg-red-50 rounded-lg transition-colors cursor-pointer" @click.stop="removeGraph(graph)">删除</span>
            </div>
            <svg viewBox="0 0 24 24" class="w-4 h-4 text-gray-300 group-hover:text-gray-400 transition-colors" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M8.25 4.5l7.5 7.5-7.5 7.5" />
            </svg>
          </div>
        </div>
      </div>

      <a-empty v-if="!loading && !graphs.length" class="mt-16" description="暂时还没有项目，先新建一个试试" />

      <!-- Pagination -->
      <div v-if="total > 0" class="flex justify-center mt-8">
        <a-pagination
          :current="searchParams.pageNum"
          :page-size="searchParams.pageSize"
          :total="total"
          show-size-changer
          @change="handlePageChange"
        />
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:open="createModalOpen"
      :title="editingGraph ? '编辑故障树项目' : '新建故障树项目'"
      :ok-text="editingGraph ? '保存修改' : '创建项目'"
      cancel-text="取消"
      :confirm-loading="submitting"
      class="!rounded-xl"
      @ok="submitGraph"
    >
      <a-form layout="vertical" :model="createForm">
        <a-form-item label="项目名称" required>
          <a-input v-model:value="createForm.graphName" maxlength="50" placeholder="例如：压缩机故障树分析" class="!rounded-lg" />
        </a-form-item>
        <a-form-item label="项目描述">
          <a-textarea v-model:value="createForm.description" :rows="4" placeholder="简单描述设备、顶事件或分析目标" class="!rounded-lg" />
        </a-form-item>
        <a-form-item label="封面链接">
          <a-input v-model:value="createForm.cover" placeholder="可选，填展示封面 URL" class="!rounded-lg" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Detail Modal -->
    <a-modal v-model:open="detailOpen" title="项目概览" :footer="null" width="680px" class="!rounded-xl">
      <a-descriptions bordered :column="2" v-if="selectedGraph">
        <a-descriptions-item label="项目名称">{{ selectedGraph.graphName }}</a-descriptions-item>
        <a-descriptions-item label="创建者">{{ selectedGraph.userId }}</a-descriptions-item>
        <a-descriptions-item label="当前版本">{{ selectedGraph.currentVersion || 'v001' }}</a-descriptions-item>
        <a-descriptions-item label="暂存状态">
          <a-tag :color="selectedGraph.hasTmp ? 'orange' : 'blue'">
            {{ selectedGraph.hasTmp ? '存在暂存版本' : '仅正式版本' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="项目描述" :span="2">
          {{ selectedGraph.description || '暂无项目描述' }}
        </a-descriptions-item>
        <a-descriptions-item v-if="isAdminView" label="创建者 ID">
          {{ selectedGraph.userId }}
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { addGraph, deleteGraph, listGraphVoByPage, updateGraph } from '@/api/graphController'

const route = useRoute()
const router = useRouter()
const isAdminView = computed(() => route.path.startsWith('/admin'))

const loading = ref(false)
const submitting = ref(false)
const total = ref(0)
const graphs = ref<API.GraphVO[]>([])
const createModalOpen = ref(false)
const detailOpen = ref(false)
const selectedGraph = ref<API.GraphVO>()
const editingGraph = ref<API.GraphVO>()

const searchParams = reactive<API.GraphQueryRequest>({
  pageNum: 1,
  pageSize: 8,
  sortField: 'createTime',
  sortOrder: 'desc',
  graphName: '',
})

const createForm = reactive<API.GraphAddRequest>({
  graphName: '',
  description: '',
  cover: '',
})

const fetchGraphs = async () => {
  loading.value = true
  try {
    const res = await listGraphVoByPage({ ...searchParams })
    if (res.data.code === 0 && res.data.data) {
      graphs.value = res.data.data.records || []
      total.value = Number(res.data.data.totalRow || 0)
    } else {
      message.error(res.data.message || '获取项目列表失败')
    }
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number, pageSize: number) => {
  searchParams.pageNum = page
  searchParams.pageSize = pageSize
  fetchGraphs()
}

const resetSearch = () => {
  searchParams.pageNum = 1
  searchParams.pageSize = 8
  searchParams.graphName = ''
  searchParams.userId = undefined
  fetchGraphs()
}

const openCreateModal = () => {
  editingGraph.value = undefined
  createForm.graphName = ''
  createForm.description = ''
  createForm.cover = ''
  createModalOpen.value = true
}

const openEditModal = (graph: API.GraphVO) => {
  editingGraph.value = graph
  createForm.graphName = graph.graphName || ''
  createForm.description = graph.description || ''
  createForm.cover = graph.cover || ''
  createModalOpen.value = true
}

const submitGraph = async () => {
  if (!createForm.graphName?.trim()) {
    message.warning('请先输入项目名称')
    return
  }
  submitting.value = true
  try {
    if (editingGraph.value?.id) {
      const res = await updateGraph({
        id: editingGraph.value.id,
        graphName: createForm.graphName.trim(),
        description: createForm.description?.trim(),
        cover: createForm.cover?.trim(),
      })
      if (res.data.code === 0) {
        message.success('项目更新成功')
        createModalOpen.value = false
        editingGraph.value = undefined
        await fetchGraphs()
      } else {
        message.error(res.data.message || '项目更新失败')
      }
    } else {
      const res = await addGraph({
        graphName: createForm.graphName.trim(),
        description: createForm.description?.trim(),
        cover: createForm.cover?.trim(),
      })
      if (res.data.code === 0 && res.data.data) {
        message.success('项目创建成功')
        createModalOpen.value = false
        await fetchGraphs()
        await router.push(`/graph/workspace/${res.data.data}`)
      } else {
        message.error(res.data.message || '项目创建失败')
      }
    }
  } finally {
    submitting.value = false
  }
}

const goWorkspace = (graph: API.GraphVO) => {
  if (graph.id) {
    router.push(`/graph/workspace/${graph.id}`)
  }
}

const showDetail = (graph: API.GraphVO) => {
  selectedGraph.value = graph
  detailOpen.value = true
}

const removeGraph = (graph: API.GraphVO) => {
  if (!graph.id) return
  Modal.confirm({
    title: '删除项目',
    content: `确定删除项目「${graph.graphName || graph.id}」吗？项目目录也会一起删除。`,
    okText: '确认删除',
    okButtonProps: { danger: true },
    cancelText: '取消',
    async onOk() {
      const res = await deleteGraph({ id: graph.id! })
      if (res.data.code === 0) {
        message.success('项目删除成功')
        await fetchGraphs()
        return
      }
      message.error(res.data.message || '项目删除失败')
    },
  })
}

onMounted(fetchGraphs)
</script>

<style scoped>
/* Most styles replaced by Tailwind utility classes */
</style>
