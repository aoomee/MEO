# 转发管理功能文档

## 概述

转发管理是MEO的原生多跳转发功能，使用 mmw-agent 的 Go 原生 relay 引擎（不依赖 xray tunnel），支持转发组、转发链、负载均衡、健康检查和 DNS 自动编排。

### 核心概念

- **转发组（Forward Group）**：一组服务器的集合，支持多种负载均衡策略
  - `round_robin`：轮询
  - `percentage`：按权重百分比分配
  - `cycle`：循环

- **转发链（Forward Chain）**：由多个转发组按顺序组成的多跳链路
  - 第一跳 = 入口组
  - 最后一跳 = 出口组
  - 中间跳 = 中转组

- **节点绑定**：将一个订阅节点绑定到转发链的指定端口
  - 入口和中转组的服务器运行 agent relay（端口 P → 下一组成员端口 P）
  - 出口组的服务器运行真实 xray 入站（监听端口 P）

## 架构设计

### 数据流向

```
用户 → 入口组服务器:P (relay)
     → 中转组服务器:P (relay)
     → 出口组服务器:P (xray inbound)
     → 目标网站
```

### 关键特性

1. **跨链聚合下发**：每台服务器承载所有绑定链的转发规则（agent Apply 是全量状态）
2. **健康检查**：agent 侧实现 TCP 健康探测 + RTT 测量 + 字节计数
3. **故障切换**：支持 failover/recover 双窗口滞回，防抖
4. **DNS 编排**：入口组可绑定域名，根据成员健康态自动增删 A 记录
5. **出口自动建站**：bind 链时自动在出口组每台服务器上创建 VLESS+TCP 入站和订阅节点

## 后端实现

### 数据库表结构

#### forward_groups（转发组）
```sql
CREATE TABLE forward_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    balance_strategy TEXT DEFAULT 'round_robin',
    failover_enabled INTEGER DEFAULT 0,
    offline_ms_threshold INTEGER DEFAULT 5000,
    dns_domain TEXT DEFAULT '',
    dns_provider_id INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### forward_group_members（组成员）
```sql
CREATE TABLE forward_group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    server_id INTEGER NOT NULL,
    weight INTEGER DEFAULT 1,
    FOREIGN KEY (group_id) REFERENCES forward_groups(id) ON DELETE CASCADE,
    FOREIGN KEY (server_id) REFERENCES remote_servers(id) ON DELETE CASCADE
);
```

#### forward_chains（转发链）
```sql
CREATE TABLE forward_chains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### forward_chain_hops（链的跳数）
```sql
CREATE TABLE forward_chain_hops (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chain_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    hop_order INTEGER NOT NULL,
    FOREIGN KEY (chain_id) REFERENCES forward_chains(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES forward_groups(id) ON DELETE CASCADE
);
```

#### forward_chain_nodes（链绑定的节点）
```sql
CREATE TABLE forward_chain_nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER NOT NULL UNIQUE,
    chain_id INTEGER NOT NULL,
    port INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (chain_id) REFERENCES forward_chains(id) ON DELETE CASCADE
);
```

