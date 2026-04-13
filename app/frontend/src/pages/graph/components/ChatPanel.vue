<template>
  <div class="flex flex-col h-full overflow-hidden">
    <div ref="chatListRef" class="chat-list flex flex-col gap-4 overflow-y-auto pb-4 px-2">
      <div v-for="item in messages" :key="item.id" class="flex flex-col group">
        <div 
          class="flex max-w-[85%]"
          :class="item.role === 'assistant' ? 'self-start' : 'self-end'"
        >
          <div class="flex items-start gap-3" :class="item.role === 'assistant' ? 'flex-row' : 'flex-row-reverse'">
            <div class="w-8 h-8 rounded-full flex-shrink-0 flex items-center justify-center text-xs font-bold text-white shadow-sm"
                 :class="item.role === 'assistant' ? 'bg-gradient-to-br from-indigo-500 to-purple-600' : 'bg-gradient-to-br from-emerald-500 to-teal-500'">
              {{ item.role === 'assistant' ? 'AI' : '我' }}
            </div>
            <div
              class="px-4 py-3 shadow-sm text-[14px] leading-relaxed relative"
              :class="[
                item.role === 'assistant' 
                  ? 'bg-white text-slate-700 border border-slate-100 rounded-2xl rounded-tl-sm' 
                  : 'bg-indigo-50 text-indigo-900 border border-indigo-100 rounded-2xl rounded-tr-sm'
              ]"
            >
              <MarkdownRenderer :content="item.content || ''" />
            </div>
          </div>
        </div>
      </div>
      <a-empty v-if="!messages.length" description="输入需求，AI 将辅助您构建故障树" class="mt-10" />
    </div>
    
    <div class="mt-auto pt-4 border-t border-slate-100 flex-shrink-0 bg-white">
      <div class="relative bg-white rounded-xl border border-slate-200 shadow-sm focus-within:border-indigo-500 focus-within:ring-1 focus-within:ring-indigo-500 transition-all overflow-hidden">
        <a-textarea
          v-model:value="chatInput"
          :rows="3"
          :maxlength="1000"
          placeholder="例如：帮我添加一个温度过高的基础事件..."
          class="w-full !border-none !shadow-none !bg-transparent resize-none py-3 px-4 focus:!shadow-none text-sm"
          @pressEnter="onPressEnter"
        />
        <div class="flex justify-between items-center px-3 py-2 bg-slate-50 border-t border-slate-100">
          <span class="text-xs text-slate-400">Shift + Enter 换行，Enter 发送</span>
          <a-button type="primary" shape="circle" :loading="chatting" @click="handleSend" 
                    class="!flex !items-center !justify-center !w-8 !h-8 !min-w-0 !bg-indigo-600 hover:!bg-indigo-500 border-none shadow-md">
            <svg v-if="!chatting" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-4 h-4">
              <path d="M3.478 2.404a.75.75 0 0 0-.926.941l2.432 7.905H13.5a.75.75 0 0 1 0 1.5H4.984l-2.432 7.905a.75.75 0 0 0 .926.94 60.519 60.519 0 0 0 18.445-8.986.75.75 0 0 0 0-1.218A60.517 60.517 0 0 0 3.478 2.404Z" />
            </svg>
          </a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'
import { chatToModifyGraphSSE, listGraphMessage } from '@/api/graphController'

interface Props {
  graphId: number
  selectedVersion: string
  currentVersion?: string
  isEditorDirty: boolean
  onSuccess?: () => void
}

const props = defineProps<Props>()

const messages = ref<API.GraphMessageVO[]>([])
const chatInput = ref('')
const chatting = ref(false)
const chatListRef = ref<HTMLElement>()

const scrollToBottom = async () => {
  await nextTick()
  if (chatListRef.value) {
    chatListRef.value.scrollTop = chatListRef.value.scrollHeight
  }
}

const loadMessages = async () => {
  const res = await listGraphMessage({ graphId: props.graphId, pageSize: 10 })
  if (res.data.code === 0 && res.data.data) {
    messages.value = (res.data.data.records || []).reverse()
    await scrollToBottom()
  }
}

const onPressEnter = (e: KeyboardEvent) => {
  if (!e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

const handleSend = async () => {
  if (!chatInput.value.trim() || chatting.value) {
    return
  }
  const userContent = chatInput.value.trim()
  chatInput.value = ''

  const userMsg: API.GraphMessageVO = {
    id: Date.now(),
    role: 'user',
    content: userContent,
  }
  messages.value.push(userMsg)
  await scrollToBottom()

  const aiMsgId = Date.now() + 1
  const aiMsg: API.GraphMessageVO = {
    id: aiMsgId,
    role: 'assistant',
    content: 'AI 正在思考...',
  }
  messages.value.push(aiMsg)
  chatting.value = true

  try {
    const version = props.selectedVersion || props.currentVersion || 'v001'
    let aiContent = ''
    let sseError = ''
    await chatToModifyGraphSSE(
      {
        graphId: props.graphId,
        message: userContent,
        version,
      },
      {
        onMessage: async (chunk) => {
          if (!chunk) return
          if (!aiContent) {
            const index = messages.value.findIndex((m) => m.id === aiMsgId)
            if (index > -1) {
              messages.value[index].content = ''
            }
          }
          aiContent += chunk
          const index = messages.value.findIndex((m) => m.id === aiMsgId)
          if (index > -1) {
            messages.value[index].content = aiContent
          }
          await scrollToBottom()
        },
        onError: (errorMessage) => {
          sseError = errorMessage
        },
      },
    )
    if (sseError) {
      throw new Error(sseError)
    }
    if (!aiContent) {
      throw new Error('empty ai response')
    }

    message.success('AI 修改建议已同步')
    if (props.onSuccess) props.onSuccess()
    await loadMessages()
  } catch (e) {
    console.error('对话失败:', e)
    const errorText = e instanceof Error ? e.message : 'AI 对话失败，请稍后重试'
    message.error(errorText)
    messages.value = messages.value.filter((m) => m.id !== aiMsgId)
  } finally {
    chatting.value = false
  }
}

onMounted(() => {
  loadMessages()
})

defineExpose({
  scrollToBottom,
  refresh: loadMessages
})
</script>

<style scoped>
.chat-list {
  flex: 1;
}
</style>
