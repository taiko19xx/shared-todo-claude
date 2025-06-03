import '@testing-library/jest-dom'
import { expect, afterEach } from 'vitest'
import { cleanup } from '@testing-library/vue'

// テスト実行後にコンポーネントをクリーンアップ
afterEach(() => {
  cleanup()
})