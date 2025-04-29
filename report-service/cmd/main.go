package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"report-service/internal/handlers"
	"report-service/internal/repositories"
	"report-service/internal/services"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Configurar Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Verificar conexión a Redis
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Error conectando a Redis: %v", err)
	}

	// Inicializar dependencias
	reportRepo := repositories.NewRedisReportRepository(redisClient)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// Configurar router
	mux := http.NewServeMux()
	mux.HandleFunc("/report", reportHandler.HandleReport)

	// Configurar servidor
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Manejar señales de terminación
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()
	log.Println("Servidor de reportes iniciado en :8080")

	// Esperar señal de terminación
	<-done
	log.Println("Servidor deteniéndose...")

	// Configurar contexto para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error deteniendo servidor: %v", err)
	}
	log.Println("Servidor detenido correctamente")
}
