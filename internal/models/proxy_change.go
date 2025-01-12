package models

import (
	"encoding/json"
	"time"
)

type ChangeType string

const (
	ChangeTypeTargetsUpdate   ChangeType = "targets_update"
	ChangeTypeConditionUpdate ChangeType = "condition_update"
)

type ProxyChange struct {
	ID            string          `json:"id"`
	ProxyID       string          `json:"proxy_id"`
	ChangeType    ChangeType      `json:"change_type"`
	PreviousState json.RawMessage `json:"previous_state"`
	NewState      json.RawMessage `json:"new_state"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     *string         `json:"created_by,omitempty"`
}
