CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS campaigns (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image_url TEXT NOT NULL,
    cta VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_url CHECK (image_url ~ '^https?://[^/]+')
);

CREATE INDEX IF NOT EXISTS idx_campaigns_status ON campaigns(status);

CREATE TABLE IF NOT EXISTS targeting_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id VARCHAR(255) NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    dimension VARCHAR(50) NOT NULL CHECK (dimension IN ('app', 'country', 'os')),
    operation VARCHAR(50) NOT NULL CHECK (operation IN ('include', 'exclude')),
    values JSONB NOT NULL CHECK (jsonb_typeof(values) = 'array'),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_targeting_rules_campaign ON targeting_rules(campaign_id);
CREATE INDEX IF NOT EXISTS idx_targeting_rules_dimension ON targeting_rules(dimension);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_campaign_timestamp
BEFORE UPDATE ON campaigns
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE OR REPLACE FUNCTION enforce_no_conflicting_rules()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM targeting_rules
        WHERE campaign_id = NEW.campaign_id
          AND dimension = NEW.dimension
          AND operation != NEW.operation
    ) THEN
        RAISE EXCEPTION 'Conflicting rule exists for campaign_id %, dimension %', NEW.campaign_id, NEW.dimension;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_no_conflicting_rules
BEFORE INSERT OR UPDATE ON targeting_rules
FOR EACH ROW
EXECUTE FUNCTION enforce_no_conflicting_rules();
