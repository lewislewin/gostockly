package repositories

import (
	"errors"
	"gostockly/internal/models"

	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// CreateCompany adds a new company.
func (r *CompanyRepository) CreateCompany(company *models.Company) error {
	return r.db.Create(company).Error
}

// GetCompanyBySubdomain retrieves a company by its subdomain.
func (r *CompanyRepository) GetCompanyBySubdomain(subdomain string) (*models.Company, error) {
	var company models.Company
	err := r.db.Where("subdomain = ?", subdomain).First(&company).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("company not found")
	}
	return &company, err
}
