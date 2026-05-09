package service

import (
	"errors"
	"time"

	"casadavo/internal/middleware"
	"casadavo/internal/model"
	"casadavo/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users      *repository.UserRepo
	jwtSecret  []byte
	adminEmail string
}

func NewAuthService(users *repository.UserRepo, secret, adminEmail string) *AuthService {
	return &AuthService{users: users, jwtSecret: []byte(secret), adminEmail: adminEmail}
}

func (s *AuthService) Register(name, email, phone, password, role string) (*model.User, error) {
	existing, err := s.users.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if s.adminEmail != "" && email == s.adminEmail {
		role = "gerente"
	} else if role == "" {
		role = "cliente"
	}
	u := &model.User{Name: name, Email: email, Phone: phone, Password: string(hash), Role: role}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(email, password string) (string, *model.User, error) {
	u, err := s.users.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if u == nil {
		return "", nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	claims := &middleware.Claims{
		UserID: u.ID,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}
	return signed, u, nil
}
