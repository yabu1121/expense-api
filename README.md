# Expense API

Go標準ライブラリとSQLiteを使った、個人向けツール群のバックエンドを目指す学習用REST APIです。

現在はExpense・Category・Tagを題材に、HTTP、JSON、SQL、テスト、Docker、レイヤ分離を学んでいます。将来的には機能単位で拡張できるモジュラーモノリスへ発展させ、Clean Architecture、gRPC、キャッシュ、非同期処理、認証、監視なども段階的に扱います。

## 技術スタック

- Go 1.27
- `net/http`
- `encoding/json`
- `database/sql`
- SQLite（`modernc.org/sqlite`）
- UUID v7
- Docker / Docker Compose

## 現在の設計

```text
Client
  ↓ HTTP / JSON
Handler
  ↓ Model / Request
Store interface
  ↓
SQLiteStore
  ↓ SQL
SQLite
```

- Handler：HTTP入力、JSON、ステータスコード
- Model：データ構造、Normalize、Validate、ドメインエラー
- Store：`database/sql`によるSQLite操作
- `cmd/api`：StoreとHandlerの組み立て、ルーティング、HTTPサーバー起動

Handlerは具体的なSQLiteStoreではなく、小さなStore interfaceへ依存します。本番では`SQLiteStore`、HandlerテストではFake Storeを注入します。構造体に追加のフィールドやメソッドがあっても、interfaceが要求する全メソッドのシグネチャを満たせば利用できます。

## ディレクトリ構成

```text
.
├── cmd/api/main.go          # 依存関係の組み立てとHTTPサーバー
├── internal/handler/        # HTTP HandlerとHandlerテスト
├── internal/model/          # Model、Request、validation、エラー
├── internal/store/          # SQLite実装とStoreテスト
├── sandbox/                 # APIを手動確認する簡易HTML
├── Dockerfile
└── docker-compose.yml
```

## API

| Method | Path | Description |
|---|---|---|
| GET | `/health` | ヘルスチェック |
| GET | `/version` | バージョン情報 |
| GET | `/expenses` | 支出一覧 |
| GET | `/expenses/{id}` | 支出の1件取得 |
| POST | `/expenses` | 支出作成 |
| PUT | `/expenses/{id}` | 支出更新 |
| DELETE | `/expenses/{id}` | 支出削除 |
| GET | `/expenses/summary` | 件数と合計金額 |
| GET | `/categories` | カテゴリ一覧 |
| GET | `/categories/{id}` | カテゴリの1件取得 |
| POST | `/categories` | カテゴリ作成 |
| GET | `/tags` | タグ一覧 |
| POST | `/tags` | タグ作成 |
| POST | `/expenses/{expenseID}/tags/{tagID}` | 支出とタグの関連付け |
| GET | `/sandbox/` | 簡易確認フォーム |

## Expense API

Expense、Category、TagのIDにはUUIDを使用します。Expenseの作成・更新では、先に作成したCategoryのUUIDが必要です。

リクエスト例：

```json
{
  "title": "coffee",
  "amount": 500,
  "category_id": "01900000-0000-7000-8000-000000000001"
}
```

レスポンス例：

```json
{
  "id": "01900000-0000-7000-8000-000000000002",
  "title": "coffee",
  "amount": 500,
  "category_id": "01900000-0000-7000-8000-000000000001",
  "created_at": "2026-09-09T12:00:00Z"
}
```

入力ルール：

| Field | Rule |
|---|---|
| `title` | 前後の空白を除去し、空文字は禁止 |
| `amount` | 1以上 |
| `category_id` | nil UUIDは禁止。存在確認はSQLiteの外部キー制約で行う |

一覧取得は現在、`created_at DESC`による全件取得です。以前のCategory絞り込み、limit、offsetはUUID移行を優先するため一時的に外しており、後から設計し直します。0件の場合は`null`ではなく`[]`を返します。

