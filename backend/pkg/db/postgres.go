package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type pgService struct {
	ctx  context.Context
	conn *sqlx.DB
}

// GetRuleByID implements DbService.
func (p *pgService) GetRuleByID(ruleID string) (*DbRule, error) {
	panic("unimplemented")
}

// GetRulesByClientID implements DbService.
func (p *pgService) GetRulesByClientID(clientID string) ([]DbRule, error) {
	panic("unimplemented")
}

// InsertRule implements DbService.
func (p *pgService) InsertRule(rule *DbRule) error {
	panic("unimplemented")
}

// GetAllClients implements DbService.
func (p *pgService) GetAllClients() ([]DbClient, error) {
	panic("unimplemented")
}

// GetClientsByQuery implements DbService.
func (p *pgService) GetClientsByQuery(builder *sqlbuilder.SelectBuilder) ([]DbClient, error) {
	panic("unimplemented")
}

// GetCampaignSpendGrouped implements DbService.
func (p *pgService) GetCampaignSpendGrouped(clientID string, start time.Time, end time.Time) ([]DbCampaignSpendGrouped, error) {
	panic("unimplemented")
}

// GetAccountSpendGrouped implements DbService.
func (p *pgService) GetAccountSpendGrouped(clientID string, start time.Time, end time.Time) ([]DbAccountSpendGrouped, error) {
	panic("unimplemented")
}

// GetCampaignSpend implements DbService.
func (p *pgService) GetCampaignSpend(clientID string, start time.Time, end time.Time) ([]DbCampaignSpend, error) {
	panic("unimplemented")
}

// InsertCampaignSpend implements DbService.
func (p *pgService) InsertCampaignSpend(data []DbCampaignSpend) error {
	panic("unimplemented")
}

// InsertAccountSpend implements DbService.
func (p *pgService) InsertAccountSpend(data []DbAccountSpend) error {
	panic("unimplemented")
}

// GetAccountSpend implements DbService.
func (p *pgService) GetAccountSpend(clientID string, start time.Time, end time.Time) ([]DbAccountSpend, error) {
	panic("unimplemented")
}

// GetClientByID implements DbService.
func (p *pgService) GetClientByID(clientId string) (*DbClient, error) {

	var res DbClient
	tx, err := p.conn.BeginTxx(p.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Commit()
	if err := tx.GetContext(p.ctx, &res, "select * from clients where client_id=$1", clientId); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetProviderByID implements DbService.
func (p *pgService) GetProviderByID(providerID string) (*DbProvider, error) {
	panic("unimplemented")
}

// GetProvidersByClientID implements DbService.
func (p *pgService) GetProvidersByClientID(clientId string) ([]DbProvider, error) {
	panic("unimplemented")
}

// InsertClient implements DbService.
func (p *pgService) InsertClient(client *DbClient) error {
	tx, err := p.conn.BeginTxx(p.ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		return err
	}
	q := `
	insert into clients (client_id, user_email, user_type, stripe_customer_id, stripe_subscription_id, stripe_subscription_status, notification_email, telegram_chat_id, slack_webhook_url)
	values (:client_id,:user_email,:user_type,:stripe_customer_id,:stripe_subscription_id,:stripe_subscription_status,:notification_email,:telegram_chat_id,:slack_webhook_url);
	`
	_, err = tx.NamedExecContext(p.ctx, q, client)
	if err != nil {
		return err
	}
	// if res.RowsAffected() <1 {
	// 	return fmt.Errorf("the row could not be inserted")
	// }
	return nil
}

// InsertProvider implements DbService.
func (p *pgService) InsertProvider(provider *DbProvider) error {
	panic("unimplemented")
}

// UpdateClient implements DbService.
func (p *pgService) UpdateClient(clientReq *ClientUpdate) (*DbClient, error) {
	panic("unimplemented")
}

// UpdateProvider implements DbService.
func (p *pgService) UpdateProvider(providerReq *ProviderUpdate) (*DbProvider, error) {
	panic("unimplemented")
}

func NewPostgreService() (DbService, error) {
	p := pgService{
		ctx: context.Background(),
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"172.17.0.1", 5432, "postgres", "admin_pg1997", "adszero")
	conn, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// sqlx.Connect()
	// conn, err := pgxpool.New(p.ctx, "postgres://postgres:admin_pg1997@172.17.0.1:5432/adszero?pool_max_conns=20")
	// if err != nil {
	// 	return nil, err
	// }
	p.conn = conn
	// types, err := p.conn.LoadTypes(p.ctx, []string{
	// 	"USERTYPE",
	// 	"SUB_STATUS",
	// 	"PROVIDER_TYPE",
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// p.conn.TypeMap().RegisterTypes(types)
	return &p, nil
}
