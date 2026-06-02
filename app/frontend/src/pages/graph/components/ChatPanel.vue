<template>
  <div class="flex flex-col h-full">
    <!-- Messages -->
    <div ref="chatListRef" class="flex-1 overflow-y-auto px-5 py-5 space-y-5">
      <div v-if="!messages.length" class="flex flex-col items-center justify-center text-center py-20 select-none">
        <div class="text-3xl mb-3">💬</div>
        <div class="text-gray-400 text-sm">输入需求，AI 将辅助您构建故障树</div>
      </div>

      <template v-for="item in messages" :key="item.id">
        <!-- user message -->
        <div v-if="item.role === 'user'" class="flex justify-end">
          <div class="max-w-[85%] bg-gray-100 rounded-2xl px-4 py-[10px]">
            <MarkdownRenderer class="text-sm text-gray-800 leading-6" :content="(item.content || '').trimEnd()" />
          </div>
        </div>

        <!-- assistant message -->
        <div v-else class="flex justify-start">
          <div class="max-w-[90%]">
            <MarkdownRenderer class="text-sm text-gray-800 leading-7" :content="(item.content || '').trimEnd()" />
          </div>
        </div>
      </template>
    </div>

    <!-- Input -->
    <div class="border-t border-gray-100 px-4 pt-2 pb-4">
      <div class="flex items-center gap-2 bg-white border border-gray-200/80 rounded-full shadow-sm pl-4 pr-1.5 py-1 transition-all focus-within:border-indigo-400 focus-within:ring-[1.5px] focus-within:ring-indigo-400/30">
        <textarea
          v-model="chatInput"
          rows="1"
          maxlength="1000"
          placeholder="给 AI 发送消息..."
          class="flex-1 bg-transparent border-none outline-none resize-none text-sm text-gray-800 placeholder-gray-400 leading-normal py-1.5"
          @keydown="onInputKeydown"
          ref="inputRef"
        ></textarea>
        <button
          class="w-7 h-7 flex items-center justify-center rounded-full flex-shrink-0 transition-all"
          :class="canSend ? 'bg-gray-800 text-white hover:bg-gray-700' : 'bg-gray-100 text-gray-300 cursor-default'"
          :disabled="!canSend"
          @click="handleSend"
        >
          <svg v-if="!chatting" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4">
            <line x1="12" y1="19" x2="12" y2="5" />
            <polyline points="5 12 12 5 19 12" />
          </svg>
          <LoadingOutlined v-else class="animate-spin" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, nextTick, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { LoadingOutlined } from '@ant-design/icons-vue'
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
const inputRef = ref<HTMLTextAreaElement>()

const canSend = computed(() => chatInput.value.trim().length > 0 && !chatting.value)

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

const onInputKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

const handleSend = async () => {
  if (!canSend.value) return
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
/* Clean markdown inside messages */
:deep(.custom-md) {
  background: transparent !important;
  font-size: inherit !important;
  color: inherit !important;
}

:deep(.custom-md p) {
  margin: 0;
}

:deep(.custom-md p + p) {
  margin-top: 0.5em;
}

:deep(.custom-md code) {
  font-size: 12px;
  padding: 1px 4px;
  border-radius: 4px;
  background: #e5e7eb;
}

:deep(.custom-md ul),
:deep(.custom-md ol) {
  margin: 0.25em 0;
  padding-left: 1.25em;
}

:deep(.custom-md li) {
  margin: 0.125em 0;
}
</style>
