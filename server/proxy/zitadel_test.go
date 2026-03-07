package proxy

import (
	"context"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
)

func TestVerifyTokenWithBackendIdentity(t *testing.T) {
	prepareEnv()
	accessToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjM2MTc1ODYxMjEwNTg0MjYzOSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2NlcmFtaS10NmlocmQudXMxLnppdGFkZWwuY2xvdWQiLCJzdWIiOiIzNjE5NzYwNDUyNjEzMTA5NzQiLCJhdWQiOlsiMzYxNzYxNDI5MzAyMzczMDgyIiwiMzYxOTc2MjUwMjExODQ3MTY2IiwiMzYxNzU5Nzg4ODI2MTUyOTExIl0sImV4cCI6MTc3MjI5Njk5MSwiaWF0IjoxNzcyMjUzNzkxLCJuYmYiOjE3NzIyNTM3OTEsImNsaWVudF9pZCI6IjM2MTc2MTQyOTMwMjM3MzA4MiIsImp0aSI6IlYyXzM2MTk4MTUwNzU3MDU4NDU3NC1hdF8zNjE5ODE1MDc1NzA2NTAxMTAifQ.iG0exsQ4X0c7KnNoxJagnOXloGxOyLqb36_ncPWf8GPwx5fpngSAI4-O20f2YOeR_BxEsq_OqCE3rze8VfHLwcxd4EoxSKmY16OtZ_1e0mh7ZrWvRFuL8VUPPXUvUpaAWDs3XqVEKrvBdHWCYlrvtkFYk8IollbAGmS0PpsYzKeEEgcmtp61s581QMqw7TcQnu__3FU20toQWx-lZ2wzEW2Mb5FVFnXv2J1GXjuNUng0HqjxHxED8XxWX25ztI5Zjci0nBOKGsN1cQnEeKqe0WNGwI44dB2OyybjeaHu8oPOSQP9O1NwzrXNgBqH3A8IKDfR7PTnJH8BRjSfDBYYTw"
	zitadelProxy := GetZitadelProxy()
	user, err := zitadelProxy.VerifyTokenWithBackendIdentity(context.Background(), accessToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.ZitadelSub == "" {
		t.Errorf("Expected 'sub' claim to be present")
	}
	fmt.Println(user)
}

func TestValidateJwtToken(t *testing.T) {
	accessToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjM2MTc1ODYxMjEwNTg0MjYzOSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2NlcmFtaS10NmlocmQudXMxLnppdGFkZWwuY2xvdWQiLCJzdWIiOiIzNjE5NzYwNDUyNjEzMTA5NzQiLCJhdWQiOlsiMzYxNzYxNDI5MzAyMzczMDgyIiwiMzYxOTc2MjUwMjExODQ3MTY2IiwiMzYxNzU5Nzg4ODI2MTUyOTExIl0sImV4cCI6MTc3MjI5Njk5MSwiaWF0IjoxNzcyMjUzNzkxLCJuYmYiOjE3NzIyNTM3OTEsImNsaWVudF9pZCI6IjM2MTc2MTQyOTMwMjM3MzA4MiIsImp0aSI6IlYyXzM2MTk4MTUwNzU3MDU4NDU3NC1hdF8zNjE5ODE1MDc1NzA2NTAxMTAifQ.iG0exsQ4X0c7KnNoxJagnOXloGxOyLqb36_ncPWf8GPwx5fpngSAI4-O20f2YOeR_BxEsq_OqCE3rze8VfHLwcxd4EoxSKmY16OtZ_1e0mh7ZrWvRFuL8VUPPXUvUpaAWDs3XqVEKrvBdHWCYlrvtkFYk8IollbAGmS0PpsYzKeEEgcmtp61s581QMqw7TcQnu__3FU20toQWx-lZ2wzEW2Mb5FVFnXv2J1GXjuNUng0HqjxHxED8XxWX25ztI5Zjci0nBOKGsN1cQnEeKqe0WNGwI44dB2OyybjeaHu8oPOSQP9O1NwzrXNgBqH3A8IKDfR7PTnJH8BRjSfDBYYTw"
	zitadelProxy := GetZitadelProxy()
	user, err := zitadelProxy.ValidateJwtToken(accessToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	fmt.Println(user)
}

func TestSyncMetaData(t *testing.T) {
	prepareEnv()
	zitadelProxy := GetZitadelProxy()
	user := &model.User{
		ZitadelSub: "361969871380040702",
		ID:         12,
	}
	err := zitadelProxy.SyncMeta2Zitadel(context.Background(), user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGenAccessToken(t *testing.T) {
	prepareEnv()
	zitadelProxy := GetZitadelProxy()
	token, err := zitadelProxy.getActualAccessToken()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	fmt.Println(token)
	token, err = zitadelProxy.getActualAccessToken()
	fmt.Println(token)
}

func prepareEnv() {
	config.Config = &config.Conf{
		ZitadelConfig: &config.ZitadelConfig{
			Host: "https://cerami-t6ihrd.us1.zitadel.cloud",
		},
	}
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Warning: No .env file found, relying on system env")
	}
	InitZitadel()
}
