<template>
  <div class="graph-manage-page">
    <section class="hero">
      <div>
        <p class="eyebrow">{{ isAdminView ? '管理员全局视图' : '我的故障树项目' }}</p>
        <h1>{{ isAdminView ? '统一查看全部故障树项目' : '管理自己的故障树项目与编辑进度' }}</h1>
        <p class="subtitle">
          {{ isAdminView ? '支持按创建人、项目名和状态筛选全部项目。' : '从这里进入图工作台、继续暂存编辑或新建分析项目。' }}
        </p>
      </div>
      <a-button type="primary" size="large" @click="openCreateModal">新建项目</a-button>
    </section>

    <a-card :bordered="false" class="toolbar-card">
      <a-form layout="inline" :model="searchParams" @finish="fetchGraphs">
        <a-form-item label="项目名">
          <a-input v-model:value="searchParams.graphName" placeholder="搜索项目名称" allow-clear />
        </a-form-item>
        <a-form-item v-if="isAdminView" label="用户 ID">
          <a-input-number v-model:value="searchParams.userId" placeholder="创建者 ID" style="width: 160px" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit">搜索</a-button>
            <a-button @click="resetSearch">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <div class="graph-grid">
      <GraphCard
        v-for="graph in graphs"
        :key="graph.id"
        :graph="graph"
        @open="goWorkspace"
        @detail="showDetail"
        @edit="openEditModal"
        @delete="removeGraph"
      />
    </div>

    <a-empty v-if="!loading && !graphs.length" description="暂时还没有项目，先新建一个试试" />

    <div class="pager">
      <a-pagination
        :current="searchParams.pageNum"
        :page-size="searchParams.pageSize"
        :total="total"
        show-size-changer
        @change="handlePageChange"
      />
    </div>

    <a-modal
      v-model:open="createModalOpen"
      :title="editingGraph ? '编辑故障树项目' : '新建故障树项目'"
      :ok-text="editingGraph ? '保存修改' : '创建项目'"
      cancel-text="取消"
      :confirm-loading="submitting"
      @ok="submitGraph"
    >
      <a-form layout="vertical" :model="createForm">
        <a-form-item label="项目名称" required>
          <a-input v-model:value="createForm.graphName" maxlength="50" placeholder="例如：压缩机故障树分析" />
        </a-form-item>
        <a-form-item label="项目描述">
          <a-textarea v-model:value="createForm.description" :rows="4" placeholder="简单描述设备、顶事件或分析目标" />
        </a-form-item>
        <a-form-item label="封面链接">
          <a-input v-model:value="createForm.cover" placeholder="可选，填展示封面 URL" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="detailOpen" title="项目概览" :footer="null" width="720px">
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
import GraphCard from '@/components/GraphCard.vue'
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
.graph-manage-page {
  min-height: calc(100vh - 120px);
  padding: 32px 32px 56px;
  background:
    radial-gradient(circle at top left, rgba(29, 78, 216, 0.18), transparent 25%),
    linear-gradient(180deg, #f8fbff 0%, #f8fafc 100%);
}

.hero {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;
}

.eyebrow {
  margin: 0 0 10px;
  color: #2563eb;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.hero h1 {
  margin: 0 0 10px;
  font-size: 34px;
  color: #0f172a;
}

.subtitle {
  margin: 0;
  color: #64748b;
  font-size: 15px;
}

.toolbar-card {
  margin-bottom: 24px;
  border-radius: 20px;
  box-shadow: 0 12px 36px rgba(15, 23, 42, 0.06);
}

.graph-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.pager {
  display: flex;
  justify-content: center;
  margin-top: 28px;
}
</style>
