package server

import (
	"context"
	"errors"
	"time"
	"user_service/internal/service"
	"user_service/internal/storage"

	userv1 "github.com/Mirwinli/proto_userService/gen/go/userService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService interface {
	NewUser(ctx context.Context, id int64, username string) error
	GetUser(ctx context.Context, username string) (service.User, error)
	UpdateUser(ctx context.Context, username string, firstname string, lastname string, birthDay time.Time, phone string) error
}

type serverApi struct {
	userv1.UnimplementedUserServiceServer
	userService UserService
}

func Register(grpcServer *grpc.Server, userService UserService) {
	userv1.RegisterUserServiceServer(grpcServer, &serverApi{userService: userService})
}

func (s *serverApi) GetProfileByUsername(ctx context.Context, req *userv1.UserRequest) (*userv1.UserResponse, error) {
	user, err := s.userService.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return &userv1.UserResponse{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		BirthDay:    timestamppb.New(user.Birthday),
		Username:    user.Username,
	}, nil
}

func (s *serverApi) CreateProfile(ctx context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	const op = "server.CreateProfile"

	if err := s.userService.NewUser(ctx, req.GetUserId(), req.Username); err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return &userv1.CreateResponse{}, nil
}

func (s *serverApi) UpdateProfile(ctx context.Context, req *userv1.UpdateRequest) (*userv1.UpdateResponse, error) {
	const op = "server.UpdateProfile"

	if err := s.userService.UpdateUser(
		ctx, req.GetUsername(),
		req.GetFirstName(),
		req.GetLastName(),
		req.GetBirthDay().AsTime(),
		req.GetPhoneNumber(),
	); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return &userv1.UpdateResponse{}, nil
}
