package main

import (
	"context"
	"log"
	"net/http"
	"notification-service/internal/handlers"
	"notification-service/internal/repositories"
	"notification-service/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Configurar repositorio de notificaciones
	notificationRepo := repositories.NewSMTPNotificationRepository(
		os.Getenv("SMTP_SERVER"),
		os.Getenv("FROM_EMAIL"),
		os.Getenv("TO_EMAIL"),
	)

	// Inicializar dependencias
	notificationService := services.NewNotificationService(notificationRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Configurar router
	mux := http.NewServeMux()
	mux.HandleFunc("/notify", notificationHandler.HandleErrorNotification)

	// Configurar servidor
	server := &http.Server{
		Addr:    ":8081",
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
	log.Println("Servidor de notificaciones iniciado en :8081")

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