#### forward_hop_metrics（历史监控数据）
```sql
CREATE TABLE forward_hop_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER NOT NULL,
    rule_id TEXT NOT NULL,
    upstream_addr TEXT NOT NULL,
    healthy INTEGER NOT NULL DEFAULT 0,
    rtt_ms INTEGER NOT NULL DEFAULT 0,
    bytes_up INTEGER NOT NULL DEFAULT 0,
    bytes_down INTEGER NOT NULL DEFAULT 0,
    at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 核心代码模块

#### 1. 编排引擎（forward_manage.go）

**buildForwardRules**：纯函数，将一条链计算成 `serverID → []forwardRule` 映射
```go
func buildForwardRules(
    chain *storage.ForwardChain,
    groups map[int64]*storage.ForwardGroup,
    servers map[int64]*storage.RemoteServer,
    port int,
) (map[int64][]forwardRule, error)
```

**aggregateServerForwardRules**：计算某台服务器在所有绑定链下应承载的全部规则
```go
func aggregateServerForwardRules(
    serverID int64,
    topo forwardTopology,
) ([]forwardRule, error)
```

**SyncServersForwarding**：重算并下发给定服务器的全量转发规则（幂等）
```go
func (h *RemoteManageHandler) SyncServersForwarding(
    ctx context.Context,
    serverIDs []int64,
) error
```

**SyncChainForwarding**：一条链变化后，重新同步该链涉及的所有服务器
```go
func (h *RemoteManageHandler) SyncChainForwarding(
    ctx context.Context,
    chainID int64,
    extraServers ...int64,
) error
```

#### 2. 出口节点自动化（forward_manage.go）

**ProvisionExitNodes**：在出口组每台服务器上创建落地节点
```go
func (h *RemoteManageHandler) ProvisionExitNodes(
    ctx context.Context,
    chainID int64,
    port int,
) error
```

流程：
1. 获取出口组（链的最后一跳）
2. 对每台出口服务器：
   - 构造 VLESS+TCP 入站配置
   - 通过 `/api/child/inbounds` 下发到 agent
   - 调用 `inboundToClashProxy` 生成 clash_config
   - 落订阅节点行（`NodeType: "forward-exit"`，挂在 `"admin"` 用户下）
3. 失败时回滚已创建的入站和节点

**TeardownExitNodes**：清理出口组的自动建站
```go
func (h *RemoteManageHandler) TeardownExitNodes(
    ctx context.Context,
    chainID int64,
) error
```

按 `InboundTag` 前缀 `fwd-exit-c{chainID}-` 查找自动创建的节点，逐个删除入站和节点行。

#### 3. 管理 API（forward_admin.go）

REST API 端点：

**转发组**：
- `GET /api/admin/forward/groups` - 列出所有组
- `POST /api/admin/forward/groups` - 创建组
- `GET /api/admin/forward/groups/{id}` - 获取组详情
- `PUT /api/admin/forward/groups/{id}` - 更新组
- `DELETE /api/admin/forward/groups/{id}` - 删除组
- `PUT /api/admin/forward/groups/{id}/members` - 设置组成员

**转发链**：
- `GET /api/admin/forward/chains` - 列出所有链
- `POST /api/admin/forward/chains` - 创建链
- `GET /api/admin/forward/chains/{id}` - 获取链详情
- `DELETE /api/admin/forward/chains/{id}` - 删除链
- `PUT /api/admin/forward/chains/{id}/hops` - 设置链的跳数
- `POST /api/admin/forward/chains/{id}/bind` - 绑定节点到链
- `DELETE /api/admin/forward/chains/{id}/bind?node_id={id}` - 解绑节点

**监控**：
- `GET /api/admin/forward/status` - 拉取所有转发服务器的运行时状态
- `GET /api/admin/forward/metrics?server_id={id}&rule_id={id}&hours={n}` - 查询历史监控数据

#### 4. 状态轮询器（forward_poller.go）

每 17 秒执行一次：
1. 调用 `loadForwardTopology` 加载全部绑定及其引用的链/组/服务器
2. 调用 `forwardingServerIDs` 计算真正承载规则的服务器（排除纯出口）
3. 逐台调用 `FetchForwardStatus` 拉取 agent 上报的运行时状态
4. 批量写入 `forward_hop_metrics` 表
5. 对设了 `DNSDomain` 的入口组，调用 DNS provider 的 `ReconcileRecordSet` 更新 A 记录
6. 每日调用 `CleanOldForwardHopMetrics` 清理旧数据（保留 7 天）

#### 5. DNS 编排（ddns/provider.go + 各 provider 实现）

新增 `ReconcileRecordSet` 接口方法：
```go
ReconcileRecordSet(
    ctx context.Context,
    fqdn string,
    recordType string,
    desiredContents []string,
    ttl int,
) error
```

把 `fqdn` 的 `recordType` 记录集调整成恰好 `desiredContents`（增缺删多，幂等）。

已实现的 6 个 DNS provider：
- Cloudflare（最干净的实现）
- Alidns（阿里云 DNS）
- Tencent（腾讯云 DNS）
- DNSPod
- Namesilo（注意 min TTL 3600）
- GoDaddy（天然契合，一次 PUT 即可 reconcile）

### Agent 侧实现

Agent 的转发引擎在 `mmw-agent/internal/forward/` 中实现：

- **Rule**：转发规则配置
- **Upstream**：上游目标（带权重）
- **Health**：健康检查配置
- **Picker**：负载均衡器接口（round_robin / weighted）
- **Manager**：规则管理器，处理 `/api/child/forward/apply` 和 `/api/child/forward/status`

健康检查机制：
- 周期性 TCP 连接探测
- RTT 测量
- Failover/Recover 双窗口滞回（防抖）
- 字节计数（up/down）

## 前端实现

### 路由结构

- `/forward` - 转发管理根路由（`forward.tsx`，带 admin 鉴权）
- `/forward/` - 转发管理首页（`forward.index.tsx`）

### 主组件（forward-panel.tsx）

**功能 1：实时状态卡片**
- 每 5 秒自动刷新 `/api/admin/forward/status`
- 显示每台服务器的转发规则和 upstream 健康状态
- 健康 upstream 显示绿色徽章，不健康显示红色

**功能 2：历史延迟监控**
- 三个下拉选择器：服务器 / 规则 / 时间范围
- 使用 recharts LineChart 绘制 wide format 数据
- 每条 upstream 一条折线，颜色区分
- 不健康的点显示为 -1（formatter 转换为 "unhealthy"）
- 每 10 秒自动刷新

### i18n 翻译

中文（`zh-CN/forward.json`）：
```json
{
  "title": "转发管理",
  "status": {
    "title": "转发规则状态",
    "refresh": "刷新状态"
  },
  "metrics": {
    "title": "历史延迟监控",
    "selectServer": "选择服务器",
    "selectRule": "选择规则",
    "timeRange": "时间范围",
    "hours": "{{n}} 小时",
    "latency": "延迟 (ms)",
    "noData": "请先选择服务器和规则"
  }
}
```

英文（`en/forward.json`）类似结构。

## 使用流程

### 1. 创建转发组

```bash
curl -X POST http://localhost:12889/api/admin/forward/groups \
  -H "MM-Authorization: Bearer <token>" \
  -d '{
    "name": "入口组-香港",
    "balance_strategy": "round_robin",
    "failover_enabled": true,
    "offline_ms_threshold": 5000,
    "dns_domain": "hk-entry.example.com",
    "dns_provider_id": 1,
    "members": [
      {"server_id": 1, "weight": 1},
      {"server_id": 2, "weight": 1}
    ]
  }'
