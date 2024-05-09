package custom_auth

import (
	"context"
	"errors"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrEmptyValue = errors.New("empty value")

	headerAuthorize         = "authorization"
	headerRootAuthorizeType = "root-authorization-type"
	headerRootAuthorize     = "root-authorization"
)

func RootAuthTypeFromMD(ctx context.Context, expectedRootAuthType map[string]struct{}) (string, error) {
	val := metadata.ExtractIncoming(ctx).Get(headerRootAuthorizeType)
	if val == "" {
		return "", ErrEmptyValue
	}
	if _, ok := expectedRootAuthType[val]; !ok {
		return "", status.Error(codes.Unauthenticated, "Request unauthenticated with "+val)

	}
	return val, nil
}

func RootAuthFromMD(ctx context.Context, expectedScheme string) (string, error) {
	val := metadata.ExtractIncoming(ctx).Get(headerRootAuthorize)
	if val == "" {
		return "", status.Error(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}
	scheme, token, found := strings.Cut(val, " ")
	if !found {
		return "", status.Error(codes.Unauthenticated, "Bad authorization string")
	}
	if !strings.EqualFold(scheme, expectedScheme) {
		return "", status.Error(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}
	return token, nil
}

func AuthFromMD(ctx context.Context, expectedScheme string) (string, error) {
	val := metadata.ExtractIncoming(ctx).Get(headerAuthorize)
	if val == "" {
		return "", ErrEmptyValue
	}
	scheme, token, found := strings.Cut(val, " ")
	if !found {
		return "", status.Error(codes.Unauthenticated, "Bad authorization string")
	}
	if !strings.EqualFold(scheme, expectedScheme) {
		return "", status.Error(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}
	return token, nil
}
