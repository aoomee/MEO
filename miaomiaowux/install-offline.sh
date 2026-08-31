#!/bin/bash
# MEO 面板 —— 离线安装 / 升级脚本
#
# 完全不联网:只把「你自己编译好的二进制」装成 systemd 服务，适合本地调试 / 内网环境。
# （对外发的版本用仓库根目录的 install.sh，那个会从 GitHub Release 拉二进制。）
# 用这个脚本装完之后，面板里的「检查更新」照样能用 —— 自更新走的是 GitHub Release。
#
# 用法:
#   ./install-offline.sh                          # 用 ./build/mmwx-linux-amd64
#   ./install-offline.sh /path/to/mmwx-linux-arm64
#   PORT=8080 DATA_DIR=/etc/mmwx/data ./install-offline.sh
#
# 卸载:
#   systemctl disable --now mmwx && rm -f /usr/local/bin/mmwx /etc/systemd/system/mmwx.service
#   （数据目录不会被删,确认不要了再手动删）

set -euo pipefail

BINARY_SRC="${1:-}"
PORT="${PORT:-12889}"
DATA_DIR="${DATA_DIR:-/etc/mmwx/data}"
BIN_DST="/usr/local/bin/mmwx"
UNIT="/etc/systemd/system/mmwx.service"

if [ "${EUID:-$(id -u)}" -ne 0 ]; then
    echo "[ERROR] 请用 root 运行" >&2
    exit 1
fi

if [ -z "$BINARY_SRC" ]; then
    case "$(uname -m)" in
        x86_64|amd64)  BINARY_SRC="build/mmwx-linux-amd64" ;;
        aarch64|arm64) BINARY_SRC="build/mmwx-linux-arm64" ;;
        *) echo "[ERROR] 未知架构 $(uname -m),请显式传入二进制路径" >&2; exit 1 ;;
    esac
fi

if [ ! -f "$BINARY_SRC" ]; then
    echo "[ERROR] 找不到二进制: $BINARY_SRC" >&2
    echo "        先编译:  ./build.sh" >&2
    echo "        或指定:  ./install-offline.sh /path/to/mmwx-linux-amd64" >&2
    exit 1
fi

echo "[1/4] 停止旧服务(如果有)..."
systemctl stop mmwx 2>/dev/null || true

echo "[2/4] 安装二进制 -> $BIN_DST"
install -m 0755 "$BINARY_SRC" "$BIN_DST"

echo "[3/4] 准备数据目录 -> $DATA_DIR"
mkdir -p "$DATA_DIR" "$DATA_DIR/agent-bin"

if [ ! -f "$UNIT" ]; then
    echo "[4/4] 写入 systemd 单元 -> $UNIT"
    cat > "$UNIT" <<EOF
[Unit]
Description=MiaomiaowuX Panel (offline edition)
After=network.target

[Service]
Type=simple
ExecStart=${BIN_DST}
Environment="PORT=${PORT}"
Environment="MMWX_DATA_DIR=${DATA_DIR}"
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
else
    echo "[4/4] 已存在 $UNIT，保留现有配置"
    systemctl daemon-reload
fi

systemctl enable mmwx >/dev/null 2>&1 || true
systemctl start mmwx

sleep 2
if systemctl is-active --quiet mmwx; then
    echo ""
    echo "✓ 面板已启动，端口 ${PORT}，访问路径要带 /login"
    echo "  首次打开是「创建管理员账号」，第一个注册的就是管理员"
    echo "  （根路径 / 是探针伪装页，默认没开，只会显示「探针暂时无法访问」）"
    echo ""
    echo "  ‼️ 必须用 HTTPS 或 localhost 打开 —— 面板加密通道依赖浏览器 WebCrypto，"
    echo "     裸 http://IP:${PORT} 会「全新安装却显示登录页且登不进去」。"
    echo "     应急：ssh -L ${PORT}:127.0.0.1:${PORT} root@<本机IP>  然后开 http://127.0.0.1:${PORT}/login"
    echo ""
    echo "下一步(可选):把编译好的 Agent 放进去,面板才能给子服务器分发:"
    echo "    cp mmw-agent-linux-amd64 ${DATA_DIR}/agent-bin/"
    echo "    cp mmw-agent-linux-arm64 ${DATA_DIR}/agent-bin/"
else
    echo "[ERROR] 启动失败,看日志:journalctl -u mmwx -n 50 --no-pager" >&2
    exit 1
fi
