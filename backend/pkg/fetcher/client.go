package fetcher

import (
	"errors"
	"fmt"
	"time"

	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
	"github.com/s0und0fs1lence/ads-zero/pkg/rule"
)

var (
	notValidClient    = errors.New("not_a_valid_client")
	unitializedClient = errors.New("unitialized_client")
)

type clientInfo struct {
	dbSvc              db.DbService
	user               *db.DbClient
	connectedProviders []Provider
	rules              []rule.Rule
	err                error
}

// GetEmail implements Client.
func (c *clientInfo) GetNotificationChannel() (tp string, value string) {
	if c.user != nil {
		if c.user.NotificationEmail != "" {
			return "email", c.user.NotificationEmail
		}
		if c.user.TelegramChatID != "" {
			return "telegram", c.user.TelegramChatID
		}
		if c.user.SlackWebhookURL != "" {
			return "slack", c.user.SlackWebhookURL
		}
	}
	return "", ""
}

// ExecuteRules implements Client.
func (c *clientInfo) ExecuteRules(task common.Event) ([]*common.RuleResult, error) {
	res := make([]*common.RuleResult, len(c.rules))
	for idx, r := range c.rules {
		ruleRes, err := r.Exec(task)
		if err != nil {
			return nil, err
		}
		res[idx] = &common.RuleResult{
			RuleName:  r.Name(),
			RuleID:    r.Id(),
			Result:    ruleRes,
			Threshold: r.Value(),
		}
	}
	return res, nil
}

// FetchData implements Client.
func (c *clientInfo) FetchData(start time.Time, end time.Time) (task *common.FetchTask, err error) {
	// create a global task wrapper that will have all the data about the underlying tasks
	globalTask := &common.FetchTask{
		Start:     start,
		End:       end,
		Accounts:  make([]db.DbAccountSpend, 0),
		Campaigns: make([]db.DbCampaignSpend, 0),
		Errors:    make([]common.FetchError, 0),
	}
	for _, provider := range c.connectedProviders {
		task, err := provider.FetchData(start, end)
		if err != nil {
			// if this provider has failed in a bad way, we notify the caller about it and continue
			globalTask.Errors = append(globalTask.Errors, common.FetchError{
				EntityType: common.PROVIDER,
				Err:        err,
			})
			continue

		}

		// copy the task results to the global task buffers
		globalTask.Accounts = append(globalTask.Accounts, task.Accounts...)
		globalTask.Campaigns = append(globalTask.Campaigns, task.Campaigns...)
		globalTask.Errors = append(globalTask.Errors, task.Errors...)
	}
	return globalTask, nil
}

// SaveData implements Client.
func (c *clientInfo) SaveAccountData(data []db.DbAccountSpend) error {
	return c.dbSvc.InsertAccountSpend(data)
}

func (c *clientInfo) SaveCampaignData(data []db.DbCampaignSpend) error {
	return c.dbSvc.InsertCampaignSpend(data)
}

// GetError implements Client.
func (c *clientInfo) GetError() error {
	return c.err
}

// IsValid implements Client.
func (c *clientInfo) IsValid() bool {
	return c.err == nil
}

// WithUserId implements Client.
func (c *clientInfo) WithUserId(userId string) Client {
	if c.dbSvc == nil {
		return c
	}
	dbClient, err := c.dbSvc.GetClientByID(userId)
	if err != nil {
		c.err = err
		return c
	}
	c.user = dbClient
	dbProviders, err := c.dbSvc.GetProvidersByClientID(dbClient.ClientID)
	if err != nil {
		c.err = err
		return c
	}
	for _, provider := range dbProviders {
		p := newEmptyProvider().
			withAccessToken(provider.ApiAccessToken).
			withAppID(provider.ApiClientID).
			withAppSecret(provider.ApiClientSecret).
			withType(provider.ProviderType).
			withID(provider.ProviderID).
			withClientID(c.user.ClientID)

		c.connectedProviders = append(c.connectedProviders, p)
	}

	dbRules, err := c.dbSvc.GetRulesByClientID(dbClient.ClientID)
	if err != nil {
		c.err = err
		return c
	}
	for _, r := range dbRules {
		newRule := rule.NewSimpleRule(
			rule.ColumnFromString(r.Column), rule.OperatorFromString(r.Operator),
			r.Value, r.RuleName, r.RuleID,
		)
		c.rules = append(c.rules, newRule)
	}

	//here we should load the client information from the db, and if it's invalid, popolate the error field
	c.err = nil

	return c
}

func (c *clientInfo) WithDbSvc(conn db.DbService) Client {
	if conn == nil {
		c.err = fmt.Errorf("invalid database connection provided")
	}
	c.dbSvc = conn
	return c
}

func NewClient() Client {

	return &clientInfo{
		dbSvc: nil,
		err:   unitializedClient,
	}
}
