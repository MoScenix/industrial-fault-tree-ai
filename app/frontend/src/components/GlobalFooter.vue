<template>
  <a-layout-footer class="footer">
    <div class="footer-content">
      <p class="copyright">
        <a
          href="https://github.com/MoScenix/industrial-fault-tree-ai"
          target="_blank"
          rel="noopener noreferrer"
          class="author-link"
        >
          industrial-fault-tree-ai · 故障树智能管理平台
        </a>
      </p>
    </div>
  </a-layout-footer>
</template>

<script setup lang="ts">
// 无需额外的响应式数据
</script>

<style scoped>
.footer {
  text-align: center;
  padding: 16px 20px;
  border-top: 1px solid #f0f0f0;
  background: #fff;
}
.copyright {
  margin: 0;
  color: #999;
  font-size: 13px;
}
</style>
