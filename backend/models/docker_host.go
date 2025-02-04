package models

type DockerHost struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	IP   string `json:"ip_address"`
}