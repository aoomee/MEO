# mmw-agent

部署在受管服务器上的远程代理程序，负责连接MEO主控、上报运行和流量信息，并接收经过鉴权的管理指令。

## 推荐安装方式

在面板中进入“服务器管理 → 添加服务器”，复制面板生成的一键安装命令到目标服务器执行。安装脚本会写入 `/etc/mmw-agent/config.yaml` 并配置系统服务。

Agent 二进制可以由主控的持久化目录分发：

```text
$MMWX_DATA_DIR/agent-bin/mmwx-agent-linux-amd64
$MMWX_DATA_DIR/agent-bin/mmwx-agent-linux-arm64
```

也可以设置 `MMWX_AGENT_GITHUB_REPO=owner/repo` 从私有仓库 Release 下载，并通过 `MMWX_GH_TOKEN` 提供访问凭据。

## 本地构建

Agent 的 `go.mod` 使用同级目录中的 `../xray-core-vision-limiter`：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o mmwx-agent-linux-amd64 ./cmd/mmw-agent
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o mmwx-agent-linux-arm64 ./cmd/mmw-agent
```

完整构建、发布和部署说明见仓库根目录的 [README](../README.md) 与 [部署文档](../docs/DEPLOYMENT.md)。
