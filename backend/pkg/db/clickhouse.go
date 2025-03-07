package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/huandu/go-sqlbuilder"
	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
)

const (
	clientsTableName          = "clients"
	providersTableName        = "providers"
	accountsSpendingTableName = "account_spends"
	campaignSpendingTableName = "campaigns_spend"
	rulesTableName            = "client_rules"
)

var (
	clientTable   = sqlbuilder.NewStruct(new(DbClient)).For(sqlbuilder.ClickHouse)
	providerTable = sqlbuilder.NewStruct(new(DbProvider)).For(sqlbuilder.ClickHouse)

	accSpendingTable = sqlbuilder.NewStruct(new(DbAccountSpend)).For(sqlbuilder.ClickHouse)
	rulesTable       = sqlbuilder.NewStruct(new(DbRule)).For(sqlbuilder.ClickHouse)
)

func initConn() (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{configuration.Config().GetString(configuration.ClickhouseHosts)}, //TODO: add proper params
		Auth: clickhouse.Auth{
			Database: configuration.Config().GetString(configuration.ClickhouseDb),
			Username: configuration.Config().GetString(configuration.ClickhouseUsername),
			Password: configuration.Config().GetString(configuration.ClickhousePassword),
		},

		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionZSTD,
			Level:  9,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         30,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "ads-zero-backend", Version: "0.1"},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return conn, conn.Ping(context.Background())
}

type clkService struct {
	ctx  context.Context
	conn clickhouse.Conn
}

func NewClickhouseService(conn clickhouse.Conn) (DbService, error) {

	svc := &clkService{
		ctx:  context.Background(),
		conn: conn,
	}
	if svc.conn == nil {
		conn, err := initConn()
		if err != nil {
			return nil, err
		}
		svc.conn = conn
	}
	return svc, nil
}

// GetAllClients implements DbService.
func (c *clkService) GetAllClients() ([]DbClient, error) {
	sb := sqlbuilder.NewSelectBuilder().Select("*").From(clientsTableName)
	sb.Where(
		sb.EQ("deleted", false),
	)
	q, args := sb.BuildWithFlavor(sqlbuilder.ClickHouse)
	rows, err := c.conn.Query(c.ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]DbClient, 0)
	for rows.Next() {
		var accSpend DbClient
		if err := rows.ScanStruct(&accSpend); err != nil {
			return nil, err
		}
		res = append(res, accSpend)
	}
	return res, nil
}

// GetClientsByQuery implements DbService.
func (c *clkService) GetClientsByQuery(builder *sqlbuilder.SelectBuilder) ([]DbClient, error) {
	q, args := builder.BuildWithFlavor(sqlbuilder.ClickHouse)
	rows, err := c.conn.Query(c.ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]DbClient, 0)
	for rows.Next() {
		var accSpend DbClient
		if err := rows.ScanStruct(&accSpend); err != nil {
			return nil, err
		}
		res = append(res, accSpend)
	}
	return res, nil
}

// GetCampaignSpendGrouped implements DbService.
func (c *clkService) GetCampaignSpendGrouped(clientID string, start time.Time, end time.Time) ([]DbCampaignSpendGrouped, error) {
	sb := sqlbuilder.NewSelectBuilder().Select(
		"client_id", "account_id", "anyLast(account_name) as account_name",
		"business_id", "anyLast(business_name) as business_name",
		"campaign_id", "anyLast(campaign_name) as campaign_name",
		"provider_id", "provider_type", "anyLast(status) as status", "sum(spend) as spend",
		"min(date_ref) as date_start", "max(date_ref) as date_end", "max(updated_at) as updated_at",
	).From(campaignSpendingTableName)
	sb.Where(
		sb.EQ("client_id", clientID),
		sb.GTE("date_ref", start),
		sb.LTE("date_ref", end),
	)
	sb.GroupBy(
		"client_id", "account_id",
		"business_id", "provider_id", "provider_type",
		"campaign_id",
	)
	q, args := sb.Build()
	rows, err := c.conn.Query(c.ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]DbCampaignSpendGrouped, 0)
	for rows.Next() {
		var accSpend DbCampaignSpendGrouped
		if err := rows.ScanStruct(&accSpend); err != nil {
			return nil, err
		}
		res = append(res, accSpend)
	}
	return res, nil
}

