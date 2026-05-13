package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ilhamazhar/golang-gpt/internal/domain"
	"github.com/ilhamazhar/golang-gpt/pkg/jwt"
	"github.com/ilhamazhar/golang-gpt/pkg/password"
)

type authService struct {
	users   domain.UserRepository
	access  *jwt.Manager
	refresh *jwt.Manager
}

func NewAuthService(users domain.UserRepository, access, refresh *jwt.Manager) domain.AuthService {
	return &authService{
		users:   users,
		access:  access,
		refresh: refresh,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (domain.UserResponse, error) {
	hash, err := password.Hash(req.Password, password.DefaultParams)
	if err != nil {
		return domain.UserResponse{}, err
	}

	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return domain.UserResponse{}, errors.New("email already registered")
	}

	return domain.ToUserResponse(user), nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.TokenResponse, error) {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("Invalid credentials")
	}

	match, err := password.Verify(req.Password, user.PasswordHash)
	if err != nil || !match {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.access.Generate(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refresh.Generate(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         domain.ToUserResponse(user),
	}, nil
}

func (s *authService) GetProfile(ctx context.Context, id uuid.UUID) (domain.UserResponse, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return domain.UserResponse{}, errors.New("user not found")
	}
	return domain.ToUserResponse(user), nil
}

func (s *authService) ChangePassword(ctx context.Context, id uuid.UUID, req domain.ChangePasswordRequest) error {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	match, err := password.Verify(req.CurrentPassword, user.PasswordHash)
	if err != nil || !match {
		return errors.New("current password is incorrect")
	}

	hash, err := password.Hash(req.NewPassword, password.DefaultParams)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = hash
	return s.users.Update(ctx, user)
}
