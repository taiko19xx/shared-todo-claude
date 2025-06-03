import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'

describe('App.vue', () => {
  it('should render app container with correct classes', () => {
    const router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } }
      ]
    })

    const wrapper = mount(App, {
      global: {
        plugins: [router]
      }
    })

    const appContainer = wrapper.find('#app')
    expect(appContainer.exists()).toBe(true)
    expect(appContainer.classes()).toContain('min-h-screen')
    expect(appContainer.classes()).toContain('bg-gray-50')
  })

  it('should render router-view correctly', async () => {
    const TestComponent = { template: '<div class="test-content">Test Content</div>' }
    const router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: TestComponent }
      ]
    })

    await router.push('/')
    await router.isReady()

    const wrapper = mount(App, {
      global: {
        plugins: [router]
      }
    })

    // Wait for component to render
    await wrapper.vm.$nextTick()
    
    // Check if router content is rendered
    expect(wrapper.find('.test-content').exists()).toBe(true)
    expect(wrapper.text()).toContain('Test Content')
  })
})