package usecases

import (
	"context"
	"kanban-go/models"
	"kanban-go/repositories/firestore"
)

type TaskUsecase interface {
	GetAllTasks(ctx context.Context) ([]models.Task, error)
}

type taskUsecaseImpl struct {
	repo firestore.FirestoreRepository
}

func NewTaskUsecase(repo firestore.FirestoreRepository) TaskUsecase {
	return &taskUsecaseImpl{
		repo: repo,
	}
}

func (u *taskUsecaseImpl) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return u.repo.GetAllTasks(ctx)
}
