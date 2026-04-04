<template>
  <div id="userLoginPage" class="min-h-[calc(100vh-120px)] flex flex-col items-center justify-center px-4 bg-[#f8fafc]">
    
    <div class="group w-full max-w-[440px]">
      <div class="bg-white border border-gray-100 rounded-2xl overflow-hidden transition-all duration-300 ease-[cubic-bezier(0.165,0.84,0.44,1)] group-hover:-translate-y-1.5 group-hover:shadow-xl p-8 md:p-10">
        
        <div class="mb-10 text-center">
          <h2 class="text-2xl font-bold text-gray-800 mb-2 tracking-tight">
            工业故障树平台 - 用户登录
          </h2>
          <div class="text-sm text-gray-400 font-medium">
            登录后进入你的故障树项目与工作台
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
              placeholder="请输入账号" 
              size="large"
              class="!rounded-xl !border-gray-100 focus:!border-blue-400"
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
              class="!rounded-xl !border-gray-100 focus:!border-blue-400"
            />
          </a-form-item>

          <div class="flex justify-between items-center mb-8 text-sm">
            <span class="text-gray-400 text-xs">安全加密登录中</span>
            <RouterLink to="/user/register" class="text-blue-500 hover:text-blue-600 font-medium transition-colors">
              没有账号？去注册
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
              登 录
            </a-button>
          </a-form-item>
        </a-form>
      </div>
      
      <p class="mt-8 text-center text-gray-400 text-xs tracking-widest uppercase">
        Fault Tree Workspace © 2026
      </p>
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
/* Ant Design 的样式覆盖（如果 Tailwind 没生效，可以通过 !important 或这种方式） */
:deep(.ant-input), :deep(.ant-input-password) {
  padding-top: 10px;
  padding-bottom: 10px;
}

:deep(.ant-form-item-label > label) {
  color: #6b7280; /* text-gray-500 */
  font-weight: 500;
}
</style>
