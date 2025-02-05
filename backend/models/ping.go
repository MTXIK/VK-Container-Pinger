package models

import "time"

type PingResult struct {
	ID            int       `json:"id,omitempty"`
	IPAddress     string    `json:"ip_address"`
	ContainerName string    `json:"container_name"`
	PingTime      float64   `json:"ping_time"`
	LastSuccess   time.Time `json:"last_success"`
}
