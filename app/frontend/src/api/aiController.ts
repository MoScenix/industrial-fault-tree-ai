// @ts-ignore
/* eslint-disable */
import request from '@/request'

export async function getPrompt(params: API.getPromptParams, options?: { [key: string]: any }) {
  return request<API.BaseResponsePromptVO>('/ai/prompt/get', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

export async function updatePrompt(
  body: API.UpdatePromptRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponseBoolean>('/ai/prompt/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
