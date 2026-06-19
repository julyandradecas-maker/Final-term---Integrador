package main

import (
	"log"

	"chat-system/backend/internal/config"
	"chat-system/backend/internal/handlers"
	"chat-system/backend/internal/hub"
	rdb "chat-system/backend/internal/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := rdb.Init(cfg); err != nil {
		log.Fatalf("[main] Error conectando a Redis: %v", err)
	}

	h := hub.NewHub()
	go h.Run()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	ws := r.Group("/ws")
	ws.Use(handlers.WSAuthMiddleware(cfg))
	{
		ws.GET("", handlers.ChatHandler(h))
	}

	log.Printf("[main] Servidor Go escuchando en :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("[main] Error arrancando servidor: %v", err)
	}
}