```bash
curl -i -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"food"}'

curl -i -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -d '{"title":"coffee","amount":500,"category_id":"<category UUID>"}'

curl -i http://localhost:8080/expenses
curl -i http://localhost:8080/expenses/<expense UUID>
```

作成は`201 Created`、取得・更新は`200 OK`、削除は`204 No Content`を返します。不正なUUIDは`400 Bad Request`、存在しないExpenseは`404 Not Found`です。

## データベース

起動時に`CREATE TABLE IF NOT EXISTS`を実行します。SQLite接続では外部キー検証を有効にしています。

```sql
CREATE TABLE expenses (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    amount INTEGER NOT NULL,
    category_id TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
```

ExpenseとTagは中間テーブルで多対多に関連付けます。

```sql
CREATE TABLE expense_tags (
    expense_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    PRIMARY KEY (expense_id, tag_id),
    FOREIGN KEY (expense_id) REFERENCES expenses(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
```

まだ独立したmigration機構はありません。既存の整数ID版DBは`CREATE TABLE IF NOT EXISTS`では更新されないため、開発中のUUID移行では新しいDBを使う必要があります。migration導入は次の設計課題です。

## テスト

通常のUnit・Storeテスト：

```bash
go vet ./...
go test ./...
```

Integrationテストはビルドタグで分離しています。

```bash
go test -tags=integration ./internal/handler
```

特定のテストだけ実行する例：

```bash
go test ./internal/handler -run '^TestListExpenses$' -count=1 -v
```

- `-run`：実行するテスト名を正規表現で指定（パッケージ全体のコンパイルは行う）
- `-count=1`：テストキャッシュを使わない
- `-v`：サブテストを含む実行結果を表示

Handlerテストは`httptest`とFake Store、StoreテストとIntegrationテストは`t.TempDir()`内のSQLite DBを使用します。

## ローカル起動

```bash
go mod download
go test ./...
go run ./cmd/api
```

SQLiteパスは`DB_PATH`で変更できます。未指定時は`expenses.db`です。

```bash
DB_PATH=./expenses-dev.db go run ./cmd/api
```

## Docker Compose

```bash
docker compose up --build -d
docker compose logs -f api
docker compose down
```

SQLiteは`expense-data` volumeへ保存されます。`docker compose down`だけではvolumeは削除されません。

## 学習した内容

- `http.ServeMux`のメソッド付きルーティングと`PathValue`
- JSON Encode / DecodeとRequest型
- NormalizeとValidate
- `Query` / `QueryRow` / `Exec`
- `RETURNING` / `RowsAffected` / `sql.ErrNoRows`
- UUID v7の生成、URL文字列の`uuid.Parse`、nil UUIDの検証
- Category外部キーとExpense・Tagの多対多中間テーブル
- Handler側interfaceによる依存性注入
- interfaceの暗黙的な実装とメソッドシグネチャ
- Fake Storeを使ったHandler Unit Test
- SQLite Store Testとビルドタグ付きIntegration Test
- `go test -run`が実行対象だけを絞り、パッケージ全体はコンパイルすること
- DockerマルチステージビルドとSQLite volume

## 現在の状況と次の課題

- Expense CRUD：UUID・`category_id`・`created_at`対応済み
- Category：UUIDで作成・一覧・1件取得
- Tag：UUIDで作成・一覧
- ExpenseとTagの関連付け：両方のUUIDに対応
- Summary：件数と合計金額を取得
- Category絞り込み・limit・offset：再設計のため一時停止
- Tag・ExpenseTagの自動テスト：未実装
- Graceful shutdown：未完成
- DB migration：未実装

将来は、ExpenseでRequest/Response DTO・UseCase・Repository境界を固めた後、CategoryとTagへ展開します。その上で個人向けのTask、Note、Habit、Budgetなどを機能単位で追加し、必要性を測定しながらキャッシュ、gRPC、非同期ジョブ、認証、監視へ進みます。
