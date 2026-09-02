# MEO

MEO 是一套自托管的服务器、Xray 节点、套餐、用户、订阅与流量管理面板。仓库包含面板、远程 Agent、配套 Xray Core 分支，以及 Docker Compose、systemd 和多架构构建文件。

界面采用简约的低饱和视觉风格；订阅生成器的规则分类完整保留。面板不提供支付、购买、价格展示或第三方推广入口。

## 选择部署方式

| 方式 | 适合场景 | 数据库 | 说明 |
| --- | --- | --- | --- |
| Docker Compose | 快速部署、方便迁移 | SQLite | 推荐首次体验和单机使用 |
| Compose + PostgreSQL | 长期运行、多用户 | PostgreSQL 17 | 数据库独立持久化 |
| systemd 一键安装 | 需要直接管理宿主机 Nginx/Xray | SQLite 或外部 PostgreSQL | 支持 amd64/arm64 |

## Docker Compose（直接复制）

镜像托管在 GHCR。复制下面完整内容保存为 `docker-compose.yml`，也可以直接粘贴到 1Panel 的 Compose 编辑器：

```yaml
name: meo

services:
  panel:
    image: ghcr.io/aoomee/meo:latest
    container_name: meo
    pull_policy: always
    restart: unless-stopped
    ports:
      # 只让本机反向代理访问面板，公网只开放 80/443
      - "127.0.0.1:12889:12889"
    environment:
      PORT: "12889"
      BIND_HOST: "0.0.0.0"
      TZ: "Asia/Shanghai"
      MMWX_DATA_DIR: "/app/data"
      MMWX_DATABASE_DRIVER: "sqlite"
      MMWX_DATABASE_PATH: "/app/data/mmwx.db"
      MMWX_UPDATE_REPO: "off"
      MMWX_AGENT_GITHUB_REPO: "off"
      MMWX_SETUP_TOKEN: ""
      # 填写实际反向代理的 IP/CIDR；不要填 0.0.0.0/0
      MMWX_TRUSTED_PROXIES: ""
    volumes:
      - mmwx-data:/app/data
    security_opt:
      - no-new-privileges:true

volumes:
  mmwx-data:
```

在该文件所在目录执行：

```bash
docker compose up -d
docker compose ps
```

查看实时日志：

```bash
docker compose logs -f panel
```

默认监听本机 `12889` 端口，数据保存在 Docker volume `mmwx-data`。首次使用请通过以下任一安全地址进入初始化页：

- 本机：`http://127.0.0.1:12889/login`
- 生产环境：`https://你的域名/login`

不要直接通过公网 `http://IP:12889` 初始化。浏览器加密接口需要 HTTPS 或 localhost 安全上下文。

首次初始化还受启动令牌保护：设置 `MMWX_SETUP_TOKEN`（至少 24 个字符）可固定令牌；留空时面板会生成随机令牌并写入启动日志。远程初始化时访问日志中给出的 `/api/setup/authorize?token=...` 一次即可获得短期 HttpOnly 授权 Cookie。

如果使用 Nginx、Caddy 或 Cloudflare Tunnel 反代，必须把代理的实际 IP/CIDR 写入 `MMWX_TRUSTED_PROXIES`。使用 `network_mode: host` 且代理也在宿主机时通常是：

```dotenv
MMWX_TRUSTED_PROXIES=127.0.0.1
```

桥接网络下请填写代理连接在面板容器中显示的源地址，通常是 Docker 网桥网关（例如 `172.17.0.1`）；代理在另一个 Docker 网络时可填写该网络的实际 CIDR（例如 `172.18.0.0/16`）。如果 Cloudflare 直接回源到面板，则填写 Cloudflare 官方 IP 段，并在防火墙中只允许这些网段访问源站。不要留空后直接把 `12889` 暴露到公网，也不要使用 `0.0.0.0/0`。

反向代理需要传递原始协议：

```nginx
proxy_set_header Host $host;
proxy_set_header X-Forwarded-Proto $scheme;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
```

修改环境变量后重新创建容器：

```bash
docker compose up -d --force-recreate panel
```

如果你沿用旧项目的 `network_mode: host`，请确保宿主机已有 Nginx/Caddy/Cloudflare Tunnel，并将面板绑定到本机：

```yaml
services:
  panel:
    image: ghcr.io/aoomee/meo:latest
    container_name: meo
    network_mode: host
    restart: unless-stopped
    environment:
      PORT: "12889"
      BIND_HOST: "127.0.0.1"
      MMWX_DATA_DIR: "/app/data"
      MMWX_DATABASE_DRIVER: "sqlite"
      MMWX_DATABASE_PATH: "/app/data/mmwx.db"
      MMWX_TRUSTED_PROXIES: "127.0.0.1"
    volumes:
      - ./data:/app/data
      - ./subscribes:/app/subscribes
      - ./rule_templates:/app/rule_templates
```

