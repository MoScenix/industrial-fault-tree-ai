<template>
  <div class="min-h-[calc(100vh-56px)] bg-gray-50 flex items-center justify-center px-4">
    <div class="w-full max-w-[420px]">
      <div class="bg-white rounded-2xl border border-gray-200 shadow-sm p-8 md:p-10">
        <div class="text-center mb-8">
          <h2 class="text-xl font-bold text-gray-900 mb-1">用户登录</h2>
          <p class="text-sm text-gray-400">工业故障树智能管理平台</p>
        </div>

        <a-form
          :model="formState"
          name="basic"
          layout="vertical"
          autocomplete="off"
          @finish="handleSubmit"
        >
          <a-form-item
            name="userAccount"
            label="账号"
            :rules="[{ required: true, message: '请输入账号' }]"
          >
            <a-input
              v-model:value="formState.userAccount"
              placeholder="请输入账号"
              size="large"
              class="!rounded-xl !border-gray-200"
            />
          </a-form-item>

          <a-form-item
            name="userPassword"
            label="密码"
            :rules="[
              { required: true, message: '请输入密码' },
              { min: 8, message: '密码长度不能小于 8 位' },
            ]"
          >
            <a-input-password
              v-model:value="formState.userPassword"
              placeholder="请输入密码"
              size="large"
              class="!rounded-xl !border-gray-200"
            />
          </a-form-item>

          <div class="flex justify-end mb-6">
            <RouterLink to="/user/register" class="text-sm text-indigo-600 hover:text-indigo-500">
              没有账号？去注册
            </RouterLink>
          </div>

          <a-form-item class="mb-0">
            <a-button
              type="primary"
              html-type="submit"
              size="large"
              block
              class="!h-11 !rounded-xl !bg-indigo-600 hover:!bg-indigo-500 !border-none font-semibold"
            >
              登 录
            </a-button>
          </a-form-item>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { reactive } from 'vue'
import { userLogin } from '@/api/userController.ts'
import { useLoginUserStore } from '@/stores/loginUser.ts'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'

const formState = reactive<API.UserLoginRequest>({
  userAccount: '',
  userPassword: '',
})

const router = useRouter()
const loginUserStore = useLoginUserStore()

const handleSubmit = async (values: any) => {
  const res = await userLogin(values)
  if (res.data.code === 0 && res.data.data) {
    await loginUserStore.fetchLoginUser()
    message.success('登录成功')
    router.push({
      path: '/',
      replace: true,
    })
  } else {
    message.error('登录失败，' + res.data.message)
  }
}
</script>

<style scoped>
:deep(.ant-form-item-label > label) {
  color: #6b7280;
  font-weight: 500;
}
</style>
