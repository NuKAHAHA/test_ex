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
	// Инициализируем логгер
	logg := logger.NewAsyncLogger(os.Stdout, 2048)

	// Репозиторий и сервис
	taskRepo := repo.NewMemoryTaskRepo(logg)
	taskService := service.NewTaskService(taskRepo, logg)

	// Роутер
	mux := handler.NewRouter(taskService, logg)

	// HTTP сервер
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Запуск сервера
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
