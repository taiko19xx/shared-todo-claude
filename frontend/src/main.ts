import { createApp } from 'vue'
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import './style.css'
import App from './App.vue'
import Home from './views/Home.vue'
import TodoList from './views/TodoList.vue'

const routes: RouteRecordRaw[] = [
  { path: '/', component: Home },
  { path: '/:listId/:userId', component: TodoList, props: true }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

createApp(App).use(router).mount('#app')