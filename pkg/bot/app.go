package bot

import (
	"time"

	"ms365bot/pkg/msgraph"

	"github.com/shiw13/go-one/pkg/logger"
)

type app struct {
	// from config
	clientID     string
	clientSecret string

	// from graph api
	accessToken  string
	refreshToken string

	nextInvokeAPITime    time.Time
	nextRefreshTokenTime time.Time

	invokeAPIFailedCount int
}

func (a *app) refreshAccessToken() {
	param := msgraph.RefreshTokenParam{
		Scope:        getBot().cfg.Scope,
		RedirectURI:  getBot().cfg.RedirectURI,
		ClientID:     a.clientID,
		ClientSecret: a.clientSecret,
		RefreshToken: a.refreshToken,
	}
	res, err := msgraph.RefreshToken(getBot().ctx, param)
	if err != nil {
		logger.Errorf("get refresh token err:%v", err)
		return
	}
	a.accessToken = res.AccessToken
	a.refreshToken = res.RefreshToken
	dur := time.Duration(res.ExpiresIn/2-60) * time.Second
	if dur <= 0 {
		dur = 600 * time.Second
	}
	a.nextRefreshTokenTime = time.Now().Add(dur)
}

func (a *app) getToken(code string) error {
	param := msgraph.GetTokenParam{
		Scope:        getBot().cfg.Scope,
		RedirectURI:  getBot().cfg.RedirectURI,
		ClientID:     a.clientID,
		ClientSecret: a.clientSecret,
		Code:         code,
	}
	res, err := msgraph.GetToken(getBot().ctx, param)
	if err != nil {
		logger.Errorf("get new token err %v", err)
		return err
	}
	a.accessToken = res.AccessToken
	a.refreshToken = res.RefreshToken
	a.nextInvokeAPITime = nextAPIInvokeTime()
	dur := time.Duration(res.ExpiresIn/2-60) * time.Second
	if dur <= 0 {
		dur = 600 * time.Second
	}
	a.nextRefreshTokenTime = time.Now().Add(dur)

	return nil
}

func (a *app) invokeAPI() {
	if err := msgraph.GetMailMessages(getBot().ctx, a.accessToken); err != nil {
		a.invokeAPIFailedCount++
		logger.Errorf("invoke err:%v", err)
		return
	}
	a.nextInvokeAPITime = time.Now().Add(time.Duration(getBot().cfg.APIInvokeInterval) * time.Second)
}
