package services

import (
	"errors"
	"gostockly/internal/models"
	"gostockly/internal/repositories"
	"gostockly/pkg/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo    *repositories.UserRepository
	CompanyRepo *repositories.CompanyRepository
	JWTSecret   string
}

func NewUserService(userRepo *repositories.UserRepository, companyRepo *repositories.CompanyRepository, jwtSecret string) *UserService {
	return &UserService{
		UserRepo:    userRepo,
		CompanyRepo: companyRepo,
		JWTSecret:   jwtSecret,
	}
}

// RegisterUser creates a new user attached to a company.
func (s *UserService) RegisterUser(email, password, companyName, subdomain string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Get or create company
	company, err := s.CompanyRepo.GetCompanyBySubdomain(subdomain)
	if err != nil && err.Error() == "company not found" {
		company = &models.Company{
			ID:        uuid.New(),
			Name:      companyName,
			Subdomain: subdomain,
		}
		if err := s.CompanyRepo.CreateCompany(company); err != nil {
			return nil, errors.New("failed to create company")
		}
	} else if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  string(hashedPassword),
		CompanyID: company.ID,
	}

	if err := s.UserRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser authenticates a user and generates a JWT.
func (s *UserService) AuthenticateUser(email, password string) (string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID.String(), s.JWTSecret)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
