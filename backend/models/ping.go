package models

import "time"

type PingResult struct {
	ID           int       `json:"id,omitempty"`
	DockerHostID int       `json:"docker_host_id"`
	IPAddress    string    `json:"ip_address"`
	PingTime     int       `json:"ping_time"`
	LastSuccess  time.Time `json:"last_success"`
}
