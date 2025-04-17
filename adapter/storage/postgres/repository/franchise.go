package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/OzkrOssa/franchi-api/adapter/storage/postgres"
	"github.com/OzkrOssa/franchi-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type FranchiseRepository struct {
	db *postgres.DB
}

func NewFranchiseRepository(db *postgres.DB) *FranchiseRepository {
	return &FranchiseRepository{db}
}

func (f *FranchiseRepository) CreateFranchise(ctx context.Context, frachise *domain.Franchise) (*domain.Franchise, error) {
	query := f.db.QueryBuilder.Insert("franchises").Columns("name").Values(frachise.Name).Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = f.db.QueryRow(ctx, sql, args...).Scan(
		&frachise.ID,
		&frachise.Name,
	)
	if err != nil {
		if errCode := f.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return frachise, nil
}

func (f *FranchiseRepository) UpdateFranchise(ctx context.Context, frachise *domain.Franchise) (*domain.Franchise, error) {
	query := f.db.QueryBuilder.Update("franchises").
		Set("name", squirrel.Expr("COALESCE(?, name)", frachise.Name))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = f.db.QueryRow(ctx, sql, args...).Scan(
		&frachise.ID,
		&frachise.Name,
	)
	if err != nil {
		if errCode := f.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return frachise, nil
}

func (f *FranchiseRepository) GetFranchiseByID(ctx context.Context, id uint64) (*domain.Franchise, error) {
	var frachise domain.Franchise
	query := f.db.QueryBuilder.Select("*").From("franchises").Where(squirrel.Eq{"id": id}).Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	err = f.db.QueryRow(ctx, sql, args...).Scan(
		&frachise.ID,
		&frachise.Name,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &frachise, nil

}
