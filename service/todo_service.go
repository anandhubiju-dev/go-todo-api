package service

import (
	"context"
	"errors"

	"todo-api/model"
	"todo-api/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TodoRepository interface {
	Create(
		ctx context.Context,
		todo model.Todo,
	) (bson.ObjectID, error)

	GetByID(
		ctx context.Context,
		id bson.ObjectID,
	) (model.Todo, error)

	GetAll(
		ctx context.Context,
		filter bson.M,
	) ([]model.Todo, error)

	Update(
		ctx context.Context,
		id bson.ObjectID,
		todo model.Todo,
	) error

	Delete(
		ctx context.Context,
		id bson.ObjectID,
	) error
}

type TodoService struct {
	repo TodoRepository
}

func NewTodoService(repo TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) Create(
	ctx context.Context,
	todo model.Todo,
) (bson.ObjectID, error) {
	return s.repo.Create(ctx, todo)
}

func (s *TodoService) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (model.Todo, error) {

	todo, err := s.repo.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Todo{}, ErrTodoNotFound
		}

		return model.Todo{}, err
	}

	return todo, nil
}

func (s *TodoService) GetAll(
	ctx context.Context,
	filter bson.M,
) ([]model.Todo, error) {
	return s.repo.GetAll(ctx, filter)
}

func (s *TodoService) Update(
	ctx context.Context,
	id bson.ObjectID,
	todo model.Todo,
) error {
	err := s.repo.Update(ctx, id, todo)

	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			return ErrTodoNotFound
		}

		return err
	}

	return nil
}

func (s *TodoService) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {
	err := s.repo.Delete(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			return ErrTodoNotFound
		}

		return err
	}

	return nil
}
