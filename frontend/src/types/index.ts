// ユーザー関連の型定義
export interface User {
  id: string
  displayName: string
  listId?: string
  createdAt?: string
}

// リスト関連の型定義
export interface List {
  id: string
  memo: string
  createdAt?: string
  updatedAt?: string
}

// ToDo関連の型定義
export interface Todo {
  id: number
  listId: string
  title: string
  priority: 'high' | 'medium' | 'low'
  dueDate: string | null
  isCompleted: boolean
  createdAt?: string
  updatedAt?: string
  userStatuses?: TodoUserStatus[]
}

export interface TodoUserStatus {
  todoId: number
  userId: string
  isChecked: boolean
  checkedAt?: string | null
}

// API レスポンス関連の型定義
export interface CreateListResponse {
  listId: string
  userId: string
}

export interface GetListDataResponse {
  users: User[]
  todos: Todo[]
  memo: string
}

export interface InviteUserResponse {
  userId: string
  url: string
}

export interface CreateTodoRequest {
  title: string
  priority: 'high' | 'medium' | 'low'
  dueDate: string | null
}

export interface UpdateTodoUserStatusRequest {
  checked: boolean
}

export interface UpdateListMemoRequest {
  memo: string
}

export interface UpdateUserNameRequest {
  name: string
}

// フォーム関連の型定義
export interface TodoForm {
  title: string
  priority: 'high' | 'medium' | 'low'
  dueDate: string
}

// エラー関連の型定義
export interface ApiError {
  message: string
  status?: number
}