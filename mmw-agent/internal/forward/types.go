package forward

// Rule 是主控下发的一条转发规则。字段必须与面板 forwardRule JSON 对齐。
type Rule struct {
	ID        string     `json:"id"`
	Listen    string     `json:"listen"`
	Protocol  string     `json:"protocol"`
	Upstreams []Upstream `json:"upstreams"`
	Strategy  string     `json:"strategy"`
	Health    Health     `json:"health"`
}

type Upstream struct {
	Addr   string `json:"addr"`
	Weight int    `json:"weight"`
}

type Health struct {
	Enabled        bool `json:"enabled"`
	IntervalMs     int  `json:"interval_ms"`
	TimeoutMs      int  `json:"timeout_ms"`
	FailoverMs     int  `json:"failover_ms"`
	RecoverMs      int  `json:"recover_ms"`
	RTTThresholdMs int  `json:"rtt_threshold_ms"`
}

// RuleStatus / UpstreamStatus 必须与面板 forwardRuleStatus JSON 对齐。
type RuleStatus struct {
	RuleID    string           `json:"rule_id"`
	Listen    string           `json:"listen"`
	Upstreams []UpstreamStatus `json:"upstreams"`
}

type UpstreamStatus struct {
	Addr      string `json:"addr"`
	Healthy   bool   `json:"healthy"`
	RTTMs     int64  `json:"rtt_ms"`
	BytesUp   uint64 `json:"bytes_up"`
	BytesDown uint64 `json:"bytes_down"`
}

func (h Health) normalized() Health {
	out := h
	if out.IntervalMs <= 0 {
		out.IntervalMs = 2000
	}
	if out.TimeoutMs <= 0 {
		out.TimeoutMs = 800
	}
	if out.FailoverMs <= 0 {
		out.FailoverMs = 5000
	}
	if out.RecoverMs <= 0 {
		out.RecoverMs = 3000
	}
	return out
}
