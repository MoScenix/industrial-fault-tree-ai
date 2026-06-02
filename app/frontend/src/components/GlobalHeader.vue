<template>
  <a-layout-header class="header">
    <a-row :wrap="false">
      <!-- 左侧：Logo和标题 -->
      <a-col flex="200px">
        <RouterLink to="/">
          <div class="header-left">
            <img class="logo" src="@/assets/logo.png" alt="Logo" />
            <h1 class="site-title">故障树管理平台</h1>
          </div>
        </RouterLink>
      </a-col>
      <!-- 中间：导航菜单 -->
      <a-col flex="auto">
        <a-menu v-model:selectedKeys="selectedKeys" mode="horizontal" :items="menuItems" @click="handleMenuClick" />
      </a-col>
      <!-- 右侧：用户操作区域 -->
      <a-col>
        <div class="user-login-status">
          <div v-if="loginUserStore.loginUser.id">
            <a-space>
              <a-upload :show-upload-list="false" :before-upload="handleUserDocumentUpload">
                <a-button>上传个人文档</a-button>
              </a-upload>
              <a-dropdown>
                <a-space>
                  <a-avatar :src="loginUserStore.loginUser.userAvatar" />
                  {{ loginUserStore.loginUser.userName ?? '无名' }}
                </a-space>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="router.push('/user/center')">
                      <UserOutlined />
                      个人中心
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item @click="doLogout" danger>
                      <LogoutOutlined />
                      退出登录
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </a-space>
          </div>
          <div v-else>
            <a-button type="primary" href="/user/login">登录</a-button>
          </div>
        </div>
      </a-col>
    </a-row>
  </a-layout-header>
</template>

<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { useRouter } from 'vue-router'
import { type MenuProps, message } from 'ant-design-vue'
import { useLoginUserStore } from '@/stores/loginUser.ts'
import { userLogout } from '@/api/userController.ts'
import { uploadUserDocument } from '@/api/documentController'
import {
  LogoutOutlined,
  HomeOutlined,
  ApartmentOutlined,
  TeamOutlined,
  SettingOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'

const loginUserStore = useLoginUserStore()
const router = useRouter()
// 当前选中菜单
const selectedKeys = ref<string[]>(['/'])
// 监听路由变化，更新当前选中菜单
router.afterEach((to) => {
  if (to.path.startsWith('/graph/workspace')) {
    selectedKeys.value = ['/graph/manage']
    return
  }
  selectedKeys.value = [to.path]
})

// 菜单配置项
const originItems = [
  {
    key: '/',
    icon: () => h(HomeOutlined),
    label: '主页',
    title: '主页',
  },
  {
    key: '/graph/manage',
    icon: () => h(ApartmentOutlined),
    label: '我的项目',
    title: '我的项目',
  },
  {
    key: '/admin/graphManage',
    icon: () => h(ApartmentOutlined),
    label: '项目总览',
    title: '项目总览',
  },
  {
    key: '/admin/userManage',
    icon: () => h(TeamOutlined),
    label: '用户管理',
    title: '用户管理',
  },
  {
    key: '/admin/promptManage',
    icon: () => h(SettingOutlined),
    label: '提示词管理',
    title: '提示词管理',
  },
]

// 过滤菜单项
const filterMenus = (menus = [] as MenuProps['items']) => {
  return menus?.filter((menu) => {
    const menuKey = menu?.key as string
    if (menuKey?.startsWith('/admin')) {
      const loginUser = loginUserStore.loginUser
      if (!loginUser || loginUser.userRole !== 'admin') {
        return false
      }
    }
    return true
  })
}

// 展示在菜单的路由数组
const menuItems = computed<MenuProps['items']>(() => filterMenus(originItems))

// 处理菜单点击
const handleMenuClick: MenuProps['onClick'] = (e) => {
  const key = e.key as string
  selectedKeys.value = [key]
  // 跳转到对应页面
  if (key.startsWith('/')) {
    router.push(key)
  }
}

// 退出登录
const doLogout = async () => {
  const res = await userLogout()
  if (res.data.code === 0) {
    loginUserStore.setLoginUser({
      userName: '未登录',
    })
    message.success('退出登录成功')
    await router.push('/user/login')
  } else {
    message.error('退出登录失败，' + res.data.message)
  }
}

const handleUserDocumentUpload = async (file: File) => {
  const res = await uploadUserDocument(file)
  if (res.data.code === 0) {
    message.success('个人文档上传成功')
  } else {
    message.error(res.data.message || '个人文档上传失败')
  }
  return false
}
</script>

<style scoped>
.header {
  background: rgba(255, 255, 255, 0.95);
  padding: 0 24px;
  backdrop-filter: blur(10px);
  border-bottom: 1px solid #e5e7eb;
  height: 56px;
  line-height: 56px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo {
  height: 36px;
  width: 36px;
}

.site-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #4f46e5;
}

:deep(.ant-menu-horizontal) {
  border-bottom: none !important;
}

:deep(.ant-menu-item) {
  padding: 0 16px !important;
  font-size: 14px;
}

.user-login-status {
  display: flex;
  align-items: center;
  height: 100%;
}
</style>
