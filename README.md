# gin-todo-api

[gin-gonic/gin](https://github.com/gin-gonic/gin)を使って適当に作った TODO API。

## 利用ツール

- Live Reload: https://github.com/air-verse/air
- Database Migration: https://github.com/pressly/goose

## やること

- [x] air の導入
- [ ] .air.toml の設定
- [ ] データベースの設定

## マジで全然関係ないメモ

- Toml のリポジトリ: https://github.com/toml-lang/toml

## エンドポイントメモ

```sh
curl localhost:8080/api/v1/todos
curl localhost:8080/api/v1/todos/{id}
curl localhost:8080/api/v1/todos -X POST -H "Content-Type: application/json" -d '{"title": "title3", "content": "content3", "done": false}'
curl localhost:8080/api/v1/todos/{id} -X PUT -H "Content-Type: application/json" -d '{"title": "updated title3", "content": "updated content3", "done": true}'
curl localhost:8080/api/v1/todos/{id} -X DELETE
```
