@echo off
REM 面板打包脚本（Windows）
REM 前端产物已随源码提供在 internal\web\dist，go:embed 打进二进制，无需 node/npm。
setlocal enabledelayedexpansion

set BUILD_DIR=build

echo ========================================
echo  构建MEO 面板
echo ========================================

if exist "%BUILD_DIR%" rmdir /s /q "%BUILD_DIR%"
mkdir "%BUILD_DIR%"

set CGO_ENABLED=0

echo.
echo [1/3] linux/amd64 ...
set GOOS=linux
set GOARCH=amd64
go build -trimpath -ldflags="-s -w" -o %BUILD_DIR%\mmwx-linux-amd64 .\cmd\server
if errorlevel 1 exit /b 1

echo [2/3] linux/arm64 ...
set GOARCH=arm64
go build -trimpath -ldflags="-s -w" -o %BUILD_DIR%\mmwx-linux-arm64 .\cmd\server
if errorlevel 1 exit /b 1

echo [3/3] windows/amd64 ...
set GOOS=windows
set GOARCH=amd64
go build -trimpath -ldflags="-s -w" -o %BUILD_DIR%\mmwx-windows-amd64.exe .\cmd\server
if errorlevel 1 exit /b 1

echo.
echo ========================================
echo 构建完成，产物在 %BUILD_DIR%\
echo ========================================
