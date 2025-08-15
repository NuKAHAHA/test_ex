package service

import (
	"context"
	"time"

	"awesomeProject/internal/model"
	"awesomeProject/internal/utils/errs"
	"awesomeProject/internal/utils/id"
	"awesomeProject/internal/utils/logger"
)

type TaskService struct {
	repo   TaskRepository
	logger *logger.AsyncLogger
}

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id string) (*model.Task, error)
	List(ctx context.Context, status *string) ([]model.Task, error)
}

func NewTaskService(r TaskRepository, l *logger.AsyncLogger) *TaskService {

	return &TaskService{repo: r, logger: l}
}

func (s *TaskService) Create(ctx context.Context, title, desc, status string) (*model.Task, error) {
	if title == "" {

		return nil, errs.ErrValidation
	}

	if status == "" {
		status = model.StatusPending
	}

	if !model.IsValidStatus(status) {

		return nil, errs.ErrValidation
	}

	now := time.Now().UTC()
	task := &model.Task{
		ID:          id.New(),
		Title:       title,
		Description: desc,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, task); err != nil {

		return nil, err
	}

	s.logger.Log(logger.Event{
		Time:   now,
		Action: "task_created",
		TaskID: task.ID,
		Meta: map[string]any{
			"title":  task.Title,
			"status": task.Status,
		},
	})

	return task, nil
}

func (s *TaskService) Get(ctx context.Context, id string) (*model.Task, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {

		return nil, err
	}

	s.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "task_viewed",
		TaskID: t.ID,
		Meta:   map[string]any{"status": t.Status},
	})

	return t, nil
}

func (s *TaskService) List(ctx context.Context, status *string) ([]model.Task, error) {
	tasks, err := s.repo.List(ctx, status)
	if err != nil {

		return nil, err
	}

	s.logger.Log(logger.Event{
		Time:   time.Now().UTC(),
		Action: "tasks_listed",
		Meta: map[string]any{
			"count":  len(tasks),
			"filter": status,
		},
	})

	return tasks, nil
}
