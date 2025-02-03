package models

import (
	"time"
)

type ProxyMode string

const (
	ProxyModeReverse  ProxyMode = "reverse"
	ProxyModeRedirect ProxyMode = "redirect"
	ProxyModePath     ProxyMode = "path"
)

type Proxy struct {
	ID        string          `json:"id" db:"id"`
	ListenURL string          `json:"listen_url" db:"listen_url"`
	Mode      ProxyMode       `json:"mode" db:"mode"`
	PathKey   *string         `json:"path_key,omitempty" db:"path_key"`
	Targets   []Target        `json:"targets" db:"targets"`
	Condition *RouteCondition `json:"condition,omitempty" db:"condition"`
	Tags      []string        `json:"tags" db:"tags"`
	IsActive  bool            `json:"is_active" db:"is_active"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}
