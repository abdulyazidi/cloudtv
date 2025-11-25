package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/abdulyazidi/cloudtv/backend/internal/db/sqlc"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/argon2"
)

var (
	// Argon2id configurations
	argon2Time    uint32 = 3         // Number of iterations
	argon2Memory  uint32 = 64 * 1024 // Memory in KiB (64 MiB)
	argon2Threads uint8  = 4         // Degree of parallelism
	argon2KeyLen  uint32 = 32        // Length of the generated key (32 bytes = 256 bits)
	saltLen              = 16        // Salt length in bytes (128 bits)
)

var (
	ErrUserAlreadyExists = errors.New("user with this username or email already exists")
	ErrValidation        = errors.New("validation failed")
)

type JWTCustomClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

type SignupParams struct {
	Username        string `validate:"required,min=3,max=32,alphanum"`
	Email           string `validate:"required,email"`
	Password        string `validate:"required,min=8,max=128"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
}

type SignupResult struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

type UserStore interface {
	CreateUser(context.Context, sqlc.CreateUserParams) (sqlc.User, error)
}

type Service struct {
	db        UserStore
	jwtSecret []byte
	validate  *validator.Validate
}

func NewService(db UserStore, jwtSecret []byte) *Service {
	return &Service{
		db:        db,
		jwtSecret: jwtSecret,
		validate:  validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (s *Service) Signup(ctx context.Context, params SignupParams) (SignupResult, error) {
	fmt.Println("HIT: Auth signup service")

	if err := s.validate.Struct(params); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMessages []string
			for _, e := range validationErrors {
				errMessages = append(errMessages, formatValidationError(e))
			}
			return SignupResult{}, fmt.Errorf("%w: %v", ErrValidation, errMessages)
		}
		return SignupResult{}, fmt.Errorf("%w: %s", ErrValidation, err.Error())
	}

	passwordHash, salt, err := hashPassword([]byte(params.Password))
	if err != nil {
		return SignupResult{}, fmt.Errorf("error hashing the password: %s", err)
	}

	user, err := s.db.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        params.Email,
		Username:     params.Username,
		PasswordHash: passwordHash,
		PasswordSalt: salt,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SignupResult{}, ErrUserAlreadyExists
		}
		return SignupResult{}, fmt.Errorf("failed to create a new user: %s", err)
	}
	expiresAt := time.Now().Add(time.Hour)
	claims := JWTCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "cloudtv",
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		Username: user.Username,
	}

	signedToken, err := createJWT(s.jwtSecret, claims)
	if err != nil {
		return SignupResult{}, fmt.Errorf("error creating a JWT: %s", err)
	}

	return SignupResult{
		Token:     signedToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}, nil
}

// returns a new signed JWT -- HS256
func createJWT(key []byte, claims JWTCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("error signing jwt: %s", err)
	}
	return signedToken, nil
}

// hashes the password using argon2id.
// it returns base64 encoded password hash and salt
func hashPassword(password []byte) (passwordHash string, saltString string, err error) {
	salt := make([]byte, saltLen)
	if _, err = rand.Read(salt); err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %s", err)
	}
	hash := argon2.IDKey(password, salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	saltString = base64.StdEncoding.EncodeToString(salt)
	passwordHash = base64.StdEncoding.EncodeToString(hash)
	return passwordHash, saltString, nil
}

// formatValidationError converts a FieldError into a human-readable message
func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", e.Field())
	case "eqfield":
		return fmt.Sprintf("%s must match %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s failed %s validation", e.Field(), e.Tag())
	}
}
