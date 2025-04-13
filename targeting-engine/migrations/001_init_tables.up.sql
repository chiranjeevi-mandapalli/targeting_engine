-- Enable UUID extension if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Campaigns table
CREATE TABLE campaigns (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image_url TEXT NOT NULL,
    cta VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_url CHECK (image_url ~ '^https?://[^/]+')
);

-- Targeting rules table
CREATE TABLE targeting_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id VARCHAR(255) NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    dimension VARCHAR(50) NOT NULL CHECK (dimension IN ('app', 'country', 'os')),
    operation VARCHAR(50) NOT NULL CHECK (operation IN ('include', 'exclude')),
    values JSONB NOT NULL CHECK (jsonb_typeof(values) = 'array'),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT no_conflicting_rules CHECK (
        NOT EXISTS (
            SELECT 1 FROM targeting_rules tr 
            WHERE tr.campaign_id = targeting_rules.campaign_id 
            AND tr.dimension = targeting_rules.dimension 
            AND tr.operation != targeting_rules.operation
        )
    )
);

-- Optimize read queries
CREATE INDEX idx_campaigns_status ON campaigns(status);
CREATE INDEX idx_targeting_rules_campaign ON targeting_rules(campaign_id);
CREATE INDEX idx_targeting_rules_dimension ON targeting_rules(dimension);

-- Update timestamp trigger
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_campaign_timestamp
BEFORE UPDATE ON campaigns
FOR EACH ROW EXECUTE FUNCTION update_timestamp();