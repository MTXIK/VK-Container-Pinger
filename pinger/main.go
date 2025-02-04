package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/go-ping/ping"
)

type DockerHost struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip_address"`
}

type PingResult struct {
	DockerHostID int       `json:"docker_host_id"`
	IPAddress    string    `json:"ip_address"`
	PingTime     int       `json:"ping_time"`
	LastSuccess  time.Time `json:"last_success"`
}

func getDockerHosts(backendURL string) ([]DockerHost, error){
	resp, err := http.Get(backendURL + "/api/docker-hosts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус при получении docker-хостов: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var hosts []DockerHost
	if err := json.Unmarshal(body, &hosts); err != nil {
		return nil, err
	}
	
	return hosts, nil
}