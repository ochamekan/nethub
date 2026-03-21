package models

import "time"

type Device struct {
	ID        int64     `json:"id"`
	Hostname  string    `json:"hostname"`
	Location  string    `json:"location"`
	IP        string    `json:"ip"`
	IsActive  bool      `json:"is_active"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
}
