<template>
  <div class="graph-card">
    <div class="cover">
      <img v-if="graph.cover" :src="graph.cover" :alt="graph.graphName" />
      <div v-else class="cover-placeholder">
        <span>FT</span>
      </div>
      <div class="cover-mask">
        <a-space>
          <a-button type="primary" @click.stop="$emit('open', graph)">进入工作台</a-button>
          <a-button @click.stop="$emit('detail', graph)">查看详情</a-button>
        </a-space>
      </div>
    </div>
    <div class="content">
      <div class="title-row">
        <h3 class="title">{{ graph.graphName || '未命名项目' }}</h3>
        <a-tag v-if="graph.hasTmp" color="orange">有暂存</a-tag>
      </div>
      <p class="desc">{{ graph.description || '暂无描述，适合作为设备故障树分析项目的起点。' }}</p>
      <div class="meta">
        <span>当前版本：{{ graph.currentVersion || 'v001' }}</span>
        <span>{{ graph.updateTime || graph.createTime || '刚刚创建' }}</span>
      </div>
      <div class="action-row">
        <a-button size="small" type="link" @click.stop="$emit('edit', graph)">编辑信息</a-button>
        <a-button size="small" type="link" danger @click.stop="$emit('delete', graph)">删除项目</a-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  graph: API.GraphVO
}>()

defineEmits<{
  (e: 'open', graph: API.GraphVO): void
  (e: 'detail', graph: API.GraphVO): void
  (e: 'edit', graph: API.GraphVO): void
  (e: 'delete', graph: API.GraphVO): void
}>()
</script>

<style scoped>
.graph-card {
  overflow: hidden;
  border: 1px solid #e5eef5;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.95);
  box-shadow: 0 18px 48px rgba(15, 23, 42, 0.08);
  transition:
    transform 0.25s ease,
    box-shadow 0.25s ease;
}

.graph-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 22px 56px rgba(15, 23, 42, 0.12);
}

.cover {
  position: relative;
  height: 180px;
  overflow: hidden;
  background:
    radial-gradient(circle at top left, rgba(31, 111, 235, 0.22), transparent 42%),
    linear-gradient(135deg, #eff6ff 0%, #f8fafc 45%, #ecfeff 100%);
}

.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-placeholder {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 42px;
  font-weight: 800;
  color: #1d4ed8;
  letter-spacing: 0.08em;
}

.cover-mask {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.44);
  opacity: 0;
  transition: opacity 0.25s ease;
}

.graph-card:hover .cover-mask {
  opacity: 1;
}

.content {
  padding: 18px;
}

.title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.title {
  margin: 0;
  font-size: 17px;
  color: #0f172a;
}

.desc {
  min-height: 44px;
  margin: 10px 0 14px;
  color: #64748b;
  line-height: 1.6;
}

.meta {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 12px;
  color: #94a3b8;
}

.action-row {
  display: flex;
  justify-content: space-between;
  margin-top: 12px;
}
</style>
