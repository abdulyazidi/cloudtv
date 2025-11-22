package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/abdulyazidi/cloudtv/backend/internal/db/sqlc"
)

type SignupParams struct {
	Username        string
	Email           string
	Password        string
	confirmPassword string
}

type Service struct {
	db *sqlc.Queries
}

func NewService(db *sqlc.Queries) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Signup(ctx context.Context, params SignupParams) bool {
	fmt.Println("AUTH SERVICE SIGNUP HIT!!!!!!!")
	log.Println(params)
	return true
}
