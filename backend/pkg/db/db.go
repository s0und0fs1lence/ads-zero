package db

import (
	"time"

	"github.com/huandu/go-sqlbuilder"
)

type DbService interface {
	//client
	InsertClient(client *DbClient) error
	UpdateClient(clientReq *ClientUpdate) (*DbClient, error)
	GetClientByID(clientId string) (*DbClient, error)
	GetAllClients() ([]DbClient, error)
	GetClientsByQuery(builder *sqlbuilder.SelectBuilder) ([]DbClient, error)

	//provider
	InsertProvider(provider *DbProvider) error
	GetProvidersByClientID(clientId string) ([]DbProvider, error)
	GetProviderByID(providerID string) (*DbProvider, error)
	UpdateProvider(providerReq *ProviderUpdate) (*DbProvider, error)

	//account spend
	GetAccountSpend(clientID string, start, end time.Time) ([]DbAccountSpend, error)
	GetAccountSpendGrouped(clientID string, start, end time.Time) ([]DbAccountSpendGrouped, error)
	InsertAccountSpend(data []DbAccountSpend) error

	//campaign spend
	GetCampaignSpend(clientID string, start, end time.Time) ([]DbCampaignSpend, error)
	GetCampaignSpendGrouped(clientID string, start, end time.Time) ([]DbCampaignSpendGrouped, error)
	InsertCampaignSpend(data []DbCampaignSpend) error

	//rule
	InsertRule(rule *DbRule) error
	GetRulesByClientID(clientID string) ([]DbRule, error)
	GetRuleByID(ruleID string) (*DbRule, error)
}
