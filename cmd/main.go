package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"awesomeProject/internal/handler"
	"awesomeProject/internal/repo"
	"awesomeProject/internal/service"
	"awesomeProject/internal/utils/logger"
)

func main() {
	// Инициализация
	logg := logger.NewAsyncLogger(os.Stdout, 2048)

	taskRepo := repo.NewMemoryTaskRepo(logg)
	taskService := service.NewTaskService(taskRepo, logg)

	mux := handler.NewRouter(taskService, logg)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем HTTP сервер
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Закрываем логгер
	if err := logg.Close(ctx); err != nil {
		log.Printf("Logger close error: %v", err)
	}

	log.Println("Server stopped")
}
