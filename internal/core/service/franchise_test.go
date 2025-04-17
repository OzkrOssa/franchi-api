package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/OzkrOssa/franchi-api/internal/core/domain"
	"github.com/OzkrOssa/franchi-api/internal/core/port/mocks"
	"github.com/OzkrOssa/franchi-api/internal/core/service"
	"github.com/OzkrOssa/franchi-api/internal/core/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

type newFranchiseTestedInput struct {
	franchise *domain.Franchise
}
type newFranchiseExpectedOutput struct {
	franchise *domain.Franchise
	err       error
}

func TestFranchiseService_CreateNewFranchise(t *testing.T) {
	ctx := context.Background()

	fakeName := gofakeit.Name()
	fakeID := gofakeit.Uint64()

	franchiseInput := &domain.Franchise{
		Name:   fakeName,
		Branch: nil,
	}

	franchiseOutput := &domain.Franchise{
		ID:     fakeID,
		Name:   fakeName,
		Branch: nil,
	}

	cacheKey := utils.GenerateCacheKey("franchise", franchiseOutput.ID)

	serialzedFranchise, _ := utils.Serialize(franchiseOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc     string
		mocks    func(repo *mocks.MockFranchiseRepository, cache *mocks.MockCacheRepository)
		input    newFranchiseTestedInput
		expected newFranchiseExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(repo *mocks.MockFranchiseRepository, cache *mocks.MockCacheRepository) {
				repo.EXPECT().CreateFranchise(ctx, franchiseInput).Return(franchiseOutput, nil)
				cache.EXPECT().Set(ctx, cacheKey, serialzedFranchise, ttl).Return(nil)
			},
			input:    newFranchiseTestedInput{franchise: franchiseInput},
			expected: newFranchiseExpectedOutput{franchise: franchiseOutput, err: nil},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(repo *mocks.MockFranchiseRepository, cache *mocks.MockCacheRepository) {
				repo.EXPECT().CreateFranchise(ctx, franchiseInput).Return(nil, domain.ErrConflictingData)
			},
			input:    newFranchiseTestedInput{franchise: franchiseInput},
			expected: newFranchiseExpectedOutput{franchise: nil, err: domain.ErrConflictingData},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(repo *mocks.MockFranchiseRepository, cache *mocks.MockCacheRepository) {
				repo.EXPECT().CreateFranchise(ctx, franchiseInput).Return(nil, domain.ErrInternal)
			},
			input:    newFranchiseTestedInput{franchise: franchiseInput},
			expected: newFranchiseExpectedOutput{franchise: nil, err: domain.ErrInternal},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(repo *mocks.MockFranchiseRepository, cache *mocks.MockCacheRepository) {
				repo.EXPECT().CreateFranchise(ctx, franchiseInput).Return(franchiseOutput, nil)
				cache.EXPECT().Set(ctx, cacheKey, serialzedFranchise, ttl).Return(domain.ErrInternal)
			},
			input:    newFranchiseTestedInput{franchise: franchiseInput},
			expected: newFranchiseExpectedOutput{franchise: nil, err: domain.ErrInternal},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := mocks.NewMockFranchiseRepository(t)
			cache := mocks.NewMockCacheRepository(t)

			tc.mocks(repo, cache)

			svc := service.NewFranchiseService(repo, cache)

			franchise, err := svc.CreateNewFranchise(ctx, tc.input.franchise)

			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.franchise, franchise, "Franchise mismatch")

		})
	}
}
