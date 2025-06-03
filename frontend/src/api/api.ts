import axios, { type AxiosResponse } from 'axios'
import type {
  CreateListResponse,
  GetListDataResponse,
  InviteUserResponse,
  CreateTodoRequest,
  UpdateTodoUserStatusRequest,
  UpdateListMemoRequest,
  UpdateUserNameRequest,
  Todo
} from '@/types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000
})

// レスポンスインターセプターでエラーハンドリングを統一
api.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error) => {
    if (error.response) {
      // サーバーエラーレスポンス
      const message = error.response.data?.error || 'サーバーエラーが発生しました'
      throw new Error(`${error.response.status}: ${message}`)
    } else if (error.request) {
      // ネットワークエラー
      throw new Error('ネットワークエラー: サーバーに接続できません')
    } else {
      // その他のエラー
      throw new Error('リクエストエラーが発生しました')
    }
  }
)

export const createList = async (): Promise<CreateListResponse> => {
  const response = await api.post<CreateListResponse>('/lists')
  return response.data
}

export const getListData = async (listId: string, userId: string): Promise<GetListDataResponse> => {
  const response = await api.get<GetListDataResponse>(`/lists/${listId}/users/${userId}`)
  return response.data
}

export const createTodo = async (listId: string, todo: CreateTodoRequest): Promise<Todo> => {
  const response = await api.post<Todo>(`/lists/${listId}/todos`, todo)
  return response.data
}

export const updateTodoUserStatus = async (
  todoId: number,
  userId: string,
  checked: boolean
): Promise<{ checked: boolean }> => {
  const requestData: UpdateTodoUserStatusRequest = { checked }
  const response = await api.put<{ checked: boolean }>(`/todos/${todoId}/status/${userId}`, requestData)
  return response.data
}

export const updateListMemo = async (listId: string, memo: string): Promise<{ memo: string }> => {
  const requestData: UpdateListMemoRequest = { memo }
  const response = await api.put<{ memo: string }>(`/lists/${listId}/memo`, requestData)
  return response.data
}

export const inviteUser = async (listId: string): Promise<InviteUserResponse> => {
  const response = await api.post<InviteUserResponse>(`/lists/${listId}/users`)
  return response.data
}

export const updateUserName = async (
  listId: string,
  userId: string,
  name: string
): Promise<{ name: string }> => {
  const requestData: UpdateUserNameRequest = { name }
  const response = await api.put<{ name: string }>(`/lists/${listId}/users/${userId}/name`, requestData)
  return response.data
}