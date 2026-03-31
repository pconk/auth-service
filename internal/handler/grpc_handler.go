package handler

import (
	"auth-service/internal/middleware"
	"auth-service/internal/pb" // Pastikan package ini ada setelah generate protobuf
	"auth-service/internal/service"
	"context"
	"fmt"
)

type AuthGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	Service service.AuthService
}

func NewAuthGRPCHandler(s service.AuthService) *AuthGRPCHandler {
	return &AuthGRPCHandler{Service: s}
}

func (h *AuthGRPCHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.UserResponse, error) {
	claims, err := h.Service.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	middleware.AddUserToLog(ctx, claims.UserID, claims.Username, claims.Role)

	// Get detail user jika perlu, atau cukup return dari claims
	return &pb.UserResponse{
		Id:       claims.UserID,
		Username: claims.Username, // Menambahkan username
		Role:     claims.Role,
	}, nil
}

func (h *AuthGRPCHandler) GetUserProfile(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.Service.GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}

	middleware.AddUserToLog(ctx, req.Id, user.Username, user.Role)

	return &pb.UserResponse{
		Id:       req.Id,
		Username: user.Username,
		Role:     user.Role,
		Position: user.Position,
		Email:    user.Email,
	}, nil
}

func (h *AuthGRPCHandler) GetUsersByIds(ctx context.Context, req *pb.GetUsersRequest) (*pb.UsersResponse, error) {
	if len(req.Ids) == 0 {
		return nil, fmt.Errorf("Invalid Ids - empty")
	}

	users, err := h.Service.GetUsersByIDs(req.Ids)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("Ids not found")
	}

	// middleware.AddUserToLog(ctx, req.Id, user.Username, user.Role)
	results := make([]*pb.UserResponse, len(users))
	for i, user := range users {
		results[i] = &pb.UserResponse{
			Id:       user.ID,
			Username: user.Username,
			Role:     user.Role,
			Position: user.Position,
			Email:    user.Email,
		}
	}

	return &pb.UsersResponse{
		Users: results,
	}, nil
}
