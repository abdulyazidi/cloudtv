package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"

	"github.com/abdulyazidi/cloudtv/backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/argon2"
)

var (
	// Argon2id configuration: m=19456 (19 MiB), t=2, p=1 -- OWASP recc
	// This provides a balanced trade-off between CPU and RAM usage
	argon2Time    uint32 = 3         // Number of iterations
	argon2Memory  uint32 = 64 * 1024 // Memory in KiB (64 MiB)
	argon2Threads uint8  = 4         // Degree of parallelism
	argon2KeyLen  uint32 = 32        // Length of the generated key (32 bytes = 256 bits)
	saltLen              = 16        // Salt length in bytes (128 bits)
)

var (
	ErrUserAlreadyExists = errors.New("user with this username or email already exists")
	ErrInvalidInput      = errors.New("invalid input parameters")
)

type SignupParams struct {
	Username string
	Email    string
	Password string
}

type Service struct {
	db        *sqlc.Queries
	jwtSecret []byte
}

func NewService(db *sqlc.Queries, jwtSecret []byte) *Service {
	return &Service{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Signup(ctx context.Context, params SignupParams) (string, error) {
	log.Println("AUTH SERVICE SIGNUP HIT!!!!!!!")
	log.Println(params)
	if params.Email == "" || params.Username == "" || params.Password == "" {
		return "", ErrInvalidInput
	}

	user, err := s.db.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        params.Email,
		Username:     params.Username,
		PasswordHash: params.Password,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("Error creating a new user, username or email already exists:  %s", err)
			return "", ErrUserAlreadyExists
		}
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	log.Println(user, "<- create user return value")
	passwordHash, salt, err := hashPassword([]byte(params.Password))
	if err != nil {
		fmt.Println("LMFAO ERROR BRO")
	} else {
		log.Printf("passwordhash: %s\nSalt: %s\n", (passwordHash), (salt))
	}

	return user.ID, nil
}

func hashPassword(password []byte) (passwordHash []byte, salt []byte, err error) {
	salt = make([]byte, saltLen)
	if _, err = rand.Read(salt); err != nil {
		return nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	log.Println("salt after rand", salt)
	passwordHash = argon2.IDKey(password, salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	return passwordHash, salt, nil
}
