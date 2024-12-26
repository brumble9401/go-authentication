package role

import (
	"context"

	"github.com/brumble9401/golang-authentication/models"
	"github.com/brumble9401/golang-authentication/repository"
	"github.com/gocql/gocql"
)

type Service interface {
    CreateRole(ctx context.Context, role *models.Role) error
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	GetRoleByID(ctx context.Context, id gocql.UUID) (*models.Role, error)
}

type service struct {
    roleRepo repository.RoleRepository
}

func NewService(roleRepo repository.RoleRepository) Service {
    return &service{roleRepo: roleRepo}
}

func (s *service) CreateRole(ctx context.Context, role *models.Role) error {
	return s.roleRepo.CreateRole(ctx, role)
}

func (s *service) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	return s.roleRepo.GetRoleByName(ctx, name)
}

func (s *service) GetRoleByID(ctx context.Context, id gocql.UUID) (*models.Role, error) {
	return s.roleRepo.GetRoleByID(ctx, id)
}