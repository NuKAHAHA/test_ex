package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"awesomeProject/internal/handler"
	"awesomeProject/internal/repo"
	"awesomeProject/internal/service"
	"awesomeProject/internal/utils/logger"
	"awesomeProject/internal/utils/shuwdown"
)

func main() {
	sm := shutdown.New(5 * time.Second)

	logg := logger.NewAsyncLogger(os.Stdout, 2048)
	sm.Register(func(ctx context.Context) error {
		return logg.Close(ctx)
	})

	taskRepo := repo.NewMemoryTaskRepo(logg)
	taskService := service.NewTaskService(taskRepo, logg)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler.NewRouter(taskService, logg),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	sm.Register(func(ctx context.Context) error {
		log.Println("Stopping HTTP server...")
		return srv.Shutdown(ctx)
	})

	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	sm.Wait()
}
