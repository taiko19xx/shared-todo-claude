# 共有ToDoリストアプリケーション

## プロジェクト概要

複数人で利用できるToDoリスト管理アプリケーションを作成します。複雑な会員登録を不要とし、URLベースでのアクセスにより簡単に共有できる使い捨て型のToDoリストです。

## 主な特徴

- URLにリストIDとユーザーIDを含めた簡単アクセス
- 複数ユーザーによる協調的なToDo管理
- 全ユーザーが完了チェックして初めてToDoが完了となる仕組み
- リアルタイム共有メモ機能（全ユーザーで共有）
- ユーザー招待機能とカスタム表示名設定
- レスポンシブデザイン（モバイル・デスクトップ対応）
- 認証が弱い代わりにシンプルで迅速な利用が可能

## 技術スタック

### フロントエンド
- **フレームワーク**: Vue.js 3 (Composition API + TypeScript)
- **言語**: TypeScript
- **CSSフレームワーク**: Tailwind CSS
- **ビルドツール**: Vite
- **型チェック**: vue-tsc

### バックエンド
- **言語**: Go
- **フレームワーク**: Gin
- **データベース**: SQLite（glebarez/sqlite - 純粋Go実装）
- **ORM**: GORM

### 開発環境
- **コンテナ**: Docker + Docker Compose
- **開発サーバー**: ホットリロード対応

## 機能仕様

### URL構造
- **形式**: `/{listId}/{userId}`
- **ID生成**: UUIDv4を使用
- **アクセス制御**: listIdとuserIdの組み合わせが存在しない場合は404

### 画面構成

#### トップページ（`/`）
- 新しいリスト作成ボタンのみ表示
- ボタンクリックで新しいリスト・ユーザーを生成してリダイレクト

#### ToDoリスト画面（`/{listId}/{userId}`）
1. **新規ToDo追加フォーム**（上部）
   - タイトル入力（必須）
   - 優先度選択（高・中・低、デフォルト：中）
   - 期限設定（日付のみ、任意）

2. **進行中ToDo一覧**（テーブル形式）
   - 並び順: 優先度（高→低）→期限（近い→遠い）
   - カラム: タイトル | 優先度 | 期限 | ユーザー1 | ユーザー2 | ...
   - 各ユーザー列にはチェックボックス表示
   - 自分のチェックボックスのみ操作可能

3. **完了済みToDo一覧**（テーブル形式）
   - 全ユーザーがチェック完了したToDoを表示
   - 表示のみ（編集不可）

4. **共有メモ欄**（下部）
   - リスト共通のフリーテキストエリア
   - 全ユーザーが編集・閲覧可能
   - 保存ボタンによる明示的な保存
   - 30秒ごとの自動更新で他ユーザーの変更を反映

5. **ユーザー管理機能**
   - 表示名設定: ユーザーが自分の表示名をカスタマイズ可能
   - 招待機能: 新しいユーザーを招待してURLを生成・共有
   - 招待されたユーザーは初回アクセス時に表示名を入力可能

### データモデル

#### Lists テーブル
```sql
CREATE TABLE lists (
    id TEXT PRIMARY KEY,
    memo TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Users テーブル
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    list_id TEXT NOT NULL,
    display_name TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (list_id) REFERENCES lists(id)
);
```

#### Todos テーブル
```sql
CREATE TABLE todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    list_id TEXT NOT NULL,
    title TEXT NOT NULL,
    priority TEXT CHECK(priority IN ('high', 'medium', 'low')) DEFAULT 'medium',
    due_date DATE,
    is_completed BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (list_id) REFERENCES lists(id)
);
```

#### TodoUserStatus テーブル
```sql
CREATE TABLE todo_user_status (
    todo_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    is_checked BOOLEAN DEFAULT FALSE,
    checked_at DATETIME,
    PRIMARY KEY (todo_id, user_id),
    FOREIGN KEY (todo_id) REFERENCES todos(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### API仕様

#### ベースURL
- 開発環境: `http://localhost:8080/api`

#### エンドポイント

##### リスト関連
- `POST /lists` - 新しいリストとユーザーを作成
- `GET /lists/{listId}/users/{userId}` - リスト情報とユーザー情報を取得
- `PUT /lists/{listId}/memo` - メモを更新

##### ユーザー関連
- `POST /lists/{listId}/users` - 新しいユーザーを招待（URLを生成）
- `PUT /lists/{listId}/users/{userId}/name` - ユーザー表示名を設定

##### ToDo関連
- `GET /lists/{listId}/todos` - ToDo一覧を取得
- `POST /lists/{listId}/todos` - 新しいToDoを作成
- `PUT /todos/{todoId}` - ToDoを更新（タイトル、優先度、期限）
- `PUT /todos/{todoId}/status/{userId}` - ユーザーのチェック状態を更新
- `DELETE /todos/{todoId}` - ToDoを削除

### ビジネスロジック

#### ToDo完了判定
- 全ユーザーがチェック完了した場合、`todos.is_completed = true`に更新
- 一人でもチェックを外した場合、`todos.is_completed = false`に戻す

#### アクセス制御
- 存在しないlistId/userIdの組み合わせは404エラー
- 他ユーザーのチェックボックス操作は403エラー

#### データ整合性
- ユーザーがリストに追加された際、既存の全ToDoに対してtodo_user_statusレコードを作成
- ToDoが追加された際、そのリストの全ユーザーに対してtodo_user_statusレコードを作成

### UI/UX要件

