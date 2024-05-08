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
	GetAllTasks(ctx context.Context) ([]models.TaskResponse, error)
	GetTask(ctx context.Context, id string) (models.TaskResponse, error)
	AddTask(ctx context.Context, task models.TaskBody) (string, error)
	UpdateTask(ctx context.Context, id string, task models.TaskBody) error
	DeleteTask(ctx context.Context, id string) error
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
func (r *FirestoreRepositoryImpl) GetAllTasks(ctx context.Context) ([]models.TaskResponse, error) {
	docs := r.client.Collection("tasks").Documents(ctx)
	tasks := []models.TaskResponse{}

	for {
		doc, err := docs.Next()
		if err != nil {
			break
		}

		statusRef, ok := doc.Data()["Status"].(*firestore.DocumentRef)
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

		task := models.TaskResponse{
			ID:          doc.Ref.ID,
			Title:       doc.Data()["Title"].(string),
			Description: doc.Data()["Description"].(string),
			Status:      status,
			CreatedAt:   doc.Data()["CreatedAt"].(time.Time),
			UpdatedAt:   doc.Data()["UpdatedAt"].(time.Time),
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *FirestoreRepositoryImpl) GetTask(ctx context.Context, id string) (models.TaskResponse, error) {
	doc, err := r.client.Collection("tasks").Doc(id).Get(ctx)
	if err != nil {
		return models.TaskResponse{}, err
	}

	statusRef, ok := doc.Data()["Status"].(*firestore.DocumentRef)
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

	task := models.TaskResponse{
		ID:          doc.Ref.ID,
		Title:       doc.Data()["Title"].(string),
		Description: doc.Data()["Description"].(string),
		Status:      status,
		CreatedAt:   doc.Data()["CreatedAt"].(time.Time),
		UpdatedAt:   doc.Data()["UpdatedAt"].(time.Time),
	}

	return task, nil
}

func (r *FirestoreRepositoryImpl) GetStatus(ctx context.Context, id *firestore.DocumentRef) (models.Status, error) {
	doc, err := r.client.Collection("statuses").Doc(id.ID).Get(ctx)
	if err != nil {
		return models.Status{}, err
	}

	status := models.Status{
		ID:   doc.Ref.ID,
		Name: doc.Data()["name"].(string),
	}

	return status, nil
}

func (r *FirestoreRepositoryImpl) AddTask(ctx context.Context, task models.TaskBody) (string, error) {
	_, err := r.GetStatus(ctx, task.Status)
	if err != nil {
		return "", err
	}
	ref := r.client.Collection("tasks").NewDoc()

	newTask := models.TaskBody{
		ID:          ref.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ref.Set(ctx, newTask)

	return ref.ID, nil
}

func (r *FirestoreRepositoryImpl) UpdateTask(ctx context.Context, id string, newTask models.TaskBody) error {
	currentTask, err := r.GetTask(ctx, id)
	if err != nil {
		return err
	}
	var updatedStatus *firestore.DocumentRef
	if newTask.Status == nil {
		updatedStatus = r.client.Doc("statuses/" + currentTask.Status.ID)
	} else {
		updatedStatus = r.client.Doc("statuses/" + newTask.Status.ID)
	}

	updatingTask := models.TaskBody{
		ID:          currentTask.ID,
		Title:       replaceIfEmpty(currentTask.Title, newTask.Title),
		Description: replaceIfEmpty(currentTask.Description, newTask.Description),
		Status:      updatedStatus,
		CreatedAt:   replaceIfEmptyTime(currentTask.CreatedAt, newTask.CreatedAt),
		UpdatedAt:   time.Now(),
	}

	_, err = r.client.Collection("tasks").Doc(id).Set(ctx, updatingTask)

	return err
}

func (r *FirestoreRepositoryImpl) DeleteTask(ctx context.Context, id string) error {
	_, err := r.client.Collection("tasks").Doc(id).Delete(ctx)
	return err
}

func replaceIfEmpty(current string, updated string) string {
	if updated == "" {
		return current
	}
	return updated
}

func replaceIfEmptyTime(current time.Time, updated time.Time) time.Time {
	if updated.IsZero() {
		return current
	}
	return updated
}
