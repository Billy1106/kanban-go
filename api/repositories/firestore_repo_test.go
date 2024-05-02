package repositories

import (
	"log"
	"testing"

	"context"
	"kanban-go/config"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func TestFirestoreRepository(t *testing.T) {

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON([]byte(config.FB_SECRET_CREDENTIAL())))
	if err != nil {
		panic("error initializing app: " + err.Error())
	}

	repo, err := NewFirestoreRepository(ctx, app)
	if err != nil {
		panic(err)
	}

	tasks, err := repo.GetAllTasks(ctx, app)
	if err != nil {
		panic("error getting tasks: " + err.Error())
	}

	log.Println(tasks)

}
