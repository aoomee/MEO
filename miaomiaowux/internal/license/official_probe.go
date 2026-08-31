package license

import (
	"context"
	"errors"
)

var ErrOfficialProbeUnavailable = errors.New("官方探测未配置，请使用自建测速端")

type OfficialProbeResult struct {
	Target string `json:"target"`
	OK     bool   `json:"ok"`
}

func (m *Manager) ProbeReachability(context.Context, []string) (map[string]bool, error) {
	return nil, ErrOfficialProbeUnavailable
}
