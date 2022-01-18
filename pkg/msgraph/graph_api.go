package msgraph

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/shiw13/go-one/pkg/net/httpx"
)

type GetTokenParam struct {
	Scope        string
	RedirectURI  string
	ClientID     string
	ClientSecret string
	Code         string
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GetToken(ctx context.Context, in GetTokenParam) (TokenResponse, error) {
	const uri = "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	params := make(url.Values)
	params.Add("scope", in.Scope)
	params.Add("redirect_uri", in.RedirectURI)
	params.Add("code", in.Code)
	params.Add("client_id", in.ClientID)
	params.Add("client_secret", in.ClientSecret)
	params.Add("grant_type", "authorization_code")

	var res TokenResponse

	if err := httpClient().PostForm(ctx, uri, params, nil, &res); err != nil {
		return res, err
	}

	return res, nil
}

type RefreshTokenParam struct {
	Scope        string
	RedirectURI  string
	ClientID     string
	ClientSecret string
	RefreshToken string
}

func RefreshToken(ctx context.Context, in RefreshTokenParam) (TokenResponse, error) {
	const uri = "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	params := make(url.Values)
	params.Add("scope", in.Scope)
	params.Add("redirect_uri", in.RedirectURI)
	params.Add("client_id", in.ClientID)
	params.Add("client_secret", in.ClientSecret)
	params.Add("refresh_token", in.RefreshToken)
	params.Add("grant_type", "refresh_token")

	var res TokenResponse

	if err := httpClient().PostForm(ctx, uri, params, nil, &res); err != nil {
		return res, err
	}

	return res, nil
}

func GetMailMessages(ctx context.Context, accessToken string) error {
	const uri = "https://graph.microsoft.com/v1.0/me/messages?$select=sender,subject"

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := httpClient().Do(ctx, req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code %d", resp.StatusCode)
	}

	return nil
}

type agent struct {
	client *httpx.Client
}

var _a *agent

func getAgent() *agent {
	return _a
}

func init() {
	cfg := httpx.ClientConfig{
		Timeout:      10 * time.Second,
		DebugEnabled: true,
	}

	_a = &agent{
		client: httpx.NewClient(cfg),
	}
}

func httpClient() *httpx.Client {
	return getAgent().client
}
