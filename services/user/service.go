package user

import (
	"context"
	"errors"

	"github.com/brumble9401/golang-authentication/models"
	"github.com/brumble9401/golang-authentication/repository"
	"github.com/brumble9401/golang-authentication/services/session"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2/log"
)

type Service interface {
    Register(ctx context.Context, payload *models.RegisterUserPayload) error
    Login(ctx context.Context, payload *models.LoginPayload) (string, error)
    LoginByGoogle(ctx context.Context, userData map[string]interface{}) (string, error)
    UpdateUsernameAndPassword(ctx context.Context, payload *models.UsernamePasswordPayload) error
    GetUserByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error)
}

type service struct {
    userRepo repository.UserRepository
    roleRepo repository.RoleRepository
    authProviderRepo repository.AuthProviderRepository
    redisService session.Service
}

func NewService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, authProviderRepo repository.AuthProviderRepository, redisService session.Service) Service {
    return &service{userRepo: userRepo, roleRepo: roleRepo, redisService: redisService, authProviderRepo: authProviderRepo}
}

func (s *service) Register(ctx context.Context, payload *models.RegisterUserPayload) error {
    role, err := s.roleRepo.GetRoleByName(ctx, payload.Role)
    if err != nil {
        return err
    }
    if role == nil {
        return errors.New("role not found")
    }

    userRetrieve, err := s.userRepo.GetUserByEmailOrUsername(ctx, payload.Username, "")
    if userRetrieve != nil || err == nil {
        return errors.New("username already taken")
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

func (s *service) GetUserByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error) {
    return s.userRepo.GetUserByEmailOrUsername(ctx, email, username)
}

func (s *service) LoginByGoogle(ctx context.Context, userData map[string]interface{}) (string, error) {
    role, err := s.roleRepo.GetRoleByName(ctx, "USER")
    if err != nil {
        return "", err
    }
    if role == nil {
        return "", errors.New("role not found")
    }
    user, err := s.userRepo.GetUserByEmailOrUsername(ctx, userData["email"].(string), "")
    if err != nil {
        return "", err
    }
    if user == nil {
        user = &models.User{
            Username: "NONE",
            Email:    userData["email"].(string),
            FullName: userData["family_name"].(string) + " " + userData["given_name"].(string),
            RoleID:   role.RoleID,
        }
        if err := s.userRepo.Register(ctx, user, ""); err != nil {
            return "", err
        }
    }

    userData["user_id"] = user.UserID
    return s.authProviderRepo.CheckAndSaveProviderUserData(ctx, userData)
}

func (s *service) UpdateUsernameAndPassword(ctx context.Context, payload *models.UsernamePasswordPayload) error {
    userID, ok := ctx.Value("userID").(gocql.UUID)
    if !ok {
        return errors.New("user_id not found or invalid")
    }
    if err, ok := ctx.Value("error").(error); ok && err != nil {
        return err
    }
    log.Debugf("Updating username and password for user %s", userID)

    return s.userRepo.UpdateUsernameAndPassword(ctx, payload, userID)
}