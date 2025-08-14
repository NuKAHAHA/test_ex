package handler

import (
	"awesomeProject/internal/utils/logger"
	"net/http"

	"awesomeProject/internal/service"
	"awesomeProject/internal/utils/middleware"
)

func NewRouter(svc *service.TaskService, log *logger.AsyncLogger) http.Handler {

	h := NewTaskHandler(svc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", h.List)
	mux.HandleFunc("GET /tasks/", h.GetByID)
	mux.HandleFunc("POST /tasks", h.Create)

	return middleware.RecoverJSON(mux)
}
