package auth

import (
	"context"
	"fmt"
	"time"

	pb "github.com/abdulyazidi/cloudtv/backend/pb/auth"
)

type Handler struct {
	pb.UnimplementedAuthServiceServer
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	fmt.Println("HIT: Auth signup handler")

	signupResponse, err := h.service.Signup(ctx, SignupParams{
		Username:        req.Username,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating a new user: %s", err)
	}

	return &pb.SignupResponse{
		AccessToken: signupResponse.Token,
		ExpiresIn:   int64(time.Until(signupResponse.ExpiresAt).Seconds())}, nil
}
