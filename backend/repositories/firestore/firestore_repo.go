package firestore

import (
	"context"
	"kanban-go/config"
	"kanban-go/models"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"time"
)

type FirestoreRepository interface {
	GetAllTasks(ctx context.Context) ([]models.Task, error)
}

type FirestoreRepositoryImpl struct {
	client *firestore.Client
}

var repository *FirestoreRepositoryImpl

func NewFirestoreRepository(ctx context.Context) (FirestoreRepository, error) {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON([]byte(config.FB_SECRET_CREDENTIAL())))
	if err != nil {
		panic(err)
	}
	client, err := app.Firestore(ctx)
	repository = &FirestoreRepositoryImpl{
		client: client,
	}

	return repository, err
}

// Note r is a pointer receiver
func (r *FirestoreRepositoryImpl) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	docs := r.client.Collection("tasks").Documents(ctx)
	tasks := []models.Task{}

	for {
		doc, err := docs.Next()
		if err != nil {
			break
		}

		statusRef, ok := doc.Data()["status"].(*firestore.DocumentRef)
		if !ok {
			panic("status is not a document ref")
		}

		// Get the status document
		docStatus, err := r.client.Doc("statuses/" + statusRef.ID).Get(ctx)
		if err != nil {
			panic(err)
		}

		status := models.Status{
			ID:   docStatus.Ref.ID,
			Name: docStatus.Data()["name"].(string),
		}

		task := models.Task{
			ID:          doc.Ref.ID,
			Title:       doc.Data()["title"].(string),
			Description: doc.Data()["description"].(string),
			Status:      status,
			CreatedAt:   doc.Data()["created_at"].(time.Time),
			UpdatedAt:   doc.Data()["updated_at"].(time.Time),
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
