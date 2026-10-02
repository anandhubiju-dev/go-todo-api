package repository

import (
	"context"
	"errors"
	"fmt"
	"todo-api/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TodoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(collection *mongo.Collection) *TodoRepository {
	return &TodoRepository{
		collection: collection,
	}
}

func (r *TodoRepository) Create(ctx context.Context, todo model.Todo) (bson.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, todo)
	if err != nil {
		return bson.NilObjectID, err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NewObjectID(), fmt.Errorf("inserted ID is not an ObjectID")
	}

	return id, nil
}

func (r *TodoRepository) GetByID(ctx context.Context, id bson.ObjectID) (model.Todo, error) {
	var todo model.Todo

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&todo)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Todo{}, mongo.ErrNoDocuments
		}

		return model.Todo{}, err
	}

	return todo, nil
}

func (r *TodoRepository) GetAll(ctx context.Context, filter bson.M) ([]model.Todo, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var todos []model.Todo

	err = cursor.All(ctx, &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	todo model.Todo,
) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"title":       todo.Title,
				"description": todo.Description,
				"completed":   todo.Completed,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrTodoNotFound
	}

	return nil
}

func (r *TodoRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {
	result, err := r.collection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrTodoNotFound
	}

	return nil
}
