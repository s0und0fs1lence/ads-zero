package db

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Status uint8

const (
	UnknownStatus Status = iota
	Active
	Inactive
)

func StatusFromString(val string) Status {
	switch strings.ToUpper(val) {
	case "ACTIVE":
		return Active
	case "INACTIVE":
		return Inactive
	default:
		return Inactive //TODO: keep an eye on this
	}
}

func (s Status) String() string {
	switch s {
	case UnknownStatus:
		return "UNKNOWN"
	case Active:
		return "ACTIVE"
	case Inactive:
		return "INACTIVE"

	default:
		panic("unreachable: ProviderEnum.ToString()") //TODO: should not be reachable
	}
}

func (s *Status) Scan(src any) error {
	if t, ok := src.(string); ok {
		switch strings.ToUpper(t) {
		case "UNKNOWN":
			*s = UnknownStatus
		case "ACTIVE":
			*s = Active
		case "INACTIVE":
			*s = Inactive
		default:
			return fmt.Errorf("cannot scan %s into customStr", t)
		}
		return nil
	}
	if t, ok := src.([]uint8); ok {
		switch strings.ToUpper(string(t)) {
		case "UNKNOWN":
			*s = UnknownStatus
		case "ACTIVE":
			*s = Active
		case "INACTIVE":
			*s = Inactive
		default:
			return fmt.Errorf("cannot scan %s into customStr", t)
		}
		return nil

	}
	return fmt.Errorf("cannot scan %T into customStr", src)
}

func (s *Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var unescaped string
	if err := json.Unmarshal(data, &unescaped); err != nil {
		return err
	}
	*s = StatusFromString(unescaped)
	return nil
}

type ProviderEnum uint8

const (
	Empty ProviderEnum = iota
	Facebook
	Google
	TikTok
	Taboola
	Invalid
)

func (s ProviderEnum) String() string {
	switch s {
	case Empty:
		return "EMPTY"
	case Facebook:
		return "FACEBOOK"
	case Google:
		return "GOOGLE"
	case TikTok:
		return "TIKTOK"
	case Taboola:
		return "TABOOLA"
	default:
		panic("unreachable: ProviderEnum.ToString()") //TODO: should not be reachable
	}
}

func ProviderFromString(val string) ProviderEnum {
	switch strings.ToUpper(val) {
	case "EMTPY":
		return Empty
	case "FACEBOOK":
		return Facebook
	case "GOOGLE":
		return Google
	case "TIKTOK":
		return TikTok
	case "TABOOLA":
		return Taboola
	default:
		return Invalid //TODO: keep an eye on this
	}
}

func (s *ProviderEnum) Scan(src any) error {
	if t, ok := src.(string); ok {
		*s = ProviderFromString(t)
		return nil
	}
	return fmt.Errorf("cannot scan %T into customStr", src)
}

func (s *ProviderEnum) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *ProviderEnum) UnmarshalJSON(data []byte) error {
	var unescaped string
	if err := json.Unmarshal(data, &unescaped); err != nil {
		return err
	}
	*s = ProviderFromString(unescaped)
	return nil
}

type DbClient struct {
	ClientID          string    `ch:"client_id" json:"client_id" db:"client_id" `
	UserEmail         string    `ch:"user_email" json:"user_email" db:"user_email"`
	NotificationEmail string    `ch:"notification_email" json:"notification_email" db:"notification_email"`
	TelegramChatID    string    `ch:"telegram_chat_id" json:"telegram_chat_id" db:"telegram_chat_id"`
	SlackWebhookURL   string    `ch:"slack_webhook_url" json:"slack_webhook_url" db:"slack_webhook_url"`
	InsertedAt        time.Time `ch:"inserted_at" json:"inserted_at" db:"inserted_at"`
	UpdatedAt         time.Time `ch:"updated_at" json:"updated_at" db:"updated_at"`
	Deleted           bool      `ch:"deleted" json:"-"`
}

