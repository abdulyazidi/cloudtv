package auth

import (
	"context"
	"fmt"

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
	fmt.Println("Signup handler is hit.....")
	h.service.Signup(ctx, SignupParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	return &pb.SignupResponse{UserId: "001", Token: "lmfao"}, nil
}
