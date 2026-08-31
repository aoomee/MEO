package storage

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"
)

// WireGuard 数据模型 —— 与官方 v0.5.2 的 wg_devices / wg_leases 表逐字段对齐
// (schema 由跑官方二进制读它自己的 SQLite 逆出来,见 00-项目总说明 v0.5.2 逆向记录)。
//
// wg_devices:一台服务器上的一个 WireGuard 入站(= 一个地址池)。
//   ipv4_cidr / ipv6_cidr + first_index..last_index 划定可分配的主机号区间;
//   server_*_key 是 WG 服务端密钥对;probe_*_key 是「健康探测」专用密钥对(发合法握手拿 RTT)。
// wg_leases:池里分配出去的一个 peer(= 一个用户/套餐分配)。host_index 是池内序号,
//   ipv4/ipv6 是据此算出的地址;released_at 非空=已回收(序号可复用)。

// WGDevice 一个 WireGuard 入站地址池。
type WGDevice struct {
	ID               int64     `json:"id"`
	ServerID         int64     `json:"server_id"`
	InboundTag       string    `json:"inbound_tag"`
	IPv4CIDR         string    `json:"ipv4_cidr"`
	IPv6CIDR         string    `json:"ipv6_cidr"`
	FirstIndex       int       `json:"first_index"`
	LastIndex        int       `json:"last_index"`
	ListenPort       int       `json:"listen_port"` // WG 监听的 UDP 端口(官方把它放 xray 配置里;我们多存一列便于热更新时重建入站)
	ServerPrivateKey string    `json:"server_private_key"`
	ServerPublicKey  string    `json:"server_public_key"`
	ProbePrivateKey  string    `json:"probe_private_key"`
	ProbePublicKey   string    `json:"probe_public_key"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// WGLease 池里分给某用户的一个 peer。
type WGLease struct {
	ID           int64        `json:"id"`
	DeviceID     int64        `json:"device_id"`
	HostIndex    int          `json:"host_index"`
	Username     string       `json:"username"`
	AssignmentID int64        `json:"assignment_id"`
	Email        string       `json:"email"`
	PrivateKey   string       `json:"private_key"`
	PublicKey    string       `json:"public_key"`
	IPv4         string       `json:"ipv4"`
	IPv6         string       `json:"ipv6"`
	ReleasedAt   sql.NullTime `json:"-"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

const wireguardTablesSchema = `
CREATE TABLE IF NOT EXISTS wg_devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER NOT NULL,
    inbound_tag TEXT NOT NULL,
    ipv4_cidr TEXT NOT NULL,
    ipv6_cidr TEXT NOT NULL,
    first_index INTEGER NOT NULL,
    last_index INTEGER NOT NULL,
    listen_port INTEGER NOT NULL DEFAULT 0,
    server_private_key TEXT NOT NULL,
    server_public_key TEXT NOT NULL,
    probe_private_key TEXT NOT NULL DEFAULT '',
    probe_public_key TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(server_id, inbound_tag),
    UNIQUE(server_id, ipv4_cidr)
);
CREATE TABLE IF NOT EXISTS wg_leases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id INTEGER NOT NULL,
    host_index INTEGER NOT NULL,
    username TEXT NOT NULL,
    assignment_id INTEGER NOT NULL DEFAULT 0,
    email TEXT NOT NULL,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    ipv4 TEXT NOT NULL,
    ipv6 TEXT NOT NULL,
    released_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(device_id, host_index),
    UNIQUE(device_id, email)
);
CREATE INDEX IF NOT EXISTS idx_wg_leases_device ON wg_leases(device_id);
CREATE INDEX IF NOT EXISTS idx_wg_leases_user ON wg_leases(username);
`

func (r *TrafficRepository) ensureWireGuardTables(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, wireguardTablesSchema); err != nil {
		return fmt.Errorf("migrate wireguard tables: %w", err)
	}
	return nil
}

// ---------------- wg_devices ----------------

// CreateWGDevice 落一个 WG 入站地址池。
func (r *TrafficRepository) CreateWGDevice(ctx context.Context, d *WGDevice) (int64, error) {
	if err := r.ensureWireGuardTables(ctx); err != nil {
		return 0, err
	}
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO wg_devices
		  (server_id, inbound_tag, ipv4_cidr, ipv6_cidr, first_index, last_index, listen_port,
		   server_private_key, server_public_key, probe_private_key, probe_public_key)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		d.ServerID, d.InboundTag, d.IPv4CIDR, d.IPv6CIDR, d.FirstIndex, d.LastIndex, d.ListenPort,
		d.ServerPrivateKey, d.ServerPublicKey, d.ProbePrivateKey, d.ProbePublicKey)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func scanWGDevice(s interface{ Scan(...any) error }) (*WGDevice, error) {
	var d WGDevice
	if err := s.Scan(&d.ID, &d.ServerID, &d.InboundTag, &d.IPv4CIDR, &d.IPv6CIDR,
		&d.FirstIndex, &d.LastIndex, &d.ListenPort, &d.ServerPrivateKey, &d.ServerPublicKey,
		&d.ProbePrivateKey, &d.ProbePublicKey, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return nil, err
	}
	return &d, nil
}

