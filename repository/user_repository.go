package repository

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/brumble9401/golang-authentication/interfaces"
	"github.com/brumble9401/golang-authentication/models"
	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("5bpehDpA1N0Hj1o+4piTXnRiJVosa9ND7n3QhBZR/cw=")

type UserRepository interface {
    Register(ctx context.Context, user *models.User, password string) error
    Login(ctx context.Context, username, password string) (string, error)
    GetUserByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error)
}

type userRepository struct {
    session      *gocql.Session
    queryBuilder interfaces.QueryBuilder
}

func NewUserRepository(session *gocql.Session, queryBuilder interfaces.QueryBuilder) UserRepository {
    return &userRepository{session: session, queryBuilder: queryBuilder}
}

func (r *userRepository) Register(ctx context.Context, user *models.User, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.PasswordHash = string(hashedPassword)
    user.UserID = gocql.TimeUUID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    dataMap := structToMap(user)
    query := r.queryBuilder.InsertQuery("users", dataMap)
    return query.Exec()
}

func (r *userRepository) Login(ctx context.Context, username, password string) (string, error) {
    query := r.queryBuilder.SelectConditionQuery("users", "username", username)
    iter := query.Iter()
    defer iter.Close()

    row := make(map[string]interface{})
    if iter.MapScan(row) {
        user := mapToUser(row)
        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
            return "", errors.New("invalid credentials")
        }

        expirationTime := time.Now().Add(24 * time.Hour)
        claims := &models.Claims{
            UserID:   user.UserID,
            Username: user.Username,
            RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            return "", err
        }

        return tokenString, nil
    }
    return "", errors.New("user not found")
}

func (r *userRepository) GetUserByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error) {
    var query *gocql.Query
    if email == ""  {
        query = r.queryBuilder.SelectConditionQuery("users", "username", username)
    } else if username == "" {
        query = r.queryBuilder.SelectConditionQuery("users", "email", email)
    }
    iter := query.Iter()
    var user models.User
    if iter.Scan(&user) {
        return &user, nil
    }
    return nil, iter.Close()
}

func structToMap(data interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    val := reflect.ValueOf(data).Elem()
    for i := 0; i < val.NumField(); i++ {
        field := val.Type().Field(i)
        fieldName := field.Tag.Get("db")
        if fieldName == "" {
            fieldName = field.Name
        }
        result[fieldName] = val.Field(i).Interface()
    }
    return result
}

func mapToUser(row map[string]interface{}) models.User {
    return models.User{
        UserID:       row["user_id"].(gocql.UUID),
        Username:     row["username"].(string),
        Email:        row["email"].(string),
        PasswordHash: row["password_hash"].(string),
        CreatedAt:    row["created_at"].(time.Time),
        UpdatedAt:    row["updated_at"].(time.Time),
    }
}