```

### 2. 创建转发链

```bash
curl -X POST http://localhost:12889/api/admin/forward/chains \
  -H "MM-Authorization: Bearer <token>" \
  -d '{
    "name": "香港-新加坡-美国",
    "group_ids": [1, 2, 3]
  }'
```

其中 `group_ids` 按顺序为：入口组、中转组（可多个）、出口组。

### 3. 绑定节点到链

```bash
curl -X POST http://localhost:12889/api/admin/forward/chains/1/bind \
  -H "MM-Authorization: Bearer <token>" \
  -d '{
    "node_id": 100,
    "port": 8443
  }'
```

此操作会：
1. 在数据库中创建绑定关系
2. 计算并下发转发规则到入口组和中转组的所有服务器
3. 在出口组的每台服务器上自动创建 VLESS+TCP 入站（端口 8443）
4. 生成订阅节点行（挂在 `admin` 用户下）

### 4. 查看运行状态

```bash
curl http://localhost:12889/api/admin/forward/status \
  -H "MM-Authorization: Bearer <token>"
```

返回示例：
```json
{
  "success": true,
  "servers": [
    {
      "server_id": 1,
      "name": "hk-entry-1",
      "rules": [
        {
          "rule_id": "fwd-c1-p8443-h0",
          "listen": ":8443",
          "upstreams": [
            {
              "addr": "sg-relay-1.example.com:8443",
              "healthy": true,
              "rtt_ms": 45,
              "bytes_up": 1048576,
              "bytes_down": 2097152
            }
          ]
        }
      ]
    }
  ]
}
```

### 5. 查询历史监控

```bash
curl "http://localhost:12889/api/admin/forward/metrics?server_id=1&rule_id=fwd-c1-p8443-h0&hours=24" \
  -H "MM-Authorization: Bearer <token>"
```

### 6. 解绑节点

```bash
curl -X DELETE "http://localhost:12889/api/admin/forward/chains/1/bind?node_id=100" \
  -H "MM-Authorization: Bearer <token>"
