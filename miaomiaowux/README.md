# MEO 面板

这是项目的 Go 面板服务，内嵌已构建的 Web 前端，支持服务器、节点、套餐、用户、订阅、流量统计及通知管理。

本分支采用自托管模式：不包含支付、购买、价格展示、拼车市场或许可证销售入口；界面使用简约主题。完整部署、构建及升级说明见仓库根目录的 [README](../README.md) 和 [部署文档](../docs/DEPLOYMENT.md)。

## 本地构建

需要 Go 1.26 或更高版本：

```bash
go build -trimpath -o mmwx ./cmd/server
```

运行时默认将数据保存到当前目录的 `data/`，可通过环境变量覆盖：

```bash
MMWX_DATA_DIR=./data PORT=12889 ./mmwx
```

首次打开页面后按向导创建管理员账号。

## 主要环境变量

| 变量 | 说明 |
| --- | --- |
| `MMWX_DATA_DIR` | 持久化数据目录，默认 `./data` |
| `PORT` | HTTP 监听端口，默认 `12889` |
| `BIND_HOST` | 监听地址，默认 `0.0.0.0` |
| `TZ` | 时区，例如 `Asia/Shanghai` |
| `JWT_SECRET` | JWT 密钥，生产环境必须设置随机值 |
| `MMWX_UPDATE_REPO` | 面板 Release 仓库，格式 `owner/repo`；默认关闭 |
| `MMWX_AGENT_GITHUB_REPO` | Agent Release 仓库；默认关闭 |
| `MMWX_GH_TOKEN` | 私有仓库读取 Release 所需的 GitHub Token |

数据库默认使用 SQLite。PostgreSQL 配置见根目录部署文档。

## 目录说明

- `cmd/server/`：服务入口
- `internal/handler/`：HTTP 接口与业务逻辑
- `internal/storage/`：SQLite/PostgreSQL 数据层
- `internal/web/dist/`：内嵌前端产物
- `configs/`、`templates/`：内置配置与模板

第三方许可证与声明保留在本目录及 `THIRD_PARTY_NOTICES.md` 中。
