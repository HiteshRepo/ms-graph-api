package token

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/internal/base"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

func GetTokenByPublicClientApp() (string, error) {
	clientId := os.Getenv("AZURE_CLIENT_ID")
	scopeArr := []string{"User.Read"} // Scopes/permissions required
	tenantId := os.Getenv("TENANT_ID")

	publicClient, err := public.New(clientId, public.WithAuthority(fmt.Sprintf("https://login.microsoftonline.com/%s", tenantId)))
	if err != nil {
		fmt.Println("Error creating public client application:", err)
		return "", err
	}

	var result base.AuthResult

	accounts, err := publicClient.Accounts(context.TODO())
	if err != nil {
		fmt.Println("Error fetching account:", err)
		return "", err
	}
	if len(accounts) > 0 {
		result, err = publicClient.AcquireTokenSilent(context.TODO(), scopeArr, public.WithSilentAccount(accounts[0]))
	}

	if err != nil || len(accounts) == 0 {
		result, err = publicClient.AcquireTokenInteractive(context.TODO(), scopeArr)
		if err != nil {
			fmt.Println("Error creating public client application:", err)
			return "", err
		}
	}

	userAccount := result.Account
	accessToken := result.AccessToken

	_ = userAccount

	return accessToken, nil
}

func GetTokenByConfidentialClientApp() (string, error) {
	clientId := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	scopeArr := []string{"User.Read"} // Scopes/permissions required
	tenantId := os.Getenv("TENANT_ID")

	cred, err := confidential.NewCredFromSecret(clientSecret)
	if err != nil {
		fmt.Println("Error creating secret:", err)
		return "", err
	}
	confidentialClient, err := confidential.New(fmt.Sprintf("https://login.microsoftonline.com/%s", tenantId), clientId, cred)
	if err != nil {
		fmt.Println("Error while creating confidentialClient:", err)
		return "", err
	}

	result, err := confidentialClient.AcquireTokenSilent(context.TODO(), scopeArr)
	if err != nil {
		result, err = confidentialClient.AcquireTokenByCredential(context.TODO(), scopeArr)
		if err != nil {
			fmt.Println("Error creating public client application:", err)
			return "", err
		}
	}
	accessToken := result.AccessToken

	fmt.Println("Access Token:", accessToken)
	return accessToken, nil
}

func SaveTokenToFile(token, filePath string) {
	err := ioutil.WriteFile(filePath, []byte(token), 0644)
	if err != nil {
		fmt.Println("Error saving token to file:", err)
		return
	}

	fmt.Println("Token saved to", filePath)
}

func ReadTokenFromFile(filePath string) (string, error) {
	accessTokenBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading token from file:", err)
		return "", err
	}

	accessToken := string(accessTokenBytes)
	fmt.Println("Access Token:", accessToken)
	return accessToken, nil
}
