package firestore

import (
	"log"
	"testing"

	"context"

	"kanban-go/models"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestFirestoreRepository(t *testing.T) {

	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	repo, err := NewFirestoreRepository(ctx)
	if err != nil {
		panic(err)
	}

	client := repo.(*FirestoreRepositoryImpl).client
	title := "New Task"
	description := "New Task Description"
	status := client.Doc("statuses/9aOxL77OwnFtVg5XcLJs")
	newTask := models.TaskBody{
		Title:       title,
		Description: description,
		Status:      status,
	}

	id, err := repo.AddTask(ctx, newTask)
	if err != nil {
		panic("error adding task: " + err.Error())
	}

	task, err := repo.GetTask(ctx, id)
	if err != nil {
		panic("error getting task: " + err.Error())

	}
	assert.Equal(t, newTask.Description, task.Description)

	newTitle := "Updated Task"

	updatedTask := models.TaskBody{
		ID:    id,
		Title: newTitle,
	}

	err = repo.UpdateTask(ctx, id, updatedTask)
	if err != nil {
		panic("error updating task: " + err.Error())
	}

	task, err = repo.GetTask(ctx, id)
	if err != nil {
		panic("error getting task: " + err.Error())
	}

	log.Println(task)
	assert.Equal(t, updatedTask.Title, newTitle)

	repo.DeleteTask(ctx, id)

	_, err = repo.GetTask(ctx, id)
	if err == nil {
		log.Println("Task not deleted")
	}

	// tasks, err := repo.GetAllTasks(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// if err != nil {
	// 	panic("error getting tasks: " + err.Error())
	// }

	// for _, task := range tasks {
	// 	log.Println(task)
	// }

}
