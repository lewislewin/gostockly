package services

import (
	"gostockly/internal/models"
	"gostockly/internal/repositories"

	"github.com/google/uuid"
)

type StoreService struct {
	Repo *repositories.StoreRepository
}

func NewStoreService(repo *repositories.StoreRepository) *StoreService {
	return &StoreService{Repo: repo}
}

// CreateStore adds a store to a company.
func (s *StoreService) CreateStore(companyID, shopifyStoreStub, accessToken string) (*models.Store, error) {
	store := &models.Store{
		ID:               uuid.New(),
		CompanyID:        uuid.MustParse(companyID),
		ShopifyStoreStub: shopifyStoreStub,
		AccessToken:      accessToken,
	}

	err := s.Repo.CreateStore(store)
	if err != nil {
		return nil, err
	}

	return store, nil
}

// GetStoresByCompany retrieves all stores for a company.
func (s *StoreService) GetStoresByCompany(companyID string) ([]models.Store, error) {
	return s.Repo.GetStoresByCompany(companyID)
}
