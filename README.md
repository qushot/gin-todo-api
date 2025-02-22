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

## opanapi-generator メモ

参考

- [OpenAPI Generator に適した OpenAPI の書き方](https://techblog.zozo.com/entry/how-to-write-openapi-for-openapi-generator)
  - 「tags、operationId を 1 エンドポイントにつき 1 つ設定する」
- [【アジャイル系男子】Go Gin Server + OpenAPI Generator 爆速サイクル戦線で生き抜く ⚔⚔](https://tech-blog.optim.co.jp/entry/2020/10/20/110000)
- https://github.com/OpenAPITools/openapi-generator/blob/master/docs/generators/go-gin-server.md
- https://openapi-generator.tech/docs/configuration/
