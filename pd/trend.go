package pd

import (
	"time"

	"github.com/pingcap/pd/pkg/typeutil"
)

// Copy from github.com/pingcap/pd/server/api/trend.go
// Trend describes the cluster's schedule trend.
type Trend struct {
	Stores  []TrendStore  `json:"stores"`
	History *trendHistory `json:"history"`
}

type TrendStore struct {
	ID              uint64             `json:"id"`
	Address         string             `json:"address"`
	StateName       string             `json:"state_name"`
	Capacity        uint64             `json:"capacity"`
	Available       uint64             `json:"available"`
	RegionCount     int                `json:"region_count"`
	LeaderCount     int                `json:"leader_count"`
	StartTS         *time.Time         `json:"start_ts,omitempty"`
	LastHeartbeatTS *time.Time         `json:"last_heartbeat_ts,omitempty"`
	Uptime          *typeutil.Duration `json:"uptime,omitempty"`

	HotWriteFlow        uint64   `json:"hot_write_flow"`
	HotWriteRegionFlows []uint64 `json:"hot_write_region_flows"`
	HotReadFlow         uint64   `json:"hot_read_flow"`
	HotReadRegionFlows  []uint64 `json:"hot_read_region_flows"`
}

type trendHistory struct {
	StartTime int64               `json:"start"`
	EndTime   int64               `json:"end"`
	Entries   []trendHistoryEntry `json:"entries"`
}

type trendHistoryEntry struct {
	From  uint64 `json:"from"`
	To    uint64 `json:"to"`
	Kind  string `json:"kind"`
	Count int    `json:"count"`
}
