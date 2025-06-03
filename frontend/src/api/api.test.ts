import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import type { MockedFunction } from 'vitest'

// Mockする axios instance
const mockAxiosInstance = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
  interceptors: {
    response: {
      use: vi.fn()
    }
  }
}

// axiosをモック
vi.mock('axios', () => ({
  default: {
    create: vi.fn(() => mockAxiosInstance)
  }
}))

// apiモジュールのインポート（axiosモック後）
const { createList, getListData, createTodo, updateTodoUserStatus, updateListMemo, inviteUser, updateUserName } = await import('./api')

describe('API Functions', () => {
  beforeEach(() => {
    // Clear all mocks before each test
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('createList', () => {
    it('should create a new list successfully', async () => {
      const mockResponse = {
        data: {
          listId: 'test-list-id',
          userId: 'test-user-id'
        }
      }
      ;(mockAxiosInstance.post as MockedFunction<any>).mockResolvedValue(mockResponse)

      const result = await createList()

      expect(mockAxiosInstance.post).toHaveBeenCalledWith('/lists')
      expect(result).toEqual(mockResponse.data)
    })

    it('should handle API errors', async () => {
      const errorMessage = 'Network Error'
      ;(mockAxiosInstance.post as MockedFunction<any>).mockRejectedValue(new Error(errorMessage))

      await expect(createList()).rejects.toThrow(errorMessage)
    })
  })

  describe('getListData', () => {
    it('should get list data successfully', async () => {
      const mockResponse = {
        data: {
          users: [{ id: 'user1', displayName: 'User 1' }],
          todos: [{ id: 1, title: 'Test Todo' }],
          memo: 'Test memo'
        }
      }
      ;(mockAxiosInstance.get as MockedFunction<any>).mockResolvedValue(mockResponse)

      const result = await getListData('list-id', 'user-id')

      expect(mockAxiosInstance.get).toHaveBeenCalledWith('/lists/list-id/users/user-id')
      expect(result).toEqual(mockResponse.data)
    })
  })

  describe('createTodo', () => {
    it('should create a todo successfully', async () => {
      const todoData = {
        title: 'New Todo',
        priority: 'high' as const,
        dueDate: '2025-06-10'
      }
      const mockResponse = {
        data: { id: 1, ...todoData }
      }
      ;(mockAxiosInstance.post as MockedFunction<any>).mockResolvedValue(mockResponse)

      const result = await createTodo('list-id', todoData)

      expect(mockAxiosInstance.post).toHaveBeenCalledWith('/lists/list-id/todos', todoData)
      expect(result).toEqual(mockResponse.data)
    })
  })

  describe('updateTodoUserStatus', () => {
    it('should update todo status successfully', async () => {
      const mockResponse = { data: { checked: true } }
      ;(mockAxiosInstance.put as MockedFunction<any>).mockResolvedValue(mockResponse)

      const result = await updateTodoUserStatus(1, 'user-id', true)

      expect(mockAxiosInstance.put).toHaveBeenCalledWith('/todos/1/status/user-id', { checked: true })
      expect(result).toEqual(mockResponse.data)
    })
  })

  describe('updateListMemo', () => {
    it('should update list memo successfully', async () => {
      const memo = 'Updated memo'
      const mockResponse = { data: { memo } }
      ;(mockAxiosInstance.put as MockedFunction<any>).mockResolvedValue(mockResponse)

      const result = await updateListMemo('list-id', memo)

      expect(mockAxiosInstance.put).toHaveBeenCalledWith('/lists/list-id/memo', { memo })
      expect(result).toEqual(mockResponse.data)
    })
  })

  describe('inviteUser', () => {
    it('should invite a user successfully', async () => {
      const mockResponse = {
        data: {
          userId: 'new-user-id',
          url: '/list-id/new-user-id'
        }
      }
      ;(mockAxiosInstance.post as MockedFunction<any>).mockResolvedValue(mockResponse)

      const result = await inviteUser('list-id')

      expect(mockAxiosInstance.post).toHaveBeenCalledWith('/lists/list-id/users')
      expect(result).toEqual(mockResponse.data)
    })
  })

  describe('updateUserName', () => {
    it('should update user name successfully', async () => {
      const name = 'New Name'
      const mockResponse = { data: { name } }
      ;(mockAxiosInstance.put as MockedFunction<any>).mockResolvedValue(mockResponse)

      const result = await updateUserName('list-id', 'user-id', name)

      expect(mockAxiosInstance.put).toHaveBeenCalledWith('/lists/list-id/users/user-id/name', { name })
      expect(result).toEqual(mockResponse.data)
    })
  })
})