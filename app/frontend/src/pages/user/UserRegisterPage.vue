<template>
  <div class="min-h-[calc(100vh-56px)] bg-gray-50 flex items-center justify-center px-4">
    <div class="w-full max-w-[420px]">
      <div class="bg-white rounded-2xl border border-gray-200 shadow-sm p-8 md:p-10">
        <div class="text-center mb-8">
          <h2 class="text-xl font-bold text-gray-900 mb-1">用户注册</h2>
          <p class="text-sm text-gray-400">创建账号后开始管理你的故障树项目</p>
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
              placeholder="请设置登录账号"
              size="large"
              class="!rounded-xl !border-gray-200"
            />
          </a-form-item>

          <a-form-item
            name="userPassword"
            label="密码"
            :rules="[
              { required: true, message: '请输入密码' },
              { min: 8, message: '密码不能小于 8 位' },
            ]"
          >
            <a-input-password
              v-model:value="formState.userPassword"
              placeholder="请输入 8 位以上密码"
              size="large"
              class="!rounded-xl !border-gray-200"
            />
          </a-form-item>

          <a-form-item
            name="checkPassword"
            label="确认密码"
            :rules="[
              { required: true, message: '请确认密码' },
              { min: 8, message: '密码不能小于 8 位' },
              { validator: validateCheckPassword },
            ]"
          >
            <a-input-password
              v-model:value="formState.checkPassword"
              placeholder="请再次输入密码"
              size="large"
              class="!rounded-xl !border-gray-200"
            />
          </a-form-item>

          <div class="flex justify-end mb-6">
            <RouterLink to="/user/login" class="text-sm text-indigo-600 hover:text-indigo-500">
              已有账号？去登录
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
              注 册
            </a-button>
          </a-form-item>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { userRegister } from '@/api/userController.ts'
import { message } from 'ant-design-vue'
import { reactive } from 'vue'

const router = useRouter()

const formState = reactive<API.UserRegisterRequest>({
  userAccount: '',
  userPassword: '',
  checkPassword: '',
})

/**
 * 验证确认密码
 */
const validateCheckPassword = (rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (value && value !== formState.userPassword) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

/**
 * 提交表单
 */
const handleSubmit = async (values: API.UserRegisterRequest) => {
  const res = await userRegister(values)
  if (res.data.code === 0) {
    message.success('注册成功')
    // 按照您的要求，直接跳转首页
    router.push({
      path: '/',
      replace: true,
    })
  } else {
    message.error('注册失败，' + res.data.message)
  }
}
</script>

<style scoped>
:deep(.ant-form-item-label > label) {
  color: #6b7280;
  font-weight: 500;
}
</style>
