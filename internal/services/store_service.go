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

func (s *StoreService) CreateStore(companyID, shopifyStoreStub, accessToken, webhookSignature, locationID string) (*models.Store, error) {
	store := &models.Store{
		ID:               uuid.New(),
		CompanyID:        uuid.MustParse(companyID),
		ShopifyStoreStub: shopifyStoreStub,
		AccessToken:      accessToken,
		WebhookSignature: webhookSignature,
		LocationID:       locationID,
	}

	err := s.Repo.CreateStore(store)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *StoreService) GetStoresByCompany(companyID string) ([]models.Store, error) {
	return s.Repo.GetStoresByCompany(companyID)
}

func (s *StoreService) GetStoreByID(storeID string) (*models.Store, error) {
	return s.Repo.GetStoreByID(storeID)
}

func (s *StoreService) UpdateStore(storeID, shopifyStoreStub, accessToken, webhookSignature, locationID string) (*models.Store, error) {
	store, err := s.Repo.GetStoreByID(storeID)
	if err != nil {
		return nil, err
	}

	store.ShopifyStoreStub = shopifyStoreStub
	store.AccessToken = accessToken
	store.WebhookSignature = webhookSignature
	store.LocationID = locationID

	if err := s.Repo.UpdateStore(store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *StoreService) DeleteStore(storeID string) error {
	return s.Repo.DeleteStore(storeID)
}