`network_mode: host` 下不要再写 `ports`；如果没有本机反向代理，请使用上面的默认 Compose 配置。PostgreSQL 建议使用仓库提供的 `deploy/compose.postgres.yml` 覆盖文件，避免 host 网络下的服务名解析问题。

## Docker Compose：PostgreSQL

先在 `.env` 中设置强随机密码：

```dotenv
POSTGRES_PASSWORD=CHANGE_ME_TO_A_STRONG_RANDOM_PASSWORD
```

然后启动面板与 PostgreSQL：

```bash
docker compose \
  -f docker-compose.yml \
  -f deploy/compose.postgres.yml \
  up -d
```

PostgreSQL 数据保存在 `postgres-data` volume。不要直接把已有 SQLite 实例切换成 PostgreSQL；请先通过面板完成数据备份或迁移。

## HTTPS 反向代理

生产环境建议只开放 80/443，并将面板端口限制在本机。Compose 中可将端口映射改为：

```yaml
ports:
  - "127.0.0.1:12889:12889"
```

Caddy 示例：

```caddyfile
panel.example.com {
    reverse_proxy 127.0.0.1:12889
}
```

保存配置并让域名解析到服务器后，Caddy 会自动申请和续期 HTTPS 证书。

## systemd 一键安装

先构建二进制，或从私有 GitHub Release 下载对应架构的文件。使用本地二进制安装：

```bash
sudo ./deploy/install.sh --binary ./build/mmwx-linux-amd64
```

从私有 GitHub Release 安装：

```bash
export MMWX_UPDATE_REPO=OWNER/PRIVATE_REPOSITORY
export MMWX_GH_TOKEN=YOUR_FINE_GRAINED_TOKEN
sudo -E ./deploy/install.sh
```

私有仓库 Token 只需要目标仓库的 Contents/Metadata 读取权限。避免把 Token 直接写进命令行、镜像或 Git 历史。

安装后的路径：

| 内容 | 路径 |
| --- | --- |
| 面板程序 | `/usr/local/bin/mmwx` |
| 环境配置 | `/etc/mmwx/mmwx.env` |
| 持久化数据 | `/var/lib/mmwx/data` |
| systemd 服务 | `/etc/systemd/system/mmwx.service` |

常用命令：

```bash
systemctl status mmwx
journalctl -u mmwx -f
sudo systemctl restart mmwx
```

## 构建 Release 文件

需要 Go 1.27+：

```bash
./scripts/build-all.sh
```

产物写入 `build/`：

- `mmwx-linux-amd64`
- `mmwx-linux-arm64`
- `mmwx-agent-linux-amd64`
- `mmwx-agent-linux-arm64`
- `SHA256SUMS`

推送 `v*` 标签后，GitHub Actions 会构建两个 Linux 架构并创建 Release。工作流同时维护 `mmwx` 和 `mmwx-agent` 两个滚动 Release tag，供安装脚本和面板下载。

## Agent 分发

Agent 建议直接安装在受管服务器宿主机上，因为它需要管理 Xray、网络和系统服务。将对应架构的 Agent 放到面板数据目录：

```text
agent-bin/mmwx-agent-linux-amd64
agent-bin/mmwx-agent-linux-arm64
```

然后在“服务器管理 → 添加服务器”中复制安装命令。

## 备份与升级

SQLite Compose 备份：

```bash
docker compose stop panel
docker run --rm \
  -v meo_mmwx-data:/source \
  -v "$PWD/backup:/backup" \
  alpine tar -czf /backup/mmwx-data.tar.gz -C /source .
docker compose start panel
```

升级前先备份数据，再执行：

```bash
docker compose pull
docker compose up -d
```

更完整的 PostgreSQL、裸机安装、升级、回滚与排障说明见 [部署文档](docs/DEPLOYMENT.md)。

## 仓库结构

- `miaomiaowux/`：Go 面板服务与内嵌 Web 前端
- `mmw-agent/`：远程服务器 Agent
- `xray-core-vision-limiter/`：Agent 构建使用的 Xray Core 分支
- `deploy/`：Compose 覆盖、安装脚本与 systemd 服务
- `scripts/`：面板和 Agent 的多架构构建脚本

## 发布前安全检查

- 不要提交 `.env`、数据库、日志、`machine-id`、Token、私钥或运行时配置。
- 会话令牌由服务端使用加密安全随机数生成；生产环境必须通过 HTTPS 访问。
- 在线更新默认关闭；启用私有仓库更新时，请使用最小权限 Token。
- 升级前备份 SQLite/PostgreSQL 数据和 Agent 配置。

上游项目及第三方组件的许可证和声明保留在对应子目录中。部署或分发前，请确认你的使用方式符合所有适用许可条款。
