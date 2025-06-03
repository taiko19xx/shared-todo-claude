<template>
  <div class="container mx-auto px-4 py-8">
    <div class="max-w-md mx-auto bg-white rounded-lg shadow-md p-6">
      <h1 class="text-2xl font-bold text-center mb-6">共有ToDoリスト</h1>
      <button
        @click="createNewList"
        :disabled="loading"
        class="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 text-white font-bold py-2 px-4 rounded transition duration-200"
      >
        {{ loading ? '作成中...' : '新しいリストを作成' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createList } from '../api/api'

const router = useRouter()
const loading = ref<boolean>(false)

const createNewList = async (): Promise<void> => {
  loading.value = true
  try {
    const response = await createList()
    await router.push(`/${response.listId}/${response.userId}`)
  } catch (error) {
    console.error('Failed to create list:', error)
    alert('リストの作成に失敗しました')
  } finally {
    loading.value = false
  }
}
</script>