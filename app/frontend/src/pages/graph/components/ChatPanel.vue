<template>
  <div class="flex flex-col h-full">
    <!-- Empty State -->
    <div v-if="!messages.length" class="flex-1 flex flex-col items-center justify-center text-center px-5 select-none">
      <div class="text-3xl mb-3">💬</div>
      <div class="text-gray-400 text-sm">输入需求，AI 将辅助您构建故障树</div>
    </div>

    <!-- Messages -->
    <div v-else ref="chatListRef" class="flex-1 overflow-y-auto px-5 py-5 space-y-5">
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

    <!-- Input Area -->
    <div class="border-t border-gray-100 px-4 py-3 bg-white">
      <div
        ref="inputAreaRef"
        class="flex flex-col bg-white border border-gray-200 rounded-3xl shadow-sm transition-all duration-200 focus-within:border-gray-300 focus-within:shadow-md"
        @dragover.prevent
        @dragleave.prevent
        @drop.prevent="handleDrop"
      >
        <!-- File Previews -->
        <div v-if="files.length > 0" class="flex flex-wrap gap-2 px-4 pt-3 pb-0">
          <div v-for="(file, idx) in files" :key="idx" class="relative group">
            <div
              v-if="file.isImage"
              class="w-16 h-16 rounded-xl overflow-hidden cursor-pointer"
              @click="openImagePreview(file.url)"
            >
              <img :src="file.url" :alt="file.name" class="h-full w-full object-cover" />
            </div>
            <div v-else class="w-16 h-16 rounded-xl bg-gray-50 border border-gray-200 flex flex-col items-center justify-center text-center p-1">
              <svg viewBox="0 0 24 24" class="w-6 h-6 text-amber-500 mb-0.5" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
                <polyline points="14 2 14 8 20 8" />
                <line x1="16" y1="13" x2="8" y2="13" />
                <line x1="16" y1="17" x2="8" y2="17" />
              </svg>
              <span class="text-[9px] text-gray-500 truncate w-full leading-tight">{{ file.name.split('.').pop()?.toUpperCase() }}</span>
            </div>
            <button class="absolute -top-1.5 -right-1.5 w-5 h-5 rounded-full bg-gray-800 text-white flex items-center justify-center shadow-sm opacity-0 group-hover:opacity-100 transition-opacity" @click="removeFile(idx)">
              <CloseOutlined class="text-xs" />
            </button>
          </div>
        </div>

        <!-- Uploading indicator -->
        <div v-if="uploadingFiles.length > 0" class="px-4 pt-3 pb-0">
          <div v-for="(item, idx) in uploadingFiles" :key="idx" class="flex items-center gap-2 text-xs text-gray-500 mb-1">
            <LoadingOutlined class="animate-spin text-blue-500" />
            <span>上传中：{{ item }}</span>
          </div>
        </div>

        <!-- Textarea -->
        <div class="px-4 pt-3 pb-1">
          <textarea
            ref="inputRef"
            v-model="chatInput"
            rows="1"
            maxlength="2000"
            :placeholder="currentPlaceholder"
            class="w-full bg-transparent border-none outline-none resize-none text-sm text-gray-800 placeholder-gray-400 leading-relaxed scrollbar-thin"
            @input="autoResize"
            @keydown="onInputKeydown"
          ></textarea>
        </div>

        <!-- Actions Bar -->
        <div class="flex items-center justify-between px-3 pb-2">
          <!-- Left: Toolbar -->
          <div class="flex items-center gap-0.5">
            <!-- Upload -->
            <button
              class="w-8 h-8 flex items-center justify-center rounded-full text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors"
              @click="fileInputRef?.click()"
              :title="'上传项目文档'"
            >
              <PaperClipOutlined class="text-sm" />
            </button>
            <input ref="fileInputRef" type="file" accept=".pdf" class="hidden" @change="handleFileSelect" />

            <div class="w-px h-5 bg-gray-200 mx-1"></div>

            <!-- Search mode toggle -->
            <button
              class="flex items-center gap-1 px-2.5 py-1.5 rounded-full text-xs transition-all whitespace-nowrap"
              :class="mode === 'search'
                ? 'bg-blue-50 text-blue-600 border border-blue-200'
                : 'text-gray-400 hover:text-gray-600'"
              @click="toggleMode('search')"
            >
              <GlobalOutlined class="text-sm flex-shrink-0" />
              <Transition name="mode-fade">
                <span v-if="mode === 'search'" class="text-xs overflow-hidden">搜索</span>
              </Transition>
            </button>

            <!-- Think mode toggle -->
            <button
              class="flex items-center gap-1 px-2.5 py-1.5 rounded-full text-xs transition-all whitespace-nowrap"
              :class="mode === 'think'
                ? 'bg-purple-50 text-purple-600 border border-purple-200'
                : 'text-gray-400 hover:text-gray-600'"
              @click="toggleMode('think')"
            >
              <BulbOutlined class="text-sm flex-shrink-0" />
              <Transition name="mode-fade">
                <span v-if="mode === 'think'" class="text-xs overflow-hidden">思考</span>
              </Transition>
            </button>

            <!-- Canvas mode toggle -->
            <button
              class="flex items-center gap-1 px-2.5 py-1.5 rounded-full text-xs transition-all whitespace-nowrap"
              :class="mode === 'canvas'
                ? 'bg-orange-50 text-orange-600 border border-orange-200'
                : 'text-gray-400 hover:text-gray-600'"
              @click="toggleMode('canvas')"
            >
              <FolderOpenOutlined class="text-sm flex-shrink-0" />
              <Transition name="mode-fade">
                <span v-if="mode === 'canvas'" class="text-xs overflow-hidden">画布</span>
              </Transition>
            </button>
          </div>

          <!-- Right: Send Button -->
          <button
            class="w-8 h-8 flex items-center justify-center rounded-full transition-all flex-shrink-0"
            :class="
              chatting
                ? 'bg-gray-700 text-white'
                : canSend
                ? 'bg-gray-800 text-white hover:bg-gray-700 shadow-sm'
                : 'bg-gray-100 text-gray-300 cursor-default'
            "
            :disabled="!canSend && !chatting"
            @click="handleSendButton"
          >
            <ArrowUpOutlined v-if="!chatting" class="text-sm" />
            <svg v-else viewBox="0 0 24 24" class="w-3.5 h-3.5" fill="currentColor"><rect x="6" y="6" width="12" height="12" rx="1.5" /></svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Image Preview Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="previewImage"
          class="fixed inset-0 z-[9999] bg-black/60 backdrop-blur-sm flex items-center justify-center"
          @click="previewImage = null"
        >
          <div class="relative max-w-[90vw] max-h-[85vh]" @click.stop>
            <img :src="previewImage" class="max-w-full max-h-[85vh] object-contain rounded-2xl shadow-2xl" />
            <button
              class="absolute top-3 right-3 w-8 h-8 rounded-full bg-black/50 text-white flex items-center justify-center hover:bg-black/70 transition-colors"
              @click="previewImage = null"
            >
              <CloseOutlined class="text-sm" />
            </button>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, nextTick, onMounted, onBeforeUnmount, watch } from 'vue'
