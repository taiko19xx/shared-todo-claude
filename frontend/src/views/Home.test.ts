import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import { createRouter, createWebHistory, type Router } from 'vue-router'
import type { MockedFunction } from 'vitest'
import Home from './Home.vue'
import * as api from '../api/api'

// API関数をモック
vi.mock('../api/api')
const mockedApi = vi.mocked(api)

describe('Home.vue', () => {
  let router: Router
  let wrapper: VueWrapper<any>

  beforeEach(() => {
    // ルーターのセットアップ
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: Home },
        { path: '/:listId/:userId', component: { template: '<div>TodoList</div>' } }
      ]
    })

    // コンポーネントのマウント
    wrapper = mount(Home, {
      global: {
        plugins: [router]
      }
    })
  })

  it('should render home page correctly', () => {
    expect(wrapper.find('h1').text()).toBe('共有ToDoリスト')
    expect(wrapper.find('button').text()).toBe('新しいリストを作成')
  })

  it('should create new list when button is clicked', async () => {
    const mockResponse = {
      listId: 'test-list-id',
      userId: 'test-user-id'
    }
    ;(mockedApi.createList as MockedFunction<any>).mockResolvedValue(mockResponse)

    const pushSpy = vi.spyOn(router, 'push')
    const button = wrapper.find('button')

    await button.trigger('click')
    
    expect(mockedApi.createList).toHaveBeenCalled()
    
    // APIレスポンスとDOM更新の完了を待機
    await vi.waitFor(() => {
      expect(pushSpy).toHaveBeenCalledWith('/test-list-id/test-user-id')
    })
  })

  it('should handle API error when creating list', async () => {
    const errorMessage = 'API Error'
    ;(mockedApi.createList as MockedFunction<any>).mockRejectedValue(new Error(errorMessage))

    // window.alertをモック
    const alertSpy = vi.spyOn(window, 'alert').mockImplementation(() => {})
    
    const button = wrapper.find('button')
    await button.trigger('click')
    
    // エラーが処理されるまで待機
    await vi.waitFor(() => {
      expect(alertSpy).toHaveBeenCalledWith('リストの作成に失敗しました')
    })
    
    alertSpy.mockRestore()
  })

  it('should show loading state when creating list', async () => {
    let resolvePromise: (value: any) => void
    const createListPromise = new Promise(resolve => {
      resolvePromise = resolve
    })
    ;(mockedApi.createList as MockedFunction<any>).mockReturnValue(createListPromise)

    const button = wrapper.find('button')
    await button.trigger('click')
    
    // ローディング状態を確認
    await wrapper.vm.$nextTick()
    expect(button.text()).toBe('作成中...')
    expect(button.attributes('disabled')).toBeDefined()
    
    // プロミスを解決
    resolvePromise!({ listId: 'test', userId: 'test' })
    
    // ローディング解除を待機
    await vi.waitFor(() => {
      expect(button.text()).toBe('新しいリストを作成')
      expect(button.attributes('disabled')).toBeUndefined()
    })
  })
})