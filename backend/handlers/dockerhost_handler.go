package handlers

import (
	"net/http"
	"strconv"

	"github.com/VK-Container-Pinger/backend/models"
	"github.com/VK-Container-Pinger/backend/repository"
	"github.com/gin-gonic/gin"
)

type DockerHostHandler struct {
	Repo *repository.DockerHostRepository
}

func NewDockerHostHandler(repo *repository.DockerHostRepository) *DockerHostHandler {
	return &DockerHostHandler{Repo: repo}
}


func (h *DockerHostHandler) GetDockerHosts(c *gin.Context) {
	hosts, err := h.Repo.GetHosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения docker-хостов"})
		return
	}
	c.JSON(http.StatusOK, hosts)
}

func (h *DockerHostHandler) AddDockerHost(c *gin.Context) {
	var host models.DockerHost
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.Repo.InsertHost(host.Name, host.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления docker-хоста"})
		return
	}

	host.ID = id
	c.JSON(http.StatusCreated, host)
}


func (h *DockerHostHandler) DeleteDockerHost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор"})
		return
	}
	if err := h.Repo.DeleteHost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления docker-хоста"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Docker-хост удалён"})
}