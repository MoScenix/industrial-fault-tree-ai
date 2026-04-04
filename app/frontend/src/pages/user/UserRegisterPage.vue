<template>
  <div id="userRegisterPage" class="min-h-[calc(100vh-120px)] flex flex-col items-center justify-center px-4 bg-[#f8fafc]">
    
    <div class="group w-full max-w-[440px]">
      <div class="bg-white border border-gray-100 rounded-2xl overflow-hidden transition-all duration-300 ease-[cubic-bezier(0.165,0.84,0.44,1)] group-hover:-translate-y-1.5 group-hover:shadow-xl p-8 md:p-10">
        
        <div class="mb-10 text-center">
          <h2 class="text-2xl font-bold text-gray-800 mb-2 tracking-tight">
            工业故障树平台 - 用户注册
          </h2>
          <div class="text-sm text-gray-400 font-medium">
            创建账号后开始管理你的故障树项目
          </div>
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
              class="!rounded-xl !border-gray-100 focus:!border-blue-400"
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
              class="!rounded-xl !border-gray-100 focus:!border-blue-400"
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
              class="!rounded-xl !border-gray-100 focus:!border-blue-400"
            />
          </a-form-item>

          <div class="flex justify-end mb-8 text-sm">
            <RouterLink to="/user/login" class="text-blue-500 hover:text-blue-600 font-medium transition-colors">
              已有账号？去登录
            </RouterLink>
          </div>

          <a-form-item class="mb-0">
            <a-button 
              type="primary" 
              html-type="submit" 
              size="large"
              block
              class="!h-12 !rounded-xl !bg-blue-600 hover:!bg-blue-700 !border-none font-bold text-white shadow-lg shadow-blue-100 transition-all duration-300 active:scale-95"
            >
              注 册
            </a-button>
          </a-form-item>
        </a-form>
      </div>
      
      <p class="mt-8 text-center text-gray-400 text-xs tracking-widest uppercase">
        Fault Tree Intelligence Workspace
      </p>
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
/* 移除原有 CSS，使用 Tailwind 控制 */
:deep(.ant-input), :deep(.ant-input-password) {
  padding-top: 10px;
  padding-bottom: 10px;
}

:deep(.ant-form-item-label > label) {
  color: #6b7280;
  font-weight: 500;
}
</style>
