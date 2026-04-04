<template>
  <div class="prompt-page">
    <section class="hero">
      <div>
        <p class="eyebrow">Admin Prompt Console</p>
        <h1>提示词管理</h1>
        <p class="subtitle">集中维护修改模式与建议模式的系统提示词，当前页面已经直接通过 BFF 读写后端配置。</p>
      </div>
      <a-tag color="gold">管理员专用</a-tag>
    </section>

    <a-card :bordered="false" class="prompt-card">
      <a-tabs v-model:activeKey="activeKey">
        <a-tab-pane key="modify" tab="修改模式">
          <a-textarea
            v-model:value="modifyPrompt"
            :rows="18"
            placeholder="用于 AI 修改图时的系统提示词"
          />
        </a-tab-pane>
        <a-tab-pane key="log" tab="建议模式">
          <a-textarea
            v-model:value="logPrompt"
            :rows="18"
            placeholder="用于校验与生成建议时的系统提示词"
          />
        </a-tab-pane>
      </a-tabs>
      <div class="actions">
        <a-space>
          <a-button @click="resetDraft">重置草稿</a-button>
          <a-button type="primary" :loading="loading" @click="saveDraft">保存提示词</a-button>
        </a-space>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import { getPrompt, updatePrompt } from '@/api/aiController'

const activeKey = ref('modify')
const modifyPrompt = ref('')
const logPrompt = ref('')
const loading = ref(false)

const MODIFY_MODE = 1
const LOG_MODE = 2

const loadPrompt = async (mode: number) => {
  const res = await getPrompt({ mode })
  if (res.data.code === 0 && res.data.data) {
    if (mode === MODIFY_MODE) {
      modifyPrompt.value = res.data.data.content || ''
    } else {
      logPrompt.value = res.data.data.content || ''
    }
  } else {
    message.error(res.data.message || '获取提示词失败')
  }
}

const resetDraft = () => {
  loadPrompt(MODIFY_MODE)
  loadPrompt(LOG_MODE)
  message.success('已恢复为服务端当前内容')
}

const saveDraft = async () => {
  loading.value = true
  try {
    const mode = activeKey.value === 'modify' ? MODIFY_MODE : LOG_MODE
    const content = activeKey.value === 'modify' ? modifyPrompt.value : logPrompt.value
    const res = await updatePrompt({ mode, content })
    if (res.data.code === 0) {
      message.success('提示词保存成功')
      await loadPrompt(mode)
    } else {
      message.error(res.data.message || '提示词保存失败')
    }
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadPrompt(MODIFY_MODE), loadPrompt(LOG_MODE)])
})
</script>

<style scoped>
.prompt-page {
  min-height: calc(100vh - 120px);
  padding: 32px;
  background:
    radial-gradient(circle at top left, rgba(37, 99, 235, 0.14), transparent 24%),
    linear-gradient(180deg, #f8fbff 0%, #f8fafc 100%);
}

.hero {
  display: flex;
  justify-content: space-between;
  align-items: end;
  gap: 24px;
  margin-bottom: 24px;
}

.eyebrow {
  margin: 0 0 10px;
  color: #2563eb;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.hero h1 {
  margin: 0 0 10px;
  color: #0f172a;
  font-size: 34px;
}

.subtitle {
  margin: 0;
  color: #64748b;
}

.prompt-card {
  border-radius: 24px;
  box-shadow: 0 20px 52px rgba(15, 23, 42, 0.06);
}

.actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 18px;
}
</style>
