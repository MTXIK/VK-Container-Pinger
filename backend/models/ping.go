package models

import "time"

type PingResult struct {
	IPAddress   string    `json:"ip_address"`
	PingTime    int       `json:"ping_time"`
	LastSuccess time.Time `json:"last_success"`
}
