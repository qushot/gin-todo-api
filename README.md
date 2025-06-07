# gin-todo-api

[gin-gonic/gin](https://github.com/gin-gonic/gin)を使って適当に作った TODO API。

## 利用ツール

- Live Reload: https://github.com/air-verse/air
- Database Migration: https://github.com/pressly/goose

## やること

- [x] air の導入
- [x] .air.toml の設定
- [ ] データベースの設定
- [ ] SQL 文 をファイルに切り出す
- [ ] StackTrace のライブラリをいくつか試したい
- [ ] ビルドタグで on-memory と PostgreSQL を切り替えたい
- [ ] MCP サーバーで API を操作したい
  - API をラップした MCP は作ってみた。そのうち cmd/mcp 的な感じでディレクトリを掘っていい感じにしたい。。

## マジで全然関係ないメモ

- Toml のリポジトリ: https://github.com/toml-lang/toml
- PostgreSQL 用のサンプル DB: https://qiita.com/pea1205/items/979ec46bff5bd7c1134f#world

## エンドポイントメモ

```sh
curl localhost:8080/api/v1/todos -H "traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
curl localhost:8080/api/v1/todos/{id} -H "traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
curl localhost:8080/api/v1/todos -X POST --json '{"title": "title3", "content": "content3", "done": false}' -H "traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
curl localhost:8080/api/v1/todos/{id} -X PUT --json '{"title": "updated title3", "content": "updated content3", "done": true}' -H "traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
curl localhost:8080/api/v1/todos/{id} -X DELETE -H "traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
```

## opanapi-generator メモ

参考

- [OpenAPI Generator に適した OpenAPI の書き方](https://techblog.zozo.com/entry/how-to-write-openapi-for-openapi-generator)
  - 「tags、operationId を 1 エンドポイントにつき 1 つ設定する」
- [【アジャイル系男子】Go Gin Server + OpenAPI Generator 爆速サイクル戦線で生き抜く ⚔⚔](https://tech-blog.optim.co.jp/entry/2020/10/20/110000)
- https://github.com/OpenAPITools/openapi-generator/blob/master/docs/generators/go-gin-server.md
- https://openapi-generator.tech/docs/configuration/
