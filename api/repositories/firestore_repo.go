package repositories

import (
	"context"
	"kanban-go/models"
	"log"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"time"
)

// type FirestoreRepository struct {
// 	GetAllTasks func() ([]models.Task, error)
// 	GetTask     func(id int) (models.Task, error)
// 	CreateTask  func(task models.Task) (models.Task, error)
// 	UpdateTask  func(task models.Task) (models.Task, error)
// 	DeleteTask  func(id int) error
// }

type FirestoreRepositoryImpl struct {
	client *firestore.Client
}

var repository *FirestoreRepositoryImpl

func NewFirestoreRepository(ctx context.Context, app *firebase.App) (*FirestoreRepositoryImpl, error) {
	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	repository = &FirestoreRepositoryImpl{
		client: client,
	}

	return repository, err
}

// Note r is a pointer receiver
func (r *FirestoreRepositoryImpl) GetAllTasks(ctx context.Context, client *firebase.App) ([]models.Task, error) {
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

		docStatus, err := r.client.Doc("statuses/" + statusRef.ID).Get(ctx)
		if err != nil {
			panic(err)
		}

		log.Print(docStatus.Data())

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
