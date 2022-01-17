package bot

import (
	"context"
	"time"
)

type bot struct {
	ctx           context.Context
	cancelCtxFunc context.CancelFunc
	cfg           *config
	apps          []*app
}

var _b *bot

func InitBot() error {
	c, err := loadConfig()
	if err != nil {
		return err
	}

	ctx, cancelCtxFunc := context.WithCancel(context.Background())
	_b = &bot{
		ctx:           ctx,
		cancelCtxFunc: cancelCtxFunc,
		cfg:           c,
	}

	_b.recoverAppFromFile()

	return nil
}

func StartBot() {
	go getBot().daemonLoop()
}

func StopBot() {
	getBot().cancelCtxFunc()
}

func (b *bot) daemonLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.refreshToken()
			b.invokeGraphAPI()
			b.saveAppToFile()
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *bot) refreshToken() {
	now := time.Now()

	for _, a := range b.apps {
		if a.nextRefreshTokenTime.After(now) {
			continue
		}
		a.refreshAccessToken()
	}
}

func (b *bot) invokeGraphAPI() {
	now := time.Now()

	for _, a := range b.apps {
		if a.nextInvokeAPITime.After(now) {
			continue
		}
		a.invokeAPI()
	}
}

func nextAPIInvokeTime() time.Time {
	return time.Now().Add(time.Duration(_b.cfg.APIInvokeInterval))
}

func getBot() *bot {
	return _b
}
