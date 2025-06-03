import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import { createRouter, createWebHistory, type Router } from 'vue-router'
import type { MockedFunction } from 'vitest'
import TodoList from './TodoList.vue'
import * as api from '../api/api'
import type { GetListDataResponse, Todo, User } from '../types'

// API関数をモック
vi.mock('../api/api')
const mockedApi = vi.mocked(api)

describe('TodoList.vue', () => {
  let router: Router
  let wrapper: VueWrapper<any>

  const mockData: GetListDataResponse = {
    users: [
      { id: 'user1', displayName: 'User 1' },
      { id: 'user2', displayName: 'User 2' }
    ] as User[],
    todos: [
      {
        id: 1,
        listId: 'test-list',
        title: 'Test Todo 1',
        priority: 'high',
        dueDate: '2025-06-10',
        isCompleted: false,
        userStatuses: [
          { todoId: 1, userId: 'user1', isChecked: true },
          { todoId: 1, userId: 'user2', isChecked: false }
        ]
      },
      {
        id: 2,
        listId: 'test-list',
        title: 'Completed Todo',
        priority: 'medium',
        dueDate: null,
        isCompleted: true,
        userStatuses: [
          { todoId: 2, userId: 'user1', isChecked: true },
          { todoId: 2, userId: 'user2', isChecked: true }
        ]
      }
    ] as Todo[],
    memo: 'Test memo'
  }

  beforeEach(async () => {
    // ルーターのセットアップ
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/:listId/:userId', component: TodoList }
      ]
    })

    // APIモックの設定
    ;(mockedApi.getListData as MockedFunction<any>).mockResolvedValue(mockData)

    // ルートを設定
    await router.push('/test-list/user1')

    // コンポーネントのマウント
    wrapper = mount(TodoList, {
      props: {
        listId: 'test-list',
        userId: 'user1'
      },
      global: {
        plugins: [router]
      }
    })

    // データロードの完了を待機
    await wrapper.vm.$nextTick()
  })

  afterEach(() => {
    vi.clearAllMocks()
    if (wrapper) {
      wrapper.unmount()
    }
  })

  it('should render TodoList page correctly', () => {
    expect(wrapper.find('h1').text()).toBe('ToDo リスト')
    expect(wrapper.find('h2').text()).toBe('新しいToDoを追加')
  })

  it('should load data on mount', () => {
    expect(mockedApi.getListData).toHaveBeenCalledWith('test-list', 'user1')
    expect(wrapper.vm.users).toEqual(mockData.users)
    expect(wrapper.vm.todos).toEqual(mockData.todos)
    expect(wrapper.vm.memo).toBe(mockData.memo)
  })

  it('should render active todos correctly', () => {
    const activeTodos = wrapper.vm.activeTodos
    expect(activeTodos).toHaveLength(1)
    expect(activeTodos[0].title).toBe('Test Todo 1')
  })

  it('should render completed todos correctly', () => {
    const completedTodos = wrapper.vm.completedTodos
    expect(completedTodos).toHaveLength(1)
    expect(completedTodos[0].title).toBe('Completed Todo')
  })

  it('should create new todo when form is submitted', async () => {
    ;(mockedApi.createTodo as MockedFunction<any>).mockResolvedValue({ id: 3, title: 'New Todo' })
    ;(mockedApi.getListData as MockedFunction<any>).mockResolvedValue(mockData) // リロード用

    // フォームに入力
    const titleInput = wrapper.find('input[placeholder="タイトル"]')
    await titleInput.setValue('New Todo')

    const form = wrapper.find('form')
    await form.trigger('submit')

    expect(mockedApi.createTodo).toHaveBeenCalledWith('test-list', {
      title: 'New Todo',
      priority: 'medium',
      dueDate: null
    })
  })

  it('should save memo when save button is clicked', async () => {
    ;(mockedApi.updateListMemo as MockedFunction<any>).mockResolvedValue({ memo: 'Updated memo' })
    ;(mockedApi.getListData as MockedFunction<any>).mockResolvedValue(mockData)

    // window.alertをモック
    const alertSpy = vi.spyOn(window, 'alert').mockImplementation(() => {})

    // メモを更新
    const textarea = wrapper.find('textarea')
    await textarea.setValue('Updated memo')

    const saveButton = wrapper.findAll('button').find(btn => btn.text().includes('メモを保存'))
    await saveButton!.trigger('click')

    expect(mockedApi.updateListMemo).toHaveBeenCalledWith('test-list', 'Updated memo')
    
    // アラートが表示されるまで待機
    await vi.waitFor(() => {
      expect(alertSpy).toHaveBeenCalledWith('メモを保存しました')
    })

    alertSpy.mockRestore()
  })

  it('should show invite modal when invite button is clicked', async () => {
    const inviteButton = wrapper.findAll('button').find(btn => btn.text().includes('ユーザーを招待'))
    await inviteButton!.trigger('click')

    expect(wrapper.vm.showInviteModal).toBe(true)
    expect(wrapper.find('.fixed').exists()).toBe(true)
  })

  it('should generate invite URL when button is clicked in modal', async () => {
    ;(mockedApi.inviteUser as MockedFunction<any>).mockResolvedValue({
      userId: 'new-user',
      url: '/test-list/new-user'
    })

    // モーダルを開く
    wrapper.vm.showInviteModal = true
    await wrapper.vm.$nextTick()

    const generateButton = wrapper.findAll('button').find(btn => btn.text().includes('URL生成'))
    await generateButton!.trigger('click')

    expect(mockedApi.inviteUser).toHaveBeenCalledWith('test-list')
    expect(wrapper.vm.inviteUrl).toBe('http://localhost:3000/test-list/new-user')
  })

  it('should show name modal when name button is clicked', async () => {
    const nameButton = wrapper.findAll('button').find(btn => btn.text().includes('表示名を設定'))
    await nameButton!.trigger('click')

    expect(wrapper.vm.showNameModal).toBe(true)
  })

  it('should update user name when form is submitted', async () => {
    ;(mockedApi.updateUserName as MockedFunction<any>).mockResolvedValue({ name: 'New Name' })
    ;(mockedApi.getListData as MockedFunction<any>).mockResolvedValue(mockData)

    // window.alertをモック
    const alertSpy = vi.spyOn(window, 'alert').mockImplementation(() => {})

    // モーダルを開く
    wrapper.vm.showNameModal = true
    wrapper.vm.newDisplayName = 'New Name'
    await wrapper.vm.$nextTick()

    const updateButton = wrapper.findAll('button').find(btn => btn.text().includes('更新'))
    await updateButton!.trigger('click')

    expect(mockedApi.updateUserName).toHaveBeenCalledWith('test-list', 'user1', 'New Name')
    
    // アラートが表示されるまで待機
    await vi.waitFor(() => {
      expect(alertSpy).toHaveBeenCalledWith('表示名を更新しました')
    })

    alertSpy.mockRestore()
  })

  it('should update todo status when checkbox is clicked', async () => {
    ;(mockedApi.updateTodoUserStatus as MockedFunction<any>).mockResolvedValue({ checked: true })
    ;(mockedApi.getListData as MockedFunction<any>).mockResolvedValue(mockData)

    await wrapper.vm.updateTodoStatus(1, 'user1', false)

    expect(mockedApi.updateTodoUserStatus).toHaveBeenCalledWith(1, 'user1', false)
  })

  it('should handle API errors gracefully', async () => {
    ;(mockedApi.getListData as MockedFunction<any>).mockRejectedValue(new Error('404: User not found in this list'))

    // window.alertをモック
    const alertSpy = vi.spyOn(window, 'alert').mockImplementation(() => {})
    const pushSpy = vi.spyOn(router, 'push')

    await wrapper.vm.loadData()

    expect(alertSpy).toHaveBeenCalledWith('指定されたリストまたはユーザーが見つかりません')
    expect(pushSpy).toHaveBeenCalledWith('/')

    alertSpy.mockRestore()
  })
})