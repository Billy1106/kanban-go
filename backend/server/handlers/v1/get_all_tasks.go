package v1

import (
	"encoding/json"
	usecase "kanban-go/usecases"
	"net/http"
)

type taskHandler struct {
	taskUsecase usecase.TaskUsecase
}

func NewTaskHandler(tu usecase.TaskUsecase) *taskHandler {
	return &taskHandler{
		taskUsecase: tu,
	}
}

func (taskHandler *taskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := taskHandler.taskUsecase.GetAllTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
