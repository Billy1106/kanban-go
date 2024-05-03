package repositories

import (
	"log"
	"testing"

	"context"
)

func TestFirestoreRepository(t *testing.T) {

	ctx := context.Background()

	repo, err := NewFirestoreRepository(ctx)
	if err != nil {
		panic(err)
	}

	tasks, err := repo.GetAllTasks(ctx)
	if err != nil {
		panic("error getting tasks: " + err.Error())
	}

	log.Println(tasks)

}
