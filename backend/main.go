package main

import (
	"context"
	"net/http"

	"kanban-go/repositories/firestore"

	v1 "kanban-go/server/handlers/v1"

	usecases "kanban-go/usecases"
)

func main() {

	ctx := context.Background()

	fr, err := firestore.NewFirestoreRepository(ctx)
	if err != nil {
		panic(err)
	}
	fu := usecases.NewTaskUsecase(fr)
	fh := v1.NewTaskHandler(fu)

	// http://localhost:8080/v1/tasks
	http.HandleFunc("/v1/tasks", fh.GetAllTasks)

	http.ListenAndServe(":8080", nil)

}
