package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/VK-Container-Pinger/backend/cache"
	"github.com/VK-Container-Pinger/backend/config"
	"github.com/VK-Container-Pinger/backend/handlers"
	"github.com/VK-Container-Pinger/backend/kafka"
	"github.com/VK-Container-Pinger/backend/repository"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.KafkaBroker == "" {
		log.Fatal("Не установлены необходимые переменные окружения (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, KAFKA_BROKER)")
	}
	if cfg.RedisAddr == "" {
		cfg.RedisAddr = "redis:6379"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка пинга БД: %v", err)
	}
	log.Println("Подключение к PostgreSQL установлено.")

	pingRepo, err := repository.NewPingRepository(db)
	if err != nil {
		log.Fatalf("Ошибка создания репозитория пингов: %v", err)
	}

	redisClient := cache.NewRedisClient(cfg.RedisAddr)
	if _, err := redisClient.Ping(cache.Ctx).Result(); err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	log.Println("Подключение к Redis установлено.")

	go kafka.StartKafkaConsumer(cfg.KafkaBroker, "backend-group", []string{"ping-results"}, &kafka.Consumer{
		Repo:        pingRepo,
		RedisClient: redisClient,
	})

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			log.Println("Запуск удаления старых записей по IP...")
			if err := pingRepo.DeleteOldRecordsForIps(); err != nil {
				log.Printf("Ошибка удаления старых записей: %v", err)
			} else {
				log.Println("Старые записи успешно удалены.")
			}
		}
	}()

	router := gin.Default()
	handler := handlers.NewHandler(pingRepo, redisClient)

	router.GET("/api/pings", handler.GetPings)
	router.POST("/api/ping", handler.PostPing)
	router.DELETE("/api/pings/old", handler.DeleteOldPings)

	log.Printf("Запуск backend-сервиса на порту %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
