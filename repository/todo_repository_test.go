package repository

import (
	"context"
	"errors"
	"testing"

	"todo-api/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupTestCollection(t *testing.T) (*mongo.Collection, func()) {
	t.Helper()

	client, err := mongo.Connect(
		options.Client().ApplyURI("mongodb://host.docker.internal:27017"),
	)

	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}

	ctx := context.Background()

	err = client.Ping(ctx, nil)

	if err != nil {
		t.Fatalf("failed to ping MongoDB: %v", err)
	}

	db := client.Database("todo_test_db")
	collection := db.Collection("todos")

	_, err = collection.DeleteMany(ctx, bson.M{})

	if err != nil {
		t.Fatalf("failed to clean test collection: %v", err)
	}

	cleanup := func() {
		err := client.Disconnect(ctx)

		if err != nil {
			t.Errorf("failed to disconnect MongoDB: %v", err)
		}
	}

	return collection, cleanup
}

func TestTodoRepository_Create(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := NewTodoRepository(collection)

	todo := model.Todo{
		Title:       "Learn integration testing",
		Description: "Test repositorywith MongoDB",
		Completed:   false,
	}

	ctx := context.Background()

	id, err := repo.Create(ctx, todo)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id.IsZero() {
		t.Fatal("expected a valid ObjectID, got zero id")
	}

	var savedTodo model.Todo

	err = collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&savedTodo)

	if err != nil {
		t.Fatalf("failed to find inserted todo: %v", err)
	}

	if savedTodo.Title != todo.Title {
		t.Fatalf("expected title %q, got %q",
			todo.Title,
			savedTodo.Title,
		)
	}

	if savedTodo.Description != todo.Description {
		t.Fatalf(
			"expected description %q, got %q",
			todo.Description,
			savedTodo.Description,
		)
	}

	if savedTodo.Completed != todo.Completed {
		t.Fatalf(
			"expected completed%v, got %v",
			todo.Completed,
			savedTodo.Completed,
		)
	}
}

func TestTodoRepository_GetByID(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := NewTodoRepository(collection)

	ctx := context.Background()

	todo := model.Todo{
		Title:       "Learn integration testing",
		Description: "Test repository with MongoDB",
		Completed:   false,
	}

	result, err := collection.InsertOne(ctx, todo)

	if err != nil {
		t.Fatalf("failed to insert todo: %v", err)
	}

	id := result.InsertedID.(bson.ObjectID)

	savedTodo, err := repo.GetByID(ctx, id)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if savedTodo.ID != id {
		t.Fatalf("expected ID %v, got %v",
			id,
			savedTodo.ID)
	}

	if savedTodo.Title != todo.Title {
		t.Fatalf("expected title %q, got %q",
			todo.Title,
			savedTodo.Title)
	}

	if savedTodo.Description != todo.Description {
		t.Fatalf("expected description %q, got %q",
			todo.Description,
			savedTodo.Description)
	}

	if savedTodo.Completed != todo.Completed {
		t.Fatalf("expected completed %v, got %v",
			todo.Completed,
			savedTodo.Completed,
		)
	}
}

func TestTodoRepository_GetByID_NotFound(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := NewTodoRepository(collection)

	ctx := context.Background()

	id := bson.NewObjectID()

	_, err := repo.GetByID(ctx, id)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		t.Fatalf(
			"expected mongo.ErrNoDocuments, got %v",
			err,
		)
	}
}

