package repository

import (
	"context"
	"time"

	"github.com/brumble9401/golang-authentication/interfaces"
	"github.com/brumble9401/golang-authentication/models"
	"github.com/brumble9401/golang-authentication/services/session"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
)

type AuthProviderRepository interface {
	CreateAuthProvider(ctx context.Context, authProvider *models.AuthProvider) error
	GetAuthProviderByProviderUserID(ctx context.Context, providerUserID string) (*models.AuthProvider, error)
	CheckAndSaveProviderUserData(ctx context.Context, userData map[string]interface{}) (string, error)
}

func parseTime(timeStr string) time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, timeStr)
	return parsedTime
}

type authProviderRepository struct {
	session *gocql.Session
	queryBuilder interfaces.QueryBuilder
    redisService session.Service
}

func NewAuthProviderRepository(session *gocql.Session, queryBuilder interfaces.QueryBuilder, redisService session.Service) AuthProviderRepository {
	return &authProviderRepository{session: session, queryBuilder: queryBuilder, redisService: redisService}
}

func (r *authProviderRepository) CreateAuthProvider(ctx context.Context, authProvider *models.AuthProvider) error {
	dataMap := structToMap(authProvider)
	query := r.queryBuilder.InsertQuery("auth_providers", dataMap)
	return query.Exec()
}

func (r *authProviderRepository) GetAuthProviderByProviderUserID(ctx context.Context, providerUserID string) (*models.AuthProvider, error) {
	query := r.queryBuilder.SelectConditionQuery("auth_providers", "provider_user_id", providerUserID)
	iter := query.Iter()
	defer iter.Close()

	row := make(map[string]interface{})
	if iter.MapScan(row) {
		return mapToAuthProvider(row), nil
	}
	return nil, nil
}

func (r *authProviderRepository) CheckAndSaveProviderUserData(ctx context.Context, userData map[string]interface{}) (string, error) {
	query := r.queryBuilder.SelectConditionQuery("auth_providers", "provider_user_id", userData["id"].(string))
	iter := query.Iter()
	defer iter.Close()

	var authProvider *models.AuthProvider
	row := make(map[string]interface{})
	log.Debug("Checking if user exists")
	if iter.MapScan(row) {
		log.Debug("Getting auth provider")
		authProvider = mapToAuthProvider(row)
	} else {
		authProvider = &models.AuthProvider{
			AuthProviderID: gocql.TimeUUID(),
			UserID:         userData["user_id"].(gocql.UUID),
			Provider:       "GOOGLE",
			ProviderUserID: userData["id"].(string),
			FamilyName:     userData["family_name"].(string),
			GivenName:      userData["given_name"].(string),
			Email:          userData["email"].(string),
			Picture:        userData["picture"].(string),
			VerifiedEmail:  userData["verified_email"].(bool),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		log.Debug("Creating new auth provider")
		err := r.CreateAuthProvider(ctx, authProvider)
		if err != nil {
			return "", err
		}
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID:   authProvider.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	log.Debug("Creating JWT token")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	log.Debug("Creating session")
	err = r.redisService.CreateSession(ctx, authProvider.UserID.String(), tokenString, time.Hour * 24)
	if err != nil {
		return "", err
	}
	
	log.Debug("Returning token")
	return tokenString, nil
}

func mapToAuthProvider(row map[string]interface{}) *models.AuthProvider {
	return &models.AuthProvider{
		AuthProviderID: row["auth_provider_id"].(gocql.UUID),
		UserID:         row["user_id"].(gocql.UUID),
		Provider:       row["provider"].(string),
		ProviderUserID: row["provider_user_id"].(string),
		FamilyName:     row["family_name"].(string),
		GivenName:      row["given_name"].(string),
		Email:          row["email"].(string),
		Picture:        row["picture"].(string),
		VerifiedEmail:  row["verified_email"].(bool),
		CreatedAt:      parseTime(row["created_at"].(string)),
		UpdatedAt:      parseTime(row["updated_at"].(string)),
	}
}