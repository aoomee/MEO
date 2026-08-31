# 部署说明

## 选择方式

| 方式 | 适合场景 | 数据库 | 宿主机服务管理 |
| --- | --- | --- | --- |
| Docker Compose | 快速部署、隔离运行 | SQLite | 受容器边界限制 |
| Compose + PostgreSQL | 多用户或长期运行 | PostgreSQL 17 | 受容器边界限制 |
| systemd 一键安装 | 需要面板管理本机 Nginx/Xray | SQLite 或外部 PostgreSQL | 完整 |

Agent 建议安装在受管服务器宿主机上，不建议放入容器，因为它需要管理 Xray、网络和系统服务。

## Docker Compose（SQLite）

```bash
cp .env.example .env
openssl rand -hex 32
# 将输出写入 .env 的 MMWX_JWT_SECRET
docker compose up -d --build
docker compose logs -f panel
```

默认端口是 `12889`，数据保存在 Docker volume `mmwx-data`。修改 `.env` 的 `MMWX_PORT` 可以改变宿主机端口。

备份：

```bash
docker compose stop panel
docker run --rm -v mmwx_mmwx-data:/source -v "$PWD/backup:/backup" alpine \
  tar -czf /backup/mmwx-data.tar.gz -C /source .
docker compose start panel
```

## Docker Compose（PostgreSQL）

先在 `.env` 设置强密码 `POSTGRES_PASSWORD`，再运行：

```bash
docker compose -f docker-compose.yml -f deploy/compose.postgres.yml up -d --build
```

PostgreSQL 数据位于 `postgres-data` volume。不要在已运行的 SQLite 实例上直接切换变量；应先通过面板备份/迁移功能完成数据库迁移。

## systemd 裸机部署

支持 Linux amd64/arm64。使用本地二进制：

```bash
sudo ./deploy/install.sh --binary ./build/mmwx-linux-amd64
```

服务位置：

- 程序：`/usr/local/bin/mmwx`
- 环境变量：`/etc/mmwx/mmwx.env`
- 数据：`/var/lib/mmwx/data`
- 服务：`/etc/systemd/system/mmwx.service`

常用命令：

```bash
systemctl status mmwx
journalctl -u mmwx -f
sudo systemctl restart mmwx
```

更新本地二进制：

```bash
sudo ./deploy/install.sh update --binary ./build/mmwx-linux-amd64
```

卸载程序但保留数据：

```bash
sudo ./deploy/install.sh uninstall
```

## 私有 GitHub Release

推荐在私有仓库中维护两个滚动 tag：

- `mmwx`：面板资产 `mmwx-linux-amd64`、`mmwx-linux-arm64`
- `mmwx-agent`：Agent 资产 `mmwx-agent-linux-amd64`、`mmwx-agent-linux-arm64`

服务器的 `/etc/mmwx/mmwx.env` 可配置：

```dotenv
MMWX_UPDATE_REPO=OWNER/PRIVATE_REPOSITORY
MMWX_AGENT_GITHUB_REPO=OWNER/PRIVATE_REPOSITORY
MMWX_GH_TOKEN=YOUR_FINE_GRAINED_TOKEN
```

不需要在线更新时保持 `off`，以减少外部访问和密钥暴露面。

## HTTPS 反向代理

以 Caddy 为例：

```caddyfile
panel.example.com {
    reverse_proxy 127.0.0.1:12889
}
```

生产环境只开放 80/443，将 12889 限制为本机或可信网络访问。若容器映射仅供本机反代，可把 Compose 端口改为：

```yaml
ports:
  - "127.0.0.1:12889:12889"
```

## Agent 分发

将对应架构的 Agent 放到面板数据目录：

```text
agent-bin/mmwx-agent-linux-amd64
agent-bin/mmwx-agent-linux-arm64
```

容器部署时可将宿主机目录额外挂载到 `/app/data/agent-bin`。裸机部署对应 `/var/lib/mmwx/data/agent-bin`。然后在面板的“服务器管理 → 添加服务器”复制安装命令。

## 升级与回滚

- Compose：先备份 volume，再执行 `docker compose build --pull && docker compose up -d`。
- systemd：`deploy/install.sh update` 在启动失败时会恢复旧二进制。
- 数据库：升级前通过面板备份，或停服务后复制 SQLite 数据目录。

检查项：

```bash
curl -fsS http://127.0.0.1:12889/api/setup/status
systemctl is-active mmwx
docker compose ps
```
