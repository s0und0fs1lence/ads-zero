package fetcher

import (
	"fmt"
	"time"

	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

type provider struct {
	accessToken  string
	appID        string
	appSecret    string
	providerID   string
	clientID     string
	providerType db.ProviderEnum
	baseEndpoint string
	fetcher      fetchFunc
}

// withClientID implements Provider.
func (p *provider) withClientID(id string) Provider {
	p.clientID = id
	return p
}

func (p *provider) withAppID(id string) Provider {
	p.appID = id
	return p
}

func (p *provider) withAppSecret(secret string) Provider {
	p.appSecret = secret
	return p
}

// withID implements Provider.
func (p *provider) withID(id string) Provider {
	p.providerID = id
	return p
}

// GetAmountSpent implements Provider.
func (p *provider) FetchData(start time.Time, end time.Time) (task *common.FetchTask, err error) {
	if p.fetcher == nil {
		return nil, fmt.Errorf("there isn't any implementation for the current provider")
	}
	task, err = p.fetcher(p.accessToken, p.appID, p.appSecret, start, end)
	if err != nil {
		return nil, err
	}
	// enrich the data
	for idx := range task.Accounts {
		task.Accounts[idx].ProviderID = p.providerID
		task.Accounts[idx].ClientID = p.clientID
	}
	for idx := range task.Campaigns {
		task.Campaigns[idx].ProviderID = p.providerID
		task.Campaigns[idx].ClientID = p.clientID
	}
	return task, nil
}

// WithAccessToken implements Provider.
func (p *provider) withAccessToken(token string) Provider {
	p.accessToken = token
	return p
}

// WithType implements Provider.
func (p *provider) withType(pType db.ProviderEnum) Provider {
	p.providerType = pType
	switch p.providerType {
	case db.Facebook:
		p.baseEndpoint = "https://graph.facebook.com/v21.0/"
		p.fetcher = facebookFetcher
	case db.Google:
		panic("define google base endpoint")
	case db.Taboola:
		panic("define taboola base endpoint")
	}
	return p
}

func newEmptyProvider() Provider {
	return &provider{
		accessToken:  "",
		baseEndpoint: "",
		providerID:   "",
		clientID:     "",
		providerType: db.Empty,
		fetcher:      nil,
	}
}
