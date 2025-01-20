package services

import (
	"gostockly/internal/models"
	"gostockly/internal/repositories"

	"github.com/google/uuid"
)

type CompanyService struct {
	Repo *repositories.CompanyRepository
}

func NewCompanyService(repo *repositories.CompanyRepository) *CompanyService {
	return &CompanyService{Repo: repo}
}

// CreateCompany creates a new company.
func (s *CompanyService) CreateCompany(name, subdomain string) (*models.Company, error) {
	company := &models.Company{
		ID:        uuid.New(),
		Name:      name,
		Subdomain: subdomain,
	}

	err := s.Repo.CreateCompany(company)
	if err != nil {
		return nil, err
	}

	return company, nil
}

// GetCompanyBySubdomain retrieves a company by subdomain.
func (s *CompanyService) GetCompanyBySubdomain(subdomain string) (*models.Company, error) {
	return s.Repo.GetCompanyBySubdomain(subdomain)
}
