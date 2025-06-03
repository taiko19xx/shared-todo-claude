<template>
  <div class="container mx-auto px-4 py-8">
    <div class="bg-white rounded-lg shadow-md p-6">
      <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-6 gap-4">
        <h1 class="text-2xl font-bold">ToDo リスト</h1>
        <div class="flex flex-col sm:flex-row gap-2">
          <button
            @click="showInviteModal = true"
            class="bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded transition duration-200"
          >
            ユーザーを招待
          </button>
          <button
            @click="showNameModal = true"
            class="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded transition duration-200"
          >
            表示名を設定
          </button>
        </div>
      </div>
      
      <!-- 新規ToDo追加フォーム -->
      <div class="mb-8 p-4 bg-gray-50 rounded-lg">
        <h2 class="text-lg font-semibold mb-4">新しいToDoを追加</h2>
        <form @submit.prevent="addTodo" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <input
            v-model="newTodo.title"
            type="text"
            placeholder="タイトル"
            required
            class="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <select
            v-model="newTodo.priority"
            class="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="high">高</option>
            <option value="medium">中</option>
            <option value="low">低</option>
          </select>
          <input
            v-model="newTodo.dueDate"
            type="date"
            class="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            type="submit"
            class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded transition duration-200"
          >
            追加
          </button>
        </form>
      </div>

      <!-- 進行中ToDo一覧 -->
      <div class="mb-8">
        <h2 class="text-lg font-semibold mb-4">進行中のToDo</h2>
        <!-- デスクトップ表示 -->
        <div class="hidden md:block overflow-x-auto">
          <table class="min-w-full bg-white border border-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-2 text-left">タイトル</th>
                <th class="px-4 py-2 text-left">優先度</th>
                <th class="px-4 py-2 text-left">期限</th>
                <th v-for="user in users" :key="user.id" class="px-4 py-2 text-center">
                  {{ user.displayName || user.id.slice(0, 8) }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="todo in activeTodos" :key="todo.id" class="border-t">
                <td class="px-4 py-2">{{ todo.title }}</td>
                <td class="px-4 py-2">
                  <span :class="getPriorityClass(todo.priority)">
                    {{ getPriorityText(todo.priority) }}
                  </span>
                </td>
                <td class="px-4 py-2">{{ formatDate(todo.dueDate) }}</td>
                <td v-for="user in users" :key="user.id" class="px-4 py-2 text-center">
                  <input
                    type="checkbox"
                    :checked="getUserStatus(todo, user.id)"
                    :disabled="user.id !== userId"
                    @change="updateTodoStatus(todo.id, user.id, ($event.target as HTMLInputElement).checked)"
                    class="w-4 h-4"
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <!-- モバイル表示 -->
        <div class="md:hidden space-y-4">
          <div v-for="todo in activeTodos" :key="todo.id" class="bg-white border border-gray-200 rounded-lg p-4">
            <div class="flex justify-between items-start mb-3">
              <h3 class="font-medium text-lg">{{ todo.title }}</h3>
              <span :class="getPriorityClass(todo.priority) + ' text-sm px-2 py-1 rounded'">
                {{ getPriorityText(todo.priority) }}
              </span>
            </div>
            <div v-if="todo.dueDate" class="text-sm text-gray-600 mb-3">
              期限: {{ formatDate(todo.dueDate) }}
            </div>
            <div class="space-y-2">
              <div v-for="user in users" :key="user.id" class="flex justify-between items-center">
                <span class="text-sm">{{ user.displayName || user.id.slice(0, 8) }}</span>
                <input
                  type="checkbox"
                  :checked="getUserStatus(todo, user.id)"
                  :disabled="user.id !== userId"
                  @change="updateTodoStatus(todo.id, user.id, ($event.target as HTMLInputElement).checked)"
                  class="w-4 h-4"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 完了済みToDo一覧 -->
      <div class="mb-8">
        <h2 class="text-lg font-semibold mb-4">完了済みのToDo</h2>
        <div class="overflow-x-auto">
          <table class="min-w-full bg-white border border-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-2 text-left">タイトル</th>
                <th class="px-4 py-2 text-left">優先度</th>
                <th class="px-4 py-2 text-left">期限</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="todo in completedTodos" :key="todo.id" class="border-t">
                <td class="px-4 py-2">{{ todo.title }}</td>
                <td class="px-4 py-2">
                  <span :class="getPriorityClass(todo.priority)">
                    {{ getPriorityText(todo.priority) }}
                  </span>
                </td>
                <td class="px-4 py-2">{{ formatDate(todo.dueDate) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- メモ欄 -->
      <div>
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-lg font-semibold">共有メモ</h2>
          <span class="text-sm text-gray-500">
            📝 全ユーザーで共有されます
          </span>
        </div>
        <textarea
          v-model="memo"
          placeholder="全ユーザーで共有されるメモを入力してください..."
          class="w-full h-32 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 mb-3"
        ></textarea>
        <div class="flex justify-between items-center">
          <span class="text-xs text-gray-400">
            💡 他のユーザーの変更は30秒ごとに自動更新されます
          </span>
          <button
            @click="saveMemo"
            :disabled="memoSaving"
            class="bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 text-white font-bold py-2 px-4 rounded transition duration-200"
          >
            {{ memoSaving ? '保存中...' : 'メモを保存' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 招待モーダル -->
    <div v-if="showInviteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 px-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h3 class="text-lg font-semibold mb-4">ユーザーを招待</h3>
        <div v-if="inviteUrl" class="mb-4">
          <p class="text-sm text-gray-600 mb-2">以下のURLを共有してください：</p>
          <div class="flex">
            <input
              :value="inviteUrl"
              readonly
              class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md bg-gray-50"
            />
            <button
              @click="copyToClipboard(inviteUrl)"
              class="px-4 py-2 bg-blue-500 text-white rounded-r-md hover:bg-blue-600"
            >
              コピー
            </button>
          </div>
        </div>
        <div class="flex justify-end gap-2">
          <button
            @click="closeInviteModal"
            class="px-4 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
          >
            閉じる
          </button>
          <button
            v-if="!inviteUrl"
            @click="generateInviteUrl"
            :disabled="inviteLoading"
            class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:bg-gray-400"
          >
            {{ inviteLoading ? '生成中...' : 'URL生成' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 表示名設定モーダル -->
    <div v-if="showNameModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 px-4">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h3 class="text-lg font-semibold mb-4">表示名を設定</h3>
        <input
          v-model="newDisplayName"
          type="text"
          placeholder="表示名を入力"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 mb-4"
        />
        <div class="flex justify-end gap-2">
          <button
            @click="closeNameModal"
            class="px-4 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
          >
            キャンセル
          </button>
          <button
            @click="updateDisplayName"
            :disabled="nameLoading"
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
          >
            {{ nameLoading ? '更新中...' : '更新' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { 
  getListData, 
  createTodo, 
  updateTodoUserStatus, 
  updateListMemo, 
  inviteUser, 
  updateUserName 
} from '../api/api'
import type { User, Todo, TodoForm } from '../types'

// Props
interface Props {
  listId: string
  userId: string
}
const props = defineProps<Props>()

// Router
const router = useRouter()

// State
const users = ref<User[]>([])
const todos = ref<Todo[]>([])
const memo = ref<string>('')
const newTodo = ref<TodoForm>({
  title: '',
  priority: 'medium',
  dueDate: ''
})
const showInviteModal = ref<boolean>(false)
const showNameModal = ref<boolean>(false)
const inviteUrl = ref<string>('')
const inviteLoading = ref<boolean>(false)
const newDisplayName = ref<string>('')
const nameLoading = ref<boolean>(false)
const memoSaving = ref<boolean>(false)
let autoRefreshInterval: NodeJS.Timeout | null = null

// Computed
const activeTodos = computed(() => {
  return todos.value.filter(todo => !todo.isCompleted)
    .sort((a, b) => {
      // 優先度順（高→中→低）
      const priorityOrder: Record<string, number> = { high: 3, medium: 2, low: 1 }
      const priorityDiff = priorityOrder[b.priority] - priorityOrder[a.priority]
      if (priorityDiff !== 0) return priorityDiff
      
      // 同じ優先度の場合は期限順（早い→遅い）
      const dueDateA = a.dueDate ? new Date(a.dueDate) : null
      const dueDateB = b.dueDate ? new Date(b.dueDate) : null
      
      if (dueDateA && dueDateB) {
        return dueDateA.getTime() - dueDateB.getTime()
      }
      if (dueDateA) return -1  // 期限ありが優先
      if (dueDateB) return 1   // 期限ありが優先
      
      // 両方とも期限なしの場合は作成日時順
      const createdAtA = a.createdAt ? new Date(a.createdAt) : new Date()
      const createdAtB = b.createdAt ? new Date(b.createdAt) : new Date()
      return createdAtA.getTime() - createdAtB.getTime()
    })
})

const completedTodos = computed(() => {
  return todos.value.filter(todo => todo.isCompleted)
    .sort((a, b) => {
      const updatedAtA = a.updatedAt ? new Date(a.updatedAt) : new Date()
      const updatedAtB = b.updatedAt ? new Date(b.updatedAt) : new Date()
      return updatedAtB.getTime() - updatedAtA.getTime() // 完了日時の新しい順
    })
})

// Methods
const loadData = async (): Promise<void> => {
  try {
    const data = await getListData(props.listId, props.userId)
    users.value = data.users
    todos.value = data.todos
    memo.value = data.memo
  } catch (error) {
    console.error('Failed to load data:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    if (errorMessage.includes('404')) {
      alert('指定されたリストまたはユーザーが見つかりません')
      await router.push('/')
    } else {
      alert(`データの読み込みに失敗しました: ${errorMessage}`)
    }
  }
}

const addTodo = async (): Promise<void> => {
  try {
    const todoData = {
      title: newTodo.value.title,
      priority: newTodo.value.priority,
      dueDate: newTodo.value.dueDate || null
    }
    await createTodo(props.listId, todoData)
    newTodo.value = { title: '', priority: 'medium', dueDate: '' }
    await loadData()
  } catch (error) {
    console.error('Failed to add todo:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    alert(`ToDoの追加に失敗しました: ${errorMessage}`)
  }
}

const updateTodoStatus = async (todoId: number, userId: string, checked: boolean): Promise<void> => {
  try {
    await updateTodoUserStatus(todoId, userId, checked)
    await loadData()
  } catch (error) {
    console.error('Failed to update todo status:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    alert(`ステータスの更新に失敗しました: ${errorMessage}`)
  }
}

const saveMemo = async (): Promise<void> => {
  memoSaving.value = true
  try {
    await updateListMemo(props.listId, memo.value)
    alert('メモを保存しました')
    // 最新のデータを再読み込み
    await loadData()
  } catch (error) {
    console.error('Failed to save memo:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    alert(`メモの保存に失敗しました: ${errorMessage}`)
  } finally {
    memoSaving.value = false
  }
}

const getUserStatus = (todo: Todo, userId: string): boolean => {
  const status = todo.userStatuses?.find(s => s.userId === userId)
  return status?.isChecked || false
}

const getPriorityClass = (priority: string): string => {
  const classes: Record<string, string> = {
    high: 'text-red-600 font-semibold',
    medium: 'text-yellow-600 font-semibold',
    low: 'text-green-600 font-semibold'
  }
  return classes[priority] || ''
}

const getPriorityText = (priority: string): string => {
  const texts: Record<string, string> = { high: '高', medium: '中', low: '低' }
  return texts[priority] || priority
}

const formatDate = (dateString: string | null): string => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString('ja-JP')
}

const generateInviteUrl = async (): Promise<void> => {
  inviteLoading.value = true
  try {
    const response = await inviteUser(props.listId)
    inviteUrl.value = `${window.location.origin}${response.url}`
  } catch (error) {
    console.error('Failed to generate invite URL:', error)
    alert('招待URLの生成に失敗しました')
  } finally {
    inviteLoading.value = false
  }
}

const closeInviteModal = (): void => {
  showInviteModal.value = false
  inviteUrl.value = ''
}

const copyToClipboard = async (text: string): Promise<void> => {
  try {
    await navigator.clipboard.writeText(text)
    alert('URLをコピーしました')
  } catch (error) {
    console.error('Failed to copy to clipboard:', error)
    alert('クリップボードへのコピーに失敗しました')
  }
}

const closeNameModal = (): void => {
  showNameModal.value = false
  const currentUser = users.value.find(u => u.id === props.userId)
  newDisplayName.value = currentUser?.displayName || ''
}

const updateDisplayName = async (): Promise<void> => {
  if (!newDisplayName.value.trim()) {
    alert('表示名を入力してください')
    return
  }
  
  nameLoading.value = true
  try {
    await updateUserName(props.listId, props.userId, newDisplayName.value.trim())
    await loadData()
    showNameModal.value = false
    alert('表示名を更新しました')
  } catch (error) {
    console.error('Failed to update display name:', error)
    alert('表示名の更新に失敗しました')
  } finally {
    nameLoading.value = false
  }
}

// Lifecycle
onMounted(async () => {
  await loadData()
  // 現在のユーザーの表示名をセット
  const currentUser = users.value.find(u => u.id === props.userId)
  if (currentUser) {
    newDisplayName.value = currentUser.displayName || ''
  }
  
  // 30秒ごとにデータを自動更新（他のユーザーの変更を反映）
  autoRefreshInterval = setInterval(() => {
    loadData()
  }, 30000)
})

onBeforeUnmount(() => {
  // インターバルをクリアしてメモリリークを防ぐ
  if (autoRefreshInterval) {
    clearInterval(autoRefreshInterval)
  }
})
</script>