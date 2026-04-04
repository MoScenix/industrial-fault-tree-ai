// @ts-ignore
/* eslint-disable */
import request from '@/request'

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
    params: {
      ...params,
    },
    ...(options || {}),
  })
}
