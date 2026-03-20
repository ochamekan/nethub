package models

import "time"

type Device struct {
	ID        int64
	Hostname  string
	IP        string
	IsActive  bool
	IsDeleted bool
	CreatedAt time.Time
}
