package user

import (
	"context"
	"errors"

	"github.com/brumble9401/golang-authentication/models"
	"github.com/brumble9401/golang-authentication/repository"
)

type Service interface {
    Register(ctx context.Context, payload *models.RegisterUserPayload) error
    Login(ctx context.Context, payload *models.LoginPayload) (string, error)
}

type service struct {
    userRepo repository.UserRepository
    roleRepo repository.RoleRepository
}

func NewService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) Service {
    return &service{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *service) Register(ctx context.Context, payload *models.RegisterUserPayload) error {
    role, err := s.roleRepo.GetRoleByName(ctx, payload.Role)
    if err != nil {
        return err
    }
    if role == nil {
        return errors.New("role not found")
    }
    user := &models.User{
        Username: payload.Username,
        Email:    payload.Email,
        FullName: payload.FullName,
        RoleID:   role.RoleID,
    }
    return s.userRepo.Register(ctx, user, payload.Password)
}

func (s *service) Login(ctx context.Context, payload *models.LoginPayload) (string, error) {
    return s.userRepo.Login(ctx, payload.Username, payload.Password)
}