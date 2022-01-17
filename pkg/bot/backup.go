package bot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/shiw13/go-one/pkg/logger"
)

type backupData struct {
	Apps []*backupApp `json:"apps"`
}

type backupApp struct {
	ClientID             string    `json:"clientID"`
	ClientSecret         string    `json:"clientSecret"`
	AccessToken          string    `json:"accessToken"`
	RefreshToken         string    `json:"refreshToken"`
	NextInvokeAPITime    time.Time `json:"nextInvokeAPITime"`
	NextRefreshTokenTime time.Time `json:"nextRefreshTokenTime"`
	InvokeAPIFailedCount int       `json:"invokeAPIFailedCount"`
}

func (b *bot) recoverAppFromFile() {
	bs, err := os.ReadFile(b.cfg.DataFileName)
	if err != nil {
		logger.Infof("read data file err:%v", err)
		return
	}

	var data backupData
	if err = json.Unmarshal(bs, &data); err != nil {
		logger.Warnf("unmarshal data file err:%v", err)
		return
	}

	apps := make([]*app, 0, len(data.Apps))
	for _, a := range data.Apps {
		v := &app{
			clientID:             a.ClientID,
			clientSecret:         a.ClientSecret,
			accessToken:          a.AccessToken,
			refreshToken:         a.RefreshToken,
			nextInvokeAPITime:    a.NextInvokeAPITime,
			nextRefreshTokenTime: a.NextRefreshTokenTime,
			invokeAPIFailedCount: a.InvokeAPIFailedCount,
		}
		apps = append(apps, v)
	}

	b.apps = apps
}

func (b *bot) saveAppToFile() {
	var data backupData

	for _, a := range b.apps {
		v := &backupApp{
			ClientID:             a.clientID,
			ClientSecret:         a.clientSecret,
			AccessToken:          a.accessToken,
			RefreshToken:         a.refreshToken,
			NextInvokeAPITime:    a.nextInvokeAPITime,
			NextRefreshTokenTime: a.nextRefreshTokenTime,
			InvokeAPIFailedCount: a.invokeAPIFailedCount,
		}
		data.Apps = append(data.Apps, v)
	}

	bs, err := json.Marshal(&data)
	if err != nil {
		logger.Errorf("marshal data file err:%v", err)
		return
	}

	if err = os.WriteFile(b.cfg.DataFileName, bs, 0644); err != nil {
		logger.Errorf("write data file %s err:%v", err)
		return
	}
}