// GetAccountSpendGrouped implements DbService.
func (c *clkService) GetAccountSpendGrouped(clientID string, start time.Time, end time.Time) ([]DbAccountSpendGrouped, error) {

	sb := sqlbuilder.NewSelectBuilder().Select(
		"client_id", "account_id", "anyLast(account_name) as account_name",
		"anyLast(account_image) as account_image",
		"business_id", "anyLast(business_name) as business_name",
		"provider_id", "provider_type", "anyLast(status) as status", "sum(spend) as spend",
		"anyLast(number_of_campaigns) as number_of_campaigns",
		"min(date_ref) as date_start", "max(date_ref) as date_end", "max(updated_at) as updated_at",
	).From(accountsSpendingTableName)
	sb.Where(
		sb.EQ("client_id", clientID),
		sb.GTE("date_ref", start),
		sb.LTE("date_ref", end),
	)
	sb.GroupBy(
		"client_id", "account_id",
		"business_id", "provider_id", "provider_type",
	)
	q, args := sb.Build()
	rows, err := c.conn.Query(c.ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]DbAccountSpendGrouped, 0)
	for rows.Next() {
		var accSpend DbAccountSpendGrouped
		if err := rows.ScanStruct(&accSpend); err != nil {
			return nil, err
		}
		res = append(res, accSpend)
	}
	return res, nil
}

// GetCampaignSpend implements DbService.
func (c *clkService) GetCampaignSpend(clientID string, start time.Time, end time.Time) ([]DbCampaignSpend, error) {
	sb := sqlbuilder.NewSelectBuilder().Select("*").From(campaignSpendingTableName)
	sb.Where(
		sb.EQ("client_id", clientID),
		sb.GTE("date_ref", start),
		sb.LTE("date_ref", end),
	)
	q, args := sb.Build()
	rows, err := c.conn.Query(c.ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]DbCampaignSpend, 0)
	for rows.Next() {
		var accSpend DbCampaignSpend
		if err := rows.ScanStruct(&accSpend); err != nil {
			return nil, err
		}
		res = append(res, accSpend)
	}
	return res, nil
}

// InsertCampaignSpend implements DbService.
func (c *clkService) InsertCampaignSpend(data []DbCampaignSpend) error {
	batch, err := c.conn.PrepareBatch(
		c.ctx, fmt.Sprintf("INSERT INTO %s ", campaignSpendingTableName),
		driver.WithCloseOnFlush(), driver.WithReleaseConnection(),
	)
	if err != nil {
		return err
	}
	for _, e := range data {
		if err := batch.AppendStruct(&e); err != nil {
			return err
		}
	}
	return batch.Send()
}

// InsertAccountSpend implements DbService.
func (c *clkService) InsertAccountSpend(data []DbAccountSpend) error {
	// create ephemeral batch
	batch, err := c.conn.PrepareBatch(
		c.ctx, fmt.Sprintf("INSERT INTO %s ", accountsSpendingTableName),
		driver.WithCloseOnFlush(), driver.WithReleaseConnection(),
	)
	if err != nil {
		return err
	}
	for _, e := range data {
		if err := batch.AppendStruct(&e); err != nil {
			return err
		}
	}
	return batch.Send()
}

// GetAccountSpend implements DbService.
func (c *clkService) GetAccountSpend(clientID string, start time.Time, end time.Time) ([]DbAccountSpend, error) {
	sb := sqlbuilder.NewSelectBuilder().Select("*").From(accountsSpendingTableName)
	sb.Where(
		sb.EQ("client_id", clientID),
		sb.GTE("date_ref", start),
		sb.LTE("date_ref", end),
	)
	q, args := sb.Build()
	rows, err := c.conn.Query(c.ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]DbAccountSpend, 0)
	for rows.Next() {
		var accSpend DbAccountSpend
		if err := rows.ScanStruct(&accSpend); err != nil {
			return nil, err
		}
		res = append(res, accSpend)
	}
	return res, nil
}

