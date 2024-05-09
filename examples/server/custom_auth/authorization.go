package custom_auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SampleAuthType = "sample"
)

var (
	rootAuthTypes = map[string]struct{}{}
)

func init() {
	rootAuthTypes = make(map[string]struct{})
	rootAuthTypes["sample"] = struct{}{}
}

func Authorization(ctx context.Context) (context.Context, error) {
	rootAuthType, err := RootAuthTypeFromMD(ctx, rootAuthTypes)
	if err != nil {
		return nil, err
	}

	switch rootAuthType {
	case SampleAuthType:
		return sampleAuth(ctx)
	default:
		return nil, status.Error(codes.Unauthenticated, "invalid root auth type")
	}

}

func sampleAuth(ctx context.Context) (context.Context, error) {
	rootToken, err := RootAuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	token, err := AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	// TODO: This is example only, perform proper Oauth/OIDC verification!
	if rootToken != "yoohoo" && token != "yolo" {
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	}
	// NOTE: You can also pass the token in the context for further interceptors or gRPC service code.
	return ctx, nil
}
