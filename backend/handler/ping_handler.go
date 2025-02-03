package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/VK-Container-Pinger/backend/cache"
	"github.com/VK-Container-Pinger/backend/models"
	"github.com/VK-Container-Pinger/backend/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Repo *repository.PingRepository
	RedisClient *redis.Client
}

func NewHandler(repo *repository.PingRepository, redisClient *redis.Client) *Handler {
	return &Handler{
		Repo: repo,
		RedisClient: redisClient,
	}
}

func (h *Handler) GetPings(c *gin.Context){
	
}

func (h *Handler) PostPing(c *gin.Context) {
	
}

func (h *Handler) DeleteOldPings(c *gin.Context) {
	
}