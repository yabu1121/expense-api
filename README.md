# Expense API

Goの標準ライブラリとSQLiteを使った、支出管理用のREST APIです。

WebフレームワークやORMを使う前に、HTTP、JSON、SQL、テスト、Dockerの仕組みを理解することを目的とした学習用プロジェクトです。

## 技術スタック

- Go
- `net/http`
- `encoding/json`
- `database/sql`
- SQLite（`modernc.org/sqlite`ドライバ）
- Docker / Docker Compose

## 設計

HTTP処理とDB処理を分離するため、Handler / Store / Modelの3層構成を採用しています。

```text
Client
  ↓ HTTP / JSON
Handler
  ↓ Model / Filter
Store
  ↓ SQL
SQLite
```

- Handler：ルーティング、HTTP入力、JSON、ステータスコード
- Store：`database/sql`とSQLiteによるデータ操作
- Model：Expense、集計結果、validation、検索条件

Handler側に小さなStore interfaceを置き、テストではFake Storeを注入します。

## ディレクトリ構成

```text
.
├── cmd/api/main.go                  # 依存関係の組み立てとHTTPサーバー
├── internal/handler/                # HTTP・JSON・ルーティング
├── internal/model/                  # Model・validation・検索条件
├── internal/store/                  # SQLiteとSQL
├── Dockerfile                       # マルチステージビルド
├── docker-compose.yml               # APIと永続volume
└── .env.sample                      # 環境変数の例
```

`cmd/api/main.go`でSQLite Storeを生成してHandlerへ注入し、Go 1.22以降のメソッド付きパターンで各ルートを`http.ServeMux`へ登録します。

## API

| Method | Path | Description |
|---|---|---|
| GET | `/health` | ヘルスチェック |
| GET | `/version` | バージョン情報 |
| GET | `/expenses` | 支出一覧 |
| GET | `/expenses/{id}` | 支出の1件取得 |
| POST | `/expenses` | 支出の作成 |
| PUT | `/expenses/{id}` | 支出の更新 |
| DELETE | `/expenses/{id}` | 支出の削除 |
| GET | `/expenses/summary` | 件数と合計金額の取得 |

### 支出一覧の検索条件

`GET /expenses`では、次のクエリパラメータを組み合わせられます。

| Parameter | Example | Description |
|---|---|---|
| `category` | `food` | カテゴリの完全一致 |
| `limit` | `10` | 最大取得件数（1以上） |
| `offset` | `20` | 先頭から飛ばす件数（0以上、limit必須） |

```bash
curl "http://localhost:8080/expenses?category=food&limit=10&offset=0"
```

不正な`limit`・`offset`には`400 Bad Request`を返します。

### Expenseの例

```json
{
  "id": 1,
  "title": "coffee",
  "amount": 500,
  "category": "food"
}
```

### 入力ルール

POST・PUTでは、JSONをDecodeした後にtitleとcategoryの前後空白を除去してからvalidationを行います。

| Field | Rule |
|---|---|
| `title` | 空文字・空白のみは不可 |
| `amount` | 1以上の整数 |
| `category` | 空文字・空白のみは不可 |

### CRUDの使用例

作成：

```bash
curl -i -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -d '{"title":"coffee","amount":500,"category":"food"}'
```

一覧・1件取得：

```bash
curl -i http://localhost:8080/expenses
curl -i http://localhost:8080/expenses/1
```

更新：

```bash
curl -i -X PUT http://localhost:8080/expenses/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"latte","amount":550,"category":"food"}'
```

削除：

```bash
curl -i -X DELETE http://localhost:8080/expenses/1
```

作成は`201 Created`、取得・更新は`200 OK`、削除は`204 No Content`を返します。

### Summaryの例

```json
{
  "count": 2,
  "total_amount": 1050
}
```

集計SQLでは`COUNT`、`SUM`、`COALESCE`を使用し、Expenseが0件の場合も合計金額を`0`として返します。

## データベース

アプリケーション起動時に`CREATE TABLE IF NOT EXISTS`を実行し、テーブルがなければ自動作成します。現時点では独立したmigrationツールは使用していません。

```sql
CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    amount INTEGER,
    category TEXT
);
```

一覧取得では`ExpenseFilter`を使い、指定された条件だけをSQLへ追加します。

```text
Categoryあり → WHERE category = ?
Limitあり    → LIMIT ?
Offsetあり   → OFFSET ?
```

値は文字列連結せずプレースホルダーの引数として渡します。結果は常にID昇順です。

## エラーハンドリング

| Error | Status |
|---|---|
| Invalid JSON | 400 |
| Invalid ID | 400 |
| Invalid filter | 400 |
| Expense Not Found | 404 |
| Method Not Allowed | 405 |
| Internal Error | 500 |

存在しないExpenseの取得では`sql.ErrNoRows`をアプリケーションのNot Foundエラーへ変換します。更新・削除では`RowsAffected()`が0の場合にNot Foundとして扱います。

## ローカル起動

必要なもの：

- Go（`go.mod`に記載されたバージョン）
- curlなどのHTTPクライアント

最短手順：

```bash
go mod download
go test ./...
go run ./cmd/api
```

SQLiteのパスは`DB_PATH`環境変数で変更できます。未指定の場合は`expenses.db`を使用します。

```bash
DB_PATH=./expenses-dev.db go run ./cmd/api
```

## テスト

```bash
go vet ./...
go test ./...
```

- `httptest`とFake StoreによるHandlerテスト
- 一時SQLite DBを使ったStoreテスト
- HandlerとSQLite Storeを接続したIntegrationテスト
- `t.TempDir()`によるテストケースごとのDB分離

## Docker Compose

Dockerfileは、Goイメージでバイナリをビルドし、実行用の`debian:stable-slim`へバイナリだけをコピーするマルチステージ構成です。

起動：

```bash
docker compose up --build -d
```

ログ：

```bash
docker compose logs -f api
```

終了：

```bash
docker compose down
```

SQLiteデータは`expense-data`という名前付きvolumeへ保存されます。通常の`docker compose down`ではvolumeは削除されません。

動作確認：

```bash
curl -i http://localhost:8080/health
curl -i http://localhost:8080/expenses
```

## 学習した内容

- `http.ServeMux`のメソッド付きルーティング
- JSONのEncode / Decode
- NormalizeとValidate
- `Query` / `QueryRow` / `Exec`
- `LastInsertId` / `RowsAffected` / `sql.ErrNoRows`
- SQL集計と動的な検索条件
- interfaceと依存性注入
- Unit Test / Integration Test
- Dockerのマルチステージビルド
- Docker volumeによるSQLite永続化

## 現在の開発状況

- Expense CRUD：実装済み
- カテゴリ絞り込み：実装済み
- limit・offset：実装済み
- Summary API：実装済み
- Handler / Store / Integrationテスト：実装済み
- Docker ComposeとSQLite永続化：実装済み
- Graceful shutdown：未完成

Graceful shutdownは今後の課題です。現在はHTTPサーバーをgoroutineで起動してmain goroutineを待機させていますが、SIGINT・SIGTERMの受信、`server.Shutdown`、終了タイムアウトはまだ実装していません。

## Future Work

- Graceful shutdown
- PostgreSQL対応
- GitHub Actions
- Kubernetes
- Terraform
