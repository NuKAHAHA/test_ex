package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"awesomeProject/internal/utils/errs"
	"awesomeProject/internal/utils/helper"
	"awesomeProject/internal/utils/logger"

	"awesomeProject/internal/model"
	"awesomeProject/internal/service"
)

type TaskHandler struct {
	svc    *service.TaskService
	logger *logger.AsyncLogger
}

func NewTaskHandler(svc *service.TaskService, log *logger.AsyncLogger) *TaskHandler {

	return &TaskHandler{svc: svc, logger: log}
}

// Маппинг ошибок
func httpStatusFromError(err error) int {
	switch err {
	case nil:

		return http.StatusOK
	case errs.ErrNotFound:

		return http.StatusNotFound
	case errs.ErrBadRequest:

		return http.StatusBadRequest
	default:

		return http.StatusInternalServerError
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "http_request_start",
		Meta:   map[string]any{"method": r.Method, "path": r.URL.Path},
	})

	var req model.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logAndRespondError(w, r, errs.ErrBadRequest, "invalid_json")

		return
	}

	task, err := h.svc.Create(r.Context(), req.Title, req.Description, req.Status)
	if err != nil {
		h.logAndRespondError(w, r, err, "create_task_error")

		return
	}

	h.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "task_created",
		Meta:   map[string]any{"task_id": task.ID},
	})

	helper.JSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "http_request_start",
		Meta:   map[string]any{"method": r.Method, "path": r.URL.Path},
	})

	id := r.URL.Path[len("/tasks/"):]
	if id == "" {
		helper.ErrorJSON(w, http.StatusBadRequest, "missing task ID")

		return
	}
	task, err := h.svc.Get(r.Context(), id)
	if err != nil {
		h.logAndRespondError(w, r, err, "get_task_error", "task_id", id)

		return
	}

	h.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "task_retrieved",
		Meta:   map[string]any{"task_id": task.ID},
	})

	helper.JSON(w, http.StatusOK, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	h.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "http_request_start",
		Meta:   map[string]any{"method": r.Method, "path": r.URL.Path},
	})

	status := r.URL.Query().Get("status")

	var filter *string
	if status != "" {
		filter = &status
	}

	tasks, err := h.svc.List(r.Context(), filter)
	if err != nil {
		h.logAndRespondError(w, r, err, "list_tasks_error")

		return
	}

	h.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "tasks_listed",
		Meta:   map[string]any{"count": len(tasks)},
	})

	helper.JSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) logAndRespondError(w http.ResponseWriter, r *http.Request, err error, action string, extraMeta ...any) {
	meta := map[string]any{
		"method": r.Method,
		"path":   r.URL.Path,
		"error":  err.Error(),
	}

	for i := 0; i < len(extraMeta)-1; i += 2 {
		key, ok := extraMeta[i].(string)

		if ok {
			meta[key] = extraMeta[i+1]
		}
	}

	status := httpStatusFromError(err)

	h.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: action,
		Meta:   meta,
	})

	helper.ErrorJSON(w, status, err.Error())
}