const wgDeviceCols = `id, server_id, inbound_tag, ipv4_cidr, ipv6_cidr, first_index, last_index, listen_port,
	server_private_key, server_public_key, probe_private_key, probe_public_key, created_at, updated_at`

// GetWGDevice 按 server+inbound_tag 取(每台服务器每个入站唯一)。
func (r *TrafficRepository) GetWGDevice(ctx context.Context, serverID int64, inboundTag string) (*WGDevice, error) {
	if err := r.ensureWireGuardTables(ctx); err != nil {
		return nil, err
	}
	row := r.db.QueryRowContext(ctx, `SELECT `+wgDeviceCols+` FROM wg_devices WHERE server_id=? AND inbound_tag=?`, serverID, inboundTag)
	d, err := scanWGDevice(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return d, err
}

// GetWGDeviceByID 按主键取。
func (r *TrafficRepository) GetWGDeviceByID(ctx context.Context, id int64) (*WGDevice, error) {
	if err := r.ensureWireGuardTables(ctx); err != nil {
		return nil, err
	}
	row := r.db.QueryRowContext(ctx, `SELECT `+wgDeviceCols+` FROM wg_devices WHERE id=?`, id)
	d, err := scanWGDevice(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return d, err
}

// ListWGDevicesByServer 某台服务器上的全部 WG 入站。
func (r *TrafficRepository) ListWGDevicesByServer(ctx context.Context, serverID int64) ([]WGDevice, error) {
	if err := r.ensureWireGuardTables(ctx); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT `+wgDeviceCols+` FROM wg_devices WHERE server_id=? ORDER BY id`, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []WGDevice
	for rows.Next() {
		d, err := scanWGDevice(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *d)
	}
	return out, rows.Err()
}

// DeleteWGDevice 删入站池(连同其 leases,由调用方先回收/下发)。
func (r *TrafficRepository) DeleteWGDevice(ctx context.Context, id int64) error {
	if err := r.ensureWireGuardTables(ctx); err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, `DELETE FROM wg_leases WHERE device_id=?`, id); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM wg_devices WHERE id=?`, id)
	return err
}

// ---------------- wg_leases ----------------

// ListActiveLeasesByDevice 某入站池当前有效(未回收)的 peer,按序号排序 —— 下发 xray 用。
func (r *TrafficRepository) ListActiveLeasesByDevice(ctx context.Context, deviceID int64) ([]WGLease, error) {
	return r.listLeases(ctx, `WHERE device_id=? AND released_at IS NULL ORDER BY host_index`, deviceID)
}

// ListActiveLeasesByUser 某用户名下所有有效 peer —— 生成订阅用。
func (r *TrafficRepository) ListActiveLeasesByUser(ctx context.Context, username string) ([]WGLease, error) {
	return r.listLeases(ctx, `WHERE username=? AND released_at IS NULL ORDER BY device_id, host_index`, username)
}

const wgLeaseCols = `id, device_id, host_index, username, assignment_id, email,
	private_key, public_key, ipv4, ipv6, released_at, created_at, updated_at`

func (r *TrafficRepository) listLeases(ctx context.Context, where string, args ...any) ([]WGLease, error) {
	if err := r.ensureWireGuardTables(ctx); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT `+wgLeaseCols+` FROM wg_leases `+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []WGLease
	for rows.Next() {
		var l WGLease
		if err := rows.Scan(&l.ID, &l.DeviceID, &l.HostIndex, &l.Username, &l.AssignmentID, &l.Email,
			&l.PrivateKey, &l.PublicKey, &l.IPv4, &l.IPv6, &l.ReleasedAt, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// AllocateWGLease 在池 [first_index, last_index] 里挑一个空闲序号,建 peer。
// 复用被回收(released_at 非空)的序号:先把回收的记录物理删掉再插新的,保持 UNIQUE(device,host_index)。
// 找不到空位返回 (nil, ErrWGPoolExhausted)。IP 由序号 + CIDR 算出。
func (r *TrafficRepository) AllocateWGLease(ctx context.Context, l *WGLease) (int64, error) {
	if err := r.ensureWireGuardTables(ctx); err != nil {
		return 0, err
	}
	dev, err := r.GetWGDeviceByID(ctx, l.DeviceID)
	if err != nil {
		return 0, err
	}
	if dev == nil {
		return 0, fmt.Errorf("wg: device %d 不存在", l.DeviceID)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 占用中的序号
	used := map[int]bool{}
	rows, err := tx.QueryContext(ctx, `SELECT host_index FROM wg_leases WHERE device_id=? AND released_at IS NULL`, l.DeviceID)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var idx int
		if err := rows.Scan(&idx); err != nil {
			rows.Close()
			return 0, err
		}
		used[idx] = true
	}
	rows.Close()

	idx := -1
	// first_index 是服务端占的池首址,客户端从 +1 起分,避免和入站 address 撞号。
	start := dev.FirstIndex + 1
	if start < 1 {
		start = 1
	}
	for i := start; i <= dev.LastIndex; i++ {
		if !used[i] {
			idx = i
			break
		}
	}
	if idx < 0 {
		return 0, ErrWGPoolExhausted
	}
	ipv4, err := hostIP(dev.IPv4CIDR, idx)
	if err != nil {
		return 0, fmt.Errorf("wg: 算 ipv4 失败: %w", err)
	}
	ipv6 := ""
	if dev.IPv6CIDR != "" {
		if ipv6, err = hostIP(dev.IPv6CIDR, idx); err != nil {
			return 0, fmt.Errorf("wg: 算 ipv6 失败: %w", err)
		}
	}
	// 该序号若有历史回收记录,物理删掉腾位(UNIQUE(device,host_index))
	if _, err := tx.ExecContext(ctx, `DELETE FROM wg_leases WHERE device_id=? AND host_index=?`, l.DeviceID, idx); err != nil {
		return 0, err
	}
	// UNIQUE(device_id, email) 含已回收行,同 email 再分配前先清掉旧行。
	if _, err := tx.ExecContext(ctx, `DELETE FROM wg_leases WHERE device_id=? AND email=? AND released_at IS NOT NULL`, l.DeviceID, l.Email); err != nil {
		return 0, err
	}
	res, err := tx.ExecContext(ctx, `
		INSERT INTO wg_leases
		  (device_id, host_index, username, assignment_id, email, private_key, public_key, ipv4, ipv6)
		VALUES (?,?,?,?,?,?,?,?,?)`,
		l.DeviceID, idx, l.Username, l.AssignmentID, l.Email, l.PrivateKey, l.PublicKey, ipv4, ipv6)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	l.ID, l.HostIndex, l.IPv4, l.IPv6 = id, idx, ipv4, ipv6
	return id, nil
}

// GetActiveWGLease 取某入站池上某 email 的有效 lease;没有则 (nil, nil)。
func (r *TrafficRepository) GetActiveWGLease(ctx context.Context, deviceID int64, email string) (*WGLease, error) {
	leases, err := r.listLeases(ctx, `WHERE device_id=? AND email=? AND released_at IS NULL LIMIT 1`, deviceID, email)
	if err != nil {
		return nil, err
	}
	if len(leases) == 0 {
		return nil, nil
	}
	l := leases[0]
	return &l, nil
}

// ListActiveLeasesByAssignment 某套餐实例下所有有效 peer —— 解绑回收用。
func (r *TrafficRepository) ListActiveLeasesByAssignment(ctx context.Context, assignmentID int64) ([]WGLease, error) {
	if assignmentID <= 0 {
		return nil, nil
	}
	return r.listLeases(ctx, `WHERE assignment_id=? AND released_at IS NULL ORDER BY device_id, host_index`, assignmentID)
}

// UpdateWGDeviceListenPort 入站改端口时同步池上的监听口。
func (r *TrafficRepository) UpdateWGDeviceListenPort(ctx context.Context, id int64, port int) error {
	if err := r.ensureWireGuardTables(ctx); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE wg_devices SET listen_port=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, port, id)
	return err
}

// ReleaseWGLease 回收一个 peer(标记 released_at,序号后续可复用)。
func (r *TrafficRepository) ReleaseWGLease(ctx context.Context, deviceID int64, email string) error {
	if err := r.ensureWireGuardTables(ctx); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE wg_leases SET released_at=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP
		 WHERE device_id=? AND email=? AND released_at IS NULL`, deviceID, email)
	return err
}

// ErrWGPoolExhausted 地址池已满。
var ErrWGPoolExhausted = fmt.Errorf("wg: 地址池已耗尽")

// hostIP 把 CIDR 基址 + 主机号 idx 算成具体 IP(idx 从网络地址往后数)。
// 例: 10.6.0.0/16 + idx=2 → 10.6.0.2;fd00::/64 + idx=2 → fd00::2。
func hostIP(cidr string, idx int) (string, error) {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	base := ip.To4()
	if base == nil {
		base = ip.To16()
	}
	b := make(net.IP, len(base))
	copy(b, base)
	// 从最低字节往高位加 idx(带进位)
	carry := idx
	for i := len(b) - 1; i >= 0 && carry > 0; i-- {
		v := int(b[i]) + (carry & 0xff)
		b[i] = byte(v & 0xff)
		carry = (carry >> 8) + (v >> 8)
	}
	return b.String(), nil
}