// UpdateProvider implements DbService.
func (c *clkService) UpdateProvider(providerReq *ProviderUpdate) (*DbProvider, error) {
	provider, err := c.GetProviderByID(providerReq.ProviderID)
	if err != nil {
		return nil, err
	}
	if providerReq.APIAccessToken != nil && *providerReq.APIAccessToken != "" {
		provider.ApiAccessToken = *providerReq.APIAccessToken
	}
	if providerReq.APIClientID != nil && *providerReq.APIClientID != "" {
		provider.ApiClientID = *providerReq.APIClientID
	}
	if providerReq.APIClientSecret != nil && *providerReq.APIClientSecret != "" {
		provider.ApiClientSecret = *providerReq.APIClientSecret
	}

	if err := c.InsertProvider(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

// InsertProvider implements DbService.
func (c *clkService) InsertProvider(provider *DbProvider) error {
	sb := providerTable.InsertInto("providers", provider)
	query, args := sb.Build()

	if err := c.conn.AsyncInsert(c.ctx, query, true, args...); err != nil {
		return err
	}

	return nil
}

// GetProviderByID implements DbService.
func (c *clkService) GetProviderByID(providerID string) (*DbProvider, error) {
	if providerID == "" {
		return nil, fmt.Errorf("invalid provider provided")
	}
	var provider DbProvider
	if err := c.conn.Select(c.ctx, &provider, "select * from providers FINAL where provider_id = ?", providerID); err != nil {
		return nil, err
	}
	return &provider, nil

}

// GetProvidersByClientID implements DbService.

// UpdateClient implements DbService.
func (c *clkService) UpdateClient(clientReq *ClientUpdate) (*DbClient, error) {
	client, err := c.GetClientByID(clientReq.ClientID)
	if err != nil {

		return nil, err
	}

	if clientReq.Email != nil && IsValidEmail(*clientReq.Email) {
		client.UserEmail = *clientReq.Email
	}
	if (clientReq.NotificationEmail != nil) && IsValidEmail(*clientReq.NotificationEmail) {
		client.NotificationEmail = *clientReq.NotificationEmail
	}
	if clientReq.TelegramChatID != nil && *clientReq.TelegramChatID != "" {
		client.TelegramChatID = *clientReq.TelegramChatID
	}
	if clientReq.SlackWebhookURL != nil && *clientReq.SlackWebhookURL != "" {
		client.SlackWebhookURL = *clientReq.SlackWebhookURL
	}

	client.UpdatedAt = time.Now().UTC()

	if err := c.InsertClient(client); err != nil {

		return nil, err
	}
	return client, nil
}

// GetProviders implements DbService.
func (c *clkService) GetProvidersByClientID(clientID string) ([]DbProvider, error) {
	if clientID == "" {
		return nil, fmt.Errorf("invalid client provided")
	}
	var providers []DbProvider
	rows, err := c.conn.Query(c.ctx, "select * from ? FINAL where client_id = ?", providersTableName, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var provider DbProvider
		if err := rows.ScanStruct(&provider); err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

// GetClient implements DbService.
func (c *clkService) GetClientByID(clientId string) (*DbClient, error) {
	var r DbClient
	if err := c.conn.QueryRow(c.ctx, "select * from ? where client_id = ? order by updated_at DESC limit 1", clientsTableName, clientId).ScanStruct(&r); err != nil {
		return nil, err
	}
	return &r, nil

}

// InsertClient implements DbService.
func (c *clkService) InsertClient(client *DbClient) error {
	sb := clientTable.InsertInto(clientsTableName, client)
	query, args := sb.BuildWithFlavor(sqlbuilder.ClickHouse)

	if err := c.conn.AsyncInsert(c.ctx, query, true, args...); err != nil {
		return err
	}

	return nil

}

// GetRuleByID implements DbService.
func (c *clkService) GetRuleByID(ruleID string) (*DbRule, error) {
	var r DbRule
	if err := c.conn.QueryRow(c.ctx, "select * from ? where client_id = ? order by updated_at DESC limit 1", rulesTableName, ruleID).ScanStruct(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetRulesByClientID implements DbService.
func (c *clkService) GetRulesByClientID(clientID string) ([]DbRule, error) {

	if clientID == "" {
		return nil, fmt.Errorf("invalid client provided")
	}
	var rules []DbRule
	rows, err := c.conn.Query(c.ctx, "select * from ? FINAL where client_id = ?", rulesTableName, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rule DbRule
		if err := rows.ScanStruct(&rule); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// InsertRule implements DbService.
func (c *clkService) InsertRule(rule *DbRule) error {
	sb := rulesTable.InsertInto(clientsTableName, rule)
	query, args := sb.BuildWithFlavor(sqlbuilder.ClickHouse)

	if err := c.conn.AsyncInsert(c.ctx, query, true, args...); err != nil {
		return err
	}

	return nil
}
