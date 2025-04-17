package port

import (
	"context"

	"github.com/OzkrOssa/franchi-api/internal/core/domain"
)

type FranchiseRepository interface {
	CreateFranchise(ctx context.Context, franchise *domain.Franchise) (*domain.Franchise, error)
	GetFranchiseByID(ctx context.Context, id uint64) (*domain.Franchise, error)
	UpdateFranchise(ctx context.Context, franchise *domain.Franchise) (*domain.Franchise, error)
}

type FranchiseService interface {
	CreateNewFranchise(ctx context.Context, franchise *domain.Franchise) (*domain.Franchise, error)
	UpdateFranchise(ctx context.Context, branch *domain.Branch) (*domain.Franchise, error)
}
