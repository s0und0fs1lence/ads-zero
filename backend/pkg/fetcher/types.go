package fetcher

import (
	"time"

	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

type Client interface {
	WithDbSvc(conn db.DbService) Client
	WithUserId(userId string) Client
	SaveAccountData(data []db.DbAccountSpend) error
	SaveCampaignData(data []db.DbCampaignSpend) error
	IsValid() bool
	GetError() error
	FetchData(start, end time.Time) (task *common.FetchTask, err error)
	ExecuteRules(task common.Event) ([]*common.RuleResult, error)
	GetNotificationChannel() (tp string, value string)
}

type Provider interface {
	//initializers
	withAccessToken(token string) Provider
	withType(pType db.ProviderEnum) Provider
	withID(id string) Provider
	withClientID(id string) Provider
	withAppID(id string) Provider
	withAppSecret(secret string) Provider

	//GetAmountSpent return the amount of spending in the given period, in CENTS, or the encountered error
	FetchData(start, end time.Time) (task *common.FetchTask, err error)
}

type fetchFunc func(accessToken, appId, appSecret string, start, end time.Time) (taskResults *common.FetchTask, err error)
