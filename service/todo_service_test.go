package service

import (
	"context"
	"errors"
	"testing"

	"todo-api/model"
	"todo-api/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type fakeTodoRepository struct {
	createdID    bson.ObjectID
	createdErr   error
	todo         model.Todo
	getByIDErr   error
	todos        []model.Todo
	getAllFilter bson.M
	getAllErr    error
	updateErr    error
	updatedTodo  model.Todo
	updatedID    bson.ObjectID
	deleteErr    error
	deletedID    bson.ObjectID
}

func (f *fakeTodoRepository) Create(
	ctx context.Context,
	todo model.Todo,
) (bson.ObjectID, error) {
	return f.createdID, f.createdErr
}

func (f *fakeTodoRepository) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (model.Todo, error) {
	return f.todo, f.getByIDErr
}

func (f *fakeTodoRepository) GetAll(
	ctx context.Context,
	filter bson.M,
) ([]model.Todo, error) {
	f.getAllFilter = filter

	return f.todos, f.getAllErr
}

func (f *fakeTodoRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	todo model.Todo,
) error {
	f.updatedID = id
	f.updatedTodo = todo
	return f.updateErr
}

func (f *fakeTodoRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {
	f.deletedID = id

	return f.deleteErr
}

func TestTodoService_Create(t *testing.T) {
	expectedID := bson.NewObjectID()
	fakeRepo := &fakeTodoRepository{
		createdID: expectedID,
	}

	todoService := NewTodoService(fakeRepo)

	todo := model.Todo{
		Title:       "Learn Go interface",
		Description: "Understand dependancy injection",
		Completed:   false,
	}

	id, err := todoService.Create(
		context.Background(),
		todo,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id == bson.NilObjectID {
		t.Fatalf("expected a valid ObjectID, ot NilObjectID")
	}

	if id != expectedID {
		t.Fatalf("expected ID %v, got %v", expectedID, id)
	}
}

func TestTodoService_Create_Error(t *testing.T) {
	repoErr := errors.New("databse error")

	fakeRepo := &fakeTodoRepository{
		createdErr: repoErr,
	}

	todoService := NewTodoService(fakeRepo)

	todo := model.Todo{
		Title: "Test error",
	}

	_, err := todoService.Create(
		context.Background(),
		todo,
	)

	if err != repoErr {
		t.Fatalf("expected error %v, but got %v", repoErr, err)
	}
}

func TestTodoService_GetByID(t *testing.T) {
	expectedID := bson.NewObjectID()

	expectedTodo := model.Todo{
		ID:          expectedID,
		Title:       "Learn Go testing",
		Description: "Practice service unit testing",
		Completed:   false,
	}

	fakeRepo := &fakeTodoRepository{
		todo: expectedTodo,
	}

	todoService := NewTodoService(fakeRepo)

	todo, err := todoService.GetByID(
		context.Background(),
		expectedID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if todo.ID != expectedTodo.ID {
		t.Fatalf("expected ID %v, got %v", expectedTodo.ID, todo.ID)
	}

	if todo.Title != expectedTodo.Title {
		t.Fatalf("expected title %q, got %q", expectedTodo.Title, todo.Title)
	}

	if todo.Description != expectedTodo.Description {
		t.Fatalf(
			"expected description %q, got %q",
			expectedTodo.Description,
			todo.Description,
		)
	}

	if todo.Completed != expectedTodo.Completed {
		t.Fatalf(
			"expected completed %v, got %v",
			expectedTodo.Completed,
			todo.Completed,
		)
	}
}

func TestTodoService_GetByID_NotFound(t *testing.T) {
	expectedID := bson.NewObjectID()

	fakeRepo := &fakeTodoRepository{
		getByIDErr: mongo.ErrNoDocuments,
	}

	todoService := NewTodoService(fakeRepo)

	_, err := todoService.GetByID(
		context.Background(),
		expectedID,
	)

	if err != ErrTodoNotFound {
		t.Fatalf(
			"expected ErrTodoNotFound. got %v",
			err,
		)
	}
}

func TestTodoService_GetByID_Error(t *testing.T) {
	repoErr := errors.New("datatbase connection failed")

	fakeRepo := &fakeTodoRepository{
		getByIDErr: repoErr,
	}

	todoService := NewTodoService(fakeRepo)

	_, err := todoService.GetByID(
		context.Background(),
		bson.NewObjectID(),
	)

	if err != repoErr {
		t.Fatalf(
			"expected error %v, got %v",
			repoErr,
			err,
		)
	}
}

func TestTodoService_GetAll(t *testing.T) {
	expectedTodos := []model.Todo{
		{
			ID:          bson.NewObjectID(),
			Title:       "Learn Go",
			Description: "Study interfaces",
			Completed:   false,
		},
		{
			ID:          bson.NewObjectID(),
			Title:       "Build API",
			Description: "Build Todo API",
			Completed:   true,
		},
	}

	fakeRepo := &fakeTodoRepository{
		todos: expectedTodos,
	}

	todoService := NewTodoService(fakeRepo)

	todos, err := todoService.GetAll(
		context.Background(),
		bson.M{},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(todos) != len(expectedTodos) {
		t.Fatalf(
			"expected %dtodos, got %d",
			len(expectedTodos),
			len(todos),
		)
	}

	if todos[1].Title != expectedTodos[1].Title {
		t.Fatalf(
			"expected title %q, got %q",
			expectedTodos[1].Title,
			todos[1].Title,
		)
	}
}

func TestTodoServivce_GetAll_filter(t *testing.T) {
	expectedFilter := bson.M{
		"completed": true,
	}

	fakeRepo := &fakeTodoRepository{
		todos: []model.Todo{
			{
				ID:        bson.NewObjectID(),
				Title:     "Completed Todo",
				Completed: true,
			},
		},
	}

	todoService := NewTodoService(fakeRepo)

	_, err := todoService.GetAll(
		context.Background(),
		expectedFilter,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fakeRepo.getAllFilter["completed"] != true {
		t.Fatalf(
			"expected completed filter to be true, got %v",
			fakeRepo.getAllFilter["completed"],
		)
	}
}

func TestTodoService_GetAll_Error(t *testing.T) {
	repoErr := errors.New("database connection failed")

	fakeRepo := &fakeTodoRepository{
		getAllErr: repoErr,
	}

	todoService := NewTodoService(fakeRepo)

	_, err := todoService.GetAll(
		context.Background(),
		bson.M{},
	)

	if err != repoErr {
		t.Fatalf(
			"expeceted error %v, got %v",
			repoErr,
			err,
		)
	}
}

func TestTodoService_Update(t *testing.T) {
	expectedID := bson.NewObjectID()

	expectedTodo := model.Todo{
		ID:          expectedID,
		Title:       "Updtaed Todo",
		Description: "Updated description",
		Completed:   true,
	}

	fakeRepo := &fakeTodoRepository{}

	todoService := NewTodoService(fakeRepo)

	err := todoService.Update(
		context.Background(),
		expectedID,
		expectedTodo,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fakeRepo.updatedID != expectedID {
		t.Fatalf(
			"expected ID %v, got %v",
			expectedID,
			fakeRepo.updatedID,
		)
	}

	if fakeRepo.updatedTodo.Title != expectedTodo.Title {
		t.Fatalf(
			"expected title%q, got %q",
			expectedTodo.Title,
			fakeRepo.updatedTodo.Title,
		)
	}

	if fakeRepo.updatedTodo.Description != expectedTodo.Description {
		t.Fatalf(
			"expected description %q, got %q",
			expectedTodo.Description,
			fakeRepo.updatedTodo.Description,
		)
	}

	if fakeRepo.updatedTodo.Completed != expectedTodo.Completed {
		t.Fatalf(
			"expected completed%v, got %v",
			expectedTodo.Completed,
			fakeRepo.updatedTodo.Completed,
		)
	}
}

func TestTodoService_Update_NotFound(t *testing.T) {
	fakeRepo := &fakeTodoRepository{
		updateErr: repository.ErrTodoNotFound,
	}

	todoService := NewTodoService(fakeRepo)

	err := todoService.Update(
		context.Background(),
		bson.NewObjectID(),
		model.Todo{
			Title: "Updated Todo",
		},
	)

	if err != ErrTodoNotFound {
		t.Fatalf(
			"expected ErrTodoNotFound, got %v",
			err,
		)
	}
}

func TestTodoService_Update_Error(t *testing.T) {
	repoErr := errors.New("database connection failed")

	fakeRepo := &fakeTodoRepository{
		updateErr: repoErr,
	}

	todoService := NewTodoService(fakeRepo)

	err := todoService.Update(
		context.Background(),
		bson.NewObjectID(),
		model.Todo{
			Title: "Updated Todo",
		},
	)

	if err != repoErr {
		t.Fatalf(
			"expected error %v, got %v",
			repoErr,
			err,
		)
	}
}

func TestTodoService_Delete(t *testing.T) {
	expectedID := bson.NewObjectID()

	fakeRepo := &fakeTodoRepository{}

	todoService := NewTodoService(fakeRepo)

	err := todoService.Delete(
		context.Background(),
		expectedID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fakeRepo.deletedID != expectedID {
		t.Fatalf(
			"expected ID %v, got %v",
			expectedID,
			fakeRepo.deletedID,
		)
	}
}

func TestTodoService_Delete_NotFound(t *testing.T) {
	fakeRepo := &fakeTodoRepository{
		deleteErr: repository.ErrTodoNotFound,
	}

	todoService := NewTodoService(fakeRepo)

	err := todoService.Delete(
		context.Background(),
		bson.NewObjectID(),
	)

	if err != ErrTodoNotFound {
		t.Fatalf(
			"expected ErrTodoNotFound, got %v",
			err,
		)
	}
}

func TestTodoService_Delete_Error(t *testing.T) {
	repoErr := errors.New("database connection failed")

	fakeRepo := &fakeTodoRepository{
		deleteErr: repoErr,
	}

	todoService := NewTodoService(fakeRepo)

	err := todoService.Delete(
		context.Background(),
		bson.NewObjectID(),
	)

	if err != repoErr {
		t.Fatalf(
			"expected error %v, got %v",
			repoErr,
			err,
		)
	}
}

var _ TodoRepository = (*fakeTodoRepository)(nil)
