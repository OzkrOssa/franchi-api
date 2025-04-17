package service

import (
	"context"
	"strings"
	"time"

	"github.com/OzkrOssa/franchi-api/internal/core/domain"
	"github.com/OzkrOssa/franchi-api/internal/core/port"
	"github.com/OzkrOssa/franchi-api/internal/core/utils"
)

type FranchiseService struct {
	repo  port.FranchiseRepository
	cache port.CacheRepository
}

func NewFranchiseService(repo port.FranchiseRepository, cache port.CacheRepository) *FranchiseService {
	return &FranchiseService{
		repo,
		cache,
	}
}

func (s *FranchiseService) CreateNewFranchise(ctx context.Context, franchise *domain.Franchise) (*domain.Franchise, error) {
	if strings.TrimSpace(franchise.Name) == "" {
		return nil, domain.ErrInvalidData
	}

	franchise, err := s.repo.CreateFranchise(ctx, franchise)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	cacheKey := utils.GenerateCacheKey("franchise", franchise.ID)

	serialzedFranchise, err := utils.Serialize(franchise)
	if err != nil {
		return nil, domain.ErrInternal
	}
	err = s.cache.Set(ctx, cacheKey, serialzedFranchise, time.Duration(0))
	if err != nil {
		return nil, domain.ErrInternal
	}

	return franchise, nil
}

func (s *FranchiseService) UpdateFranchise(ctx context.Context, franchise *domain.Franchise) (*domain.Franchise, error) {
	existingFranchise, err := s.repo.GetFranchiseByID(ctx, franchise.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}

		return nil, domain.ErrInternal
	}

	emptyData := franchise.Name == ""
	sameData := franchise.Name == existingFranchise.Name

	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	_, err = s.repo.UpdateFranchise(ctx, franchise)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	cacheKey := utils.GenerateCacheKey("franchise", franchise.ID)

	err = s.cache.Delete(ctx, cacheKey)
	if err != nil {
		return nil, domain.ErrInternal
	}

	serialzedFranchise, err := utils.Serialize(franchise)
	if err != nil {
		return nil, domain.ErrInternal
	}
	err = s.cache.Set(ctx, cacheKey, serialzedFranchise, time.Duration(0))
	if err != nil {
		return nil, domain.ErrInternal
	}

	return nil, nil
}
