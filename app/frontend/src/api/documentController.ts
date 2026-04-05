// @ts-ignore
/* eslint-disable */
import request from '@/request'

export async function uploadUserDocument(file: File, options?: { [key: string]: any }) {
  const formData = new FormData()
  formData.append('file', file)
  return request<API.BaseResponseBoolean>('/document/user/upload', {
    method: 'POST',
    data: formData,
    ...(options || {}),
  })
}

export async function uploadProjectDocument(
  graphId: number,
  file: File,
  options?: { [key: string]: any },
) {
  const formData = new FormData()
  formData.append('graphId', String(graphId))
  formData.append('file', file)
  return request<API.BaseResponseBoolean>('/document/project/upload', {
    method: 'POST',
    data: formData,
    ...(options || {}),
  })
}