type DbProvider struct {
	ProviderID      string       `ch:"provider_id" json:"provider_id" db:"provider_id" fieldtag:"provider_id"`
	ProviderType    ProviderEnum `ch:"provider_type" json:"provider_type" db:"provider_type"`
	ClientID        string       `ch:"client_id" json:"client_id" db:"client_id"`
	InsertedAt      time.Time    `ch:"inserted_at" json:"inserted_at" db:"inserted_at"`
	ApiClientID     string       `ch:"api_client_id" json:"api_client_id" db:"api_client_id"`
	ApiClientSecret string       `ch:"api_client_secret" json:"api_client_secret" db:"api_client_secret"`
	ApiAccessToken  string       `ch:"api_access_token" json:"api_access_token" db:"api_access_token"`
}

type ClientCreate struct {
	ClientID string `json:"client_id"`
	Email    string `json:"email"`
}

func (c *ClientCreate) IsValid() bool {
	return c.ClientID != "" && IsValidEmail(c.Email)
}

type ClientUpdate struct {
	ClientID                 string  `json:"client_id"`
	Email                    *string `json:"email"`
	StripeCustomerID         *string `json:"stripe_customer_id"`
	StripeSubscriptionID     *string `json:"stripe_subscription_id"`
	StripeSubscriptionStatus *Status `json:"stripe_subscription_status"`
	NotificationEmail        *string `json:"notification_email"`
	TelegramChatID           *string `json:"telegram_chat_id"`
	SlackWebhookURL          *string `json:"slack_webhook_url"`
}

