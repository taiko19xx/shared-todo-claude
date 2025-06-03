# 共有ToDoリストアプリケーション

複数人で利用できるToDoリスト管理アプリケーションです。複雑な会員登録を不要とし、URLベースでのアクセスにより簡単に共有できる使い捨て型のToDoリストを提供します。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Vue.js](https://img.shields.io/badge/Vue.js-3.x-green.svg)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)

## ✨ 主な特徴

- 🔗 **URLベースアクセス** - リストIDとユーザーIDを含めた簡単アクセス
- 👥 **協調的ToDo管理** - 複数ユーザーによるリアルタイム共有
- ✅ **全員完了システム** - 全ユーザーがチェックして初めてToDoが完了
- 📝 **共有メモ機能** - リスト参加者全員で編集可能なメモ
- 🎫 **ユーザー招待機能** - 専用URLでの簡単招待
- 📱 **レスポンシブデザイン** - モバイル・デスクトップ対応
- 🚀 **シンプル設計** - 認証が弱い代わりに迅速で簡単な利用

## 🛠 技術スタック

### フロントエンド
- **フレームワーク**: Vue.js 3 (Composition API)
- **言語**: TypeScript
- **CSSフレームワーク**: Tailwind CSS
- **ビルドツール**: Vite
- **型チェック**: vue-tsc
- **テスト**: Vitest + Vue Test Utils

### バックエンド
- **言語**: Go
- **フレームワーク**: Gin
- **データベース**: SQLite (glebarez/sqlite - 純粋Go実装)
- **ORM**: GORM
- **テスト**: Go標準testing + Testify

### 開発環境
- **コンテナ**: Docker + Docker Compose
- **開発サーバー**: ホットリロード対応

## 🚀 クイックスタート

### 前提条件

- Docker
- Docker Compose

### 起動方法

1. **リポジトリをクローン**
   ```bash
   git clone <repository-url>
   cd shared-todo
   ```

2. **Docker Composeで起動**
   ```bash
   docker-compose up -d
   ```

3. **アプリケーションにアクセス**
   - フロントエンド: http://localhost:3000
   - バックエンドAPI: http://localhost:8080

### 開発環境での起動

#### バックエンド
```bash
cd backend
go mod download
go run main.go
```

#### フロントエンド
```bash
cd frontend
npm install
npm run dev
```

## 📖 使い方

### 基本的な流れ

1. **新しいリスト作成**
   - トップページで「新しいリストを作成」をクリック
   - 自動的に新しいリストとユーザーIDが生成されます

2. **ToDoの追加**
   - タイトル、優先度、期限を設定してToDoを追加
   - 優先度: 高・中・低
   - 期限: 任意設定

3. **ユーザーの招待**
   - 「ユーザーを招待」ボタンで招待URL生成
   - URLを共有して他のユーザーを招待

4. **協調作業**
   - 各ユーザーが個別にToDoをチェック
   - 全員がチェックすると完了状態に移行
   - 共有メモでコミュニケーション

### URL構造

- **フォーマット**: `/{listId}/{userId}`
- **例**: `/a1b2c3d4-e5f6-7890-abcd-ef1234567890/u9v8w7x6-y5z4-3210-9876-543210fedcba`

## 🧪 テスト

### 全テスト実行
```bash
./test.sh
```

### 個別テスト実行

#### バックエンドテスト
```bash
cd backend
go test -v ./...
```

#### フロントエンドテスト
```bash
cd frontend
npm test
```

#### TypeScript型チェック
```bash
cd frontend
npm run typecheck
```

## 📊 API仕様

### エンドポイント一覧

| メソッド | エンドポイント | 説明 |
|---------|---------------|------|
| `POST` | `/api/lists` | 新しいリストとユーザーを作成 |
| `GET` | `/api/lists/{listId}/users/{userId}` | リスト情報を取得 |
| `PUT` | `/api/lists/{listId}/memo` | メモを更新 |
| `POST` | `/api/lists/{listId}/users` | ユーザーを招待 |
| `PUT` | `/api/lists/{listId}/users/{userId}/name` | ユーザー表示名を設定 |
| `POST` | `/api/lists/{listId}/todos` | 新しいToDoを作成 |
| `PUT` | `/api/todos/{todoId}/status/{userId}` | ユーザーのチェック状態を更新 |

### データ構造

#### Todo
```typescript
interface Todo {
  id: number
  listId: string
  title: string
  priority: 'high' | 'medium' | 'low'
  dueDate: string | null
  isCompleted: boolean
  userStatuses?: TodoUserStatus[]
}
```

#### User
```typescript
interface User {
  id: string
  displayName: string
  listId?: string
}
```

## 🗄 データベーススキーマ

### テーブル構造

- **lists**: リスト情報とメモ
- **users**: ユーザー情報と表示名
- **todos**: ToDo項目
- **todo_user_status**: ユーザー別チェック状態

### 外部キー制約

- `users.list_id` → `lists.id`
- `todos.list_id` → `lists.id`
- `todo_user_status.todo_id` → `todos.id`
- `todo_user_status.user_id` → `users.id`

## 🔧 設定

### 環境変数

#### バックエンド
- `PORT`: サーバーポート（デフォルト: 8080）
- `DB_PATH`: SQLiteファイルパス
- `CORS_ORIGIN`: CORS許可オリジン

#### フロントエンド
- `VITE_API_BASE_URL`: APIベースURL（デフォルト: http://localhost:8080/api）

## 📂 プロジェクト構造

```
shared-todo/
├── docker-compose.yml          # Docker Compose設定
├── frontend/                   # Vue.js + TypeScript フロントエンド
│   ├── Dockerfile
│   ├── package.json
│   ├── tsconfig.json          # TypeScript設定
│   ├── env.d.ts               # 型定義
│   ├── vite.config.ts         # Vite設定
│   ├── tailwind.config.js     # Tailwind CSS設定
│   └── src/
│       ├── main.ts            # エントリーポイント
│       ├── App.vue
│       ├── types/             # 型定義
│       │   └── index.ts
│       ├── views/             # ページコンポーネント
│       │   ├── Home.vue
│       │   └── TodoList.vue
│       └── api/               # API層
│           └── api.ts
├── backend/                   # Go + Gin バックエンド
│   ├── Dockerfile
│   ├── go.mod
│   ├── main.go
│   ├── models/               # データモデル
│   ├── handlers/             # APIハンドラ
│   ├── middleware/           # ミドルウェア
│   └── database/             # データベース操作
├── test.sh                   # テスト実行スクリプト
├── coverage.sh               # カバレッジ測定スクリプト
└── README.md
```

## 🤝 コントリビューション

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### 開発ガイドライン

- **コードスタイル**: `go fmt`、ESLint + Prettier
- **テスト**: 新機能には必ずテストを追加
- **型安全性**: TypeScriptの型チェックを必ず通す
- **コミットメッセージ**: 分かりやすい日本語または英語

## 📝 ライセンス

このプロジェクトは [MIT License](LICENSE) の下で公開されています。

## 🙏 謝辞

- [Vue.js](https://vuejs.org/) - プログレッシブJavaScriptフレームワーク
- [Go](https://golang.org/) - シンプルで効率的なプログラミング言語
- [Tailwind CSS](https://tailwindcss.com/) - ユーティリティファーストCSSフレームワーク
- [GORM](https://gorm.io/) - GoのファンタスティックORMライブラリ

## 📞 サポート

問題や質問がある場合は、[Issues](https://github.com/your-username/shared-todo/issues) でお知らせください。

---

Made with ❤️ by [Your Name]