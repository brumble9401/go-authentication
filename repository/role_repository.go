package repository

import (
	"context"
	"time"

	"github.com/brumble9401/golang-authentication/interfaces"
	"github.com/brumble9401/golang-authentication/models"
	"github.com/gocql/gocql"
)

type RoleRepository interface {
	CreateRole(ctx context.Context, role *models.Role) error;
	GetRoleByName(ctx context.Context, name string) (*models.Role, error);
	GetRoleByID(ctx context.Context, id gocql.UUID) (*models.Role, error);
}

type roleRepository struct {
    session      *gocql.Session
    queryBuilder interfaces.QueryBuilder
}

func NewRoleRepository(session *gocql.Session, queryBuilder interfaces.QueryBuilder) RoleRepository {
    return &roleRepository{session: session, queryBuilder: queryBuilder}
}

func (r *roleRepository) CreateRole(ctx context.Context, role *models.Role) error {
	dataMap := structToMap(role)
	query := r.queryBuilder.InsertQuery("roles", dataMap)
	return query.Exec()
}

func (r *roleRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	query := r.queryBuilder.SelectConditionQuery("roles", "role_name", name)
	iter := query.Iter()
	defer iter.Close()

	row := make(map[string]interface{})
	if iter.MapScan(row) {
		return mapToRole(row), nil
	}
	return nil, nil
}

func (r *roleRepository) GetRoleByID(ctx context.Context, id gocql.UUID) (*models.Role, error) {
	query := r.queryBuilder.SelectQuery("roles", id)
	iter := query.Iter()
	defer iter.Close()

	row := make(map[string]interface{})
	if iter.MapScan(row) {
		return mapToRole(row), nil
	}
	return nil, nil
}

func mapToRole(row map[string]interface{}) *models.Role {
	return &models.Role{
		RoleID:       row["role_id"].(gocql.UUID),
		RoleName:     row["role_name"].(string),
		CreatedAt:    row["created_at"].(time.Time),
		UpdatedAt:    row["updated_at"].(time.Time),
	}
}