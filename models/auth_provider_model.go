package models

import (
	"time"

	"github.com/gocql/gocql"
)

type AuthProvider struct {
    AuthProviderID gocql.UUID `json:"auth_provider_id" db:"auth_provider_id"`
    UserID         gocql.UUID `json:"user_id" db:"user_id"`
    Provider       string     `json:"provider" db:"provider"`
    ProviderUserID string     `json:"provider_user_id" db:"provider_user_id"`
    FamilyName     string     `json:"family_name" db:"family_name"`
    GivenName      string     `json:"given_name" db:"given_name"`
    Email          string     `json:"email" db:"email"`
    Picture        string     `json:"picture" db:"picture"`
    VerifiedEmail  bool       `json:"verified_email" db:"verified_email"`
    CreatedAt      time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

