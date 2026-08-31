#!/bin/bash
# 面板打包脚本
#
# 前端产物已随源码提供在 internal/web/dist/，由 go:embed 打进二进制，
# 因此这里只编译 Go，不需要 node / npm，也不需要任何许可证公钥。
#
# 用法:
#   ./build.sh                 # 编译 linux/amd64 + linux/arm64 + windows/amd64
#   TARGETS="linux/amd64" ./build.sh
set -e

BUILD_DIR="build"
TARGETS="${TARGETS:-linux/amd64 linux/arm64 windows/amd64}"

cd "$(dirname "$0")"

if [ ! -d "internal/web/dist" ] || [ -z "$(ls -A internal/web/dist 2>/dev/null)" ]; then
    echo "❌ internal/web/dist 为空 —— 面板前端是 go:embed 进二进制的，缺了它编译出来没有界面。"
    exit 1
fi

VERSION=$(sed -n 's/.*Version = "\(.*\)".*/\1/p' internal/version/version.go | head -1)
echo "========================================"
echo " 构建MEO 面板 v${VERSION}"
echo "========================================"

rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

for target in $TARGETS; do
    goos="${target%%/*}"
    goarch="${target##*/}"
    out="${BUILD_DIR}/mmwx-${goos}-${goarch}"
    [ "$goos" = "windows" ] && out="${out}.exe"
    echo ""
    echo "→ ${goos}/${goarch}"
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
        go build -trimpath -ldflags="-s -w" -o "$out" ./cmd/server
    echo "  ✓ $out"
done

echo ""
echo "========================================"
echo "构建完成，产物在 ${BUILD_DIR}/"
echo "========================================"
echo ""
echo "部署（Linux）:"
echo "  install -m 0755 ${BUILD_DIR}/mmwx-linux-amd64 /usr/local/bin/mmwx"
echo "  PORT=12889 MMWX_DATA_DIR=/etc/mmwx/data /usr/local/bin/mmwx"
echo ""
echo "提示：如果要让面板能给子服务器分发 Agent，"
echo "      把编译好的 mmw-agent-linux-{amd64,arm64} 放进 \$MMWX_DATA_DIR/agent-bin/。"
