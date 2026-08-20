# 年度公益补助结算规则平台

这是一个离线可运行的补助规则、费用申报、试算、复核和年度台账服务。后端按 domain/application/repository/service/transport 分层，默认使用内存仓储，生产部署可切换 PostgreSQL 适配器。

## 本地运行

```bash
go mod tidy
go test ./...
go test -race ./...
go vet ./...
go run ./cmd/server
```

HTTP API 位于 `/api/v1`，健康检查为 `/healthz` 与 `/readyz`。前端目录 `web/` 是 Vue 3 + TypeScript + Pinia 的离线工作台原型。

## 业务规则

规则按项目、年度和方案版本生效；试算只读，确认结算会在同一应用事务中追加结算记录并更新年度累计。已发布的规则版本不可修改，撤销和更正会产生新的审计事件。金额使用整数分保存，避免浮点误差。

## 配置与数据

复制 `.env.example` 可配置 HTTP 地址、默认年度和附件根目录。迁移位于 `migrations/`，演示种子位于 `scripts/seed.go`，不会覆盖已有数据。
