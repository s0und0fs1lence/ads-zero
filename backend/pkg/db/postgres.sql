CREATE TYPE USERTYPE AS ENUM ('UNKNOWN', 'FREE', 'TIER1', 'TIER2', 'TIER3');
CREATE TYPE SUB_STATUS AS ENUM ('UNKNOWN','ACTIVE', 'INACTIVE');
CREATE TYPE PROVIDER_TYPE AS ENUM ('INVALID','FACEBOOK','GOOGLE','TIKTOK','TABOOLA');


CREATE TABLE clients (
    client_id VARCHAR(64) PRIMARY KEY,
    user_email VARCHAR(255) NOT NULL CHECK (user_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    user_type USERTYPE NOT NULL DEFAULT 'UNKNOWN',
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    stripe_subscription_status SUB_STATUS NOT NULL DEFAULT 'UNKNOWN' ,
    notification_email VARCHAR(255)  CHECK (user_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    telegram_chat_id TEXT,
    slack_webhook_url TEXT,
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE providers (
    provider_id varchar(26) PRIMARY KEY,
    provider_type PROVIDER_TYPE default 'INVALID',
    client_id varchar(64) NOT NULL,
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    api_client_id text,
    api_client_secret text,
    api_access_token text
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_clients_updated_at
BEFORE UPDATE ON clients
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();