import { message } from 'ant-design-vue'
import {
  LoadingOutlined,
  ArrowUpOutlined,
  PaperClipOutlined,
  CloseOutlined,
  GlobalOutlined,
  BulbOutlined,
  FolderOpenOutlined,
} from '@ant-design/icons-vue'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'
import { chatToModifyGraphSSE, listGraphMessage } from '@/api/graphController'
import { uploadProjectDocument } from '@/api/documentController'

interface Props {
  graphId: number
  selectedVersion: string
  currentVersion?: string
  isEditorDirty: boolean
  onSuccess?: () => void
}

const props = defineProps<Props>()

type InputMode = '' | 'search' | 'think' | 'canvas'

interface FileItem {
  name: string
  url: string
  file: File
  isImage: boolean
}

const messages = ref<API.GraphMessageVO[]>([])
const chatInput = ref('')
const chatting = ref(false)
const mode = ref<InputMode>('')
const files = ref<FileItem[]>([])
const uploadingFiles = ref<string[]>([])
const previewImage = ref<string | null>(null)

const chatListRef = ref<HTMLElement>()
const inputRef = ref<HTMLTextAreaElement>()
const fileInputRef = ref<HTMLInputElement>()
const inputAreaRef = ref<HTMLElement>()

const canSend = computed(() => (chatInput.value.trim().length > 0 || files.value.length > 0) && !chatting.value)

const currentPlaceholder = computed(() => {
  if (mode.value === 'search') return '搜索网络信息...'
  if (mode.value === 'think') return '深入推理思考...'
  if (mode.value === 'canvas') return '生成或更新故障树图...'
  return '给 AI 发送消息...'
})

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

const autoResize = () => {
  if (!inputRef.value) return
  inputRef.value.style.height = 'auto'
  inputRef.value.style.height = Math.min(inputRef.value.scrollHeight, 200) + 'px'
}

const onInputKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

const toggleMode = (newMode: InputMode) => {
  mode.value = mode.value === newMode ? '' : newMode
  nextTick(() => inputRef.value?.focus())
}

// --- File Handling ---

const processFile = (file: File) => {
  if (file.size > 20 * 1024 * 1024) {
    message.warning('文件不能超过 20MB')
    return
  }
  if (file.type.startsWith('image/')) {
    const reader = new FileReader()
    reader.onload = (e) => {
      files.value = [{
        name: file.name,
        url: e.target?.result as string,
        file,
        isImage: true,
      }]
    }
    reader.readAsDataURL(file)
  } else {
    // Document file — only store metadata, no preview URL
    files.value = [{
      name: file.name,
      url: '',
      file,
      isImage: false,
    }]
  }
}

const handleFileSelect = () => {
  const input = fileInputRef.value
  if (!input?.files?.length) return
  processFile(input.files[0])
  input.value = ''
}

const handleDrop = (e: DragEvent) => {
  const droppedFiles = Array.from(e.dataTransfer?.files || [])
  if (droppedFiles.length > 0) processFile(droppedFiles[0])
}

const handlePaste = (e: ClipboardEvent) => {
  const items = e.clipboardData?.items
  if (!items) return
  for (const item of items) {
    if (item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) {
        e.preventDefault()
        processFile(file)
        break
      }
    }
  }
}

const removeFile = (idx: number) => {
  files.value.splice(idx, 1)
}

const openImagePreview = (url: string) => {
  previewImage.value = url
}

// --- Send ---

const handleSendButton = () => {
  if (chatting.value) {
    // TODO: abort SSE connection
    return
  }
  handleSend()
}

const handleSend = async () => {
  if (!canSend.value && files.value.length === 0) return

  // --- Step 1: Upload documents first ---
  const uploadedDocNames: string[] = []
  if (files.value.length > 0) {
    const pendingFiles = [...files.value]
    files.value = []
    uploadingFiles.value = pendingFiles.map(f => f.name)

    try {
      for (const item of pendingFiles) {
        const res = await uploadProjectDocument(props.graphId, item.file)
        if (res.data.code === 0) {
          uploadedDocNames.push(item.name)
        } else {
          message.error(`「${item.name}」上传失败：${res.data.message || '未知错误'}`)
        }
      }
    } catch (e) {
      message.error('文档上传出错，请稍后重试')
      uploadingFiles.value = []
      return
    }
    uploadingFiles.value = []
  }

  // --- Step 2: Build message with document prefix ---
  let userContent = chatInput.value.trim()
  if (!userContent && uploadedDocNames.length > 0) {
    userContent = '请根据新上传的项目文档分析内容'
  }

  // Mode prefix
  if (mode.value) {
    userContent = `[${mode.value === 'search' ? '搜索' : mode.value === 'think' ? '思考' : '画布'}: ${userContent}]`
  }

  // Document upload prefix (let AI know to base answer on the new docs)
  if (uploadedDocNames.length > 0) {
    const names = uploadedDocNames.join('、')
    userContent = `[已上传项目文档: ${names}] 根据新上传的项目文档：${userContent}`
  }

  chatInput.value = ''

  // Reset height to single row
  if (inputRef.value) {
    inputRef.value.style.height = 'auto'
  }

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
    if (sseError) throw new Error(sseError)
    if (!aiContent) throw new Error('empty ai response')

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

// --- Watch empty state ---
watch(() => messages.value.length, (len) => {
  if (len > 0) {
    nextTick(scrollToBottom)
  }
})

onMounted(() => {
  loadMessages()
  document.addEventListener('paste', handlePaste)
})

onBeforeUnmount(() => {
  document.removeEventListener('paste', handlePaste)
})

defineExpose({
  scrollToBottom,
  refresh: loadMessages,
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

/* Mode label transition */
.mode-fade-enter-active,
.mode-fade-leave-active {
  transition: all 0.15s ease-out;
}
.mode-fade-enter-from,
.mode-fade-leave-to {
  opacity: 0;
  width: 0;
  margin: 0;
}

/* Image preview modal transition */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

/* Thin scrollbar for textarea */
:deep(textarea::-webkit-scrollbar) {
  width: 4px;
}
:deep(textarea::-webkit-scrollbar-track) {
  background: transparent;
}
:deep(textarea::-webkit-scrollbar-thumb) {
  background-color: #d1d5db;
  border-radius: 2px;
}
</style>
