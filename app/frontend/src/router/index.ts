import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import UserLoginPage from '@/pages/user/UserLoginPage.vue'
import UserRegisterPage from '@/pages/user/UserRegisterPage.vue'
import UserManagePage from '@/pages/admin/UserManagePage.vue'
import UserCenterPage from '@/pages/user/UserCenterPage.vue'
import GraphManagePage from '@/pages/graph/GraphManagePage.vue'
import AdminGraphManagePage from '@/pages/admin/GraphManagePage.vue'
import GraphWorkspacePage from '@/pages/graph/GraphWorkspacePage.vue'
import PromptManagePage from '@/pages/admin/PromptManagePage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: '主页',
      component: HomePage,
    },
    {
      path: '/user/login',
      name: '用户登录',
      component: UserLoginPage,
    },
    {
      path: '/user/register',
      name: '用户注册',
      component: UserRegisterPage,
    },
    {
      path: '/admin/userManage',
      name: '用户管理',
      component: UserManagePage,
    },
    {
      path: '/graph/manage',
      name: '我的项目',
      component: GraphManagePage,
    },
    {
      path: '/graph/workspace/:id',
      name: '图工作台',
      component: GraphWorkspacePage,
    },
    {
      path: '/admin/graphManage',
      name: '项目总览',
      component: AdminGraphManagePage,
    },
    {
      path: '/admin/promptManage',
      name: '提示词管理',
      component: PromptManagePage,
    },
    {
      path: '/user/center',
      name: '个人中心',
      component: UserCenterPage,
    },
  ],
})

export default router
