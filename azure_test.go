package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type testTokenCredential struct{}

func (testTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "credential-token"}, nil
}

func TestGetAccessTokenFromCommand(t *testing.T) {
	command := filepath.Join(t.TempDir(), "token-command")
	if err := os.WriteFile(command, []byte("#!/bin/sh\nprintf '  command-token\\n'\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv(accessTokenCommandEnv, command)

	token, err := getAccessToken(context.Background(), nil)
	if err != nil {
		t.Fatalf("getAccessToken returned an error: %v", err)
	}
	if token != "command-token" {
		t.Fatalf("getAccessToken returned %q, want %q", token, "command-token")
	}
}

func TestGetAccessTokenCommandFailure(t *testing.T) {
	command := filepath.Join(t.TempDir(), "token-command")
	if err := os.WriteFile(command, []byte("#!/bin/sh\nexit 7\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv(accessTokenCommandEnv, command)

	if _, err := getAccessToken(context.Background(), nil); err == nil {
		t.Fatal("getAccessToken returned nil error for a failed command")
	}
}

func TestGetAccessTokenUsesCredentialWithoutCommand(t *testing.T) {
	t.Setenv(accessTokenCommandEnv, "")

	token, err := getAccessToken(context.Background(), testTokenCredential{})
	if err != nil {
		t.Fatalf("getAccessToken returned an error: %v", err)
	}
	if token != "credential-token" {
		t.Fatalf("getAccessToken returned %q, want %q", token, "credential-token")
	}
}
