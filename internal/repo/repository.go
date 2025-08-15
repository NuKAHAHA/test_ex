package repo

import (
	"context"
	"sync"
	"time"

	"awesomeProject/internal/utils/logger"

	"awesomeProject/internal/model"
	"awesomeProject/internal/utils/errs"
)

type MemoryTaskRepo struct {
	mu     sync.RWMutex
	tasks  map[string]model.Task
	logger *logger.AsyncLogger
}

func NewMemoryTaskRepo(log *logger.AsyncLogger) *MemoryTaskRepo {

	return &MemoryTaskRepo{
		tasks:  make(map[string]model.Task),
		logger: log,
	}
}

func (r *MemoryTaskRepo) Create(ctx context.Context, task *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = *task
	r.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "repo_create",
		TaskID: task.ID,
		Meta:   map[string]any{"status": task.Status},
	})

	return nil
}

func (r *MemoryTaskRepo) GetByID(ctx context.Context, id string) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tasks[id]
	if !ok {
		r.logger.Log(logger.Event{
			Time:   time.Now().UTC(),
			Action: "repo_get_not_found",
			TaskID: id,
		})

		return nil, errs.ErrNotFound
	}

	r.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "repo_get",
		TaskID: t.ID,
	})

	return &t, nil
}

func (r *MemoryTaskRepo) List(ctx context.Context, status *string) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var res []model.Task
	for _, t := range r.tasks {

		if status == nil || t.Status == *status {
			res = append(res, t)
		}
	}

	r.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "repo_list",
		Meta:   map[string]any{"count": len(res)},
	})

	return res, nil
}
