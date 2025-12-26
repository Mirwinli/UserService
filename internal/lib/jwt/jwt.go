package jwtt

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey string

const (
	emptyValue = 0
	userIDKey  = "userID"
)

func AuthInterceptor(appSecret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is empty")
		}
		authHeader := meta.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization header not found in request")
		}

		tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")

		uid, err := validateToken(tokenStr, appSecret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		newCtx := context.WithValue(ctx, userIDKey, uid)
		return handler(newCtx, req)
	}
}

func GetUserId(ctx context.Context) int64 {
	uid, _ := ctx.Value(userIDKey).(int64)
	return uid
}

func validateToken(tokenString string, appSecret string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appSecret), nil
	})

	if err != nil {
		return emptyValue, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uidRaw, ok := claims["uid"]
		if !ok {
			return emptyValue, errors.New("uid claim not found in token")
		}
		uidFloat, ok := uidRaw.(float64)
		if !ok {
			return emptyValue, errors.New("uid claim is not number")
		}
		return int64(uidFloat), nil
	}
	return emptyValue, fmt.Errorf("invalid token")
}
