/* Database creation */
CREATE DATABASE IF NOT EXISTS adszero ENGINE = Atomic;

/* Table definition */


CREATE TABLE adszero.clients (
    client_id String NOT NULL default generateULID(),
    user_email String NOT NULL,
    notification_email String,
    telegram_chat_id String,
    slack_webhook_url String,
    inserted_at DateTime64(9) default now64(9),
    updated_at DateTime64(9) default now64(9),
    deleted Bool default false
)
ENGINE=ReplacingMergeTree(updated_at)
ORDER BY (client_id);



CREATE TABLE adszero.providers (
    provider_id FixedString(26) default generateULID(),
    provider_type Enum8('INVALID'=0,'FACEBOOK'=1,'GOOGLE'=2,'TIKTOK'=3,'TABOOLA'=4) default 'INVALID',
    client_id String NOT NULL,
    inserted_at DateTime64(9) default now64(9),
    api_client_id String,
    api_client_secret String,
    api_access_token String
)
ENGINE=ReplacingMergeTree(inserted_at)
ORDER BY (provider_id,client_id);


CREATE TABLE adszero.account_spends (
    client_id String NOT NULL,
    account_id String NOT NULL,
    account_name String NOT NULL,
    account_image String,
    business_id String NOT NULL,
    business_name String NOT NULL,
    provider_id FixedString(26) NOT NULL,
    provider_type Enum8('INVALID'=0,'FACEBOOK'=1,'GOOGLE'=2,'TIKTOK'=3,'TABOOLA'=4) default 'INVALID',
    status Enum8('UNKNOWN'=0,'ACTIVE'=1,'INACTIVE'=2) default 'UNKNOWN',
    spend Float64,
    number_of_campaigns UInt16,
    date_ref Date32 default now(),
    updated_at DateTime64(9) default now64(9)
)
ENGINE=ReplacingMergeTree(updated_at)
ORDER BY (account_id,client_id,date_ref)
partition by toMonth(date_ref);

CREATE TABLE adszero.campaigns_spend (
    client_id String NOT NULL,
    account_id String NOT NULL,
    account_name String NOT NULL,
    business_id String NOT NULL,
    business_name String NOT NULL,
    campaign_id String NOT NULL,
    campaign_name String NOT NULL,
    provider_id FixedString(26) NOT NULL,
    provider_type Enum8('INVALID'=0,'FACEBOOK'=1,'GOOGLE'=2,'TIKTOK'=3,'TABOOLA'=4) default 'INVALID',
    status Enum8('UNKNOWN'=0,'ACTIVE'=1,'INACTIVE'=2) default 'UNKNOWN',
    spend Float64,
    date_ref Date32 default now(),
    updated_at DateTime64(9) default now64(9)
)
ENGINE=ReplacingMergeTree(updated_at)
ORDER BY (client_id,campaign_id,account_id,date_ref)
partition by toMonth(date_ref);



CREATE TABLE adszero.fetch_history (
    request_id FixedString(26) NOT NULL,
    client_id String NOT NULL,
    account_id String NOT NULL,
    business_id String NOT NULL,
    provider_id FixedString(26) NOT NULL,
    start_date_request Date32,
    end_date_request Date32,
    status Enum8('UNKNOWN'=0,'RUNNING'=1,'FAILED'=2,'SUCCESS'= 3) default 'UNKNOWN',
    error_message String,
    inserted_at Date32 default now(),
    updated_at DateTime64(9) default now64(9)
)
ENGINE=ReplacingMergeTree(updated_at)
ORDER BY (request_id,client_id,account_id,inserted_at);


CREATE TABLE adszero.client_rules (
    client_id String NOT NULL,
    rule_id FixedString(26) NOT NULL default generateULID(),
    rule_name String NOT NULL,
    column String NOT NULL,
    operator String NOT NULL,
    value Float64 NOT NULL,
    notification_way Enum8('EMAIL'=0,'TELEGRAM'=1,'SLACK'=2) default 'EMAIL',
    inserted_at DateTime64(9) default now64(9),
    updated_at DateTime64(9) default now64(9)
)
ENGINE=ReplacingMergeTree(updated_at)
ORDER BY (rule_id,client_id);
