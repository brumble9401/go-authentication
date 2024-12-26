package models

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
    UserID       gocql.UUID `json:"user_id" db:"user_id"`
    Username     string     `json:"username" db:"username"`
    PasswordHash string     `json:"password_hash" db:"password_hash"`
    Email        string     `json:"email" db:"email"`
    FullName     string     `json:"full_name,omitempty" db:"full_name"`
    RoleID       gocql.UUID `json:"role_id" db:"role_id"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type Role struct {
    RoleID    gocql.UUID `json:"role_id" db:"role_id"`
    RoleName  string     `json:"role_name" db:"role_name"`
    CreatedAt time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type Permission struct {
    PermissionID   gocql.UUID `json:"permission_id" db:"permission_id"`
    PermissionName string     `json:"permission_name" db:"permission_name"`
    CreatedAt      time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type RegisterUserPayload struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
    FullName string `json:"full_name"`
    Role     string `json:"role"`
}

type LoginPayload struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type Claims struct {
    UserID   gocql.UUID `json:"user_id"`
    Username string     `json:"username"`
    jwt.RegisteredClaims
}