package handlers

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
	Repo        *repository.PingRepository
	RedisClient *redis.Client
}

func NewHandler(repo *repository.PingRepository, redisClient *redis.Client) *Handler {
	return &Handler{
		Repo:        repo,
		RedisClient: redisClient,
	}
}

func (h *Handler) GetPings(c *gin.Context) {
	cacheKey := "pings_cache"
	cached, err := cache.GetCache(h.RedisClient, cacheKey)
	if err == nil {
		var results []models.PingResult
		if err := json.Unmarshal([]byte(cached), &results); err == nil {
			c.JSON(http.StatusOK, results)
			return
		}
	}

	limit := 100
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	results, err := h.Repo.GetPingResults(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка запроса к базе данных"})
		return
	}

	data, _ := json.Marshal(results)
	cache.SetCache(h.RedisClient, cacheKey, data, 10*time.Second)
	c.JSON(http.StatusOK, results)
}

func (h *Handler) PostPing(c *gin.Context) {
	var pr models.PingResult
	if err := c.ShouldBindJSON(&pr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	if err := h.Repo.InsertPingResult(pr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка вставки в базу данных"})
		return
	}

	cache.DeleteCache(h.RedisClient, "pings_cache")
	c.Status(http.StatusCreated)
}

func (h *Handler) DeleteOldPings(c *gin.Context) {
	beforeStr := c.Query("before")
	if beforeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр 'before' обязателен и должен быть в формате RFC3339"})
		return
	}

	before, err := time.Parse(time.RFC3339, beforeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметра 'before', ожидается RFC3339"})
		return
	}

	if err := h.Repo.DeleteOldPingResults(before); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления устаревших записей"})
		return
	}

	cache.DeleteCache(h.RedisClient, "pings_cache")
	c.JSON(http.StatusOK, gin.H{"message": "Старые записи успешно удалены"})
}