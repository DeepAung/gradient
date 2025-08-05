package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type BaseRepo[T any] struct {
	BaseCreate[T]
	BaseFind[T]
	BaseUpdate[T]
	BaseDelete[T]
	BasePagination[T]
}

func NewBaseRepo[T any](db bun.IDB) BaseRepo[T] {
	return BaseRepo[T]{
		BaseCreate[T]{db},
		BaseFind[T]{db},
		BaseUpdate[T]{db},
		BaseDelete[T]{db},
		BasePagination[T]{db},
	}
}

// ------------------------------------------------------------------------- //

type BaseCreate[T any] struct{ db bun.IDB }

func (b BaseCreate[T]) Create(ctx context.Context, model *T) error {
	_, err := b.db.NewInsert().Model(model).Returning("*").Exec(ctx)
	return err
}

// ------------------------------------------------------------------------- //

type BaseFind[T any] struct{ db bun.IDB }

func (b BaseFind[T]) FindOneById(ctx context.Context, id uuid.UUID) (*T, error) {
	model := new(T)
	err := b.db.NewSelect().Model(model).Where("id = ?", id).Scan(ctx, model)
	return model, err
}

// ------------------------------------------------------------------------- //

type BaseUpdate[T any] struct{ db bun.IDB }

func (b BaseUpdate[T]) UpdateOne(ctx context.Context, model *T) error {
	_, err := b.db.NewUpdate().Model(model).WherePK().Exec(ctx)
	return err
}

// ------------------------------------------------------------------------- //

type BaseDelete[T any] struct{ db bun.IDB }

func (b BaseDelete[T]) DeleteOneById(ctx context.Context, id uuid.UUID) error {
	model := new(T)
	_, err := b.db.NewDelete().Model(&model).Where("id = ?", id).Exec(ctx)
	return err
}

// ------------------------------------------------------------------------- //

type PaginationInput struct {
	Offset int
	Limit  int
	// the input is db.NewSelect().Model(&models), you can add extra functionality like sort, etc.
	ExtraFn func(q *bun.SelectQuery) *bun.SelectQuery
}

type BasePagination[T any] struct {
	db bun.IDB
}

func (b BasePagination[T]) FindManyByPagination(
	ctx context.Context,
	input PaginationInput,
) ([]*T, error) {
	var models []*T
	query := b.db.NewSelect().Model(&models)
	if input.ExtraFn != nil {
		query = input.ExtraFn(query)
	}
	err := query.Offset(input.Offset).Limit(input.Limit).Scan(ctx, &models)

	return models, err
}