```

此操作会：
1. 删除数据库中的绑定关系
2. 重新同步规则（入口和中转组的服务器不再转发到此链）
3. 删除出口组服务器上的自动建站（入站 + 订阅节点）

## DNS 自动编排

### 配置入口组 DNS

在创建或更新转发组时设置：
```json
{
  "dns_domain": "hk-entry.example.com",
  "dns_provider_id": 1
}
```

其中 `dns_provider_id` 对应 `ddns_providers` 表中的 DNS provider 配置。

### 自动翻转机制

状态轮询器每 17 秒检查一次：

1. 获取入口组的所有成员服务器
2. 检查每台服务器的健康状态（基于其承载的转发规则的 upstream 健康态）
3. 计算 `desiredContents` = 所有健康成员的对外 IP（优先 IPv4）
4. 调用 DNS provider 的 `ReconcileRecordSet` 方法
5. Provider 对比当前 DNS 记录与 `desiredContents`：
   - 缺少的 IP → 添加 A 记录
   - 多余的 IP → 删除 A 记录
   - 已存在的 IP → 保持不变

### 滞回防抖

健康检查的滞回在 agent 侧实现（failover_ms / recover_ms 双窗口），DNS 层直接采用 agent 上报的 `Healthy` 状态，不再叠加额外防抖。

## 运维指南

### 监控指标

1. **转发规则数量**：`forward_chain_nodes` 表的行数
2. **活跃服务器数量**：承载转发规则的服务器（`forwardingServerIDs` 返回）
3. **历史监控数据量**：`forward_hop_metrics` 表的行数（每 7 天自动清理）
4. **DNS 同步频率**：每 17 秒（poller 间隔）

### 日志关键字

- `[ForwardProvision]` - 出口节点自动建站日志
- `[ForwardTeardown]` - 出口节点清理日志
- `[forward]` - 规则聚合和下发日志

### 故障排查

**问题 1：规则下发失败**
- 检查 agent 连接状态（`/api/admin/remote-servers`）
- 检查 agent 日志中的 `/api/child/forward/apply` 请求
- 确认服务器的 IP 地址或域名配置正确

**问题 2：出口节点创建失败**
- 检查日志中的 `[ForwardProvision]` 错误
- 确认出口组的所有服务器在线
- 检查端口是否被占用
- 验证 agent 版本支持 `/api/child/inbounds` 接口

**问题 3：DNS 自动编排不生效**
- 检查入口组的 `dns_domain` 和 `dns_provider_id` 配置
- 验证 DNS provider 的凭据（`ddns_providers` 表）
- 检查 DNS provider 的 API 限流
- 查看 poller 日志中的 DNS 调用错误

**问题 4：健康检查误报**
- 调整组的 `offline_ms_threshold`（默认 5000ms）
- 检查网络延迟是否超过阈值
- 查看 agent 日志中的健康检查结果

### 性能优化

1. **批量操作**：使用 `SyncServersForwarding` 批量同步多台服务器，而非逐台调用
2. **跨链聚合**：每台服务器只下发一次全量规则，避免多链互相覆盖
3. **监控数据清理**：定期清理 7 天前的 `forward_hop_metrics` 数据
4. **DNS 节流**：poller 间隔 17 秒（避开整点抖动），避免频繁调用 DNS API

## 限制和注意事项

1. **端口唯一性**：同一条链上的节点必须使用相同的端口 P
2. **最少两跳**：转发链至少需要入口组和出口组（中转组可选）
3. **出口组无规则**：出口组的服务器不承载转发规则，只运行 xray 入站
4. **DNS provider 支持**：目前支持 6 家，新增需实现 `ReconcileRecordSet` 接口
5. **Action Guard**：出口节点创建绕过 Action Guard（与 tunnel_chains/routed_outbound 一致）
6. **系统节点**：自动创建的出口节点挂在 `admin` 用户下，普通用户不可见

## 待办事项（TODO）

- [ ] 支持 UDP 转发（当前只支持 TCP）
- [ ] 支持更多负载均衡策略（consistent hashing、least connections）
- [ ] 支持转发链的流量统计（按链聚合）
- [ ] 支持转发链的访问控制（按用户/套餐限制）
- [ ] 支持转发链的配置模板（快速创建常用拓扑）
- [ ] 前端增加转发组/链的可视化拓扑图
- [ ] 前端增加转发链的创建向导
- [ ] 支持 DNS 记录的 TTL 自定义（当前使用默认值）

## 参考资料

- 计划文档：`/root/.claude/plans/optimized-popping-kurzweil.md`
- Agent 转发引擎：`mmw-agent/internal/forward/`
- DNS provider 接口：`internal/ddns/provider.go`
- 状态轮询器：`internal/handler/forward_poller.go`
- 前端组件：`miaomiaowux-frontend/src/components/forward-panel.tsx`