func TestTodoRepository_GetAll(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := NewTodoRepository(collection)

	ctx := context.Background()

	todos := []model.Todo{
		{
			Title:       "Learn Go",
			Description: "Study repository patterns",
			Completed:   false,
		},
		{
			Title:       "Build API",
			Description: "Build Todo API",
			Completed:   true,
		},
	}

	_, err := collection.InsertMany(ctx, todos)

	if err != nil {
		t.Fatalf("failed to insert test todos: %v", err)
	}

	result, err := repo.GetAll(ctx, bson.M{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) < 2 {
		t.Fatalf(
			"expected at least 2 todos, got %d",
			len(result),
		)
	}
}

func TestTodoRepository_GetAll_Filter(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := NewTodoRepository(collection)

	ctx := context.Background()

	todos := []model.Todo{
		{
			Title:     "Incomplete Todo",
			Completed: false,
		},
		{
			Title:     "Completed Todo 1",
			Completed: true,
		},
		{
			Title:     "Completed Todo 2",
			Completed: true,
		},
	}

	_, err := collection.InsertMany(ctx, todos)

	if err != nil {
		t.Fatalf("failed to insert test todos: %v", err)
	}

	filter := bson.M{
		"completed": true,
	}

	result, err := repo.GetAll(ctx, filter)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 completed todos, got %d",
			len(result),
		)
	}

	for _, todo := range result {
		if !todo.Completed {
			t.Fatalf(
				"expected only completed todos, got incomplete todo %q",
				todo.Title,
			)
		}
	}
}

func TestTodoRepository_Update(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := NewTodoRepository(collection)

	ctx := context.Background()

	todo := model.Todo{
		Title:       "Original Todo",
		Description: "Original description",
		Completed:   false,
	}

	result, err := collection.InsertOne(ctx, todo)

	if err != nil {
		t.Fatalf("failed to insert test todo: %v", err)
	}

	id := result.InsertedID.(bson.ObjectID)

	updatedTodo := model.Todo{
		Title:       "Updated Todo",
		Description: "Updated description",
		Completed:   true,
	}

	err = repo.Update(ctx, id, updatedTodo)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var savedTodo model.Todo

	err = collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&savedTodo)

	if err != nil {
		t.Fatalf("failed to find updated todo: %v", err)
	}

	if savedTodo.Title != updatedTodo.Title {
		t.Fatalf(
			"expected title %q, got %q",
			updatedTodo.Title,
			savedTodo.Title,
		)
	}

	if savedTodo.Description != updatedTodo.Description {
		t.Fatalf(
			"expected description %q, got %q",
			updatedTodo.Description,
			savedTodo.Description,
		)
	}

	if savedTodo.Completed != updatedTodo.Completed {
		t.Fatalf(
			"expected completed %v, got %v",
			updatedTodo.Completed,
			savedTodo.Completed,
		)
	}
}

func TestTodoRepository_Update_NotFound(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := NewTodoRepository(collection)

	ctx := context.Background()

	id := bson.NewObjectID()

	err := repo.Update(
		ctx,
		id,
		model.Todo{
			Title:       "Updated Todo",
			Description: "Updated description",
			Completed:   true,
		},
	)

	if err != ErrTodoNotFound {
		t.Fatalf(
			"expected ErrTodoNotFound, got %v",
			err,
		)
	}
}

func TestTodoRepository_Delete(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := NewTodoRepository(collection)

	ctx := context.Background()

	todo := model.Todo{
		Title:       "Todo to delete",
		Description: "This todo should be deleted",
		Completed:   false,
	}

	result, err := collection.InsertOne(ctx, todo)

	if err != nil {
		t.Fatalf("failed to insert test todo: %v", err)
	}

	id := result.InsertedID.(bson.ObjectID)

	err = repo.Delete(ctx, id)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var deletedTodo model.Todo

	err = collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&deletedTodo)

	if !errors.Is(err, mongo.ErrNoDocuments) {
		t.Fatalf(
			"expected mongo.ErrNoDocuments, got %v",
			err,
		)
	}
}

func TestTodoRepository_Delete_NotFound(t *testing.T) {
	collection, cleanup := setupTestCollection(t)
	defer cleanup()

	repo := NewTodoRepository(collection)

	ctx := context.Background()

	id := bson.NewObjectID()

	err := repo.Delete(ctx, id)

	if err != ErrTodoNotFound {
		t.Fatalf(
			"expected ErrTodoNotFound, got %v",
			err,
		)
	}
}
