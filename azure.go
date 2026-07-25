package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const armScope = "https://management.azure.com/.default"

const accessTokenCommandEnv = "AZURE_BASTION_TUNNEL_ACCESS_TOKEN_COMMAND"

func getAccessToken(ctx context.Context, cred azcore.TokenCredential) (string, error) {
	if command := os.Getenv(accessTokenCommandEnv); command != "" {
		output, err := exec.CommandContext(ctx, command).Output()
		if err != nil {
			return "", fmt.Errorf("get Azure access token with %s: %w", accessTokenCommandEnv, err)
		}

		token := strings.TrimSpace(string(output))
		if token == "" {
			return "", fmt.Errorf("get Azure access token with %s: command returned empty output", accessTokenCommandEnv)
		}
		return token, nil
	}

	if cred == nil {
		return "", fmt.Errorf("get Azure access token: no token credential configured")
	}

	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{armScope}})
	if err != nil {
		return "", fmt.Errorf("get Azure access token: %w", err)
	}
	return token.Token, nil
}