#### レスポンシブデザイン
- モバイルフレンドリーなテーブルレイアウト
- スマートフォンでは横スクロール対応とカード型レイアウト
- 各画面要素の適切なブレークポイント設定
- モーダルのモバイル対応

#### リアルタイム性・データ同期
- リアルタイム通信は不使用（シンプル性重視）
- 30秒間隔での自動データ更新によるコラボレーション支援
- 手動リフレッシュボタンでの即座更新
- メモ保存時の自動データ再読み込み
- 更新競合は後勝ち方式

#### エラーハンドリング
- 適切なHTTPステータスコードの返却
- フロントエンドでのエラー表示
- バリデーションエラーの詳細表示

### 開発環境構成

#### ディレクトリ構造
```
project-root/
├── docker-compose.yml
├── frontend/
│   ├── Dockerfile
│   ├── package.json
│   ├── tsconfig.json
│   ├── env.d.ts
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── src/
│       ├── main.ts
│       ├── App.vue
│       ├── components/
│       ├── views/
│       ├── types/
│       │   └── index.ts
│       └── api/
│           └── api.ts
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── models/
│   ├── handlers/
│   ├── middleware/
│   └── database/
└── README.md
```

#### Docker Compose設定
- フロントエンド: ポート3000でViteデブサーバー
- バックエンド: ポート8080でGinサーバー
- データベース: SQLiteファイルをボリュームマウント
- ホットリロード: 両環境で開発時の自動リロード対応

### デプロイメント想定

#### コンテナ化
- 本番用のマルチステージDockerfile
- 静的ファイルの最適化
- ヘルスチェック機能

#### 環境変数
- `PORT`: サーバーポート（デフォルト: 8080）
- `DB_PATH`: SQLiteファイルパス
- `CORS_ORIGIN`: CORS許可オリジン

## 実装完了状況

### ✅ Phase 1: 基本機能（完了）
1. ✅ Docker環境のセットアップ
2. ✅ バックエンド基盤（Gin + GORM + SQLite）
3. ✅ データベーススキーマとマイグレーション
4. ✅ フロントエンド基盤（Vue3 + Vite + Tailwind）

### ✅ Phase 2: コア機能（完了）
1. ✅ リスト作成とユーザー管理
2. ✅ ToDo CRUD操作
3. ✅ チェック状態管理とToDo完了判定
4. ✅ 基本的なUI実装

### ✅ Phase 3: 追加機能（完了）
1. ✅ 共有メモ機能（保存ボタン付き）
2. ✅ 招待機能とURL生成・コピー
3. ✅ 優先度・期限による並び替え
4. ✅ エラーハンドリングとバリデーション

### ✅ Phase 4: 品質向上（完了）
1. ✅ レスポンシブデザイン対応（モバイル・デスクトップ）
2. ✅ 操作性向上（ユーザー表示名設定、コピー機能等）
3. ✅ 30秒間隔での自動データ更新
4. ✅ 詳細なエラーメッセージ表示

### ✅ Phase 5: テスト実装（完了）
1. ✅ バックエンドユニットテスト（APIハンドラー、モデル、データベース）
2. ✅ フロントエンドコンポーネントテスト（Vue.js、ユーザーインタラクション）
3. ✅ テスト自動化環境（GitHub Actions、Docker統合）
4. ✅ カバレッジ測定とレポート生成

### 🎯 現在の状況
全ての主要機能が実装完了し、本格的な利用が可能な状態です。

## テスト

### テスト戦略
- **ユニットテスト**: 各コンポーネントと関数の単体テスト
- **統合テスト**: API・データベース連携テスト
- **フロントエンドテスト**: Vueコンポーネントとユーザーインタラクションテスト
- **カバレッジ測定**: コードカバレッジの可視化と品質管理

### バックエンドテスト
- **フレームワーク**: Go標準テストパッケージ + Testify
- **テスト対象**:
  - APIハンドラー（handlers_test.go）
  - データモデル（models_test.go）
  - データベース操作（database_test.go）
- **テストデータベース**: インメモリSQLiteを使用
- **カバレッジ**: HTMLレポート生成対応

### フロントエンドテスト
- **フレームワーク**: Vitest + Vue Test Utils + Testing Library
- **テスト対象**:
  - Vueコンポーネント（Home.vue, TodoList.vue, App.vue）
  - API関数（api.test.js）
  - ユーザーインタラクション
- **モック**: API呼び出し、ルーター、DOM操作
- **環境**: jsdom（ブラウザ環境をシミュレート）

### テスト実行方法

#### ローカル実行
```bash
# 全テスト実行
./test.sh

# バックエンドのみ
cd backend && go test -v ./...

# フロントエンドのみ
cd frontend && npm test

# カバレッジ付き実行
./coverage.sh
```

#### Docker環境でのテスト
```bash
# テスト用プロファイルで実行
docker compose --profile test up --abort-on-container-exit

# 個別実行
docker compose run backend-test
docker compose run frontend-test
```

#### CI/CD（GitHub Actions）
- プッシュ・プルリクエスト時に自動実行
- バックエンド・フロントエンド・統合テストを並列実行
- カバレッジレポートをアーティファクトとして保存

### テストカバレッジ
- **目標**: 80%以上のコードカバレッジ
- **レポート**: HTML形式でブラウザで確認可能
- **対象外**: テストファイル、node_modules

## 注意事項

- セキュリティは最小限（URLベースアクセスのみ）
- 使い捨て利用を想定した設計
- 管理機能は現段階では不要
- スケーラビリティよりもシンプルさを重視