func IsValidEmail(email string) bool {
	// Regular expression for validating an email
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

type ProviderCreate struct {
	ProviderType    string `json:"provider_type"`
	ClientID        string `json:"client_id"`
	APIClientID     string `json:"api_client_id"`
	APIClientSecret string `json:"api_client_secret"`
	APIAccessToken  string `json:"api_access_token"`
}

func (p *ProviderCreate) IsValid() bool {
	return p.ProviderType != "" && p.ClientID != "" &&
		p.APIAccessToken != "" && p.APIClientID != "" && p.APIClientSecret != ""
}

func (p ProviderCreate) AsDbProvider() DbProvider {
	return DbProvider{
		ProviderType:    ProviderFromString(p.ProviderType),
		ClientID:        p.ClientID,
		ApiClientID:     p.APIClientID,
		ApiClientSecret: p.APIClientSecret,
		ApiAccessToken:  p.APIAccessToken,
	}
}

type ProviderUpdate struct {
	ProviderID      string  `json:"provider_id"`
	APIClientID     *string `json:"api_client_id"`
	APIClientSecret *string `json:"api_client_secret"`
	APIAccessToken  *string `json:"api_access_token"`
}

type DbAccountSpend struct {
	ClientID          string       `ch:"client_id" json:"client_id"`
	AccountID         string       `ch:"account_id" json:"account_id"`
	AccountName       string       `ch:"account_name" json:"account_name"`
	AccountImage      string       `ch:"account_image" json:"account_image"`
	BusinessID        string       `ch:"business_id" json:"business_id"`
	BusinessName      string       `ch:"business_name" json:"business_name"`
	ProviderID        string       `ch:"provider_id" json:"provider_id"`
	ProviderType      ProviderEnum `ch:"provider_type" json:"provider_type"`
	Status            string       `ch:"status" json:"status"`
	Spend             float64      `ch:"spend" json:"spend"`
	NumberOfCampaigns uint16       `ch:"number_of_campaigns" json:"number_of_campaigns"`
	DateRef           time.Time    `ch:"date_ref" json:"date_ref"`
	UpdatedAt         time.Time    `ch:"updated_at" json:"updated_at"`
}

type DbAccountSpendGrouped struct {
	ClientID          string       `ch:"client_id" json:"client_id"`
	AccountID         string       `ch:"account_id" json:"account_id"`
	AccountName       string       `ch:"account_name" json:"account_name"`
	AccountImage      string       `ch:"account_image" json:"account_image"`
	BusinessID        string       `ch:"business_id" json:"business_id"`
	BusinessName      string       `ch:"business_name" json:"business_name"`
	ProviderID        string       `ch:"provider_id" json:"provider_id"`
	ProviderType      ProviderEnum `ch:"provider_type" json:"provider_type"`
	Status            string       `ch:"status" json:"status"`
	Spend             float64      `ch:"spend" json:"spend"`
	NumberOfCampaigns uint16       `ch:"number_of_campaigns" json:"number_of_campaigns"`
	DateStart         time.Time    `ch:"date_start" json:"date_start"`
	DateEnd           time.Time    `ch:"date_end" json:"date_end"`
	UpdatedAt         time.Time    `ch:"updated_at" json:"updated_at"`
}

type ClientSpendRequest struct {
	ClientID string    `form:"client_id"`
	Start    time.Time `form:"start" time_format:"2006-01-02"`
	End      time.Time `form:"end" time_format:"2006-01-02"`
}

type DbCampaignSpend struct {
	ClientID     string       `ch:"client_id" json:"client_id"`
	AccountID    string       `ch:"account_id" json:"account_id"`
	AccountName  string       `ch:"account_name" json:"account_name"`
	AccountImage string       `ch:"account_image" json:"account_image"`
	BusinessID   string       `ch:"business_id" json:"business_id"`
	BusinessName string       `ch:"business_name" json:"business_name"`
	CampaignID   string       `ch:"campaign_id" json:"campaign_id"`
	CampaignName string       `ch:"campaign_name" json:"campaign_name"`
	ProviderID   string       `ch:"provider_id" json:"provider_id"`
	ProviderType ProviderEnum `ch:"provider_type" json:"provider_type"`
	Status       string       `ch:"status" json:"status"`
	Spend        float64      `ch:"spend" json:"spend"`
	DateRef      time.Time    `ch:"date_ref" json:"date_ref"`
	UpdatedAt    time.Time    `ch:"updated_at" json:"updated_at"`
}

type DbCampaignSpendGrouped struct {
	ClientID     string       `ch:"client_id" json:"client_id"`
	AccountID    string       `ch:"account_id" json:"account_id"`
	AccountName  string       `ch:"account_name" json:"account_name"`
	AccountImage string       `ch:"account_image" json:"account_image"`
	BusinessID   string       `ch:"business_id" json:"business_id"`
	BusinessName string       `ch:"business_name" json:"business_name"`
	CampaignID   string       `ch:"campaign_id" json:"campaign_id"`
	CampaignName string       `ch:"campaign_name" json:"campaign_name"`
	ProviderID   string       `ch:"provider_id" json:"provider_id"`
	ProviderType ProviderEnum `ch:"provider_type" json:"provider_type"`
	Status       string       `ch:"status" json:"status"`
	Spend        float64      `ch:"spend" json:"spend"`
	DateStart    time.Time    `ch:"date_start" json:"date_start"`
	DateEnd      time.Time    `ch:"date_end" json:"date_end"`
	UpdatedAt    time.Time    `ch:"updated_at" json:"updated_at"`
}

type DbRule struct {
	RuleID          string    `ch:"rule_id" json:"rule_id"`
	ClientID        string    `ch:"client_id" json:"client_id"`
	RuleName        string    `ch:"rule_name" json:"rule_name"`
	Column          string    `ch:"column" json:"column"`
	Operator        string    `ch:"operator" json:"operator"`
	Value           float64   `ch:"value" json:"value"`
	NotificationWay string    `ch:"notification_way" json:"notification_way"`
	InsertedAt      time.Time `ch:"inserted_at" json:"inserted_at"`
	UpdatedAt       time.Time `ch:"updated_at" json:"updated_at"`
}
