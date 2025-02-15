package services

import (
	"errors"
	"gostockly/internal/models"
	"gostockly/internal/repositories"
	"gostockly/pkg/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo    *repositories.UserRepository
	CompanyRepo *repositories.CompanyRepository
	JWTSecret   string
}

func NewAuthService(userRepo *repositories.UserRepository, companyRepo *repositories.CompanyRepository, jwtSecret string) *AuthService {
	return &AuthService{
		UserRepo:    userRepo,
		CompanyRepo: companyRepo,
		JWTSecret:   jwtSecret,
	}
}

// RegisterUser creates a new user account attached to a company.
func (s *AuthService) RegisterUser(email, password, companyName, subdomain string) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Check if the company exists or create a new one
	company, err := s.CompanyRepo.GetCompanyBySubdomain(subdomain)
	if err != nil && err.Error() == "company not found" {
		company = &models.Company{
			ID:        uuid.New(),
			Name:      companyName,
			Subdomain: subdomain,
			CreatedAt: time.Now(),
		}
		if err := s.CompanyRepo.CreateCompany(company); err != nil {
			return nil, errors.New("failed to create company")
		}
	} else if err != nil {
		return nil, err
	}

	// Create the user
	user := &models.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  string(hashedPassword),
		CompanyID: company.ID,
		CreatedAt: time.Now(),
	}

	if err := s.UserRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser checks user credentials and generates a JWT.
func (s *AuthService) AuthenticateUser(email, password string) (string, error) {
	// Retrieve the user by email
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate a JWT token
	token, err := utils.GenerateJWT(user.ID.String(), s.JWTSecret)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
