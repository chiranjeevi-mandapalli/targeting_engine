-- Insert sample campaigns
INSERT INTO campaigns (id, name, image_url, cta, status) VALUES
('spotify', 'Spotify - Music for everyone', 'https://somelink', 'Download', 'ACTIVE'),
('duolingo', 'Duolingo: Best way to learn', 'https://somelink2', 'Install', 'ACTIVE'),
('subwaysurfer', 'Subway Surfer', 'https://somelink3', 'Play', 'ACTIVE');

-- Insert targeting rules
INSERT INTO targeting_rules (campaign_id, dimension, operation, values) VALUES
-- Spotify: Include US and Canada
('spotify', 'country', 'include', '["US", "Canada"]'),

-- Duolingo: Include Android/iOS, Exclude US
('duolingo', 'os', 'include', '["Android", "iOS"]'),
('duolingo', 'country', 'exclude', '["US"]'),

-- Subway Surfer: Include Android and specific app
('subwaysurfer', 'os', 'include', '["Android"]'),
('subwaysurfer', 'app', 'include', '["com.gametion.ludokinggame"]');

-- Create a read-only user for the application
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'app_user') THEN
        CREATE ROLE app_user WITH LOGIN PASSWORD 'securepassword';
    END IF;
END
$$;

GRANT CONNECT ON DATABASE targeting TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO app_user;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;