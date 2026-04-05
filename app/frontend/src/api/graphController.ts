// @ts-ignore
/* eslint-disable */
import request from '@/request'
import { API_BASE_URL } from '@/config/env'

export async function addGraph(body: API.GraphAddRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseLong>('/graph/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function deleteGraph(body: API.DeleteRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/graph/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function updateGraph(body: API.GraphUpdateRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/graph/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function getGraphVoById(
  params: API.getGraphVOByIdParams,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponseGraphVO>('/graph/get/vo', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

export async function listGraphVoByPage(
  body: API.GraphQueryRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponsePageGraphVO>('/graph/list/page/vo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function startEdit(body: API.StartEditRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseGraphEditState>('/graph/edit/start', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function getWorkingGraph(
  params: API.getWorkingGraphParams,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponseWorkingGraph>('/graph/edit/working', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

export async function getCurrentSuggestion(
  params: API.getCurrentSuggestionParams,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponseGraphSuggestion>('/graph/suggestion/current', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

export async function discardWorkingGraph(
  body: API.DiscardWorkingGraphRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponseBoolean>('/graph/edit/discard', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function saveGraph(body: API.SaveGraphRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseSaveResult>('/graph/save', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function listGraphVersion(
  params: API.listGraphVersionParams,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponsePageGraphVersionVO>('/graph/version/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

export async function createGraphVersion(
  body: API.CreateGraphVersionRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponseString>('/graph/version/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function deleteGraphVersion(
  body: API.DeleteGraphVersionRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponseBoolean>('/graph/version/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function renameGraphVersion(
  body: API.RenameGraphVersionRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponseBoolean>('/graph/version/rename', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

export async function listGraphMessage(
  params: API.listGraphMessageParams,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponsePageGraphMessageVO>('/graph/message/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

export async function downloadGraph(
  params: API.downloadGraphParams,
  options?: { [key: string]: any },
) {
  const { graphId, ...queryParams } = params
  return request<API.BaseResponseBytes>(`/graph/download/${graphId}`, {
    method: 'GET',
    params: {
      ...queryParams,
    },
    responseType: 'blob',
    ...(options || {}),
  })
}

export async function chatToModifyGraph(
  params: API.chatToModifyGraphParams,
  options?: { [key: string]: any },
) {
  return request<API.ServerSentEventString>('/graph/chat', {
    method: 'GET',
    params,
    ...(options || {}),
  })
}

/**
 * AI 对话流式接口 (SSE)
 */
export async function chatToModifyGraphSSE(
  params: API.chatToModifyGraphParams,
  handlers: {
    onMessage?: (chunk: string) => void
    onError?: (error: string) => void
    onDone?: () => void
  },
) {
  const url = new URL(`${API_BASE_URL}/graph/chat`, window.location.origin)
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      url.searchParams.set(key, String(value))
    }
  })

  try {
    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: {
        Accept: 'text/event-stream',
      },
      credentials: 'include',
    })

    if (!response.ok) {
      const errorText = await response.text()
      handlers.onError?.(`HTTP ${response.status}: ${errorText}`)
      return
    }

    const reader = response.body?.getReader()
    if (!reader) {
      handlers.onError?.('SSE reader unavailable')
      return
    }

    const decoder = new TextDecoder()
    let buffer = ''

    const flushEvent = (eventBlock: string) => {
      const lines = eventBlock.split('\n')
      let eventName = 'message'
      const dataLines: string[] = []

      for (const rawLine of lines) {
        const line = rawLine.trim()
        if (!line) continue
        if (line.startsWith('event:')) {
          eventName = line.slice(6).trim()
          continue
        }
        if (line.startsWith('data:')) {
          dataLines.push(line.slice(5).trim())
        }
      }

      const dataText = dataLines.join('\n')
      if (eventName === 'done') {
        handlers.onDone?.()
        return
      }
      if (!dataText) return

      try {
        const payload = JSON.parse(dataText) as API.ServerSentEventString
        if (eventName === 'business-error') {
          handlers.onError?.(payload.message || 'AI 对话失败')
          return
        }
        if (payload.d) {
          handlers.onMessage?.(payload.d)
        }
      } catch {
        if (eventName === 'business-error') {
          handlers.onError?.(dataText || 'AI 对话失败')
          return
        }
        handlers.onMessage?.(dataText)
      }
    }

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const parts = buffer.split('\n\n')
      buffer = parts.pop() || ''
      for (const part of parts) {
        flushEvent(part)
      }
    }

    if (buffer.trim()) {
      flushEvent(buffer)
    }
    handlers.onDone?.()
  } catch (err: any) {
    handlers.onError?.(err.message || 'Network error')
  }
}

export async function validateGraph(body: API.ValidateGraphRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/graph/validate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
