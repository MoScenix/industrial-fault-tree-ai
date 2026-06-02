<template>
  <div class="min-h-[calc(100vh-56px)] bg-gray-50">
    <!-- Hero -->
    <div class="max-w-6xl mx-auto px-6 pt-16 pb-12 md:pt-20 md:pb-16">
      <div class="flex flex-col lg:flex-row items-start gap-10">
        <div class="flex-1">
          <p class="text-xs font-bold tracking-widest text-indigo-600 uppercase mb-4">Industrial Fault Tree AI</p>
          <h1 class="text-3xl md:text-4xl lg:text-5xl font-bold text-gray-900 leading-tight mb-4">
            用项目、版本和建议<br/>协同管理工业故障树
          </h1>
          <p class="text-base text-gray-500 leading-relaxed mb-8 max-w-lg">
            围绕项目文档解析、图版本管理、AI 辅助修改和校验建议构建的一体化工作台。
          </p>
          <div class="flex gap-3">
            <a-button type="primary" size="large" class="!rounded-xl !h-11 !px-6" @click="goMyProjects">进入我的项目</a-button>
            <a-button size="large" class="!rounded-xl !h-11 !px-6" @click="goWorkspaceDemo">了解更多</a-button>
          </div>
        </div>
        <div class="w-full lg:w-80 flex-shrink-0">
          <div class="bg-gray-900 rounded-2xl p-5 space-y-3">
            <div class="rounded-xl bg-white/10 p-4">
              <p class="text-[11px] text-indigo-300 uppercase tracking-widest font-semibold mb-2">项目管理</p>
              <p class="text-white text-sm font-semibold">Graph + Version + Tmp</p>
            </div>
            <div class="rounded-xl bg-white/10 p-4">
              <p class="text-[11px] text-indigo-300 uppercase tracking-widest font-semibold mb-2">文档入口</p>
              <p class="text-white text-sm font-semibold">用户文档 / 项目文档</p>
            </div>
            <div class="rounded-xl bg-white/10 p-4">
              <p class="text-[11px] text-indigo-300 uppercase tracking-widest font-semibold mb-2">AI 模式</p>
              <p class="text-white text-sm font-semibold">修改模式 / 建议模式</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Features -->
    <div class="max-w-6xl mx-auto px-6 pb-16">
      <div class="mb-8">
        <p class="text-xs font-bold tracking-widest text-indigo-600 uppercase mb-2">Main Flow</p>
        <h2 class="text-2xl font-bold text-gray-900">当前前端信息架构</h2>
      </div>
      <div class="grid md:grid-cols-3 gap-5">
        <div class="bg-white rounded-xl border border-gray-200 p-6 shadow-sm">
          <h3 class="text-base font-semibold text-gray-900 mb-2">我的项目页</h3>
          <p class="text-sm text-gray-500 leading-relaxed">管理故障树项目，进入工作台、查看暂存状态、继续版本编辑。</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-6 shadow-sm">
          <h3 class="text-base font-semibold text-gray-900 mb-2">项目总览页</h3>
          <p class="text-sm text-gray-500 leading-relaxed">管理员查看全部项目，按项目名和创建人筛选。</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-6 shadow-sm">
          <h3 class="text-base font-semibold text-gray-900 mb-2">图工作台</h3>
          <p class="text-sm text-gray-500 leading-relaxed">集中承载版本、建议、聊天、文档上传和画布编辑能力。</p>
        </div>
      </div>
    </div>

    <!-- Roles -->
    <div class="max-w-6xl mx-auto px-6 pb-16">
      <div class="mb-8">
        <p class="text-xs font-bold tracking-widest text-indigo-600 uppercase mb-2">Roles</p>
        <h2 class="text-2xl font-bold text-gray-900">角色能力预览</h2>
      </div>
      <div class="grid md:grid-cols-2 gap-5">
        <div class="bg-white rounded-xl border border-gray-200 p-6 shadow-sm">
          <h3 class="text-base font-semibold text-gray-900 mb-3">工程师</h3>
          <ul class="text-sm text-gray-500 space-y-2 pl-5 list-disc">
            <li>管理自己的图项目</li>
            <li>进入工作台查看当前展示版本</li>
            <li>上传项目文档并结合 AI 建议调整方案</li>
          </ul>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-6 shadow-sm">
          <h3 class="text-base font-semibold text-gray-900 mb-3">管理员</h3>
          <ul class="text-sm text-gray-500 space-y-2 pl-5 list-disc">
            <li>查看全部图项目和用户信息</li>
            <li>从平台侧检查项目进展与暂存状态</li>
            <li>继续复用现有用户管理模块</li>
          </ul>
        </div>
      </div>
    </div>
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
