package user

import (
	"context"
	"errors"
	"time"

	"github.com/brumble9401/golang-authentication/models"
	"github.com/brumble9401/golang-authentication/repository"
	"github.com/brumble9401/golang-authentication/services/session"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
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
    // queryBuilder interfaces.QueryBuilder
}

var jwtKey = []byte("5bpehDpA1N0Hj1o+4piTXnRiJVosa9ND7n3QhBZR/cw=")

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

// func (s *service) LoginByGoogle(ctx context.Context, userData map[string]interface{}) (string, error) {
//     role, err := s.roleRepo.GetRoleByName(ctx, "USER")
//     if err != nil {
//         return "", err
//     }
//     if role == nil {
//         return "", errors.New("role not found")
//     }

//     user, err := s.userRepo.GetUserByEmailOrUsername(ctx, userData["email"].(string), "")
//     if err != nil {
//         return "", err
//     }
//     if user == nil {
//         user = &models.User{
//             Username: "NONE",
//             Email:    userData["email"].(string),
//             FullName: userData["family_name"].(string) + " " + userData["given_name"].(string),
//             RoleID:   role.RoleID,
//         }
//         if err := s.userRepo.Register(ctx, user, ""); err != nil {
//             return "", err
//         }
//     }

//     userData["user_id"] = user.UserID
//     return s.authProviderRepo.CheckAndSaveProviderUserData(ctx, userData)
// }
func (s *service) LoginByGoogle(ctx context.Context, userData map[string]interface{}) (string, error) {
    role, err := s.roleRepo.GetRoleByName(ctx, "USER")
    if err != nil {
        return "", err
    }
    if role == nil {
        return "", errors.New("role not found")
    }

    // Start a batch
    batch := s.userRepo.NewBatch(gocql.LoggedBatch)

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
        // Add user registration to the batch
        s.userRepo.RegisterBatch(ctx, batch, user, "")
    }

    userData["user_id"] = user.UserID
    // Add auth provider data to the batch
    s.authProviderRepo.CheckAndSaveProviderUserDataBatch(ctx, batch, userData)

    // Apply the batch
    if err := s.userRepo.ExecuteBatch(batch); err != nil {
        return "", err
    }

    // Generate token
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &models.Claims{
        UserID:   user.UserID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    // Create session
    err = s.redisService.CreateSession(ctx, user.UserID.String(), tokenString, time.Hour * 24)
    if err != nil {
        return "", err
    }

    return tokenString, nil
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