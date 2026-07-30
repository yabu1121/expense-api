# Expense API

REST API for managing expenses built with Go standard libraries and SQLite.

---

# Context

## Why this project?

このプロジェクトは、GoのWebフレームワークやORMを利用する前に、

- HTTPサーバがどのように動作するのか
- SQLをどのように実行するのか
- REST APIがどのように構成されるのか

を理解することを目的として開発した。

そのため、Gin・Echo・GORMなどは利用せず、

- net/http
- database/sql
- SQLite

のみを使用して実装している。

---

## Design Goals

設計時には以下を意識した。

### 1. Responsibility Separation

HTTP処理とDB処理を分離するため、

- Handler
- Store
- Model

の3層構成を採用した。

HandlerはHTTP通信のみを担当し、
SQLの知識を持たない。

StoreはSQLのみを担当し、
HTTPの知識を持たない。

Modelはデータ構造のみを保持する。

この分離により、
各レイヤーを独立して変更しやすい構成となっている。

---

### 2. Standard Library Only

フレームワークに依存しないよう、
Go標準ライブラリのみを利用した。

HTTPサーバは

```go
net/http
```

DBアクセスは

```go
database/sql
```

JSON処理は

```go
encoding/json
```

を利用している。

---

### 3. SQL First

ORMを利用せず、
SQLを直接記述している。

そのため、

- SELECT
- INSERT
- UPDATE
- DELETE

それぞれについて、
database/sqlの

- Query
- QueryRow
- Exec

を使い分けている。

---

## Error Handling

HTTPステータスコードを適切に返すことを意識した。

| Error | Status |
|-------|--------|
| Invalid JSON | 400 |
| Invalid ID | 400 |
| Expense Not Found | 404 |
| Internal Error | 500 |

存在しないレコードを更新・削除した場合は、

RowsAffected()

を利用し、

```go
sql.ErrNoRows
```

を返すことで
Handler側で404を返却している。

---

## Request Flow

```
Client
   │
   ▼
HTTP Request
   │
   ▼
Handler
   │
JSON Decode
   │
   ▼
Store
   │
SQL
   │
   ▼
SQLite
   │
Rows
   ▼
Store
   │
Model
   ▼
Handler
   │
JSON Encode
   ▼
Client
```

---

## API Design

REST APIとして

GET
POST
PUT
DELETE

を実装した。

各エンドポイントでは

- JSON Request
- JSON Response
- HTTP Status

を統一している。

---

## What I Learned

このプロジェクトでは

- HTTPサーバの構築
- REST API設計
- SQL実装
- database/sql
- JSONエンコード・デコード
- Layered Architecture
- Error Handling

について理解を深めた。

---

## Future Work

今後は

- Docker
- PostgreSQL
- Unit Test
- GitHub Actions
- Kubernetes
- Terraform

へ発展させる予定である。
