<template>
  <div class="home-page">
    <section class="hero">
      <div class="hero-copy">
        <p class="eyebrow">Industrial Fault Tree AI</p>
        <h1>用项目、版本和建议协同管理工业故障树</h1>
        <p class="description">
          围绕项目文档解析、图版本管理、AI 辅助修改和校验建议构建的一体化工作台。用户管理自己的项目，管理员查看全部项目状态。
        </p>
        <a-space wrap size="large">
          <a-button type="primary" size="large" @click="goMyProjects">进入我的项目</a-button>
          <a-button size="large" @click="goWorkspaceDemo">查看工作台设计</a-button>
        </a-space>
      </div>
      <div class="hero-panel">
        <div class="metric-card">
          <span>项目管理</span>
          <strong>Graph + Version + Tmp</strong>
        </div>
        <div class="metric-card">
          <span>文档入口</span>
          <strong>用户文档 / 项目文档</strong>
        </div>
        <div class="metric-card">
          <span>AI 模式</span>
          <strong>修改模式 / 建议模式</strong>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <div>
          <p class="eyebrow">Main Flow</p>
          <h2>当前前端信息架构</h2>
        </div>
      </div>
      <div class="feature-grid">
        <div class="feature-card">
          <h3>我的项目页</h3>
          <p>普通用户管理自己的故障树项目，进入工作台、查看暂存状态、继续版本编辑。</p>
        </div>
        <div class="feature-card">
          <h3>项目总览页</h3>
          <p>管理员查看全部项目，按项目名和创建人筛选，并保留用户管理入口。</p>
        </div>
        <div class="feature-card">
          <h3>图工作台</h3>
          <p>集中承载版本、建议、聊天、文档上传和后续 Vue Flow 画布编辑能力。</p>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <div>
          <p class="eyebrow">Roles</p>
          <h2>角色能力预览</h2>
        </div>
      </div>
      <div class="role-grid">
        <a-card class="role-card" :bordered="false">
          <template #title>工程师</template>
          <ul>
            <li>管理自己的图项目</li>
            <li>进入工作台查看当前展示版本</li>
            <li>上传项目文档并结合 AI 建议调整方案</li>
          </ul>
        </a-card>
        <a-card class="role-card" :bordered="false">
          <template #title>管理员</template>
          <ul>
            <li>查看全部图项目和用户信息</li>
            <li>从平台侧检查项目进展与暂存状态</li>
            <li>继续复用现有用户管理模块</li>
          </ul>
        </a-card>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goMyProjects = () => {
  router.push('/graph/manage')
}

const goWorkspaceDemo = () => {
  router.push('/graph/manage')
}
</script>

<style scoped>
.home-page {
  min-height: calc(100vh - 120px);
  padding: 40px 32px 56px;
  background:
    radial-gradient(circle at top left, rgba(37, 99, 235, 0.15), transparent 20%),
    radial-gradient(circle at right 20%, rgba(8, 145, 178, 0.12), transparent 18%),
    linear-gradient(180deg, #f8fbff 0%, #f8fafc 100%);
}

.hero {
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) minmax(320px, 0.9fr);
  gap: 28px;
  align-items: stretch;
}

.hero-copy,
.hero-panel,
.feature-card,
.role-card {
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid #e2e8f0;
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.08);
}

.hero-copy {
  padding: 40px;
}

.eyebrow {
  margin: 0 0 12px;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.14em;
  color: #2563eb;
  text-transform: uppercase;
}

.hero-copy h1,
.section-head h2 {
  margin: 0 0 14px;
  color: #0f172a;
}

.hero-copy h1 {
  font-size: 40px;
  line-height: 1.15;
}

.description {
  margin: 0 0 28px;
  color: #475569;
  font-size: 16px;
  line-height: 1.8;
}

.hero-panel {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: linear-gradient(180deg, #0f172a 0%, #1e293b 100%);
}

.metric-card {
  padding: 18px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.08);
  color: #e2e8f0;
}

.metric-card span {
  display: block;
  margin-bottom: 8px;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #93c5fd;
}

.metric-card strong {
  font-size: 18px;
}

.section {
  margin-top: 32px;
}

.section-head {
  margin-bottom: 18px;
}

.feature-grid,
.role-grid {
  display: grid;
  gap: 18px;
}

.feature-grid {
  grid-template-columns: repeat(3, 1fr);
}

.role-grid {
  grid-template-columns: repeat(2, 1fr);
}

.feature-card {
  padding: 24px;
}

.feature-card h3 {
  margin-top: 0;
  color: #0f172a;
}

.feature-card p,
.role-card ul {
  color: #475569;
  line-height: 1.7;
}

.role-card ul {
  margin: 0;
  padding-left: 18px;
}
</style>
